package reactnative

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type DevResult struct {
	Success   bool      `json:"success"`
	Action    string    `json:"action"`
	PID       int       `json:"pid,omitempty"`
	Device    string    `json:"device,omitempty"`
	MetroPort int       `json:"metro_port,omitempty"`
	Started   time.Time `json:"started_at"`
	Error     string    `json:"error,omitempty"`
}

func (s *Session) DevMode() (*DevResult, error) {
	s.mu.Lock()
	s.state = SessionStarting
	s.mu.Unlock()

	s.log.Info("starting react native dev mode")

	if s.projectDir == "" {
		return nil, fmt.Errorf("no project directory set — run from a React Native project")
	}

	s.log.Info("resolving dependencies")
	if err := s.ResolveDependencies(); err != nil {
		return nil, fmt.Errorf("dependency resolution failed: %w", err)
	}

	if s.cfg.AutoStartMetro {
		if err := s.StartMetro(); err != nil {
			return nil, fmt.Errorf("metro start failed: %w", err)
		}
	}

	deviceID := s.deviceID
	entry := s.cfg.Entry
	if entry == "" {
		entry = "index.js"
	}

	args := []string{"react-native", "run-ios"}
	if deviceID != "" {
		args = append(args, "--device", deviceID)
	}

	s.log.Info("launching react native app", "device", deviceID)

	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		s.mu.Lock()
		s.state = SessionInactive
		if s.cfg.AutoStartMetro {
			s.stopMetro()
		}
		s.mu.Unlock()
		return nil, fmt.Errorf("react-native run-ios failed: %s", strings.TrimSpace(string(output)))
	}

	pid := 0
	device := deviceID
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "pid:") || strings.Contains(line, "PID:") {
			fmt.Sscanf(line, "%*s %d", &pid)
		}
		if strings.Contains(line, "device:") || strings.Contains(line, "Device:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) > 1 {
				device = strings.TrimSpace(parts[1])
			}
		}
	}

	s.mu.Lock()
	s.state = SessionActive
	s.startedAt = time.Now()
	if device != "" {
		s.deviceID = device
	}
	s.mu.Unlock()

	result := &DevResult{
		Success:   true,
		Action:    "dev",
		PID:       pid,
		Device:    device,
		MetroPort: s.metroPort,
		Started:   s.startedAt,
	}

	s.log.Info("react native dev mode active", "pid", pid, "device", device, "metroPort", s.metroPort)

	if s.cfg.AutoAttach {
		s.mu.Lock()
		s.state = SessionAttached
		s.mu.Unlock()
		s.log.Info("auto-attached to react native session")
	}

	return result, nil
}

func (s *Session) BuildAndInstall() error {
	s.log.Info("building and installing react native app")

	deviceID := s.deviceID
	args := []string{"react-native", "run-ios", "--mode", "Release"}
	if deviceID != "" {
		args = append(args, "--device", deviceID)
	}

	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("react-native build and install failed: %s", strings.TrimSpace(string(output)))
	}

	s.log.Info("react native app installed")
	return nil
}
