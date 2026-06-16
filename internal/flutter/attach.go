package flutter

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type AttachResult struct {
	Success   bool      `json:"success"`
	PID       int       `json:"pid,omitempty"`
	Device    string    `json:"device,omitempty"`
	URI       string    `json:"uri,omitempty"`
	AttachedAt time.Time `json:"attached_at"`
	Error     string    `json:"error,omitempty"`
}

func (s *Session) Attach(deviceID string) (*AttachResult, error) {
	s.mu.Lock()
	s.state = SessionStarting
	s.mu.Unlock()

	s.log.Info("attaching to flutter app", "device", deviceID)

	args := []string{"attach", "--machine"}
	if deviceID != "" {
		args = append(args, "-d", deviceID)
	}

	cmd := exec.Command("flutter", args...)
	cmd.Dir = s.projectDir
	output, err := cmd.Output()
	if err != nil {
		s.mu.Lock()
		s.state = SessionInactive
		s.mu.Unlock()
		return nil, fmt.Errorf("flutter attach failed: %w", err)
	}

	pid := 0
	uri := ""
	device := deviceID
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "pid:") {
			fmt.Sscanf(line, "pid: %d", &pid)
		}
		if strings.Contains(line, "uri:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) > 1 {
				uri = strings.TrimSpace(parts[1])
			}
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
	s.flutterPID = pid
	s.startedAt = time.Now()
	if device != "" {
		s.deviceID = device
	}
	s.mu.Unlock()

	result := &AttachResult{
		Success:    true,
		PID:        pid,
		Device:     device,
		URI:        uri,
		AttachedAt: s.startedAt,
	}

	s.log.Info("attached to flutter app", "pid", pid, "uri", uri)
	return result, nil
}

func (s *Session) Detach() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state == SessionInactive {
		return nil
	}

	s.state = SessionInactive
	s.flutterPID = 0
	s.log.Info("detached from flutter app")
	return nil
}
