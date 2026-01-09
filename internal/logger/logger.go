// Package logger provides structured logging for the application.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Level represents the logging level.
type Level int

const (
	// LevelDebug is for debug messages.
	LevelDebug Level = iota
	// LevelInfo is for informational messages.
	LevelInfo
	// LevelWarn is for warning messages.
	LevelWarn
	// LevelError is for error messages.
	LevelError
)

// Logger provides structured logging capabilities.
type Logger struct {
	mu     sync.Mutex
	level  Level
	logger *log.Logger
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// New creates a new logger with the specified level and output.
func New(level Level, out io.Writer) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(out, "", log.LstdFlags),
	}
}

// Default returns the default logger.
func Default() *Logger {
	once.Do(func() {
		defaultLogger = New(LevelInfo, os.Stdout)
	})
	return defaultLogger
}

// SetLevel sets the logging level.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// Debug logs a debug message.
func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(LevelDebug, "DEBUG", format, v...)
}

// Info logs an informational message.
func (l *Logger) Info(format string, v ...interface{}) {
	l.log(LevelInfo, "INFO", format, v...)
}

// Warn logs a warning message.
func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(LevelWarn, "WARN", format, v...)
}

// Error logs an error message.
func (l *Logger) Error(format string, v ...interface{}) {
	l.log(LevelError, "ERROR", format, v...)
}

func (l *Logger) log(level Level, prefix, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	msg := fmt.Sprintf(format, v...)
	l.logger.Printf("[%s] %s", prefix, msg)
}

// Convenience functions using the default logger

// Debug logs a debug message using the default logger.
func Debug(format string, v ...interface{}) {
	Default().Debug(format, v...)
}

// Info logs an informational message using the default logger.
func Info(format string, v ...interface{}) {
	Default().Info(format, v...)
}

// Warn logs a warning message using the default logger.
func Warn(format string, v ...interface{}) {
	Default().Warn(format, v...)
}

// Error logs an error message using the default logger.
func Error(format string, v ...interface{}) {
	Default().Error(format, v...)
}
