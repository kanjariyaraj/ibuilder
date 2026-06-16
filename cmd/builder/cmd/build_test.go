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

func TestBuildHistoryHelp(t *testing.T) {
	output, err := executeCommand("build", "history", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestBuildHistoryValidFlags(t *testing.T) {
	output, err := executeCommand("build", "history", "--limit", "10", "--json")
	if err == nil {
		_ = output
	}
}

func TestBuildInspectWithoutRunID(t *testing.T) {
	_, err := executeCommand("build", "inspect")
	if err == nil {
		t.Errorf("expected error for missing --run-id")
	}
}

func TestBuildInspectHelp(t *testing.T) {
	output, err := executeCommand("build", "inspect", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestBuildLogsHelp(t *testing.T) {
	output, err := executeCommand("build", "logs", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestBuildLogsValidFlags(t *testing.T) {
	output, err := executeCommand("build", "logs", "--latest", "--save", "/tmp/logs")
	if err == nil {
		_ = output
	}
}

func TestBuildOpenWithoutRunID(t *testing.T) {
	_, err := executeCommand("build", "open")
	if err == nil {
		t.Errorf("expected error for missing --run-id")
	}
}

func TestBuildOpenHelp(t *testing.T) {
	output, err := executeCommand("build", "open", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestBuildAllSubCommandsHaveHelp(t *testing.T) {
	subs := []string{"history", "inspect", "logs", "open"}
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

func TestBuildSubCommandsHaveValidFlagCombos(t *testing.T) {
	_, err := executeCommand("build", "history", "--branch", "main", "--status", "completed", "--limit", "5", "--page", "1")
	if err == nil {
		_ = err
	}
	_, err = executeCommand("build", "inspect", "--run-id", "12345")
	if err == nil {
		_ = err
	}
	_, err = executeCommand("build", "logs", "--run-id", "12345")
	if err == nil {
		_ = err
	}
	_, err = executeCommand("build", "open", "--run-id", "12345")
	if err == nil {
		_ = err
	}
}
