package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// ConsoleLoggingService implements LoggingService using console output.
type ConsoleLoggingService struct {
	serviceName string
}

// NewConsoleLoggingService creates a new console logging service.
func NewConsoleLoggingService(serviceName string) *ConsoleLoggingService {
	return &ConsoleLoggingService{
		serviceName: serviceName,
	}
}

// Info logs an info message with structured data.
func (l *ConsoleLoggingService) Info(ctx context.Context, message string, fields map[string]interface{}) {
	l.logWithLevel("INFO", message, nil, fields)
}

// Error logs an error message with structured data.
func (l *ConsoleLoggingService) Error(ctx context.Context, message string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	if err != nil {
		fields["error"] = err.Error()
	}

	l.logWithLevel("ERROR", message, err, fields)
}

// Warn logs a warning message with structured data.
func (l *ConsoleLoggingService) Warn(ctx context.Context, message string, fields map[string]interface{}) {
	l.logWithLevel("WARN", message, nil, fields)
}

// logWithLevel logs a message with the specified level.
func (l *ConsoleLoggingService) logWithLevel(level, message string, err error, fields map[string]interface{}) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"level":     level,
		"service":   l.serviceName,
		"message":   message,
	}

	// Add error if present
	if err != nil {
		logEntry["error"] = err.Error()
	}

	// Add additional fields
	for k, v := range fields {
		logEntry[k] = v
	}

	// Convert to JSON for structured logging
	jsonBytes, err := json.Marshal(logEntry)
	if err != nil {
		// Fallback to simple logging if JSON marshal fails
		log.Printf("[%s] %s: %s (JSON marshal error: %v)", level, l.serviceName, message, err)
		return
	}

	// Output the structured log using log.Printf for consistency
	log.Printf("%s", string(jsonBytes))
}
