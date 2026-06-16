package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.ProjectName != "Builder" {
		t.Errorf("expected 'Builder', got '%s'", cfg.ProjectName)
	}
	if cfg.IOS.MinimumVersion != "15.0" {
		t.Errorf("expected '15.0', got '%s'", cfg.IOS.MinimumVersion)
	}
	if len(cfg.IOS.Devices) != 2 {
		t.Errorf("expected 2 devices, got %d", len(cfg.IOS.Devices))
	}
	if cfg.Flutter.Channel != "stable" {
		t.Errorf("expected 'stable', got '%s'", cfg.Flutter.Channel)
	}
	if cfg.ReactNative.Entry != "index.js" {
		t.Errorf("expected 'index.js', got '%s'", cfg.ReactNative.Entry)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "builder.json")

	cfg := Default()
	cfg.ProjectName = "TestProject"
	cfg.Repository = "https://example.com/repo.git"

	if err := Save(path, cfg); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.ProjectName != "TestProject" {
		t.Errorf("expected 'TestProject', got '%s'", loaded.ProjectName)
	}
	if loaded.Repository != "https://example.com/repo.git" {
		t.Errorf("expected 'https://example.com/repo.git', got '%s'", loaded.Repository)
	}
}

func TestLoadNonExistent(t *testing.T) {
	cfg, err := Load("/nonexistent/path/builder.json")
	if err != nil {
		t.Fatalf("expected no error for non-existent file, got %v", err)
	}
	if cfg.ProjectName != "Builder" {
		t.Errorf("expected defaults, got '%s'", cfg.ProjectName)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "builder.json")
	os.WriteFile(path, []byte("{invalid json}"), 0644)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestValidate(t *testing.T) {
	cfg := Default()
	errs := Validate(cfg)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors for default config, got %d", len(errs))
	}
}

func TestValidateEmptyProjectName(t *testing.T) {
	cfg := Default()
	cfg.ProjectName = ""
	errs := Validate(cfg)
	hasProjectError := false
	for _, e := range errs {
		if e.Error() == "validation: project_name is required" {
			hasProjectError = true
		}
	}
	if !hasProjectError {
		t.Errorf("expected validation error for empty project name")
	}
}

func TestValidateEmptyDevices(t *testing.T) {
	cfg := Default()
	cfg.IOS.Devices = []string{}
	errs := Validate(cfg)
	hasDeviceError := false
	for _, e := range errs {
		if e.Error() == "validation: ios.devices must have at least one device" {
			hasDeviceError = true
		}
	}
	if !hasDeviceError {
		t.Errorf("expected validation error for empty devices")
	}
}
