package reactnative

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/logger"
)

func newTestSession(t *testing.T) *Session {
	t.Helper()
	cfg := &config.ReactNativeSettings{
		Enabled:        true,
		Entry:          "index.js",
		MetroPort:      8081,
		AutoStartMetro: true,
		AutoAttach:     true,
		AutoInstall:    true,
		FastRefresh:    true,
	}
	log := logger.New(logger.LevelDebug)
	return NewSession(cfg, log)
}

func TestNewSession(t *testing.T) {
	session := newTestSession(t)
	if session == nil {
		t.Fatal("expected session to be non-nil")
	}
	if session.State() != SessionInactive {
		t.Fatalf("expected inactive state, got %s", session.State())
	}
	if session.metroPort != 8081 {
		t.Fatalf("expected metro port 8081, got %d", session.metroPort)
	}
}

func TestNewSessionCustomPort(t *testing.T) {
	cfg := &config.ReactNativeSettings{MetroPort: 9090}
	log := logger.New(logger.LevelDebug)
	session := NewSession(cfg, log)
	if session.metroPort != 9090 {
		t.Fatalf("expected metro port 9090, got %d", session.metroPort)
	}
}

func TestSessionStateString(t *testing.T) {
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
		if got := tt.state.String(); got != tt.want {
			t.Errorf("SessionState(%d).String() = %q, want %q", tt.state, got, tt.want)
		}
	}
}

func TestSetDeviceID(t *testing.T) {
	session := newTestSession(t)
	if id := session.DeviceID(); id != "" {
		t.Fatalf("expected empty device id, got %s", id)
	}

	session.SetDeviceID("test-device-123")
	if id := session.DeviceID(); id != "test-device-123" {
		t.Fatalf("expected 'test-device-123', got %s", id)
	}
}

func TestConfigAccess(t *testing.T) {
	session := newTestSession(t)
	cfg := session.Config()
	if cfg == nil {
		t.Fatal("expected config to be non-nil")
	}
	if !cfg.Enabled {
		t.Fatal("expected enabled to be true")
	}
}

func TestUpdateConfig(t *testing.T) {
	session := newTestSession(t)
	newCfg := &config.ReactNativeSettings{Enabled: false}
	session.UpdateConfig(newCfg)

	if session.cfg.Enabled {
		t.Fatal("expected enabled to be false after update")
	}
}

func TestDetectRNProjectNoPackageJSON(t *testing.T) {
	session := newTestSession(t)
	dir := t.TempDir()

	valid, err := session.DetectRNProject(dir)
	if valid {
		t.Fatal("expected invalid project")
	}
	if err == nil {
		t.Fatal("expected error for missing package.json")
	}
}

func TestDetectRNProjectNoReactNative(t *testing.T) {
	session := newTestSession(t)
	dir := t.TempDir()

	pkgPath := filepath.Join(dir, "package.json")
	if err := os.WriteFile(pkgPath, []byte(`{"name": "test"}`), 0644); err != nil {
		t.Fatal(err)
	}

	valid, err := session.DetectRNProject(dir)
	if valid {
		t.Fatal("expected invalid project")
	}
	if err == nil || err.Error() != "react-native dependency not found in package.json" {
		t.Fatalf("expected react-native dependency error, got: %v", err)
	}
}

func TestDetectRNProjectNoIOSDir(t *testing.T) {
	session := newTestSession(t)
	dir := t.TempDir()

	pkgPath := filepath.Join(dir, "package.json")
	if err := os.WriteFile(pkgPath, []byte(`{"dependencies": {"react-native": "0.73.0"}}`), 0644); err != nil {
		t.Fatal(err)
	}

	valid, err := session.DetectRNProject(dir)
	if valid {
		t.Fatal("expected invalid project")
	}
	if err == nil || err.Error() != "no ios/ directory found — not a React Native iOS project" {
		t.Fatalf("expected ios/ directory error, got: %v", err)
	}
}

