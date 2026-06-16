package flutter

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type DevResult struct {
	Success bool      `json:"success"`
	Action  string    `json:"action"`
	PID     int       `json:"pid,omitempty"`
	Device  string    `json:"device,omitempty"`
	Started time.Time `json:"started_at"`
	Error   string    `json:"error,omitempty"`
}

func (s *Session) DevMode() (*DevResult, error) {
	s.mu.Lock()
	s.state = SessionStarting
	s.mu.Unlock()

	s.log.Info("starting flutter dev mode")

	if s.projectDir == "" {
		return nil, fmt.Errorf("no project directory set — run from a Flutter project")
	}

	s.log.Info("resolving dependencies")
	if err := s.ResolveDependencies(); err != nil {
		return nil, fmt.Errorf("dependency resolution failed: %w", err)
	}

	deviceID := s.deviceID
	args := []string{"run", "--machine"}
	if deviceID != "" {
		args = append(args, "-d", deviceID)
	}

	s.log.Info("building and launching flutter app", "device", deviceID)
	cmd := exec.Command("flutter", args...)
	cmd.Dir = s.projectDir

	output, err := cmd.Output()
	if err != nil {
		s.mu.Lock()
		s.state = SessionInactive
		s.mu.Unlock()
		return nil, fmt.Errorf("flutter run failed: %w", err)
	}

	pid := 0
	device := deviceID
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "pid:") {
			fmt.Sscanf(line, "pid: %d", &pid)
		}
		if strings.Contains(line, "device:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) > 1 {
				device = strings.TrimSpace(parts[1])
			}
		}
	}

	s.mu.Lock()
	s.state = SessionActive
	s.flutterPID = pid
	s.startedAt = time.Now()
	if device != "" {
		s.deviceID = device
	}
	s.mu.Unlock()

	result := &DevResult{
		Success: true,
		Action:  "dev",
		PID:     pid,
		Device:  device,
		Started: s.startedAt,
	}

	s.log.Info("flutter dev mode active", "pid", pid, "device", device)

	if s.cfg.AutoAttach {
		s.mu.Lock()
		s.state = SessionAttached
		s.mu.Unlock()
		s.log.Info("auto-attached to flutter session")
	}

	return result, nil
}

func (s *Session) BuildAndInstall() error {
	s.log.Info("building and installing flutter app")

	deviceID := s.deviceID
	args := []string{"install"}
	if deviceID != "" {
		args = append(args, "-d", deviceID)
	}

	cmd := exec.Command("flutter", args...)
	cmd.Dir = s.projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("flutter install failed: %s", strings.TrimSpace(string(output)))
	}

	s.log.Info("flutter app installed")
	return nil
}
