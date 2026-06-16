package reactnative

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
	SessionMetroRunning
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
	case SessionMetroRunning:
		return "Metro Running"
	case SessionRecovering:
		return "Recovering"
	default:
		return "Unknown"
	}
}

type RNInfo struct {
	NodeVersion  string `json:"node_version"`
	NPMVersion   string `json:"npm_version"`
	RNVersion    string `json:"rn_version"`
	ProjectValid bool   `json:"project_valid"`
	MetroReady   bool   `json:"metro_ready"`
	DeviceID     string `json:"device_id,omitempty"`
}

type Session struct {
	mu         sync.RWMutex
	cfg        *config.ReactNativeSettings
	log        *logger.Logger
	state      SessionState
	projectDir string
	deviceID   string
	metroPID   int
	metroPort  int
	startedAt  time.Time
}

func NewSession(cfg *config.ReactNativeSettings, log *logger.Logger) *Session {
	port := 8081
	if cfg.MetroPort > 0 {
		port = cfg.MetroPort
	}
	return &Session{
		cfg:       cfg,
		log:       log,
		state:     SessionInactive,
		metroPort: port,
	}
}

func (s *Session) State() SessionState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

func (s *Session) Config() *config.ReactNativeSettings {
	return s.cfg
}

func (s *Session) UpdateConfig(cfg *config.ReactNativeSettings) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cfg = cfg
}

func (s *Session) DetectRNProject(dir string) (bool, error) {
	pkg := filepath.Join(dir, "package.json")
	if _, err := os.Stat(pkg); os.IsNotExist(err) {
		return false, fmt.Errorf("no package.json found in %s", dir)
	}

	data, err := os.ReadFile(pkg)
	if err != nil {
		return false, fmt.Errorf("failed to read package.json: %w", err)
	}

	if !strings.Contains(string(data), "react-native") {
		return false, fmt.Errorf("react-native dependency not found in package.json")
	}

	iosDir := filepath.Join(dir, "ios")
	if _, err := os.Stat(iosDir); os.IsNotExist(err) {
		return false, fmt.Errorf("no ios/ directory found — not a React Native iOS project")
	}

	s.projectDir = dir
	return true, nil
}

func (s *Session) CheckNode() (string, error) {
	cmd := exec.Command("node", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("node not found: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (s *Session) CheckNPM() (string, error) {
	cmd := exec.Command("npm", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("npm not found: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (s *Session) CheckYarn() (string, error) {
	cmd := exec.Command("yarn", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(output)), nil
}

func (s *Session) CheckPNPM() (string, error) {
	cmd := exec.Command("pnpm", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(output)), nil
}

func (s *Session) CheckMetro() bool {
	cmd := exec.Command("npx", "react-native", "--help")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (s *Session) ResolveDependencies() error {
	s.log.Info("resolving npm dependencies")

	pkg := filepath.Join(s.projectDir, "package.json")
	data, err := os.ReadFile(pkg)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	usesYarn := strings.Contains(string(data), "yarn") || s.cfg.Entry == "yarn"

	var cmd *exec.Cmd
	if usesYarn {
		cmd = exec.Command("yarn", "install")
	} else {
		cmd = exec.Command("npm", "install")
	}
	cmd.Dir = s.projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("dependency resolution failed: %s: %w", strings.TrimSpace(string(output)), err)
	}
	s.log.Info("dependencies resolved")
	return nil
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