func TestDetectRNProjectValid(t *testing.T) {
	session := newTestSession(t)
	dir := t.TempDir()

	pkgPath := filepath.Join(dir, "package.json")
	if err := os.WriteFile(pkgPath, []byte(`{"dependencies": {"react-native": "0.73.0"}}`), 0644); err != nil {
		t.Fatal(err)
	}

	iosDir := filepath.Join(dir, "ios")
	if err := os.MkdirAll(iosDir, 0755); err != nil {
		t.Fatal(err)
	}

	valid, err := session.DetectRNProject(dir)
	if !valid {
		t.Fatal("expected valid project")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session.ProjectDir() != dir {
		t.Fatalf("expected project dir %s, got %s", dir, session.ProjectDir())
	}
}

func TestIsActive(t *testing.T) {
	session := newTestSession(t)
	if session.IsActive() {
		t.Fatal("expected inactive session")
	}

	session.mu.Lock()
	session.state = SessionActive
	session.mu.Unlock()

	if !session.IsActive() {
		t.Fatal("expected active session")
	}
}

func TestSessionInfo(t *testing.T) {
	session := newTestSession(t)
	info := session.SessionInfo()

	if info.State != SessionInactive {
		t.Fatalf("expected inactive state, got %s", info.State)
	}
	if info.ProjectDir != "" {
		t.Fatalf("expected empty project dir, got %s", info.ProjectDir)
	}
	if info.MetroPort != 8081 {
		t.Fatalf("expected metro port 8081, got %d", info.MetroPort)
	}
}

func TestFormatSessionInfo(t *testing.T) {
	info := SessionInfo{
		State:      SessionActive,
		ProjectDir: "/test/project",
		DeviceID:   "test-device",
		MetroPID:   12345,
		MetroPort:  8081,
	}

	output := FormatSessionInfo(info)
	if output == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestHealthStatusConstants(t *testing.T) {
	if StatusHealthy != "HEALTHY" {
		t.Fatalf("expected HEALTHY, got %s", StatusHealthy)
	}
	if StatusWarning != "WARNING" {
		t.Fatalf("expected WARNING, got %s", StatusWarning)
	}
	if StatusFailure != "FAILURE" {
		t.Fatalf("expected FAILURE, got %s", StatusFailure)
	}
}

func TestDoctorReport(t *testing.T) {
	report := DoctorReport{
		Healthy: true,
		Checks:  []HealthCheck{},
	}

	if !report.Healthy {
		t.Fatal("expected healthy report")
	}
}

func TestLogLevelDetection(t *testing.T) {
	tests := []struct {
		msg  string
		want string
	}{
		{"ERROR: something failed", "ERROR"},
		{"FATAL: crash", "ERROR"},
		{"WARN: something", "WARN"},
		{"WARNING: caution", "WARN"},
		{"DEBUG: detail", "DEBUG"},
		{"info message", "INFO"},
		{"random log", "INFO"},
	}

	for _, tt := range tests {
		if got := detectLogLevel(tt.msg); got != tt.want {
			t.Errorf("detectLogLevel(%q) = %q, want %q", tt.msg, got, tt.want)
		}
	}
}

func TestParseLogOutput(t *testing.T) {
	output := "ERROR: test error\nINFO: test info\nWARN: test warning"
	logs := parseLogOutput(output, nil)

	if len(logs) != 3 {
		t.Fatalf("expected 3 log entries, got %d", len(logs))
	}

	if logs[0].Level != "ERROR" {
		t.Fatalf("expected ERROR level, got %s", logs[0].Level)
	}
	if logs[1].Level != "INFO" {
		t.Fatalf("expected INFO level, got %s", logs[1].Level)
	}
	if logs[2].Level != "WARN" {
		t.Fatalf("expected WARN level, got %s", logs[2].Level)
	}
}

func TestParseLogOutputWithFilter(t *testing.T) {
	output := "ERROR: test error\nINFO: test info\nWARN: test warning"

	filter := &LogFilter{Level: "ERROR"}
	logs := parseLogOutput(output, filter)

	if len(logs) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(logs))
	}
	if logs[0].Level != "ERROR" {
		t.Fatalf("expected ERROR level, got %s", logs[0].Level)
	}
}

func TestParseLogOutputWithSearch(t *testing.T) {
	output := "ERROR: critical failure\nINFO: all good\nWARN: minor issue"

	filter := &LogFilter{Search: "failure"}
	logs := parseLogOutput(output, filter)

	if len(logs) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(logs))
	}
	if !containsIgnoreCase(logs[0].Message, "failure") {
		t.Fatalf("expected message containing 'failure', got %s", logs[0].Message)
	}
}

func TestParseLogOutputEmptyLines(t *testing.T) {
	output := "line1\n\n\nline2"
	logs := parseLogOutput(output, nil)

	if len(logs) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(logs))
	}
}

func TestRecoveryResult(t *testing.T) {
	result := RecoveryResult{
		Success:  true,
		Action:   "recovered",
		Attempts: 2,
	}

	if !result.Success {
		t.Fatal("expected success")
	}
	if result.Action != "recovered" {
		t.Fatalf("expected 'recovered', got %s", result.Action)
	}
}

func TestDevResult(t *testing.T) {
	result := DevResult{
		Success:   true,
		Action:    "dev",
		PID:       12345,
		Device:    "test-device",
		MetroPort: 8081,
	}

	if !result.Success {
		t.Fatal("expected success")
	}
	if result.MetroPort != 8081 {
		t.Fatalf("expected metro port 8081, got %d", result.MetroPort)
	}
}

func TestAttachResult(t *testing.T) {
	result := AttachResult{
		Success:   true,
		PID:       12345,
		Device:    "test-device",
		MetroPort: 8081,
	}

	if !result.Success {
		t.Fatal("expected success")
	}
}

func TestInstallResult(t *testing.T) {
	result := InstallResult{
		Success: true,
		Action:  "install_latest",
		Output:  "Installed successfully",
	}

	if !result.Success {
		t.Fatal("expected success")
	}
}

func TestReloadResult(t *testing.T) {
	result := ReloadResult{
		Success: true,
		Type:    "fast_refresh",
	}

	if !result.Success {
		t.Fatal("expected success")
	}
	if result.Type != "fast_refresh" {
		t.Fatalf("expected 'fast_refresh', got %s", result.Type)
	}
}

func TestMetroStatus(t *testing.T) {
	status := MetroStatus{
		Running: true,
		PID:     12345,
		Port:    8081,
		Host:    "localhost",
	}

	if !status.Running {
		t.Fatal("expected running")
	}
	if status.Port != 8081 {
		t.Fatalf("expected port 8081, got %d", status.Port)
	}
}

func TestMetroResult(t *testing.T) {
	result := MetroResult{
		Success: true,
		Action:  "start",
		PID:     12345,
		Port:    8081,
	}

	if !result.Success {
		t.Fatal("expected success")
	}
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
