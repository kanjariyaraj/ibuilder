package flutter

import (
	"fmt"
	"time"
)

type SessionInfo struct {
	State       SessionState `json:"state"`
	ProjectDir  string       `json:"project_dir"`
	DeviceID    string       `json:"device_id,omitempty"`
	FlutterPID  int          `json:"flutter_pid,omitempty"`
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
		FlutterPID: s.flutterPID,
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

	s.state = SessionInactive
	s.flutterPID = 0
	s.deviceID = ""

	select {
	case <-s.watchStop:
	default:
		close(s.watchStop)
	}
	s.watchStop = make(chan struct{})

	s.log.Info("flutter session stopped")
	return nil
}

func (s *Session) IsActive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state == SessionActive || s.state == SessionAttached || s.state == SessionWatchMode
}

func FormatSessionInfo(info SessionInfo) string {
	return fmt.Sprintf(`Flutter Session:
  State:       %s
  Project:     %s
  Device:      %s
  PID:         %d
  Uptime:      %s
  Started:     %s
`, info.State, info.ProjectDir, info.DeviceID, info.FlutterPID, info.Uptime,
		info.StartedAt.Format(time.RFC3339))
}
