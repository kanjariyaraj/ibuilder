package reactnative

import (
	"fmt"
	"time"
)

type SessionInfo struct {
	State       SessionState `json:"state"`
	ProjectDir  string       `json:"project_dir"`
	DeviceID    string       `json:"device_id,omitempty"`
	MetroPID    int          `json:"metro_pid,omitempty"`
	MetroPort   int          `json:"metro_port,omitempty"`
	Uptime      string       `json:"uptime"`
	StartedAt   time.Time    `json:"started_at,omitempty"`
}

func (s *Session) SessionInfo() SessionInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	info := SessionInfo{
		State:      s.state,
		ProjectDir: s.projectDir,
		DeviceID:   s.deviceID,
		MetroPID:   s.metroPID,
		MetroPort:  s.metroPort,
		StartedAt:  s.startedAt,
	}

	if !s.startedAt.IsZero() {
		uptime := time.Since(s.startedAt).Round(time.Second)
		info.Uptime = uptime.String()
	}

	return info
}

func (s *Session) StopSession() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.metroPID != 0 {
		s.stopMetro()
	}

	s.state = SessionInactive
	s.deviceID = ""

	s.log.Info("react native session stopped")
	return nil
}

func (s *Session) IsActive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state == SessionActive || s.state == SessionAttached || s.state == SessionMetroRunning
}

func FormatSessionInfo(info SessionInfo) string {
	return fmt.Sprintf(`React Native Session:
  State:       %s
  Project:     %s
  Device:      %s
  Metro PID:   %d
  Metro Port:  %d
  Uptime:      %s
  Started:     %s
`, info.State, info.ProjectDir, info.DeviceID, info.MetroPID, info.MetroPort,
		info.Uptime, info.StartedAt.Format(time.RFC3339))
}
