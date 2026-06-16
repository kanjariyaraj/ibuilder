package flutter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/logger"
)

func newTestSession() *Session {
	cfg := &config.FlutterSettings{
		Enabled:     true,
		Channel:     "stable",
		Watch:       true,
		HotReload:   true,
		DebounceMs:  500,
		AutoAttach:  true,
		AutoInstall: true,
	}
	log := logger.New(logger.LevelDebug)
	return NewSession(cfg, log)
}

func TestNewSession(t *testing.T) {
	s := newTestSession()
	if s == nil {
		t.Fatal("expected session, got nil")
	}
	if s.State() != SessionInactive {
		t.Errorf("expected inactive, got %v", s.State())
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
		{SessionWatchMode, "Watch Mode"},
		{SessionRecovering, "Recovering"},
		{SessionState(99), "Unknown"},
	}

	for _, tt := range tests {
		got := tt.state.String()
		if got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}

func TestConfig(t *testing.T) {
	s := newTestSession()
	cfg := s.Config()
	if cfg.Channel != "stable" {
		t.Errorf("expected stable, got %s", cfg.Channel)
	}
	if !cfg.HotReload {
		t.Error("expected hot reload enabled")
	}
}

func TestUpdateConfig(t *testing.T) {
	s := newTestSession()
	newCfg := &config.FlutterSettings{
		Enabled: false,
		Channel: "beta",
	}
	s.UpdateConfig(newCfg)
	cfg := s.Config()
	if cfg.Channel != "beta" {
		t.Errorf("expected beta, got %s", cfg.Channel)
	}
	if cfg.Enabled {
		t.Error("expected disabled")
	}
}

func TestDetectFlutterProject(t *testing.T) {
	s := newTestSession()

	dir := t.TempDir()
	pubspec := filepath.Join(dir, "pubspec.yaml")
	if err := os.WriteFile(pubspec, []byte("name: test_project"), 0644); err != nil {
		t.Fatal(err)
	}
	iosDir := filepath.Join(dir, "ios")
	if err := os.MkdirAll(iosDir, 0755); err != nil {
		t.Fatal(err)
	}

	valid, err := s.DetectFlutterProject(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !valid {
		t.Error("expected valid project")
	}
}

func TestDetectFlutterProjectNoPubspec(t *testing.T) {
	s := newTestSession()
	dir := t.TempDir()

	_, err := s.DetectFlutterProject(dir)
	if err == nil {
		t.Error("expected error for missing pubspec.yaml")
	}
}

func TestDetectFlutterProjectNoIOS(t *testing.T) {
	s := newTestSession()
	dir := t.TempDir()
	pubspec := filepath.Join(dir, "pubspec.yaml")
	if err := os.WriteFile(pubspec, []byte("name: test"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := s.DetectFlutterProject(dir)
	if err == nil {
		t.Error("expected error for missing ios/ directory")
	}
}

func TestSetDeviceID(t *testing.T) {
	s := newTestSession()
	s.SetDeviceID("test-device-123")
	if id := s.DeviceID(); id != "test-device-123" {
		t.Errorf("expected test-device-123, got %s", id)
	}
}

func TestStopSession(t *testing.T) {
	s := newTestSession()
	if err := s.StopSession(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if s.State() != SessionInactive {
		t.Errorf("expected inactive, got %v", s.State())
	}
}

func TestIsActive(t *testing.T) {
	s := newTestSession()
	if s.IsActive() {
		t.Error("expected inactive")
	}

	s.mu.Lock()
	s.state = SessionActive
	s.mu.Unlock()

	if !s.IsActive() {
		t.Error("expected active")
	}
}

func TestDetach(t *testing.T) {
	s := newTestSession()
	if err := s.Detach(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if s.State() != SessionInactive {
		t.Errorf("expected inactive, got %v", s.State())
	}
}

func TestNewFileWatcher(t *testing.T) {
	w := NewFileWatcher("/tmp", 500)
	if w == nil {
		t.Fatal("expected watcher, got nil")
	}
	if w.IsRunning() {
		t.Error("expected watcher not running initially")
	}
}

func TestFileWatcherStartStop(t *testing.T) {
	dir := t.TempDir()
	w := NewFileWatcher(dir, 500)

	if err := w.Start(); err != nil {
		t.Fatalf("expected no error starting watcher, got %v", err)
	}

	if !w.IsRunning() {
		t.Error("expected watcher running")
	}

	w.Stop()

	if w.IsRunning() {
		t.Error("expected watcher stopped")
	}
}

func TestFileWatcherDoubleStart(t *testing.T) {
	dir := t.TempDir()
	w := NewFileWatcher(dir, 500)

	if err := w.Start(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err := w.Start()
	if err == nil {
		t.Error("expected error starting already running watcher")
	}

	w.Stop()
}

func TestFormatSessionInfo(t *testing.T) {
	info := SessionInfo{
		State:      SessionActive,
		ProjectDir: "/test/project",
		DeviceID:   "test-device",
		FlutterPID: 1234,
	}

	output := FormatSessionInfo(info)
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestDoctor(t *testing.T) {
	s := newTestSession()
	report := s.Doctor()

	if len(report.Checks) == 0 {
		t.Error("expected at least one check")
	}

	hasFlutter := false
	hasDart := false
	hasProject := false
	hasDeps := false
	hasDevices := false

	for _, check := range report.Checks {
		switch check.Name {
		case "Flutter SDK":
			hasFlutter = true
		case "Dart SDK":
			hasDart = true
		case "Flutter Project":
			hasProject = true
		case "Dependencies":
			hasDeps = true
		case "Devices":
			hasDevices = true
		}
	}

	if !hasFlutter {
		t.Error("expected Flutter SDK check")
	}
	if !hasDart {
		t.Error("expected Dart SDK check")
	}
	if !hasProject {
		t.Error("expected Flutter Project check")
	}
	if !hasDeps {
		t.Error("expected Dependencies check")
	}
	if !hasDevices {
		t.Error("expected Devices check")
	}
}

func TestLogFilter(t *testing.T) {
	logs := parseLogOutput(`[INFO] App started
[ERROR] Connection failed
[WARN] Low memory
[DEBUG] Initializing`, &LogFilter{
		Level: "ERROR",
	})

	if len(logs) != 1 {
		t.Errorf("expected 1 ERROR log, got %d", len(logs))
	}
	if len(logs) > 0 && logs[0].Level != "ERROR" {
		t.Errorf("expected ERROR level, got %s", logs[0].Level)
	}
}

func TestLogFilterSearch(t *testing.T) {
	logs := parseLogOutput(`App started
Connection failed
Low memory
Initializing`, &LogFilter{
		Search: "memory",
	})

	if len(logs) != 1 {
		t.Errorf("expected 1 log matching search, got %d", len(logs))
	}
}

func TestDetectLogLevel(t *testing.T) {
	tests := []struct {
		msg  string
		want string
	}{
		{"Error: something failed", "ERROR"},
		{"Fatal: crash", "ERROR"},
		{"Warning: low disk", "WARN"},
		{"Debug: init", "DEBUG"},
		{"Info: started", "INFO"},
	}

	for _, tt := range tests {
		got := detectLogLevel(tt.msg)
		if got != tt.want {
			t.Errorf("detectLogLevel(%q) = %q, want %q", tt.msg, got, tt.want)
		}
	}
}

func TestSaveLogs(t *testing.T) {
	s := newTestSession()
	logs := []LogEntry{
		{Level: "INFO", Message: "test log"},
	}

	path, err := s.SaveLogs(logs, t.TempDir())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path == "" {
		t.Error("expected non-empty path")
	}
}

func TestRecovery(t *testing.T) {
	s := newTestSession()
	result := s.Recover()

	if result.Success {
		t.Log("recovery reported success (may succeed if flutter is available in CI)")
	} else {
		t.Log("recovery reported failure (expected if flutter not in PATH)")
	}
	if result.Attempts == 0 {
		t.Error("expected at least 1 attempt")
	}
}

func TestSnapshotFiles(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.dart"), []byte("void main() {}"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(dir, "lib")
	os.MkdirAll(subDir, 0755)
	if err := os.WriteFile(filepath.Join(subDir, "app.dart"), []byte("class App {}"), 0644); err != nil {
		t.Fatal(err)
	}

	w := NewFileWatcher(dir, 500)
	snapshot := w.snapshotFiles()

	if len(snapshot) == 0 {
		t.Error("expected at least one file in snapshot")
	}
}

func TestIgnoreDirs(t *testing.T) {
	dir := t.TempDir()
	buildDir := filepath.Join(dir, "build")
	os.MkdirAll(buildDir, 0755)
	if err := os.WriteFile(filepath.Join(buildDir, "output.dart"), []byte(""), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.dart"), []byte("void main() {}"), 0644); err != nil {
		t.Fatal(err)
	}

	w := NewFileWatcher(dir, 500)
	snapshot := w.snapshotFiles()

	for path := range snapshot {
		if filepath.HasPrefix(path, "build") {
			t.Errorf("expected build dir to be ignored, found %s", path)
		}
	}
}

func TestEventChannel(t *testing.T) {
	dir := t.TempDir()
	w := NewFileWatcher(dir, 100)

	if err := w.Start(); err != nil {
		t.Fatalf("failed to start watcher: %v", err)
	}
	defer w.Stop()

	events := w.Events()
	if events == nil {
		t.Error("expected event channel")
	}
}

func TestSessionInfo(t *testing.T) {
	s := newTestSession()
	s.mu.Lock()
	s.projectDir = "/test"
	s.deviceID = "device-1"
	s.flutterPID = 42
	s.mu.Unlock()

	info := s.SessionInfo()
	if info.ProjectDir != "/test" {
		t.Errorf("expected /test, got %s", info.ProjectDir)
	}
	if info.DeviceID != "device-1" {
		t.Errorf("expected device-1, got %s", info.DeviceID)
	}
	if info.FlutterPID != 42 {
		t.Errorf("expected 42, got %d", info.FlutterPID)
	}
}
