package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := rootCmd.Execute()

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = original

	fmt.Fprint(buf, string(out))
	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	output, err := executeCommand()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected output, got empty string")
	}
}

func TestVersionCommand(t *testing.T) {
	output, err := executeCommand("version")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected version output, got empty string")
	}
}

func TestDoctorCommand(t *testing.T) {
	output, err := executeCommand("doctor")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected doctor output, got empty string")
	}
}

func TestHelpCommand(t *testing.T) {
	output, err := executeCommand("--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output) == 0 {
		t.Errorf("expected help output, got empty string")
	}
}
