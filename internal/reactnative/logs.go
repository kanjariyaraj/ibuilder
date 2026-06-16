package reactnative

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Tag       string    `json:"tag,omitempty"`
}

type LogFilter struct {
	Level  string
	Search string
	Since  time.Duration
	Tag    string
}

func (s *Session) FetchLogs(filter *LogFilter) ([]LogEntry, error) {
	s.mu.RLock()
	projectDir := s.projectDir
	s.mu.RUnlock()

	if projectDir == "" {
		return nil, fmt.Errorf("no project directory set")
	}

	args := []string{"react-native", "log-ios", "--json"}
	cmd := exec.Command("npx", args...)
	cmd.Dir = projectDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RN logs: %w", err)
	}

	logs := parseLogOutput(string(output), filter)
	return logs, nil
}

func (s *Session) StreamLogs(stopChan chan struct{}) (chan LogEntry, error) {
	s.mu.RLock()
	projectDir := s.projectDir
	s.mu.RUnlock()

	if projectDir == "" {
		return nil, fmt.Errorf("no project directory set")
	}

	logChan := make(chan LogEntry)

	go func() {
		defer close(logChan)

		args := []string{"react-native", "log-ios", "--json", "--follow"}
		cmd := exec.Command("npx", args...)
		cmd.Dir = projectDir

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return
		}

		if err := cmd.Start(); err != nil {
			return
		}

		buf := make([]byte, 4096)
		for {
			select {
			case <-stopChan:
				cmd.Process.Kill()
				return
			default:
				n, err := stdout.Read(buf)
				if err != nil {
					return
				}
				lines := strings.Split(string(buf[:n]), "\n")
				for _, line := range lines {
					if line == "" {
						continue
					}
					logChan <- LogEntry{
						Timestamp: time.Now(),
						Level:     detectLogLevel(line),
						Message:   line,
					}
				}
			}
		}
	}()

	return logChan, nil
}

func (s *Session) SaveLogs(logs []LogEntry, outputDir string) (string, error) {
	if outputDir == "" {
		outputDir = ".build/logs/react-native"
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create log directory: %w", err)
	}

	filename := fmt.Sprintf("rn_logs_%s.txt", s.Timestamp())
	path := filepath.Join(outputDir, filename)

	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create log file: %w", err)
	}
	defer f.Close()

	for _, l := range logs {
		line := fmt.Sprintf("[%s] [%s] %s\n",
			l.Timestamp.Format(time.RFC3339), l.Level, l.Message)
		if _, err := f.WriteString(line); err != nil {
			return "", fmt.Errorf("failed to write log: %w", err)
		}
	}

	return path, nil
}

func parseLogOutput(output string, filter *LogFilter) []LogEntry {
	var logs []LogEntry

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		entry := LogEntry{
			Timestamp: time.Now(),
			Level:     detectLogLevel(line),
			Message:   line,
		}

		if filter != nil {
			if filter.Level != "" && !strings.EqualFold(entry.Level, filter.Level) {
				continue
			}
			if filter.Search != "" && !strings.Contains(strings.ToLower(entry.Message), strings.ToLower(filter.Search)) {
				continue
			}
			if filter.Since > 0 && time.Since(entry.Timestamp) > filter.Since {
				continue
			}
			if filter.Tag != "" && !strings.Contains(entry.Message, filter.Tag) {
				continue
			}
		}

		logs = append(logs, entry)
	}

	return logs
}

func detectLogLevel(msg string) string {
	upper := strings.ToUpper(msg)
	switch {
	case strings.Contains(upper, "ERROR") || strings.Contains(upper, "FATAL"):
		return "ERROR"
	case strings.Contains(upper, "WARN"):
		return "WARN"
	case strings.Contains(upper, "DEBUG"):
		return "DEBUG"
	default:
		return "INFO"
	}
}
