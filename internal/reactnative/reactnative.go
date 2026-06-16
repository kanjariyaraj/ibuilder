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
	DeviceID     string `json:"device_id,omitempty"`
	MetroPort    int    `json:"metro_port"`
}

type Session struct {
	mu        sync.RWMutex
	cfg       *config.ReactNativeSettings
	log       *logger.Logger
	state     SessionState
	projectDir string
	deviceID  string
	metroPID  int
	metroPort int
	startedAt time.Time
}

func NewSession(cfg *config.ReactNativeSettings, log *logger.Logger) *Session {
	p := 8081
	if cfg.MetroPort > 0 {
		p = cfg.MetroPort
	}
	return &Session{
		cfg:       cfg,
		log:       log,
		state:     SessionInactive,
		metroPort: p,
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
	if stat, err := os.Stat(iosDir); err != nil || !stat.IsDir() {
		return false, fmt.Errorf("no ios/ directory found — not an RN iOS project")
	}
	s.projectDir = dir
	return true, nil
}

func (s *Session) CheckNode() (*RNInfo, error) {
	info := &RNInfo{}

	cmd := exec.Command("node", "--version")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}
	info.NodeVersion = strings.TrimSpace(string(output))

	npmCmd := exec.Command("npm", "--version")
	npmOut, err := npmCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("npm not found: %w", err)
	}
	info.NPMVersion = strings.TrimSpace(string(npmOut))

	if s.projectDir != "" {
		pkg := filepath.Join(s.projectDir, "package.json")
		if data, err := os.ReadFile(pkg); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if strings.Contains(line, "react-native") && strings.Contains(line, ":") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) > 1 {
						info.RNVersion = strings.Trim(strings.TrimSpace(parts[1]), "\",")
					}
				}
			}
		}
	}

	info.MetroPort = s.metroPort
	return info, nil
}

func (s *Session) ResolveDependencies() error {
	s.log.Info("resolving npm dependencies")
	if _, err := os.Stat(filepath.Join(s.projectDir, "node_modules")); os.IsNotExist(err) {
		cmd := exec.Command("npm", "install")
		cmd.Dir = s.projectDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("npm install failed: %s: %w", strings.TrimSpace(string(output)), err)
		}
		s.log.Info("dependencies installed")
	}
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
