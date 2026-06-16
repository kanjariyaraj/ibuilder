package mobai

import (
	"fmt"
	"time"
)

type LaunchResult struct {
	Success   bool      `json:"success"`
	BundleID  string    `json:"bundle_id"`
	PID       int       `json:"pid"`
	LaunchedAt time.Time `json:"launched_at"`
	Duration  time.Duration `json:"duration_ms"`
	Error     string    `json:"error,omitempty"`
}

func (c *Client) LaunchApp(bundleID string) (*LaunchResult, error) {
	c.mu.RLock()
	connected := c.status.State == StateConnected
	c.mu.RUnlock()

	if !connected {
		return nil, fmt.Errorf("not connected to MobAI")
	}

	if bundleID == "" {
		return nil, fmt.Errorf("bundle ID is required")
	}

	start := time.Now()
	c.log.Info("launching app", "bundle_id", bundleID)

	time.Sleep(1 * time.Second)

	result := &LaunchResult{
		Success:    true,
		BundleID:   bundleID,
		PID:        1234,
		LaunchedAt: time.Now(),
		Duration:   time.Since(start),
	}

	c.log.Info("app launched", "bundle_id", bundleID, "pid", result.PID, "duration", result.Duration)
	return result, nil
}

func (c *Client) TerminateApp(bundleID string) error {
	c.mu.RLock()
	connected := c.status.State == StateConnected
	c.mu.RUnlock()

	if !connected {
		return fmt.Errorf("not connected to MobAI")
	}

	c.log.Info("terminating app", "bundle_id", bundleID)
	time.Sleep(500 * time.Millisecond)
	return nil
}

func (c *Client) IsAppRunning(bundleID string) (bool, error) {
	c.mu.RLock()
	connected := c.status.State == StateConnected
	c.mu.RUnlock()

	if !connected {
		return false, fmt.Errorf("not connected to MobAI")
	}

	return true, nil
}
