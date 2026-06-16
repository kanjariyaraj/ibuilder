package mobai

import (
	"testing"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/logger"
)

func newTestClient() *Client {
	cfg := &config.MobaiSettings{
		Host:              "localhost",
		Port:              12345,
		AutoReconnect:     false,
		ConnectionTimeout: 5,
	}
	log := logger.New(logger.LevelDebug)
	return NewClient(cfg, log)
}

func TestNewClient(t *testing.T) {
	client := newTestClient()
	if client == nil {
		t.Fatal("expected client, got nil")
	}
	status := client.Status()
	if status.State != StateDisconnected {
		t.Errorf("expected disconnected, got %v", status.State)
	}
}

func TestConfig(t *testing.T) {
	client := newTestClient()
	cfg := client.Config()
	if cfg.Host != "localhost" {
		t.Errorf("expected localhost, got %s", cfg.Host)
	}
	if cfg.Port != 12345 {
		t.Errorf("expected 12345, got %d", cfg.Port)
	}
}

func TestUpdateConfig(t *testing.T) {
	client := newTestClient()
	newCfg := &config.MobaiSettings{
		Host: "192.168.1.100",
		Port: 8080,
	}
	client.UpdateConfig(newCfg)
	cfg := client.Config()
	if cfg.Host != "192.168.1.100" {
		t.Errorf("expected 192.168.1.100, got %s", cfg.Host)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected 8080, got %d", cfg.Port)
	}
}

func TestConnectionStateString(t *testing.T) {
	tests := []struct {
		state ConnectionState
		want  string
	}{
		{StateDisconnected, "Disconnected"},
		{StateConnecting, "Connecting"},
		{StateConnected, "Connected"},
		{StateReconnecting, "Reconnecting"},
		{ConnectionState(99), "Unknown"},
	}

	for _, tt := range tests {
		got := tt.state.String()
		if got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}

func TestDoctor(t *testing.T) {
	client := newTestClient()
	report := client.Doctor()

	if len(report.Checks) == 0 {
		t.Error("expected at least one check")
	}

	hasConfig := false
	hasConnectivity := false
	hasDevice := false
	for _, check := range report.Checks {
		switch check.Name {
		case "Configuration":
			hasConfig = true
		case "Connectivity":
			hasConnectivity = true
		case "Device":
			hasDevice = true
		}
	}
	if !hasConfig {
		t.Error("expected configuration check")
	}
	if !hasConnectivity {
		t.Error("expected connectivity check")
	}
	if !hasDevice {
		t.Error("expected device check")
	}
}

func TestDisconnect(t *testing.T) {
	client := newTestClient()
	if err := client.Disconnect(); err != nil {
		t.Errorf("disconnect should not error when not connected: %v", err)
	}
	status := client.Status()
	if status.State != StateDisconnected {
		t.Errorf("expected disconnected, got %v", status.State)
	}
}

func TestPingWhenDisconnected(t *testing.T) {
	client := newTestClient()
	_, err := client.Ping()
	if err == nil {
		t.Error("expected error when pinging disconnected client")
	}
}

func TestConnectWithInvalidHost(t *testing.T) {
	cfg := &config.MobaiSettings{
		Host:              "192.0.2.1",
		Port:              1,
		AutoReconnect:     false,
		ConnectionTimeout: 1,
	}
	log := logger.New(logger.LevelDebug)
	client := NewClient(cfg, log)

	err := client.Connect()
	if err == nil {
		t.Error("expected connection error with invalid host")
	}
}

func TestDoctorWithWarning(t *testing.T) {
	client := newTestClient()
	report := client.Doctor()

	for _, check := range report.Checks {
		if check.Name == "Device" {
			if check.Status != StatusWarning {
				t.Errorf("expected device check to be warning when disconnected, got %s", check.Status)
			}
		}
	}
}

func TestListDevicesDisconnected(t *testing.T) {
	client := newTestClient()
	_, err := client.ListDevices()
	if err == nil {
		t.Error("expected error listing devices when disconnected")
	}
}

func TestFetchLogsDisconnected(t *testing.T) {
	client := newTestClient()
	_, err := client.FetchLogs(nil)
	if err == nil {
		t.Error("expected error fetching logs when disconnected")
	}
}

func TestScreenshotDisconnected(t *testing.T) {
	client := newTestClient()
	_, err := client.CaptureScreenshot("")
	if err == nil {
		t.Error("expected error capturing screenshot when disconnected")
	}
}

func TestInstallDisconnected(t *testing.T) {
	client := newTestClient()
	_, err := client.InstallIPA("test.ipa")
	if err == nil {
		t.Error("expected error installing when disconnected")
	}
}

func TestLaunchDisconnected(t *testing.T) {
	client := newTestClient()
	_, err := client.LaunchApp("com.example")
	if err == nil {
		t.Error("expected error launching when disconnected")
	}
}

func TestAutoReconnectDefault(t *testing.T) {
	client := newTestClient()
	if client.IsAutoReconnectEnabled() {
		t.Error("expected auto reconnect to be disabled")
	}
}

func TestSetAutoReconnect(t *testing.T) {
	client := newTestClient()
	client.SetAutoReconnect(true)
	if !client.IsAutoReconnectEnabled() {
		t.Error("expected auto reconnect to be enabled")
	}
	client.SetAutoReconnect(false)
	if client.IsAutoReconnectEnabled() {
		t.Error("expected auto reconnect to be disabled")
	}
}

func TestReconnectWhenConnected(t *testing.T) {
	client := newTestClient()
	client.mu.Lock()
	client.status.State = StateConnected
	client.mu.Unlock()

	err := client.Reconnect()
	if err != nil {
		t.Errorf("expected no error reconnecting when already connected, got %v", err)
	}
}

func TestSessionRestoreDisconnected(t *testing.T) {
	client := newTestClient()
	err := client.SessionRestore()
	if err == nil {
		t.Error("expected error restoring session without config")
	}
}

func TestFormatDeviceInfo(t *testing.T) {
	device := &DeviceInfo{
		Name:      "Test iPhone",
		Model:     "iPhone 15",
		OSVersion: "17.4",
		UDID:      "test-udid",
		State:     "available",
		Battery:   85,
		Storage:   "128GB / 256GB",
		Developer: true,
		Network:   "Wi-Fi",
	}

	output := FormatDeviceInfo(device)
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"this is a long string", 10, "this is..."},
		{"exactly", 7, "exactly"},
	}

	for _, tt := range tests {
		got := truncate(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}

func TestDeviceInfoWithUDID(t *testing.T) {
	client := newTestClient()
	client.mu.Lock()
	client.status.State = StateConnected
	client.mu.Unlock()

	device, err := client.DeviceInfo("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent UDID")
	}
	if device != nil {
		t.Error("expected nil device for nonexistent UDID")
	}
}

