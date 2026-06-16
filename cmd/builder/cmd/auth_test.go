package cmd

import (
	"testing"
)

func TestAuthHelpCommand(t *testing.T) {
	output, err := executeCommand("auth", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestAuthStatusWithoutToken(t *testing.T) {
	output, err := executeCommand("auth", "status")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestAuthLogoutWithoutToken(t *testing.T) {
	output, err := executeCommand("auth", "logout")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestRepoHelpCommand(t *testing.T) {
	output, err := executeCommand("repo", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestRepoInfoWithoutConfig(t *testing.T) {
	output, err := executeCommand("repo", "info")
	if err == nil && len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestRepoValidateWithoutAuth(t *testing.T) {
	output, err := executeCommand("repo", "validate")
	if err == nil && len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestRootWithAllSubCommands(t *testing.T) {
	cmds := []string{"version", "doctor", "config", "auth", "repo"}
	for _, c := range cmds {
		output, err := executeCommand(c, "--help")
		if err != nil {
			t.Fatalf("command '%s --help' failed: %v", c, err)
		}
		if len(output) == 0 {
			t.Errorf("command '%s --help' produced no output", c)
		}
	}
}
