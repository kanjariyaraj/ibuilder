package artifacts

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kanjariyaraj/Builder/internal/errors"
)

type DownloadOptions struct {
	ArtifactID int64
	Name       string
	Build      int
	Latest     bool
	DestDir    string
	Overwrite  bool
	Resume     bool
}

type DownloadResult struct {
	Path     string
	Size     int64
	Checksum string
}

func (m *ArtifactManager) DownloadArtifact(opts *DownloadOptions) (*DownloadResult, error) {
	artifacts, err := m.ListAllArtifacts()
	if err != nil {
		return nil, err
	}

	var target *Artifact
	for _, a := range artifacts {
		if opts.ArtifactID > 0 && a.ID == opts.ArtifactID {
			target = &a
			break
		}
		if opts.Name != "" && a.Name == opts.Name {
			target = &a
			break
		}
		if opts.Build > 0 && a.BuildNumber == opts.Build {
			target = &a
			break
		}
	}

	if opts.Latest && len(artifacts) > 0 {
		target = &artifacts[0]
	}

	if target == nil {
		return nil, errors.New(errors.KindNotFound, "no matching artifact found")
	}

	destDir := opts.DestDir
	if destDir == "" {
		destDir = "dist"
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, errors.Wrap(errors.KindInternal, "failed to create destination directory", err)
	}

	fileName := target.Name + ".zip"
	destPath := filepath.Join(destDir, fileName)

	if _, err := os.Stat(destPath); err == nil && !opts.Overwrite {
		return nil, errors.New(errors.KindInternal, fmt.Sprintf("file %s already exists (use --overwrite)", destPath))
	}

	fmt.Println("Downloading artifact", target.Name)

	req, err := http.NewRequest("GET", target.DownloadURL, nil)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "failed to create download request", err)
	}
	req.Header.Set("Authorization", "Bearer "+m.client.AccessToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(errors.KindNetwork, "download failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(errors.KindNetwork, fmt.Sprintf("download returned status %d", resp.StatusCode))
	}

	outFile, err := os.Create(destPath)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "failed to create output file", err)
	}
	defer outFile.Close()

	written, err := io.Copy(outFile, resp.Body)
	if err != nil {
		return nil, errors.Wrap(errors.KindNetwork, "download incomplete", err)
	}

	hash := sha256.New()
	outFile.Seek(0, 0)
	io.Copy(hash, outFile)
	sum := fmt.Sprintf("%x", hash.Sum(nil))

	fmt.Printf("Downloaded %s (%d bytes)\n", target.Name, written)

	m.storage.SaveArtifactMetadata(target.ID, target.Name, destPath, written, sum)

	return &DownloadResult{
		Path:     destPath,
		Size:     written,
		Checksum: sum,
	}, nil
}

func (m *ArtifactManager) DownloadLatest() (*DownloadResult, error) {
	return m.DownloadArtifact(&DownloadOptions{Latest: true, DestDir: "dist", Overwrite: true})
}
