package build

import (
	"os"
	"testing"
)

func TestBuildReportStruct(t *testing.T) {
	r := BuildReport{RunID: 1, Status: "completed", Conclusion: "success"}
	if r.RunID != 1 {
		t.Errorf("expected RunID 1")
	}
}

func TestBuildResultStruct(t *testing.T) {
	r := BuildResult{RunID: 42, RunNumber: 5, Status: "completed"}
	if r.RunNumber != 5 {
		t.Errorf("expected RunNumber 5")
	}
}

func TestValidatePrereqsNoToken(t *testing.T) {
	opts := &BuildOptions{Token: ""}
	err := validatePrereqs(opts)
	if err == nil {
		t.Errorf("expected error for missing token")
	}
}

func TestValidatePrereqsNoRepo(t *testing.T) {
	opts := &BuildOptions{Token: "test", Owner: "", Name: ""}
	err := validatePrereqs(opts)
	if err == nil {
		t.Errorf("expected error for missing repo")
	}
}

func TestValidatePrereqsNoWorkflow(t *testing.T) {
	opts := &BuildOptions{Token: "test", Owner: "o", Name: "r", WorkflowID: ""}
	err := validatePrereqs(opts)
	if err == nil {
		t.Errorf("expected error for missing workflow")
	}
}

func TestValidatePrereqsValid(t *testing.T) {
	opts := &BuildOptions{Token: "test", Owner: "o", Name: "r", WorkflowID: "ci.yml"}
	err := validatePrereqs(opts)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGenerateReport(t *testing.T) {
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	result := &BuildResult{
		RunID:      1,
		RunNumber:  42,
		Status:     "completed",
		Conclusion: "success",
		Artifact:   "dist/build.zip",
	}
	path := generateReport(result)
	if path == "" {
		t.Errorf("expected report path")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("report file not created: %s", path)
	}
}
