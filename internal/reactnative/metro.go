package reactnative

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type MetroResult struct {
	Success  bool      `json:"success"`
	Action   string    `json:"action"`
	PID      int       `json:"pid,omitempty"`
	Port     int       `json:"port"`
	Host     string    `json:"host,omitempty"`
	Started  time.Time `json:"started_at,omitempty"`
	Error    string    `json:"error,omitempty"`
}

func (s *Session) StartMetro() (*MetroResult, error) {
	s.mu.Lock()
	port := s.metroPort
	s.mu.Unlock()

	s.log.Info("starting metro bundler", "port", port)

	args := []string{"react-native", "start", "--port", fmt.Sprintf("%d", port), "--no-interactive"}

	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start metro: %w", err)
	}

	s.mu.Lock()
	s.state = SessionMetroRunning
	s.metroPID = cmd.Process.Pid
	s.startedAt = time.Now()
	s.mu.Unlock()

	result := &MetroResult{
		Success: true,
		Action:  "start",
		PID:     cmd.Process.Pid,
		Port:    port,
		Started: s.startedAt,
	}

	s.log.Info("metro started", "pid", cmd.Process.Pid, "port", port)
	return result, nil
}

func (s *Session) StopMetro() (*MetroResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.metroPID == 0 {
		result := &MetroResult{
			Success: true,
			Action:  "stop",
			Port:    s.metroPort,
			Error:   "metro not running",
		}
		return result, nil
	}

	proc, err := os.FindProcess(s.metroPID)
	if err == nil {
		proc.Kill()
	}

	s.metroPID = 0
	s.state = SessionInactive

	result := &MetroResult{
		Success: true,
		Action:  "stop",
		Port:    s.metroPort,
	}

	s.log.Info("metro stopped")
	return result, nil
}

func (s *Session) RestartMetro() (*MetroResult, error) {
	s.log.Info("restarting metro bundler")

	if _, err := s.StopMetro(); err != nil {
		return nil, fmt.Errorf("failed to stop metro: %w", err)
	}

	time.Sleep(1 * time.Second)

	startResult, err := s.StartMetro()
	if err != nil {
		return nil, fmt.Errorf("failed to restart metro: %w", err)
	}
	startResult.Action = "restart"

	s.log.Info("metro restarted", "port", startResult.Port)
	return startResult, nil
}

func (s *Session) MetroStatus() *MetroResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	running := true
	if s.metroPID == 0 || s.state != SessionMetroRunning {
		running = false
	}

	result := &MetroResult{
		Success: running,
		Action:  "status",
		PID:     s.metroPID,
		Port:    s.metroPort,
	}

	if running {
		result.Started = s.startedAt
	}

	return result
}

func (s *Session) CheckMetroPort() error {
	s.mu.RLock()
	port := s.metroPort
	s.mu.RUnlock()

	cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port))
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return fmt.Errorf("port %d is already in use:\n%s", port, strings.TrimSpace(string(output)))
	}
	return nil
}
