package cmd

import (
	"testing"
)

func TestInitHelpCommand(t *testing.T) {
	output, err := executeCommand("init", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestInitRunHelpCommand(t *testing.T) {
	output, err := executeCommand("init", "run", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestInitRunDryRun(t *testing.T) {
	output, err := executeCommand("init", "run", "--dry-run")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestInitAllFlags(t *testing.T) {
	flags := []string{"--force", "--dry-run", "--yes", "--json"}
	for _, f := range flags {
		output, err := executeCommand("init", "run", f)
		if err != nil {
			t.Fatalf("init run %s failed: %v", f, err)
		}
		if len(output) == 0 {
			t.Errorf("init run %s produced no output", f)
		}
	}
}
