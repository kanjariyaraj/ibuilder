package github

import (
	"os"
	"testing"
)

func TestConstants(t *testing.T) {
	if DeviceCodeURL != "https://github.com/login/device/code" {
		t.Errorf("unexpected DeviceCodeURL")
	}
	if AccessTokenURL != "https://github.com/login/oauth/access_token" {
		t.Errorf("unexpected AccessTokenURL")
	}
	if APIRoot != "https://api.github.com" {
		t.Errorf("unexpected APIRoot")
	}
	if DefaultClientID != "BuilderCLI" {
		t.Errorf("unexpected DefaultClientID")
	}
}

func TestTokenDir(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	dir, err := tokenDir()
	if err != nil {
		t.Fatalf("failed to get token dir: %v", err)
	}
	if dir != tmpHome+"/.builder" {
		t.Errorf("expected %s/.builder, got %s", tmpHome, dir)
	}
}

func TestParseRemoteURLHTTPS(t *testing.T) {
	remote, err := parseRemoteURL("https://github.com/kanjariyaraj/Builder.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if remote.Owner != "kanjariyaraj" {
		t.Errorf("expected 'kanjariyaraj', got '%s'", remote.Owner)
	}
	if remote.Name != "Builder" {
		t.Errorf("expected 'Builder', got '%s'", remote.Name)
	}
}

func TestParseRemoteURLSSH(t *testing.T) {
	remote, err := parseRemoteURL("git@github.com:kanjariyaraj/Builder.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if remote.Owner != "kanjariyaraj" {
		t.Errorf("expected 'kanjariyaraj', got '%s'", remote.Owner)
	}
	if remote.Name != "Builder" {
		t.Errorf("expected 'Builder', got '%s'", remote.Name)
	}
}

func TestParseRemoteURLInvalid(t *testing.T) {
	_, err := parseRemoteURL("invalid-url")
	if err == nil {
		t.Errorf("expected error for invalid URL")
	}
}
