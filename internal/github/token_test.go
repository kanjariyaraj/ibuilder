package github

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadToken(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	token := &TokenData{
		AccessToken: "gho_test123",
		TokenType:   "bearer",
		Scope:       "repo,workflow",
	}

	if err := SaveToken(token); err != nil {
		t.Fatalf("failed to save token: %v", err)
	}

	loaded, err := LoadToken()
	if err != nil {
		t.Fatalf("failed to load token: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected token, got nil")
	}
	if loaded.AccessToken != "gho_test123" {
		t.Errorf("expected 'gho_test123', got '%s'", loaded.AccessToken)
	}
	if loaded.Scope != "repo,workflow" {
		t.Errorf("expected 'repo,workflow', got '%s'", loaded.Scope)
	}
}

func TestLoadTokenNotExist(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	token, err := LoadToken()
	if err != nil {
		t.Fatalf("expected no error for missing token, got %v", err)
	}
	if token != nil {
		t.Errorf("expected nil token, got %v", token)
	}
}

func TestDeleteToken(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	token := &TokenData{AccessToken: "gho_test"}
	SaveToken(token)

	if err := DeleteToken(); err != nil {
		t.Fatalf("failed to delete token: %v", err)
	}

	loaded, _ := LoadToken()
	if loaded != nil {
		t.Errorf("expected nil after delete")
	}
}

func TestDeleteTokenNotExist(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	if err := DeleteToken(); err != nil {
		t.Fatalf("expected no error deleting non-existent token, got %v", err)
	}
}

func TestTokenExists(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	if TokenExists() {
		t.Errorf("expected false when no token exists")
	}

	SaveToken(&TokenData{AccessToken: "gho_test"})
	if !TokenExists() {
		t.Errorf("expected true when token exists")
	}
}

func TestTokenPermissions(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	SaveToken(&TokenData{AccessToken: "gho_test"})

	dir, _ := tokenDir()
	path := filepath.Join(dir, TokenFile)

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat token file: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("expected 0600 permissions, got %o", perm)
	}
}
