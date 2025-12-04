package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.App.Name != "Remo" {
		t.Errorf("Expected app name 'Remo', got '%s'", cfg.App.Name)
	}

	if cfg.Log.Level != "info" {
		t.Errorf("Expected log level 'info', got '%s'", cfg.Log.Level)
	}

}

func TestConfigSaveAndLoad(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")

	// 创建配置
	cfg := DefaultConfig()
	cfg.App.Language = "en-US"

	// 保存配置
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// 加载配置
	loadedCfg, err := loadConfig(cfgPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 验证加载的配置
	if loadedCfg.App.Language != "en-US" {
		t.Errorf("Expected language 'en-US', got '%s'", loadedCfg.App.Language)
	}

}

func TestConfigUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")

	cfg := DefaultConfig()
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// 更新配置
	err := cfg.Update(func(c *Config) {
		c.App.Language = "ja-JP"
		c.Log.Level = "debug"
	})

	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// 验证更新
	if cfg.App.Language != "ja-JP" {
		t.Errorf("Expected language 'ja-JP', got '%s'", cfg.App.Language)
	}

	if cfg.Log.Level != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", cfg.Log.Level)
	}
}

func TestSetMethods(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")

	cfg := DefaultConfig()
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// 测试设置方法
	if err := cfg.SetAppLanguage("fr-FR"); err != nil {
		t.Fatalf("SetAppLanguage failed: %v", err)
	}

	if cfg.App.Language != "fr-FR" {
		t.Errorf("Expected language 'fr-FR', got '%s'", cfg.App.Language)
	}

	if err := cfg.SetWindowSize(1024, 768); err != nil {
		t.Fatalf("SetWindowSize failed: %v", err)
	}

	if err := cfg.SetLogLevel("error"); err != nil {
		t.Fatalf("SetLogLevel failed: %v", err)
	}

	if cfg.Log.Level != "error" {
		t.Errorf("Expected log level 'error', got '%s'", cfg.Log.Level)
	}

}
