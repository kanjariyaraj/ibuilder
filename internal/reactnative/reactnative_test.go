package reactnative

import (
	"testing"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/logger"
)

func newTestSession() *Session {
	cfg := &config.ReactNativeSettings{
		Enabled:        true,
		Entry:          "index.js",
		MetroPort:      8081,
		AutoStartMetro: true,
		AutoAttach:     true,
		AutoInstall:    true,
		FastRefresh:    true,
	}
	log := logger.New(logger.LevelInfo)
	return NewSession(cfg, log)
}

func TestNewSession(t *testing.T) {
	s := newTestSession()
	if s == nil {
		t.Fatal("expected session, got nil")
	}
	if s.State() != SessionInactive {
		t.Errorf("expected SessionInactive, got %v", s.State())
	}
	if s.metroPort != 8081 {
		t.Errorf("expected metroPort 8081, got %d", s.metroPort)
	}
}

func TestSessionConfig(t *testing.T) {
	s := newTestSession()
	cfg := s.Config()
	if cfg == nil {
		t.Fatal("expected config, got nil")
	}
	if !cfg.AutoStartMetro {
		t.Error("expected AutoStartMetro to be true")
	}
}

func TestSessionUpdateConfig(t *testing.T) {
	s := newTestSession()
	newCfg := &config.ReactNativeSettings{MetroPort: 9090}
	s.UpdateConfig(newCfg)
	if s.cfg.MetroPort != 9090 {
		t.Errorf("expected MetroPort 9090, got %d", s.cfg.MetroPort)
	}
}

func TestSessionState(t *testing.T) {
	s := newTestSession()
	if s.State() != SessionInactive {
		t.Errorf("expected Inactive, got %v", s.State())
	}
}

func TestSessionDeviceID(t *testing.T) {
	s := newTestSession()
	if s.DeviceID() != "" {
		t.Errorf("expected empty DeviceID, got %s", s.DeviceID())
	}
	s.SetDeviceID("test-device")
	if s.DeviceID() != "test-device" {
		t.Errorf("expected test-device, got %s", s.DeviceID())
	}
}

func TestSessionIsActive(t *testing.T) {
	s := newTestSession()
	if s.IsActive() {
		t.Error("expected not active")
	}
}

func TestSessionInfo(t *testing.T) {
	s := newTestSession()
	info := s.SessionInfo()
	if info.State != SessionInactive {
		t.Errorf("expected Inactive, got %v", info.State)
	}
}

func TestFormatSessionInfo(t *testing.T) {
	info := SessionInfo{
		State:      SessionActive,
		MetroPort:  8081,
		MetroPID:   12345,
		ProjectDir: "/test/project",
	}
	formatted := FormatSessionInfo(info)
	if formatted == "" {
		t.Error("expected non-empty formatted string")
	}
}

func TestTimestamp(t *testing.T) {
	s := newTestSession()
	ts := s.Timestamp()
	if ts == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestHealthStatusConstants(t *testing.T) {
	if StatusHealthy != "HEALTHY" {
		t.Errorf("expected HEALTHY, got %s", StatusHealthy)
	}
	if StatusWarning != "WARNING" {
		t.Errorf("expected WARNING, got %s", StatusWarning)
	}
	if StatusFailure != "FAILURE" {
		t.Errorf("expected FAILURE, got %s", StatusFailure)
	}
}

func TestLogLevelDetection(t *testing.T) {
	tests := []struct {
		msg     string
		want    string
	}{
		{"this is an ERROR", "ERROR"},
		{"FATAL error occurred", "ERROR"},
		{"warning: something", "WARN"},
		{"WARN: low disk", "WARN"},
		{"debug info", "DEBUG"},
		{"DEBUG: verbose", "DEBUG"},
		{"regular info message", "INFO"},
		{"", "INFO"},
	}

	for _, tt := range tests {
		got := detectLogLevel(tt.msg)
		if got != tt.want {
			t.Errorf("detectLogLevel(%q) = %q, want %q", tt.msg, got, tt.want)
		}
	}
}

func TestParseLogOutput(t *testing.T) {
	output := "line1\nERROR: something failed\nWARN: almost there\n"
	logs := parseLogOutput(output, nil)
	if len(logs) != 3 {
		t.Errorf("expected 3 log entries, got %d", len(logs))
	}
}

func TestParseLogOutputWithFilter(t *testing.T) {
	output := "INFO: normal\nERROR: critical\nWARN: careful\n"
	filter := &LogFilter{Level: "ERROR"}
	logs := parseLogOutput(output, filter)
	if len(logs) != 1 {
		t.Errorf("expected 1 filtered log entry, got %d", len(logs))
	}
	if len(logs) > 0 && logs[0].Level != "ERROR" {
		t.Errorf("expected ERROR level, got %s", logs[0].Level)
	}
}

func TestParseLogOutputWithSearch(t *testing.T) {
	output := "everything is fine\ncritical error occurred\nall good\n"
	filter := &LogFilter{Search: "critical"}
	logs := parseLogOutput(output, filter)
	if len(logs) != 1 {
		t.Errorf("expected 1 filtered log entry, got %d", len(logs))
	}
}

func TestMetroStatusWhenInactive(t *testing.T) {
	s := newTestSession()
	status := s.MetroStatus()
	if status.Success {
		t.Error("expected metro to be stopped")
	}
}

func TestRecoveryResultDefaults(t *testing.T) {
	result := &RecoveryResult{
		Success:  false,
		Action:   "recovery_failed",
		Attempts: 3,
	}
	if result.Success {
		t.Error("expected failed recovery")
	}
	if result.Action != "recovery_failed" {
		t.Errorf("expected recovery_failed, got %s", result.Action)
	}
}

func TestSessionStateStrings(t *testing.T) {
	tests := []struct {
		state SessionState
		want  string
	}{
		{SessionInactive, "Inactive"},
		{SessionStarting, "Starting"},
		{SessionActive, "Active"},
		{SessionAttached, "Attached"},
		{SessionMetroRunning, "Metro Running"},
		{SessionRecovering, "Recovering"},
		{SessionState(99), "Unknown"},
	}
	for _, tt := range tests {
		got := tt.state.String()
		if got != tt.want {
			t.Errorf("SessionState(%d).String() = %q, want %q", tt.state, got, tt.want)
		}
	}
}
