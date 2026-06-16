package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	l := New(LevelInfo)
	if l.level != LevelInfo {
		t.Errorf("expected LevelInfo, got %v", l.level)
	}
}

func TestLevelStrings(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
	}

	for _, tt := range tests {
		if tt.level.String() != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.level.String())
		}
	}
}

func TestLogOutput(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(LevelDebug, &buf)

	l.Info("test message")
	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("expected output to contain INFO, got %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got %s", output)
	}
}

func TestLogLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(LevelError, &buf)

	l.Info("should not appear")
	if buf.Len() > 0 {
		t.Errorf("expected no output for Info log at Error level")
	}
}

func TestLogWithArgs(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(LevelInfo, &buf)

	l.Debug("debug message", "key", "value")
	if buf.Len() > 0 {
		t.Errorf("expected no output for Debug log at Info level")
	}
}
