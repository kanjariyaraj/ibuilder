package artifacts

import (
	"testing"
)

func TestArtifactStruct(t *testing.T) {
	a := Artifact{ID: 1, Name: "test.zip", Size: 1024}
	if a.ID != 1 {
		t.Errorf("expected 1, got %d", a.ID)
	}
	if a.Name != "test.zip" {
		t.Errorf("expected test.zip, got %s", a.Name)
	}
	if a.Size != 1024 {
		t.Errorf("expected 1024, got %d", a.Size)
	}
}

func TestBuildRecordStruct(t *testing.T) {
	r := BuildRecord{RunID: 1, RunNumber: 42, Status: "completed"}
	if r.RunNumber != 42 {
		t.Errorf("expected 42, got %d", r.RunNumber)
	}
	if r.Status != "completed" {
		t.Errorf("expected completed, got %s", r.Status)
	}
}

func TestBuildInspectStruct(t *testing.T) {
	b := BuildInspect{RunID: 1, Status: "completed", Branch: "main"}
	if b.Status != "completed" {
		t.Errorf("expected completed, got %s", b.Status)
	}
}

func TestJobInfoStruct(t *testing.T) {
	j := JobInfo{ID: 1, Name: "build", Status: "completed"}
	if j.Name != "build" {
		t.Errorf("expected build, got %s", j.Name)
	}
}

func TestCleanupResultStruct(t *testing.T) {
	r := CleanupResult{FilesRemoved: 5}
	if r.FilesRemoved != 5 {
		t.Errorf("expected 5, got %d", r.FilesRemoved)
	}
}

func TestLogsOptionsDefaults(t *testing.T) {
	o := &LogsOptions{}
	if o.Latest {
		t.Error("expected Latest to be false")
	}
}

func TestDownloadOptionsDefaults(t *testing.T) {
	o := &DownloadOptions{}
	if o.Overwrite {
		t.Error("expected Overwrite to be false")
	}
}

func TestHistoryOptionsDefaults(t *testing.T) {
	o := &HistoryOptions{}
	if o.Limit != 0 {
		t.Errorf("expected 0, got %d", o.Limit)
	}
}

func TestCleanupOptionsDefaults(t *testing.T) {
	o := &CleanupOptions{}
	if o.All {
		t.Error("expected All to be false")
	}
}
