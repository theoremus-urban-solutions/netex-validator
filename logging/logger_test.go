package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

type contextKey string

func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	config := LoggerConfig{
		Level:         LevelInfo,
		Format:        "json",
		Output:        &buf,
		IncludeSource: false,
		Component:     "test-component",
	}

	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}

	logger.Info("test message", "key", "value")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected log output to contain 'test message', got: %s", output)
	}

	if !strings.Contains(output, "test-component") {
		t.Errorf("Expected log output to contain component name, got: %s", output)
	}
}

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{LogLevel(999), "UNKNOWN"},
	}

	for _, test := range tests {
		if got := test.level.String(); got != test.expected {
			t.Errorf("LogLevel(%d).String() = %s, want %s", test.level, got, test.expected)
		}
	}
}

func TestNewDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger()
	if logger == nil {
		t.Fatal("NewDefaultLogger returned nil")
	}

	// Test that it doesn't panic
	logger.Info("test message")
}

func TestNewJSONLogger(t *testing.T) {
	var buf bytes.Buffer
	// Temporarily replace stdout to capture output
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	logger := NewLogger(LoggerConfig{
		Level:  LevelInfo,
		Format: "json",
		Output: &buf,
	})

	logger.Info("test json message", "key", "value")

	output := buf.String()

	// Verify it's valid JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &jsonData); err != nil {
		t.Errorf("Output is not valid JSON: %v\nOutput: %s", err, output)
	}

	if jsonData["msg"] != "test json message" {
		t.Errorf("Expected message 'test json message', got: %v", jsonData["msg"])
	}

	if jsonData["key"] != "value" {
		t.Errorf("Expected key 'value', got: %v", jsonData["key"])
	}
}

func TestNewDebugLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LoggerConfig{
		Level:         LevelDebug,
		Format:        "text",
		Output:        &buf,
		IncludeSource: true,
	})

	logger.Debug("debug message")

	output := buf.String()
	if !strings.Contains(output, "debug message") {
		t.Errorf("Expected debug message in output, got: %s", output)
	}
}

func TestLogger_WithMethods(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LoggerConfig{
		Level:  LevelInfo,
		Format: "json",
		Output: &buf,
	})

	// Test WithFile
	fileLogger := logger.WithFile("test.xml")
	fileLogger.Info("file test")

	output := buf.String()
	if !strings.Contains(output, "test.xml") {
		t.Errorf("Expected filename in output, got: %s", output)
	}

	// Reset buffer
	buf.Reset()

	// Test WithValidation
	validationLogger := logger.WithValidation("val-123", "NO")
	validationLogger.Info("validation test")

	output = buf.String()
	if !strings.Contains(output, "val-123") || !strings.Contains(output, "NO") {
		t.Errorf("Expected validation context in output, got: %s", output)
	}

	// Reset buffer
	buf.Reset()

	// Test WithRule
	ruleLogger := logger.WithRule("LINE_1", "Line missing name")
	ruleLogger.Info("rule test")

	output = buf.String()
	if !strings.Contains(output, "LINE_1") {
		t.Errorf("Expected rule code in output, got: %s", output)
	}

	// Reset buffer
	buf.Reset()

	// Test WithError
	err := errors.New("test error")
	errorLogger := logger.WithError(err)
	errorLogger.Info("error test")

	output = buf.String()
	if !strings.Contains(output, "test error") {
		t.Errorf("Expected error message in output, got: %s", output)
	}

	// Reset buffer
	buf.Reset()

	// Test WithDuration
	duration := 150 * time.Millisecond
	durationLogger := logger.WithDuration("validation", duration)
	durationLogger.Info("duration test")

	output = buf.String()
	if !strings.Contains(output, "150") {
		t.Errorf("Expected duration in output, got: %s", output)
	}

	// Reset buffer
	buf.Reset()

	// Test WithMetrics
	metricsLogger := logger.WithMetrics(5, 12, 500*time.Millisecond)
	metricsLogger.Info("metrics test")

	output = buf.String()
	if !strings.Contains(output, "\"files_processed\":5") {
		t.Errorf("Expected metrics in output, got: %s", output)
	}
}

func TestLogger_ValidationMethods(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LoggerConfig{
		Level:  LevelInfo,
		Format: "json",
		Output: &buf,
	})

	filename := "test.xml"
	codespace := "NO"
	duration := 100 * time.Millisecond

	// Test ValidationStart
	logger.ValidationStart(filename, codespace)
	output := buf.String()
	if !strings.Contains(output, "Starting validation") {
		t.Errorf("Expected validation start message, got: %s", output)
	}
	buf.Reset()

	// Test ValidationComplete
	logger.ValidationComplete(filename, duration, 5, true)
	output = buf.String()
	if !strings.Contains(output, "Validation completed") {
		t.Errorf("Expected validation complete message, got: %s", output)
	}
	buf.Reset()

	// Test ValidationError
	err := errors.New("validation failed")
	logger.ValidationError(filename, err)
	output = buf.String()
	if !strings.Contains(output, "Validation error") {
		t.Errorf("Expected validation error message, got: %s", output)
	}
	buf.Reset()

	// Test RuleViolation
	logger.RuleViolation(filename, "LINE_1", "Line missing name", "Name element is required", 42)
	output = buf.String()
	if !strings.Contains(output, "Rule violation") || !strings.Contains(output, "LINE_1") {
		t.Errorf("Expected rule violation message, got: %s", output)
	}
	buf.Reset()

	// Test PerformanceWarning
	threshold := 50 * time.Millisecond
	logger.PerformanceWarning("validation", duration, threshold)
	output = buf.String()
	if !strings.Contains(output, "Performance warning") {
		t.Errorf("Expected performance warning message, got: %s", output)
	}
	buf.Reset()

	// Test BatchValidationStart
	logger.BatchValidationStart(10)
	output = buf.String()
	if !strings.Contains(output, "Starting batch validation") {
		t.Errorf("Expected batch validation start message, got: %s", output)
	}
	buf.Reset()

	// Test BatchValidationComplete
	logger.BatchValidationComplete(10, 8, 25, time.Second)
	output = buf.String()
	if !strings.Contains(output, "Batch validation completed") {
		t.Errorf("Expected batch validation complete message, got: %s", output)
	}
	buf.Reset()

	// Test ConfigurationLoaded
	logger.ConfigurationLoaded("config.yaml", 150)
	output = buf.String()
	if !strings.Contains(output, "Configuration loaded") {
		t.Errorf("Expected configuration loaded message, got: %s", output)
	}
}

