package cmd

import (
	"testing"
)

func TestBuildHelpCommand(t *testing.T) {
	output, err := executeCommand("build", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestBuildRunWithoutArgs(t *testing.T) {
	_, err := executeCommand("build", "run")
	if err == nil {
		t.Errorf("expected error for missing args")
	}
}

func TestBuildStatusWithoutArgs(t *testing.T) {
	_, err := executeCommand("build", "status")
	if err == nil {
		t.Errorf("expected error for missing args")
	}
}

func TestBuildListHelp(t *testing.T) {
	output, err := executeCommand("build", "list", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestBuildLogWithoutArgs(t *testing.T) {
	_, err := executeCommand("build", "log")
	if err == nil {
		t.Errorf("expected error for missing args")
	}
}

func TestBuildArtifactsWithoutArgs(t *testing.T) {
	_, err := executeCommand("build", "artifacts")
	if err == nil {
		t.Errorf("expected error for missing args")
	}
}

func TestBuildAllSubCommandsHaveHelp(t *testing.T) {
	subs := []string{"run", "status", "list", "log", "artifacts"}
	for _, sub := range subs {
		output, err := executeCommand("build", sub, "--help")
		if err != nil {
			t.Fatalf("build %s --help failed: %v", sub, err)
		}
		if len(output) == 0 {
			t.Errorf("build %s --help produced no output", sub)
		}
	}
}
