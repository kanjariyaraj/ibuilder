package artifacts

import (
	"os"
	"testing"
	"time"
)

func TestStorage_EnsureDirs(t *testing.T) {
	s := NewStorage()
	s.BaseDir = t.TempDir()
	if err := s.EnsureDirs(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	dirs := []string{s.ArtifactDir(), s.LogDir(), s.ReportDir(), s.CacheDir(), s.MetadataDir()}
	for _, d := range dirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			t.Errorf("expected %s to exist", d)
		}
	}
}

func TestStorage_ArtifactDir(t *testing.T) {
	s := NewStorage()
	s.BaseDir = "/tmp/test"
	if s.ArtifactDir() != "/tmp/test/artifacts" {
		t.Errorf("unexpected artifact dir: %s", s.ArtifactDir())
	}
}

func TestStorage_LogDir(t *testing.T) {
	s := NewStorage()
	s.BaseDir = "/tmp/test"
	if s.LogDir() != "/tmp/test/logs" {
		t.Errorf("unexpected log dir: %s", s.LogDir())
	}
}

func TestStorage_ReportDir(t *testing.T) {
	s := NewStorage()
	s.BaseDir = "/tmp/test"
	if s.ReportDir() != "/tmp/test/reports" {
		t.Errorf("unexpected report dir: %s", s.ReportDir())
	}
}

func TestStorage_CacheDir(t *testing.T) {
	s := NewStorage()
	s.BaseDir = "/tmp/test"
	if s.CacheDir() != "/tmp/test/cache" {
		t.Errorf("unexpected cache dir: %s", s.CacheDir())
	}
}

func TestStorage_MetadataDir(t *testing.T) {
	s := NewStorage()
	s.BaseDir = "/tmp/test"
	if s.MetadataDir() != "/tmp/test/metadata" {
		t.Errorf("unexpected metadata dir: %s", s.MetadataDir())
	}
}

func TestStorage_CleanAll(t *testing.T) {
	s := NewStorage()
	s.BaseDir = t.TempDir()
	if err := s.EnsureDirs(); err != nil {
		t.Fatal(err)
	}
	if err := s.CleanAll(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestStorage_SaveAndGetArtifactMetadata(t *testing.T) {
	s := NewStorage()
	s.BaseDir = t.TempDir()
	s.EnsureDirs()

	if err := s.SaveArtifactMetadata(1, "test.zip", "/tmp/test.zip", 1024, "abc123"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	meta, err := s.GetArtifactMetadata(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if meta.Name != "test.zip" {
		t.Errorf("expected test.zip, got %s", meta.Name)
	}
	if meta.Size != 1024 {
		t.Errorf("expected 1024, got %d", meta.Size)
	}
	if meta.Checksum != "abc123" {
		t.Errorf("expected abc123, got %s", meta.Checksum)
	}
}

func TestStorage_ListArtifactMetadata(t *testing.T) {
	s := NewStorage()
	s.BaseDir = t.TempDir()
	s.EnsureDirs()

	s.SaveArtifactMetadata(1, "a.zip", "/tmp/a.zip", 100, "hash1")
	s.SaveArtifactMetadata(2, "b.zip", "/tmp/b.zip", 200, "hash2")

	metas, err := s.ListArtifactMetadata()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(metas) != 2 {
		t.Errorf("expected 2 metadata entries, got %d", len(metas))
	}
}

func TestCleanOldArtifacts(t *testing.T) {
	s := NewStorage()
	s.BaseDir = t.TempDir()
	s.EnsureDirs()

	os.WriteFile(s.ArtifactDir()+"/old.zip", []byte("old"), 0644)
	time.Sleep(10 * time.Millisecond)

	count, err := s.CleanOldArtifacts(1 * time.Nanosecond)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 removed, got %d", count)
	}
}

func TestDurationParsing(t *testing.T) {
	d, err := time.ParseDuration("72h")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if d != 72*time.Hour {
		t.Errorf("expected 72h, got %v", d)
	}
}

func TestDurationParsing_Empty(t *testing.T) {
	_, err := time.ParseDuration("")
	if err == nil {
		t.Error("expected error for empty duration")
	}
}
