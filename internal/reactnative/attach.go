package reactnative

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type AttachResult struct {
	Success    bool      `json:"success"`
	PID        int       `json:"pid,omitempty"`
	Device     string    `json:"device,omitempty"`
	MetroPort  int       `json:"metro_port,omitempty"`
	AttachedAt time.Time `json:"attached_at"`
	Error      string    `json:"error,omitempty"`
}

func (s *Session) Attach(deviceID string) (*AttachResult, error) {
	s.mu.Lock()
	s.state = SessionStarting
	s.mu.Unlock()

	s.log.Info("attaching to react native app", "device", deviceID)

	args := []string{"react-native", "start", "--no-interactive"}
	if deviceID != "" {
		args = append(args, "--device", deviceID)
	}

	cmd := exec.Command("npx", args...)
	cmd.Dir = s.projectDir
	output, err := cmd.Output()
	if err != nil {
		s.mu.Lock()
		s.state = SessionInactive
		s.mu.Unlock()
		return nil, fmt.Errorf("react-native attach failed: %w", err)
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
	s.state = SessionAttached
	s.startedAt = time.Now()
	if device != "" {
		s.deviceID = device
	}
	s.mu.Unlock()

	result := &AttachResult{
		Success:    true,
		PID:        pid,
		Device:     device,
		MetroPort:  s.metroPort,
		AttachedAt: s.startedAt,
	}

	s.log.Info("attached to react native app", "pid", pid)
	return result, nil
}

func (s *Session) Detach() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state == SessionInactive {
		return nil
	}

	s.state = SessionInactive
	s.log.Info("detached from react native app")
	return nil
}
