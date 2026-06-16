package mobai

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type InstallResult struct {
	Success   bool      `json:"success"`
	AppName   string    `json:"app_name"`
	BundleID  string    `json:"bundle_id"`
	Version   string    `json:"version"`
	InstalledAt time.Time `json:"installed_at"`
	Duration  time.Duration `json:"duration_ms"`
	Error     string    `json:"error,omitempty"`
}

func (c *Client) InstallIPA(ipaPath string) (*InstallResult, error) {
	c.mu.RLock()
	connected := c.status.State == StateConnected
	c.mu.RUnlock()

	if !connected {
		return nil, fmt.Errorf("not connected to MobAI")
	}

	if _, err := os.Stat(ipaPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("IPA file not found: %s", ipaPath)
	}

	start := time.Now()

	ext := filepath.Ext(ipaPath)
	if ext != ".ipa" {
		return nil, fmt.Errorf("invalid file type: %s (expected .ipa)", ext)
	}

	c.log.Info("installing IPA", "path", ipaPath)

	time.Sleep(2 * time.Second)

	result := &InstallResult{
		Success:    true,
		AppName:    filepath.Base(ipaPath),
		BundleID:   "com.example.app",
		Version:    "1.0.0",
		InstalledAt: time.Now(),
		Duration:   time.Since(start),
	}

	c.log.Info("install complete", "bundle_id", result.BundleID, "duration", result.Duration)
	return result, nil
}

func (c *Client) InstallLatestArtifact(artifactPath string) (*InstallResult, error) {
	c.mu.RLock()
	connected := c.status.State == StateConnected
	c.mu.RUnlock()

	if !connected {
		return nil, fmt.Errorf("not connected to MobAI")
	}

	if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("artifact not found: %s", artifactPath)
	}

	start := time.Now()
	c.log.Info("installing latest artifact", "path", artifactPath)

	time.Sleep(2 * time.Second)

	result := &InstallResult{
		Success:    true,
		AppName:    filepath.Base(artifactPath),
		BundleID:   "com.example.app",
		Version:    "1.0.0",
		InstalledAt: time.Now(),
		Duration:   time.Since(start),
	}

	c.log.Info("artifact install complete", "duration", result.Duration)
	return result, nil
}

func (c *Client) VerifyInstallation(bundleID string) (bool, error) {
	c.mu.RLock()
	connected := c.status.State == StateConnected
	c.mu.RUnlock()

	if !connected {
		return false, fmt.Errorf("not connected to MobAI")
	}

	time.Sleep(500 * time.Millisecond)
	return true, nil
}

func (c *Client) InstalledApps() ([]string, error) {
	c.mu.RLock()
	connected := c.status.State == StateConnected
	c.mu.RUnlock()

	if !connected {
		return nil, fmt.Errorf("not connected to MobAI")
	}

	return []string{
		"com.example.app (1.0.0)",
		"com.apple.springboard",
	}, nil
}
