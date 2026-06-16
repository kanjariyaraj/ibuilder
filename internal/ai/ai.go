package ai

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
	mu           sync.RWMutex
	log          *logger.Logger
	projectDir   string
	analyzer     *Analyzer
	knowledgeBase *KnowledgeBase
}

func NewSession(log *logger.Logger) *Session {
	return &Session{
		log:          log,
		analyzer:     NewAnalyzer(),
		knowledgeBase: NewKnowledgeBase(),
	}
}

func (s *Session) SetProjectDir(dir string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.projectDir = dir
}

func (s *Session) ProjectDir() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.projectDir
}

func (s *Session) Timestamp() string {
	return time.Now().Format("20060102_150405")
}

func (s *Session) collectLogs(dir string) ([]string, error) {
	var logs []string

	buildLogs := filepath.Join(dir, ".build", "logs")
	err := filepath.Walk(buildLogs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".log") {
			data, err := os.ReadFile(path)
			if err == nil {
				logs = append(logs, fmt.Sprintf("--- %s ---\n%s", path, string(data)))
			}
		}
		return nil
	})
	if err != nil {
		return logs, nil
	}

	return logs, nil
}
