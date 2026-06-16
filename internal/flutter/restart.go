package flutter

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type RestartResult struct {
	Success  bool          `json:"success"`
	Duration time.Duration `json:"duration_ms"`
	Output   string        `json:"output"`
	Error    string        `json:"error,omitempty"`
}

func (s *Session) Restart() (*RestartResult, error) {
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
		s.log.Error("hot restart failed", "error", err)
		return &RestartResult{
			Success:  false,
			Duration: duration,
			Output:   strings.TrimSpace(string(output)),
			Error:    err.Error(),
		}, fmt.Errorf("hot restart failed: %s", strings.TrimSpace(string(output)))
	}

	result := &RestartResult{
		Success:  true,
		Duration: duration,
		Output:   strings.TrimSpace(string(output)),
	}

	s.log.Info("hot restart succeeded", "duration", duration)
	return result, nil
}
