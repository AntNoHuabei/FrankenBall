package config

// Example 展示如何使用配置模块

/*
使用示例：

1. 初始化配置
```go
import (
	"github.com/AntNoHuabei/Remo/internal/config"
)

func main() {
	// 初始化配置，如果配置文件不存在会创建默认配置
	cfg, err := config.Init("config/app.json")
	if err != nil {
		log.Fatal("Failed to initialize config:", err)
	}

	// 使用配置
	appConfig := cfg.GetApp()
	fmt.Printf("App Name: %s, Version: %s\n", appConfig.Name, appConfig.Version)
}
```

2. 读取配置
```go
// 获取全局配置实例
cfg := config.Get()

// 读取应用配置
appCfg := cfg.GetApp()
fmt.Printf("Language: %s\n", appCfg.Language)

// 读取窗口配置
winCfg := cfg.GetWindow()
fmt.Printf("Window Size: %dx%d\n", winCfg.Width, winCfg.Height)

// 读取日志配置
logCfg := cfg.GetLog()
fmt.Printf("Log Level: %s\n", logCfg.Level)

// 读取鼠标配置
mouseCfg := cfg.GetMouse()
fmt.Printf("Hook Enabled: %v\n", mouseCfg.EnableHook)
```

3. 修改配置
```go
cfg := config.Get()

// 使用便捷方法修改
err := cfg.SetAppLanguage("en-US")
if err != nil {
	log.Printf("Failed to set language: %v", err)
}

err = cfg.SetWindowSize(1024, 768)
if err != nil {
	log.Printf("Failed to set window size: %v", err)
}

err = cfg.SetLogLevel("debug")
if err != nil {
	log.Printf("Failed to set log level: %v", err)
}

err = cfg.SetMouseHook(false)
if err != nil {
	log.Printf("Failed to set mouse hook: %v", err)
}
```

4. 自定义更新配置
```go
cfg := config.Get()

// 使用 Update 方法进行批量修改
err := cfg.Update(func(c *config.Config) {
	c.App.Language = "zh-CN"
	c.Window.AlwaysOnTop = false
	c.Log.Level = "warn"
	c.Mouse.CaptureClicks = false
})
if err != nil {
	log.Printf("Failed to update config: %v", err)
}
```

5. 重新加载配置
```go
cfg := config.Get()

// 从文件重新加载配置
err := cfg.Reload()
if err != nil {
	log.Printf("Failed to reload config: %v", err)
}
```

6. 在 Wails 应用中集成
```go
import (
	"github.com/AntNoHuabei/Remo/internal/config"
	"github.com/AntNoHuabei/Remo/internal/log"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	config *config.Config
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	// 初始化配置
	cfg, err := config.Init("config/app.json")
	if err != nil {
		log.Error("Failed to initialize config", "error", err)
		return
	}
	a.config = cfg

	// 初始化日志
	logCfg := cfg.GetLog()
	if err := log.Init(&log.Config{
		Level:      logCfg.Level,
		OutputFile: logCfg.OutputFile,
		MaxSize:    logCfg.MaxSize,
		MaxBackups: logCfg.MaxBackups,
		MaxAge:     logCfg.MaxAge,
		Compress:   logCfg.Compress,
	}); err != nil {
		log.Error("Failed to initialize logger", "error", err)
		return
	}

	log.Info("Application started", "version", cfg.GetApp().Version)
}

// 暴露给前端的方法
func (a *App) GetConfig() *config.Config {
	return config.Get()
}

func (a *App) UpdateLanguage(language string) error {
	return a.config.SetAppLanguage(language)
}

func (a *App) UpdateWindowSize(width, height int) error {
	return a.config.SetWindowSize(width, height)
}
```

7. 配置文件格式 (config/app.json)
```json
{
  "app": {
    "name": "Remo",
    "version": "1.0.0",
    "language": "zh-CN"
  },
  "window": {
    "width": 1920,
    "height": 1080,
    "always_on_top": true,
    "frameless": false,
    "dev_tools": true
  },
  "log": {
    "level": "info",
    "output_file": "logs/app.log",
    "max_size": 10,
    "max_backups": 5,
    "max_age": 30,
    "compress": true
  },
  "mouse": {
    "enable_hook": true,
    "hook_timeout": 5000,
    "ignore_events": false,
    "capture_clicks": true,
    "capture_move": true,
    "capture_scroll": true
  }
}
```

*/
