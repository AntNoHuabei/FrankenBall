package config

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

// Config 应用程序配置结构
type Config struct {
	App    AppConfig    `json:"app"`
	Log    LogConfig    `json:"log"`
	Window WindowConfig `json:"window"`
	mu     sync.RWMutex
}

// AppConfig 应用程序基本配置
type AppConfig struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Language string `json:"language"` // zh-CN, en-US
}

// WindowConfig 窗口配置
type WindowConfig struct {
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `json:"level"`       // debug, info, warn, error
	OutputFile string `json:"output_file"` // 日志文件路径
	MaxSize    int    `json:"max_size"`    // 最大文件大小(MB)
	MaxBackups int    `json:"max_backups"` // 最大备份数量
	MaxAge     int    `json:"max_age"`     // 最大保存天数
	Compress   bool   `json:"compress"`    // 是否压缩
}

var (
	globalConfig *Config
	once         sync.Once
	v            *viper.Viper
)

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:     "Remo",
			Version:  "1.0.0",
			Language: "zh-CN",
		},
		Window: WindowConfig{},
		Log: LogConfig{
			Level:      "info",
			OutputFile: "logs/app.log",
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		},
	}
}

// Init 初始化配置,如果配置文件不存在则创建默认配置
func Init(cfgPath string) (*Config, error) {
	var err error
	once.Do(func() {
		globalConfig, err = loadConfig(cfgPath)
	})
	return globalConfig, err
}

// Get 获取全局配置实例
func Get() *Config {
	if globalConfig == nil {
		panic("config not initialized, call Init() first")
	}
	return globalConfig
}

// loadConfig 从文件加载配置
func loadConfig(cfgPath string) (*Config, error) {
	// 初始化 viper
	v = viper.New()

	// 设置配置文件路径和名称
	dir := filepath.Dir(cfgPath)
	filename := filepath.Base(cfgPath)
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]

	v.SetConfigName(name)
	v.SetConfigType("json")
	v.AddConfigPath(dir)

	// 设置默认值
	setDefaults(v)

	// 尝试读取配置文件
	if err := v.ReadInConfig(); err != nil {
		// 配置文件不存在,创建默认配置
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cfg := DefaultConfig()
			if err := cfg.Save(cfgPath); err != nil {
				return nil, fmt.Errorf("failed to save default config: %w", err)
			}
			// 重新读取配置
			if err := v.ReadInConfig(); err != nil {
				return nil, fmt.Errorf("failed to read config after creation: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// 解析配置到结构体
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// Save 保存配置到文件
func (c *Config) Save(cfgPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if v == nil {
		// 如果 viper 未初始化,使用临时实例
		tmpViper := viper.New()
		dir := filepath.Dir(cfgPath)
		filename := filepath.Base(cfgPath)
		ext := filepath.Ext(filename)
		name := filename[:len(filename)-len(ext)]

		tmpViper.SetConfigName(name)
		tmpViper.SetConfigType("json")
		tmpViper.AddConfigPath(dir)

		// 将配置写入 viper
		syncToViper(tmpViper, c)

		// 写入配置文件
		if err := tmpViper.SafeWriteConfig(); err != nil {
			// 如果文件已存在,使用 WriteConfig 覆盖
			if err := tmpViper.WriteConfig(); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}
		}
		return nil
	}

	// 同步配置到 viper
	syncToViper(v, c)

	// 写入配置文件
	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Reload 重新加载配置
func (c *Config) Reload() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if v == nil {
		return fmt.Errorf("viper not initialized")
	}

	// 重新读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	// 解析到结构体
	if err := v.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// Update 更新配置并保存
func (c *Config) Update(updateFn func(*Config)) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if v == nil {
		return fmt.Errorf("viper not initialized")
	}

	// 执行更新
	updateFn(c)

	// 同步到 viper
	syncToViper(v, c)

	// 保存到文件
	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetApp 获取应用配置
func (c *Config) GetApp() AppConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.App
}

// GetWindow 获取窗口配置
func (c *Config) GetWindow() WindowConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Window
}

// GetLog 获取日志配置
func (c *Config) GetLog() LogConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Log
}

// SetAppLanguage 设置应用语言
func (c *Config) SetAppLanguage(language string) error {
	return c.Update(func(cfg *Config) {
		cfg.App.Language = language
	})
}

// SetWindowSize 设置窗口大小
func (c *Config) SetWindowSize(width, height int) error {
	return c.Update(func(cfg *Config) {
		//TODO
	})
}

// SetLogLevel 设置日志级别
func (c *Config) SetLogLevel(level string) error {
	return c.Update(func(cfg *Config) {
		cfg.Log.Level = level
	})
}

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
	defaultCfg := DefaultConfig()

	// App 默认值
	v.SetDefault("app.name", defaultCfg.App.Name)
	v.SetDefault("app.version", defaultCfg.App.Version)
	v.SetDefault("app.language", defaultCfg.App.Language)

	// Log 默认值
	v.SetDefault("log.level", defaultCfg.Log.Level)
	v.SetDefault("log.output_file", defaultCfg.Log.OutputFile)
	v.SetDefault("log.max_size", defaultCfg.Log.MaxSize)
	v.SetDefault("log.max_backups", defaultCfg.Log.MaxBackups)
	v.SetDefault("log.max_age", defaultCfg.Log.MaxAge)
	v.SetDefault("log.compress", defaultCfg.Log.Compress)
}

// syncToViper 将配置同步到 viper
func syncToViper(v *viper.Viper, cfg *Config) {
	// App 配置
	v.Set("app.name", cfg.App.Name)
	v.Set("app.version", cfg.App.Version)
	v.Set("app.language", cfg.App.Language)

	// Log 配置
	v.Set("log.level", cfg.Log.Level)
	v.Set("log.output_file", cfg.Log.OutputFile)
	v.Set("log.max_size", cfg.Log.MaxSize)
	v.Set("log.max_backups", cfg.Log.MaxBackups)
	v.Set("log.max_age", cfg.Log.MaxAge)
	v.Set("log.compress", cfg.Log.Compress)
}

// GetViper 获取 viper 实例 (用于高级用法)
func GetViper() *viper.Viper {
	return v
}