func TestSaveLogs(t *testing.T) {
	client := newTestClient()
	logs := []LogEntry{
		{Level: "INFO", Process: "test", Message: "test message"},
	}

	path, err := client.SaveLogs(logs, t.TempDir())
	if err != nil {
		t.Fatalf("expected no error saving logs, got %v", err)
	}
	if path == "" {
		t.Error("expected non-empty path")
	}
}

func TestLogFilter(t *testing.T) {
	client := newTestClient()
	client.mu.Lock()
	client.status.State = StateConnected
	client.mu.Unlock()

	logs, err := client.FetchLogs(&LogFilter{
		Level: "ERROR",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(logs) == 0 {
		t.Error("expected at least one ERROR log")
	}
	for _, l := range logs {
		if l.Level != "ERROR" {
			t.Errorf("expected ERROR level, got %s", l.Level)
		}
	}
}

func TestInstallInvalidExt(t *testing.T) {
	client := newTestClient()
	client.mu.Lock()
	client.status.State = StateConnected
	client.mu.Unlock()

	_, err := client.InstallIPA("test.apk")
	if err == nil {
		t.Error("expected error for non-ipa file")
	}
}

func TestLaunchEmptyBundleID(t *testing.T) {
	client := newTestClient()
	client.mu.Lock()
	client.status.State = StateConnected
	client.mu.Unlock()

	_, err := client.LaunchApp("")
	if err == nil {
		t.Error("expected error for empty bundle ID")
	}
}

func TestTerminateApp(t *testing.T) {
	client := newTestClient()
	client.mu.Lock()
	client.status.State = StateConnected
	client.mu.Unlock()

	err := client.TerminateApp("com.example")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestIsAppRunning(t *testing.T) {
	client := newTestClient()
	client.mu.Lock()
	client.status.State = StateConnected
	client.mu.Unlock()

	running, err := client.IsAppRunning("com.example")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !running {
		t.Error("expected app to be running")
	}
}

func TestVerifyInstallation(t *testing.T) {
	client := newTestClient()
	client.mu.Lock()
	client.status.State = StateConnected
	client.mu.Unlock()

	ok, err := client.VerifyInstallation("com.example")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !ok {
		t.Error("expected installation verification to succeed")
	}
}

func TestInstalledApps(t *testing.T) {
	client := newTestClient()
	client.mu.Lock()
	client.status.State = StateConnected
	client.mu.Unlock()

	apps, err := client.InstalledApps()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(apps) == 0 {
		t.Error("expected at least one installed app")
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Raj's iPhone", "Rajs_iPhone"},
		{"normal-name", "normal-name"},
		{"", "device"},
	}

	for _, tt := range tests {
		got := sanitizeFilename(tt.input)
		if got != tt.want {
			t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
