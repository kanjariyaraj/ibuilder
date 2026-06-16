package release

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kanjariyaraj/Builder/internal/logger"
)

type Session struct {
	mu         sync.RWMutex
	log        *logger.Logger
	projectDir string
	apiKey     string
}

func NewSession(log *logger.Logger) *Session {
	return &Session{
		log: log,
	}
}

func (s *Session) SetProjectDir(dir string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.projectDir = dir
}

func (s *Session) logInfo(msg string, args ...any) {
	if s.log == nil {
		return
	}
	s.log.Info(msg, args...)
}

func (s *Session) logWarn(msg string, args ...any) {
	if s.log == nil {
		return
	}
	s.log.Warn(msg, args...)
}

func (s *Session) ProjectDir() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.projectDir
}

func (s *Session) Timestamp() string {
	return time.Now().Format("20060102_150405")
}

func (s *Session) findIPA() (string, error) {
	dir := s.ProjectDir()
	if dir == "" {
		return "", fmt.Errorf("no project directory set")
	}

	buildDir := filepath.Join(dir, ".build")
	var ipaPath string
	err := filepath.Walk(buildDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".ipa") {
			ipaPath = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if ipaPath == "" {
		return "", fmt.Errorf("no .ipa file found in .build directory")
	}
	return ipaPath, nil
}
