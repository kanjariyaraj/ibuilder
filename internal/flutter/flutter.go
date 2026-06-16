package flutter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/logger"
)

type SessionState int

const (
	SessionInactive SessionState = iota
	SessionStarting
	SessionActive
	SessionAttached
	SessionWatchMode
	SessionRecovering
)

func (s SessionState) String() string {
	switch s {
	case SessionInactive:
		return "Inactive"
	case SessionStarting:
		return "Starting"
	case SessionActive:
		return "Active"
	case SessionAttached:
		return "Attached"
	case SessionWatchMode:
		return "Watch Mode"
	case SessionRecovering:
		return "Recovering"
	default:
		return "Unknown"
	}
}

type FlutterInfo struct {
	Version      string `json:"version"`
	DartVersion  string `json:"dart_version"`
	Channel      string `json:"channel"`
	ProjectValid bool   `json:"project_valid"`
	DeviceID     string `json:"device_id,omitempty"`
}

type Session struct {
	mu         sync.RWMutex
	cfg        *config.FlutterSettings
	log        *logger.Logger
	state      SessionState
	projectDir string
	deviceID   string
	watchStop  chan struct{}
	flutterPID int
	startedAt  time.Time
}

func NewSession(cfg *config.FlutterSettings, log *logger.Logger) *Session {
	return &Session{
		cfg:       cfg,
		log:       log,
		state:     SessionInactive,
		watchStop: make(chan struct{}),
	}
}

func (s *Session) State() SessionState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

func (s *Session) Config() *config.FlutterSettings {
	return s.cfg
}

func (s *Session) UpdateConfig(cfg *config.FlutterSettings) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cfg = cfg
}

func (s *Session) DetectFlutterProject(dir string) (bool, error) {
	pubspec := filepath.Join(dir, "pubspec.yaml")
	if _, err := os.Stat(pubspec); os.IsNotExist(err) {
		return false, fmt.Errorf("no pubspec.yaml found in %s", dir)
	}
	iosDir := filepath.Join(dir, "ios")
	iosStat, err := os.Stat(iosDir)
	if err != nil || !iosStat.IsDir() {
		return false, fmt.Errorf("no ios/ directory found — not a Flutter iOS project")
	}
	s.projectDir = dir
	return true, nil
}

func (s *Session) CheckFlutterSDK() (*FlutterInfo, error) {
	cmd := exec.Command("flutter", "--version")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("flutter SDK not found: %w", err)
	}

	info := &FlutterInfo{}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) > 1 {
			info.Version = parts[1]
		}
	}
	for _, line := range lines {
		if strings.Contains(line, "Dart") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				info.DartVersion = parts[1]
			}
		}
		if strings.Contains(line, "channel") {
			parts := strings.Fields(line)
			for i, p := range parts {
				if p == "channel" && i+1 < len(parts) {
					info.Channel = parts[i+1]
				}
			}
		}
	}

	if info.Version == "" {
		info.Version = "unknown"
	}

	return info, nil
}

func (s *Session) CheckDartSDK() error {
	cmd := exec.Command("dart", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dart SDK not found: %w", err)
	}
	return nil
}

func (s *Session) ResolveDependencies() error {
	s.log.Info("resolving flutter dependencies")
	cmd := exec.Command("flutter", "pub", "get")
	cmd.Dir = s.projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("dependency resolution failed: %s: %w", strings.TrimSpace(string(output)), err)
	}
	s.log.Info("dependencies resolved")
	return nil
}

func (s *Session) ListDevices() ([]string, error) {
	cmd := exec.Command("flutter", "devices", "--machine")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list flutter devices: %w", err)
	}
	var devices []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			devices = append(devices, line)
		}
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("no Flutter devices found")
	}
	return devices, nil
}

func (s *Session) DeviceID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.deviceID
}

func (s *Session) SetDeviceID(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deviceID = id
}

func (s *Session) ProjectDir() string {
	return s.projectDir
}

func (s *Session) Timestamp() string {
	return time.Now().Format("20060102_150405")
}
