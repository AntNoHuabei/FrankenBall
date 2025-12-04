package log

// Example 展示如何使用日志模块

/*
使用示例：

1. 基本初始化和使用
```go
import (
	"github.com/AntNoHuabei/Remo/internal/log"
)

func main() {
	// 使用默认配置初始化
	err := log.Init(nil)
	if err != nil {
		panic(err)
	}

	// 记录不同级别的日志
	log.Debug("This is a debug message")
	log.Info("Application started successfully")
	log.Warn("This is a warning")
	log.Error("An error occurred")
}
```

2. 自定义配置
```go
import (
	"github.com/AntNoHuabei/Remo/internal/log"
)

func main() {
	// 自定义日志配置
	cfg := &log.Config{
		Level:      "debug",           // 日志级别：debug, info, warn, error
		OutputFile: "logs/myapp.log",  // 日志文件路径
		MaxSize:    20,                // 单个文件最大 20MB
		MaxBackups: 10,                // 保留最多 10 个备份
		MaxAge:     60,                // 保留 60 天
		Compress:   true,              // 压缩旧日志
	}

	err := log.Init(cfg)
	if err != nil {
		panic(err)
	}

	log.Info("Logger initialized with custom config")
}
```

3. 带结构化字段的日志
```go
// 记录带有键值对的日志
log.Info("User logged in",
	"user_id", 12345,
	"username", "john_doe",
	"ip", "192.168.1.1",
)

log.Error("Failed to connect to database",
	"error", err,
	"host", "localhost",
	"port", 5432,
)

log.Debug("Processing request",
	"method", "POST",
	"path", "/api/users",
	"duration_ms", 150,
)
```

4. 使用子记录器
```go
// 创建带有固定字段的子记录器
userLogger := log.With("component", "user_service", "version", "1.0")
userLogger.Info("User services started")
userLogger.Info("Processing user request", "user_id", 123)

// 创建分组记录器
dbLogger := log.WithGroup("database")
dbLogger.Info("Connected to database", "host", "localhost")
```

5. 在 Wails 应用中集成
```go
import (
	"context"
	"github.com/AntNoHuabei/Remo/internal/config"
	"github.com/AntNoHuabei/Remo/internal/log"
)

type App struct {
	logger *slog.Logger
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	// 从配置初始化日志
	cfg := config.Get()
	logCfg := cfg.GetLog()

	err := log.Init(&log.Config{
		Level:      logCfg.Level,
		OutputFile: logCfg.OutputFile,
		MaxSize:    logCfg.MaxSize,
		MaxBackups: logCfg.MaxBackups,
		MaxAge:     logCfg.MaxAge,
		Compress:   logCfg.Compress,
	})

	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	// 创建应用专用的子记录器
	a.logger = log.With("component", "app")
	a.logger.Info("Application started", "version", cfg.GetApp().Version)
}

func (a *App) ProcessData(data string) error {
	a.logger.Debug("Processing data", "length", len(data))

	// 处理逻辑...

	a.logger.Info("Data processed successfully")
	return nil
}

func (a *App) HandleError(err error) {
	a.logger.Error("Error occurred",
		"error", err,
		"timestamp", time.Now(),
	)
}
```

6. 不同场景的日志记录
```go
// 应用启动
log.Info("Application starting",
	"version", "1.0.0",
	"environment", "production",
)

// HTTP 请求
log.Info("HTTP request",
	"method", "GET",
	"path", "/api/users",
	"status", 200,
	"duration_ms", 45,
)

// 数据库操作
log.Debug("Database query",
	"query", "SELECT * FROM users WHERE id = ?",
	"params", []int{123},
	"rows_affected", 1,
)

// 错误处理
log.Error("Failed to process request",
	"error", err,
	"retry_count", 3,
	"will_retry", false,
)

// 性能监控
log.Info("Performance metric",
	"operation", "data_sync",
	"duration_ms", 1234,
	"records_processed", 5000,
)

// 安全事件
log.Warn("Suspicious activity detected",
	"ip", "192.168.1.100",
	"attempts", 5,
	"action", "blocked",
)
```

7. 日志级别说明
- Debug: 详细的调试信息，通常只在开发环境使用
- Info:  一般信息，记录应用的正常运行状态
- Warn:  警告信息，表示潜在问题但不影响运行
- Error: 错误信息，表示发生了错误需要关注

8. 日志文件管理
日志文件会自动轮转：
- 当文件大小超过 MaxSize 时创建新文件
- 旧文件会根据时间戳重命名
- 保留最近 MaxBackups 个文件
- 删除超过 MaxAge 天的文件
- 如果 Compress=true，旧文件会被压缩为 .gz 格式

示例日志文件结构：
logs/
├── app.log           (当前日志文件)
├── app-2024-01-15.log
├── app-2024-01-14.log.gz
└── app-2024-01-13.log.gz

9. 获取原生 slog.Logger
```go
// 如果需要使用标准库的 slog.Logger
logger := log.Get()
logger.Info("Using native slog logger")
```

*/
