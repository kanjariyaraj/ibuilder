package cmd

import (
	"testing"
)

func TestSigningHelpCommand(t *testing.T) {
	output, err := executeCommand("signing", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestSigningSetupHelpCommand(t *testing.T) {
	output, err := executeCommand("signing", "setup", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestSigningValidateHelpCommand(t *testing.T) {
	output, err := executeCommand("signing", "validate", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestSigningDoctorHelpCommand(t *testing.T) {
	output, err := executeCommand("signing", "doctor", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestSigningStatusHelpCommand(t *testing.T) {
	output, err := executeCommand("signing", "status", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestSigningRemoveHelpCommand(t *testing.T) {
	output, err := executeCommand("signing", "remove", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestSigningAllSubCommandsHelp(t *testing.T) {
	subs := []string{"setup", "validate", "doctor", "status", "remove"}
	for _, sub := range subs {
		output, err := executeCommand("signing", sub, "--help")
		if err != nil {
			t.Fatalf("signing %s --help failed: %v", sub, err)
		}
		if len(output) == 0 {
			t.Errorf("signing %s --help produced no output", sub)
		}
	}
}
