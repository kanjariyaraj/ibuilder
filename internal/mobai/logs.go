package mobai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Process   string    `json:"process"`
	Message   string    `json:"message"`
}

type LogFilter struct {
	Process string
	Level   string
	Search  string
	Since   time.Duration
}

var mockLogs = []LogEntry{
	{Timestamp: time.Now().Add(-10 * time.Minute), Level: "INFO", Process: "SpringBoard", Message: "Application launched: com.example.app"},
	{Timestamp: time.Now().Add(-9 * time.Minute), Level: "DEBUG", Process: "kernel", Message: "AppleMobileFileIntegrity: staging applied"},
	{Timestamp: time.Now().Add(-8 * time.Minute), Level: "INFO", Process: "app", Message: "Session started"},
	{Timestamp: time.Now().Add(-7 * time.Minute), Level: "WARN", Process: "app", Message: "Network reachability changed"},
	{Timestamp: time.Now().Add(-5 * time.Minute), Level: "ERROR", Process: "app", Message: "Failed to fetch remote config: timeout"},
	{Timestamp: time.Now().Add(-3 * time.Minute), Level: "INFO", Process: "SpringBoard", Message: "Application entered foreground"},
	{Timestamp: time.Now().Add(-1 * time.Minute), Level: "DEBUG", Process: "app", Message: "View did load: MainViewController"},
	{Timestamp: time.Now(), Level: "INFO", Process: "app", Message: "User logged in"},
}

func (c *Client) FetchLogs(filter *LogFilter) ([]LogEntry, error) {
	c.mu.RLock()
	connected := c.status.State == StateConnected
	c.mu.RUnlock()

	if !connected {
		return nil, fmt.Errorf("not connected to MobAI")
	}

	logs := mockLogs

	if filter != nil {
		if filter.Process != "" {
			var filtered []LogEntry
			for _, l := range logs {
				if strings.EqualFold(l.Process, filter.Process) {
					filtered = append(filtered, l)
				}
			}
			logs = filtered
		}
		if filter.Level != "" {
			var filtered []LogEntry
			for _, l := range logs {
				if strings.EqualFold(l.Level, filter.Level) {
					filtered = append(filtered, l)
				}
			}
			logs = filtered
		}
		if filter.Search != "" {
			var filtered []LogEntry
			for _, l := range logs {
				if strings.Contains(strings.ToLower(l.Message), strings.ToLower(filter.Search)) {
					filtered = append(filtered, l)
				}
			}
			logs = filtered
		}
		if filter.Since > 0 {
			since := time.Now().Add(-filter.Since)
			var filtered []LogEntry
			for _, l := range logs {
				if l.Timestamp.After(since) {
					filtered = append(filtered, l)
				}
			}
			logs = filtered
		}
	}

	return logs, nil
}

func (c *Client) StreamLogs(filter *LogFilter, stopChan chan struct{}) (chan LogEntry, error) {
	c.mu.RLock()
	connected := c.status.State == StateConnected
	c.mu.RUnlock()

	if !connected {
		return nil, fmt.Errorf("not connected to MobAI")
	}

	logChan := make(chan LogEntry)

	go func() {
		defer close(logChan)
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				entry := LogEntry{
					Timestamp: time.Now(),
					Level:     "INFO",
					Process:   "app",
					Message:   "Live log entry...",
				}
				logChan <- entry
			}
		}
	}()

	return logChan, nil
}

func (c *Client) SaveLogs(logs []LogEntry, outputDir string) (string, error) {
	if outputDir == "" {
		outputDir = ".build/logs"
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create log directory: %w", err)
	}

	filename := fmt.Sprintf("device_logs_%s.txt", c.Timestamp())
	path := filepath.Join(outputDir, filename)

	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create log file: %w", err)
	}
	defer f.Close()

	for _, l := range logs {
		line := fmt.Sprintf("[%s] [%s] [%s] %s\n",
			l.Timestamp.Format(time.RFC3339), l.Level, l.Process, l.Message)
		if _, err := f.WriteString(line); err != nil {
			return "", fmt.Errorf("failed to write log: %w", err)
		}
	}

	return path, nil
}
