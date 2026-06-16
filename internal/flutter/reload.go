package flutter

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type ReloadResult struct {
	Success  bool          `json:"success"`
	Duration time.Duration `json:"duration_ms"`
	Output   string        `json:"output"`
	Error    string        `json:"error,omitempty"`
}

func (s *Session) HotReload() (*ReloadResult, error) {
	s.mu.RLock()
	if s.state != SessionActive && s.state != SessionAttached && s.state != SessionWatchMode {
		s.mu.RUnlock()
		return nil, fmt.Errorf("no active flutter session — start one with 'flutter dev' or 'flutter attach'")
	}
	s.mu.RUnlock()

	if !s.cfg.HotReload {
		return nil, fmt.Errorf("hot reload is disabled in configuration")
	}

	s.log.Info("triggering hot reload")
	start := time.Now()

	cmd := exec.Command("flutter", "reload")
	cmd.Dir = s.projectDir
	output, err := cmd.CombinedOutput()
	duration := time.Since(start)

	if err != nil {
		result := &ReloadResult{
			Success:  false,
			Duration: duration,
			Output:   string(output),
			Error:    err.Error(),
		}
		s.log.Error("hot reload failed", "error", err, "output", strings.TrimSpace(string(output)))
		return result, fmt.Errorf("hot reload failed: %s", strings.TrimSpace(string(output)))
	}

	result := &ReloadResult{
		Success:  true,
		Duration: duration,
		Output:   strings.TrimSpace(string(output)),
	}

	s.log.Info("hot reload succeeded", "duration", duration)
	return result, nil
}

func (s *Session) HotRestart() (*ReloadResult, error) {
	s.mu.RLock()
	if s.state != SessionActive && s.state != SessionAttached && s.state != SessionWatchMode {
		s.mu.RUnlock()
		return nil, fmt.Errorf("no active flutter session — start one with 'flutter dev' or 'flutter attach'")
	}
	s.mu.RUnlock()

	s.log.Info("triggering hot restart")
	start := time.Now()

	cmd := exec.Command("flutter", "restart")
	cmd.Dir = s.projectDir
	output, err := cmd.CombinedOutput()
	duration := time.Since(start)

	if err != nil {
		result := &ReloadResult{
			Success:  false,
			Duration: duration,
			Output:   string(output),
			Error:    err.Error(),
		}
		s.log.Error("hot restart failed", "error", err)
		return result, fmt.Errorf("hot restart failed: %s", strings.TrimSpace(string(output)))
	}

	result := &ReloadResult{
		Success:  true,
		Duration: duration,
		Output:   strings.TrimSpace(string(output)),
	}

	s.log.Info("hot restart succeeded", "duration", duration)
	return result, nil
}
