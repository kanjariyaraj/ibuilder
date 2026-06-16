package artifacts

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	BaseDir string
}

func NewStorage() *Storage {
	return &Storage{BaseDir: filepath.Join(".build")}
}

func (s *Storage) EnsureDirs() error {
	dirs := []string{
		filepath.Join(s.BaseDir, "artifacts"),
		filepath.Join(s.BaseDir, "logs"),
		filepath.Join(s.BaseDir, "reports"),
		filepath.Join(s.BaseDir, "cache"),
		filepath.Join(s.BaseDir, "metadata"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) ArtifactDir() string {
	return filepath.Join(s.BaseDir, "artifacts")
}

func (s *Storage) LogDir() string {
	return filepath.Join(s.BaseDir, "logs")
}

func (s *Storage) ReportDir() string {
	return filepath.Join(s.BaseDir, "reports")
}

func (s *Storage) CacheDir() string {
	return filepath.Join(s.BaseDir, "cache")
}

func (s *Storage) MetadataDir() string {
	return filepath.Join(s.BaseDir, "metadata")
}

func (s *Storage) CleanOldArtifacts(olderThan time.Duration) (int, error) {
	dir := s.ArtifactDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	cutoff := time.Now().Add(-olderThan)
	count := 0
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			path := filepath.Join(dir, e.Name())
			if err := os.RemoveAll(path); err != nil {
				fmt.Printf("Warning: failed to remove %s: %v\n", path, err)
				continue
			}
			count++
		}
	}
	return count, nil
}

func (s *Storage) CleanAll() error {
	dirs := []string{s.ArtifactDir(), s.LogDir(), s.ReportDir(), s.CacheDir(), s.MetadataDir()}
	for _, d := range dirs {
		if err := os.RemoveAll(d); err != nil {
			return err
		}
	}
	return s.EnsureDirs()
}
