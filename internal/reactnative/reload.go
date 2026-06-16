package reactnative

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ReloadResult struct {
	Success    bool          `json:"success"`
	Type       string        `json:"type"`
	Duration   time.Duration `json:"duration_ms"`
	Output     string        `json:"output"`
	Error      string        `json:"error,omitempty"`
}

func (s *Session) FastRefresh() (*ReloadResult, error) {
	s.mu.RLock()
	if s.state != SessionActive && s.state != SessionAttached && s.state != SessionMetroRunning {
		s.mu.RUnlock()
		return nil, fmt.Errorf("no active react native session — start dev mode first")
	}

	if !s.cfg.FastRefresh {
		s.mu.RUnlock()
		return nil, fmt.Errorf("fast refresh is disabled in configuration")
	}
	s.mu.RUnlock()

	s.log.Info("triggering fast refresh via metro")
	start := time.Now()

	url := fmt.Sprintf("http://localhost:%d/onchange", s.metroPort)
	resp, err := http.Post(url, "application/json", nil)
	duration := time.Since(start)

	if err != nil {
		result := &ReloadResult{
			Success:  false,
			Type:     "fast_refresh",
			Duration: duration,
			Error:    err.Error(),
		}
		s.log.Error("fast refresh failed", "error", err)
		return result, fmt.Errorf("fast refresh failed: %w", err)
	}
	defer resp.Body.Close()

	result := &ReloadResult{
		Success:  true,
		Type:     "fast_refresh",
		Duration: duration,
		Output:   fmt.Sprintf("HTTP %d", resp.StatusCode),
	}

	s.log.Info("fast refresh triggered", "duration", duration)
	return result, nil
}

func (s *Session) ManualReload() (*ReloadResult, error) {
	s.mu.RLock()
	if s.state != SessionActive && s.state != SessionAttached {
		s.mu.RUnlock()
		return nil, fmt.Errorf("no active react native session")
	}
	s.mu.RUnlock()

	s.log.Info("triggering manual reload")
	start := time.Now()

	url := fmt.Sprintf("http://localhost:%d/reload", s.metroPort)
	resp, err := http.Post(url, "application/json", strings.NewReader(`{"reload":true}`))
	duration := time.Since(start)

	if err != nil {
		result := &ReloadResult{
			Success:  false,
			Type:     "manual_reload",
			Duration: duration,
			Error:    err.Error(),
		}
		s.log.Error("manual reload failed", "error", err)
		return result, fmt.Errorf("manual reload failed: %w", err)
	}
	defer resp.Body.Close()

	result := &ReloadResult{
		Success:  true,
		Type:     "manual_reload",
		Duration: duration,
		Output:   fmt.Sprintf("HTTP %d", resp.StatusCode),
	}

	s.log.Info("manual reload triggered", "duration", duration)
	return result, nil
}

func (s *Session) DeviceRefresh() (*ReloadResult, error) {
	s.mu.RLock()
	if s.state != SessionActive && s.state != SessionAttached {
		s.mu.RUnlock()
		return nil, fmt.Errorf("no active react native session")
	}
	s.mu.RUnlock()

	s.log.Info("sending device refresh signal")
	start := time.Now()

	url := fmt.Sprintf("http://localhost:%d/device-refresh", s.metroPort)
	resp, err := http.Post(url, "application/json", nil)
	duration := time.Since(start)

	if err != nil {
		result := &ReloadResult{
			Success:  false,
			Type:     "device_refresh",
			Duration: duration,
			Error:    err.Error(),
		}
		s.log.Error("device refresh failed", "error", err)
		return result, fmt.Errorf("device refresh failed: %w", err)
	}
	defer resp.Body.Close()

	result := &ReloadResult{
		Success:  true,
		Type:     "device_refresh",
		Duration: duration,
		Output:   fmt.Sprintf("HTTP %d", resp.StatusCode),
	}

	s.log.Info("device refresh triggered", "duration", duration)
	return result, nil
}
