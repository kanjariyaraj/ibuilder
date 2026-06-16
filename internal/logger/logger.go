package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	level Level
	out   io.Writer
}

func New(level Level) *Logger {
	return &Logger{level: level, out: os.Stdout}
}

func NewWithWriter(level Level, w io.Writer) *Logger {
	return &Logger{level: level, out: w}
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) log(level Level, msg string, args ...any) {
	if level < l.level {
		return
	}
	timestamp := time.Now().Format(time.RFC3339)
	prefix := fmt.Sprintf("[%s] [%s]", timestamp, level)
	if len(args) > 0 {
		fmt.Fprintf(l.out, "%s %s %v\n", prefix, msg, args)
	} else {
		fmt.Fprintf(l.out, "%s %s\n", prefix, msg)
	}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log(LevelError, msg, args...)
}
