package reactnative

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type MetroStatus struct {
	Running  bool   `json:"running"`
	PID      int    `json:"pid,omitempty"`
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Uptime   string `json:"uptime,omitempty"`
	Error    string `json:"error,omitempty"`
}

type MetroResult struct {
	Success  bool      `json:"success"`
	Action   string    `json:"action"`
	PID      int       `json:"pid,omitempty"`
	Port     int       `json:"port"`
	Started  time.Time `json:"started_at,omitempty"`
	Error    string    `json:"error,omitempty"`
}

func (s *Session) StartMetro() error {
	if s.isMetroRunning() {
		s.log.Info("metro bundler is already running", "port", s.metroPort)
		return nil
	}

	s.log.Info("starting metro bundler", "port", s.metroPort)

	portStr := strconv.Itoa(s.metroPort)
	args := []string{
		"react-native", "start",
		"--port", portStr,
		"--no-interactive",
	}

	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start metro: %w", err)
	}

	s.mu.Lock()
	s.metroPID = cmd.Process.Pid
	s.state = SessionMetroRunning
	s.mu.Unlock()

	s.log.Info("metro bundler started", "pid", s.metroPID, "port", s.metroPort)

	time.Sleep(3 * time.Second)

	if !s.isMetroRunning() {
		return fmt.Errorf("metro exited prematurely")
	}

	return nil
}

func (s *Session) StopMetro() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.stopMetro()
}

func (s *Session) stopMetro() error {
	if s.metroPID == 0 {
		return nil
	}

	s.log.Info("stopping metro bundler", "pid", s.metroPID)

	process, err := os.FindProcess(s.metroPID)
	if err != nil {
		s.metroPID = 0
		return nil
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		if err := process.Kill(); err != nil {
			s.metroPID = 0
			return fmt.Errorf("failed to stop metro: %w", err)
		}
	}

	time.Sleep(1 * time.Second)
	s.metroPID = 0
	s.state = SessionInactive
	s.log.Info("metro bundler stopped")
	return nil
}

func (s *Session) RestartMetro() error {
	s.log.Info("restarting metro bundler")

	if err := s.StopMetro(); err != nil {
		s.log.Error("failed to stop metro for restart", "error", err)
	}

	time.Sleep(2 * time.Second)

	if err := s.StartMetro(); err != nil {
		return fmt.Errorf("metro restart failed: %w", err)
	}

	s.log.Info("metro bundler restarted")
	return nil
}

func (s *Session) MetroStatus() MetroStatus {
	s.mu.RLock()
	running := s.isMetroRunning()
	pid := s.metroPID
	port := s.metroPort
	s.mu.RUnlock()

	status := MetroStatus{
		Running: running,
		PID:     pid,
		Port:    port,
		Host:    "localhost",
	}

	if running && !s.startedAt.IsZero() {
		uptime := time.Since(s.startedAt).Round(time.Second)
		status.Uptime = uptime.String()
	}

	return status
}

func (s *Session) isMetroRunning() bool {
	if s.metroPID == 0 {
		return false
	}

	process, err := os.FindProcess(s.metroPID)
	if err != nil {
		return false
	}

	if err := process.Signal(syscall.Signal(0)); err != nil {
		return false
	}

	return true
}

func (s *Session) CheckPortConflict() (bool, int) {
	portStr := strconv.Itoa(s.metroPort)

	cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%s", portStr))
	output, err := cmd.Output()
	if err != nil {
		return false, 0
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) > 1 {
		for _, line := range lines[1:] {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				pid, _ := strconv.Atoi(fields[1])
				if pid > 0 && pid != s.metroPID {
					return true, pid
				}
			}
		}
	}

	return false, 0
}
