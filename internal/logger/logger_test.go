package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name       string
		level      Level
		logFunc    func(*Logger, string)
		message    string
		shouldLog  bool
		wantPrefix string
	}{
		{
			name:       "info level logs info",
			level:      LevelInfo,
			logFunc:    func(l *Logger, m string) { l.Info(m) },
			message:    "test info message",
			shouldLog:  true,
			wantPrefix: "[INFO]",
		},
		{
			name:       "info level skips debug",
			level:      LevelInfo,
			logFunc:    func(l *Logger, m string) { l.Debug(m) },
			message:    "test debug message",
			shouldLog:  false,
			wantPrefix: "[DEBUG]",
		},
		{
			name:       "debug level logs debug",
			level:      LevelDebug,
			logFunc:    func(l *Logger, m string) { l.Debug(m) },
			message:    "test debug message",
			shouldLog:  true,
			wantPrefix: "[DEBUG]",
		},
		{
			name:       "error level logs error",
			level:      LevelInfo,
			logFunc:    func(l *Logger, m string) { l.Error(m) },
			message:    "test error message",
			shouldLog:  true,
			wantPrefix: "[ERROR]",
		},
		{
			name:       "warn level logs warn",
			level:      LevelInfo,
			logFunc:    func(l *Logger, m string) { l.Warn(m) },
			message:    "test warn message",
			shouldLog:  true,
			wantPrefix: "[WARN]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := New(tt.level, buf)

			tt.logFunc(logger, tt.message)

			output := buf.String()

			if tt.shouldLog {
				if !strings.Contains(output, tt.wantPrefix) {
					t.Errorf("Logger output missing prefix %s: %s", tt.wantPrefix, output)
				}
				if !strings.Contains(output, tt.message) {
					t.Errorf("Logger output missing message %s: %s", tt.message, output)
				}
			} else {
				if output != "" {
					t.Errorf("Logger should not have logged anything, but got: %s", output)
				}
			}
		})
	}
}

func TestSetLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(LevelInfo, buf)

	// Should not log debug messages at info level
	logger.Debug("debug message 1")
	if buf.Len() > 0 {
		t.Errorf("Expected no output at info level, but got: %s", buf.String())
	}

	// Change to debug level
	logger.SetLevel(LevelDebug)
	buf.Reset()

	// Should now log debug messages
	logger.Debug("debug message 2")
	if buf.Len() == 0 {
		t.Errorf("Expected output at debug level, but got none")
	}

	if !strings.Contains(buf.String(), "debug message 2") {
		t.Errorf("Expected debug message 2 in output, but got: %s", buf.String())
	}
}
