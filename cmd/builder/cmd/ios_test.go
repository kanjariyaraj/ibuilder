package cmd

import (
	"testing"
)

func TestIosHelpCommand(t *testing.T) {
	output, err := executeCommand("ios", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestIosBuildHelpCommand(t *testing.T) {
	output, err := executeCommand("ios", "build", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestIosBuildMissingWorkflow(t *testing.T) {
	output, err := executeCommand("ios", "build")
	if err == nil && len(output) > 0 {
		t.Logf("got output (may be auth): %s", output)
	}
}

func TestIosBuildAllFlagsHelp(t *testing.T) {
	output, err := executeCommand("ios", "build", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}
