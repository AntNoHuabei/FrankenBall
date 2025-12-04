package log

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Level != "info" {
		t.Errorf("Expected level 'info', got '%s'", cfg.Level)
	}

	if cfg.MaxSize != 10 {
		t.Errorf("Expected max size 10, got %d", cfg.MaxSize)
	}

	if !cfg.Compress {
		t.Error("Expected compress to be true")
	}
}

func TestInitWithDefaultConfig(t *testing.T) {
	// 重置全局 logger
	logger = nil

	err := Init(nil)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	if logger == nil {
		t.Error("Logger should not be nil after initialization")
	}
}

func TestInitWithCustomConfig(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	cfg := &Config{
		Level:      "debug",
		OutputFile: logFile,
		MaxSize:    5,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
	}

	// 重置全局 logger
	logger = nil

	err := Init(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	if logger == nil {
		t.Error("Logger should not be nil after initialization")
	}

	// 写入日志
	Info("test message")

	// 检查日志文件是否创建
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"unknown", slog.LevelInfo}, // 默认值
	}

	for _, tt := range tests {
		result := parseLevel(tt.input)
		if result != tt.expected {
			t.Errorf("parseLevel(%s) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestLogFunctions(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	cfg := &Config{
		Level:      "debug",
		OutputFile: logFile,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
	}

	// 重置全局 logger
	logger = nil

	if err := Init(cfg); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// 测试各级别日志
	Debug("debug message", "key", "value")
	Info("info message", "key", "value")
	Warn("warn message", "key", "value")
	Error("error message", "key", "value")

	// 读取日志文件内容
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// 验证日志内容
	expectedMessages := []string{
		"debug message",
		"info message",
		"warn message",
		"error message",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(logContent, msg) {
			t.Errorf("Log file should contain '%s'", msg)
		}
	}
}

func TestWith(t *testing.T) {
	// 重置并初始化 logger
	logger = nil
	if err := Init(nil); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	childLogger := With("component", "test")
	if childLogger == nil {
		t.Error("With() should return a non-nil logger")
	}
}

func TestWithGroup(t *testing.T) {
	// 重置并初始化 logger
	logger = nil
	if err := Init(nil); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	groupLogger := WithGroup("mygroup")
	if groupLogger == nil {
		t.Error("WithGroup() should return a non-nil logger")
	}
}

func TestGetUninitialized(t *testing.T) {
	// 重置全局 logger
	logger = nil

	// Get 应该自动初始化 logger
	l := Get()
	if l == nil {
		t.Error("Get() should initialize logger if not already initialized")
	}
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	f()
	return buf.String()
}