func TestLogger_DebugMethods(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LoggerConfig{
		Level:  LevelDebug,
		Format: "json",
		Output: &buf,
	})

	filename := "test.xml"
	duration := 50 * time.Millisecond

	// Test SchemaValidationStart
	logger.SchemaValidationStart(filename)
	output := buf.String()
	if !strings.Contains(output, "Starting schema validation") {
		t.Errorf("Expected schema validation start message, got: %s", output)
	}
	buf.Reset()

	// Test SchemaValidationComplete
	logger.SchemaValidationComplete(filename, duration, true)
	output = buf.String()
	if !strings.Contains(output, "Schema validation completed") {
		t.Errorf("Expected schema validation complete message, got: %s", output)
	}
	buf.Reset()

	// Test XPathValidationStart
	logger.XPathValidationStart(filename, 150)
	output = buf.String()
	if !strings.Contains(output, "Starting XPath validation") {
		t.Errorf("Expected XPath validation start message, got: %s", output)
	}
	buf.Reset()

	// Test XPathValidationComplete
	logger.XPathValidationComplete(filename, duration, 3)
	output = buf.String()
	if !strings.Contains(output, "XPath validation completed") {
		t.Errorf("Expected XPath validation complete message, got: %s", output)
	}
	buf.Reset()

	// Test MemoryUsage
	logger.MemoryUsage("validation", 25.5, 128.0)
	output = buf.String()
	if !strings.Contains(output, "Memory usage") {
		t.Errorf("Expected memory usage message, got: %s", output)
	}
}

func TestLogger_IsLevelEnabled(t *testing.T) {
	logger := NewLogger(LoggerConfig{Level: LevelWarn})

	if !logger.IsLevelEnabled(LevelError) {
		t.Error("Expected ERROR level to be enabled for WARN logger")
	}

	if !logger.IsLevelEnabled(LevelWarn) {
		t.Error("Expected WARN level to be enabled for WARN logger")
	}

	if logger.IsLevelEnabled(LevelInfo) {
		t.Error("Expected INFO level to be disabled for WARN logger")
	}

	if logger.IsLevelEnabled(LevelDebug) {
		t.Error("Expected DEBUG level to be disabled for WARN logger")
	}
}

func TestGlobalLogger(t *testing.T) {
	var buf bytes.Buffer
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	// Set a test logger as default
	testLogger := NewLogger(LoggerConfig{
		Level:  LevelInfo,
		Format: "json",
		Output: &buf,
	})
	SetDefaultLogger(testLogger)

	if GetDefaultLogger() != testLogger {
		t.Error("GetDefaultLogger did not return the expected logger")
	}

	// Test global convenience functions
	Info("test info", "key", "value")
	output := buf.String()
	if !strings.Contains(output, "test info") {
		t.Errorf("Expected global Info to work, got: %s", output)
	}
	buf.Reset()

	Warn("test warning")
	output = buf.String()
	if !strings.Contains(output, "test warning") {
		t.Errorf("Expected global Warn to work, got: %s", output)
	}
	buf.Reset()

	Error("test error")
	output = buf.String()
	if !strings.Contains(output, "test error") {
		t.Errorf("Expected global Error to work, got: %s", output)
	}
	buf.Reset()

	ValidationStart("test.xml", "NO")
	output = buf.String()
	if !strings.Contains(output, "Starting validation") {
		t.Errorf("Expected global ValidationStart to work, got: %s", output)
	}
	buf.Reset()

	ValidationComplete("test.xml", 100*time.Millisecond, 2, true)
	output = buf.String()
	if !strings.Contains(output, "Validation completed") {
		t.Errorf("Expected global ValidationComplete to work, got: %s", output)
	}
	buf.Reset()

	ValidationError("test.xml", errors.New("test error"))
	output = buf.String()
	if !strings.Contains(output, "Validation error") {
		t.Errorf("Expected global ValidationError to work, got: %s", output)
	}
	buf.Reset()

	RuleViolation("test.xml", "RULE_1", "Test rule", "Test message", 10)
	output = buf.String()
	if !strings.Contains(output, "Rule violation") {
		t.Errorf("Expected global RuleViolation to work, got: %s", output)
	}
}

func TestWithContext(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LoggerConfig{
		Level:  LevelInfo,
		Format: "json",
		Output: &buf,
	})

	ctx := context.WithValue(context.Background(), contextKey("request_id"), "req-123")
	contextLogger := logger.WithContext(ctx)

	contextLogger.Info("context test")

	output := buf.String()
	// Note: context value might be nil if not properly set up, but method should not panic
	if output == "" {
		t.Error("Expected some output from context logger")
	}
}
