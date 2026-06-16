package cmd

import (
	"testing"
)

func TestArtifactHelpCommand(t *testing.T) {
	output, err := executeCommand("artifact", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestArtifactSubHelpCommands(t *testing.T) {
	subs := []string{"list", "download", "inspect", "latest", "clean"}
	for _, sub := range subs {
		output, err := executeCommand("artifact", sub, "--help")
		if err != nil {
			t.Fatalf("artifact %s --help failed: %v", sub, err)
		}
		if len(output) == 0 {
			t.Errorf("artifact %s --help produced no output", sub)
		}
	}
}

func TestArtifactListValidFlags(t *testing.T) {
	output, err := executeCommand("artifact", "list", "--limit", "10", "--json")
	if err == nil {
		_ = output
	}
}

func TestArtifactDownloadValidFlags(t *testing.T) {
	output, err := executeCommand("artifact", "download", "--dest", "dist", "--overwrite")
	if err == nil {
		_ = output
	}
}

func TestArtifactLatestValidFlags(t *testing.T) {
	output, err := executeCommand("artifact", "latest")
	if err == nil {
		_ = output
	}
}

func TestArtifactCleanValidFlags(t *testing.T) {
	output, err := executeCommand("artifact", "clean", "--all")
	if err == nil {
		_ = output
	}
}
