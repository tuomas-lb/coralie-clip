// Package logging provides structured logging for the application.
package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Logger handles structured logging.
type Logger struct {
	logFile *os.File
}

// NewLogger creates a new logger.
func NewLogger(logPath string) (*Logger, error) {
	if logPath == "" {
		return &Logger{}, nil
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &Logger{logFile: file}, nil
}

// LogRequest logs an API request.
func (l *Logger) LogRequest(logEntry interface{}) error {
	if l.logFile == nil {
		return nil
	}

	data, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	if _, err := l.logFile.WriteString(string(data) + "\n"); err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}

	return l.logFile.Sync()
}

// Close closes the log file.
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// RedactSecret redacts a secret value for logging.
func RedactSecret(secret string) string {
	if secret == "" {
		return ""
	}
	if len(secret) <= 8 {
		return "***"
	}
	return secret[:4] + "***" + secret[len(secret)-4:]
}

// FormatDuration formats a duration as a string.
func FormatDuration(d time.Duration) string {
	return d.String()
}

