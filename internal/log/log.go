package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 全局日志记录器
var (
	logger *slog.Logger
)

// Config 日志配置
type Config struct {
	Level      string // debug, info, warn, error
	OutputFile string // 日志文件路径,为空则只输出到控制台
	MaxSize    int    // 单个日志文件最大大小(MB)
	MaxBackups int    // 最多保留的旧日志文件数量
	MaxAge     int    // 最多保留天数
	Compress   bool   // 是否压缩旧日志
}

// DefaultConfig 返回默认日志配置
func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		OutputFile: "logs/app.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
}

// Init 初始化日志系统
func Init(cfg *Config) error {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// 解析日志级别
	level := parseLevel(cfg.Level)

	// 创建日志输出
	var writer io.Writer
	if cfg.OutputFile != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(cfg.OutputFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		// 配置日志轮转
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.OutputFile,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}

		// 同时输出到文件和控制台
		writer = io.MultiWriter(os.Stdout, fileWriter)
	} else {
		// 只输出到控制台
		writer = os.Stdout
	}

	// 创建带颜色的文本处理器
	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 自定义时间格式
			if a.Key == slog.TimeKey {
				return slog.String("time", time.Now().Format("2006-01-02 15:04:05"))
			}
			return a
		},
	})

	logger = slog.New(handler)
	slog.SetDefault(logger)

	return nil
}

// parseLevel 解析日志级别
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Get 获取日志记录器
func Get() *slog.Logger {
	if logger == nil {
		// 如果未初始化,使用默认配置初始化
		if err := Init(nil); err != nil {
			panic(fmt.Sprintf("failed to initialize logger: %v", err))
		}
	}
	return logger
}

// Debug 记录调试信息
func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

// Info 记录普通信息
func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

// Warn 记录警告信息
func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

// Error 记录错误信息
func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

// With 创建带有固定字段的子记录器
func With(args ...any) *slog.Logger {
	return Get().With(args...)
}

// WithGroup 创建带有分组的子记录器
func WithGroup(name string) *slog.Logger {
	return Get().WithGroup(name)
}
