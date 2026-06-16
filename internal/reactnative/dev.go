package reactnative

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type DevResult struct {
	Success   bool      `json:"success"`
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
		return nil, fmt.Errorf("no project directory set — run from an RN project")
	}

	if s.cfg.AutoStartMetro {
		s.log.Info("starting metro bundler")
		metroResult, err := s.StartMetro()
		if err != nil {
			return nil, fmt.Errorf("metro failed to start: %w", err)
		}
		s.mu.Lock()
		s.metroPID = metroResult.PID
		s.mu.Unlock()
	}

	s.log.Info("installing and launching RN app")
	args := []string{"react-native", "run-ios", "--json"}
	if s.deviceID != "" {
		args = append(args, "--device", s.deviceID)
	}
	if s.metroPort != 8081 {
		args = append(args, "--port", fmt.Sprintf("%d", s.metroPort))
	}

	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		s.mu.Lock()
		s.state = SessionInactive
		s.mu.Unlock()
		return nil, fmt.Errorf("react-native run-ios failed: %s", strings.TrimSpace(string(output)))
	}

	s.mu.Lock()
	s.state = SessionActive
	s.startedAt = time.Now()
	s.mu.Unlock()

	result := &DevResult{
		Success:   true,
		Device:    s.deviceID,
		MetroPort: s.metroPort,
		Started:   s.startedAt,
	}

	s.log.Info("react native dev mode active", "port", s.metroPort, "device", s.deviceID)

	if s.cfg.AutoAttach {
		s.mu.Lock()
		s.state = SessionAttached
		s.mu.Unlock()
		s.log.Info("auto-attached to RN session")
	}

	return result, nil
}

func (s *Session) BuildAndInstall() error {
	s.log.Info("building and installing RN app")

	args := []string{"react-native", "run-ios", "--no-launch", "--json"}
	if s.deviceID != "" {
		args = append(args, "--device", s.deviceID)
	}

	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build and install failed: %s", strings.TrimSpace(string(output)))
	}

	s.log.Info("RN app built and installed")
	return nil
}
