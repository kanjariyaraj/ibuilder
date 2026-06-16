package artifacts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ArtifactMeta struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	LocalPath string `json:"local_path"`
	Size      int64  `json:"size"`
	Checksum  string `json:"checksum"`
	DownloadedAt string `json:"downloaded_at"`
}

type BuildMeta struct {
	RunID      int64  `json:"run_id"`
	RunNumber  int    `json:"run_number"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	Branch     string `json:"branch"`
	CheckedAt  string `json:"checked_at"`
}

func (s *Storage) SaveArtifactMetadata(id int64, name, localPath string, size int64, checksum string) error {
	meta := ArtifactMeta{
		ID:        id,
		Name:      name,
		LocalPath: localPath,
		Size:      size,
		Checksum:  checksum,
		DownloadedAt: time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(s.MetadataDir(), fmt.Sprintf("artifact-%d.json", id))
	return os.WriteFile(path, data, 0644)
}

func (s *Storage) GetArtifactMetadata(id int64) (*ArtifactMeta, error) {
	path := filepath.Join(s.MetadataDir(), fmt.Sprintf("artifact-%d.json", id))
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var meta ArtifactMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

func (s *Storage) ListArtifactMetadata() ([]ArtifactMeta, error) {
	dir := s.MetadataDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var metas []ArtifactMeta
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var meta ArtifactMeta
		if err := json.Unmarshal(data, &meta); err != nil {
			continue
		}
		metas = append(metas, meta)
	}
	return metas, nil
}

func (s *Storage) SaveBuildHistory(records []BuildRecord) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(s.MetadataDir(), "build-history.json")
	return os.WriteFile(path, data, 0644)
}
