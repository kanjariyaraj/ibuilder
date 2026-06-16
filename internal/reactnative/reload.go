package reactnative

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ReloadResult struct {
	Success  bool          `json:"success"`
	Action   string        `json:"action"`
	Duration time.Duration `json:"duration_ms"`
	Error    string        `json:"error,omitempty"`
}

func (s *Session) FastRefresh() (*ReloadResult, error) {
	s.mu.RLock()
	if s.state != SessionActive && s.state != SessionAttached && s.state != SessionMetroRunning {
		s.mu.RUnlock()
		return nil, fmt.Errorf("no active RN session — start one with 'rn dev' or 'rn attach'")
	}
	port := s.metroPort
	s.mu.RUnlock()

	if !s.cfg.FastRefresh {
		return nil, fmt.Errorf("fast refresh is disabled in configuration")
	}

	s.log.Info("triggering fast refresh")
	start := time.Now()

	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/onchange", port),
		"application/json",
		strings.NewReader(`{"changes":["__fast_refresh__"]}`),
	)
	duration := time.Since(start)

	if err != nil {
		result := &ReloadResult{
			Success:  false,
			Action:   "fast_refresh",
			Duration: duration,
			Error:    err.Error(),
		}
		s.log.Error("fast refresh failed", "error", err)
		return result, fmt.Errorf("fast refresh failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result := &ReloadResult{
			Success:  false,
			Action:   "fast_refresh",
			Duration: duration,
			Error:    fmt.Sprintf("metro returned status %d", resp.StatusCode),
		}
		return result, fmt.Errorf("fast refresh failed: metro returned %d", resp.StatusCode)
	}

	result := &ReloadResult{
		Success:  true,
		Action:   "fast_refresh",
		Duration: duration,
	}

	s.log.Info("fast refresh triggered", "duration", duration)
	return result, nil
}

func (s *Session) ManualReload() (*ReloadResult, error) {
	s.mu.RLock()
	if s.state != SessionActive && s.state != SessionAttached && s.state != SessionMetroRunning {
		s.mu.RUnlock()
		return nil, fmt.Errorf("no active RN session — start one with 'rn dev' or 'rn attach'")
	}
	port := s.metroPort
	s.mu.RUnlock()

	s.log.Info("triggering manual reload")
	start := time.Now()

	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/reload", port),
		"application/json",
		nil,
	)
	duration := time.Since(start)

	if err != nil {
		result := &ReloadResult{
			Success:  false,
			Action:   "manual_reload",
			Duration: duration,
			Error:    err.Error(),
		}
		s.log.Error("manual reload failed", "error", err)
		return result, fmt.Errorf("manual reload failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result := &ReloadResult{
			Success:  false,
			Action:   "manual_reload",
			Duration: duration,
			Error:    fmt.Sprintf("metro returned status %d", resp.StatusCode),
		}
		return result, fmt.Errorf("manual reload failed: metro returned %d", resp.StatusCode)
	}

	result := &ReloadResult{
		Success:  true,
		Action:   "manual_reload",
		Duration: duration,
	}

	s.log.Info("manual reload triggered", "duration", duration)
	return result, nil
}
