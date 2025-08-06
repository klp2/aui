package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultConfig(t *testing.T) {
	cfg := NewDefault()

	if cfg == nil {
		t.Fatal("Expected default config, got nil")
	}

	if cfg.Database.Path == "" {
		t.Error("Expected default database path")
	}

	if cfg.UI.Theme != "default" {
		t.Errorf("Expected default theme, got %s", cfg.UI.Theme)
	}

	if cfg.UI.RefreshRate != 100 {
		t.Errorf("Expected refresh rate 100, got %d", cfg.UI.RefreshRate)
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("Expected log level info, got %s", cfg.Logging.Level)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `
api_keys:
  anthropic: "test-anthropic-key"
  openai: "test-openai-key"
database:
  path: "/custom/path/aui.db"
ui:
  theme: "dark"
  refresh_rate: 200
logging:
  level: "debug"
  file: "/custom/log/aui.log"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := LoadFromFile(configFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.APIKeys["anthropic"] != "test-anthropic-key" {
		t.Errorf("Expected anthropic key, got %s", cfg.APIKeys["anthropic"])
	}

	if cfg.Database.Path != "/custom/path/aui.db" {
		t.Errorf("Expected custom db path, got %s", cfg.Database.Path)
	}

	if cfg.UI.Theme != "dark" {
		t.Errorf("Expected dark theme, got %s", cfg.UI.Theme)
	}
}

func TestLoadConfigWithEnvironmentOverrides(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("AUI_API_KEY_ANTHROPIC", "env-anthropic-key")
	os.Setenv("AUI_API_KEY_OPENAI", "env-openai-key")
	os.Setenv("AUI_DATABASE_PATH", "/env/path/aui.db")
	os.Setenv("AUI_UI_THEME", "terminal")
	defer func() {
		os.Unsetenv("AUI_API_KEY_ANTHROPIC")
		os.Unsetenv("AUI_API_KEY_OPENAI")
		os.Unsetenv("AUI_DATABASE_PATH")
		os.Unsetenv("AUI_UI_THEME")
	}()

	cfg.LoadFromEnv()

	if cfg.APIKeys["anthropic"] != "env-anthropic-key" {
		t.Errorf("Expected env anthropic key, got %s", cfg.APIKeys["anthropic"])
	}

	if cfg.Database.Path != "/env/path/aui.db" {
		t.Errorf("Expected env db path, got %s", cfg.Database.Path)
	}

	if cfg.UI.Theme != "terminal" {
		t.Errorf("Expected terminal theme, got %s", cfg.UI.Theme)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
	}{
		{
			name: "valid config with API keys",
			config: &Config{
				APIKeys: map[string]string{
					"anthropic": "key1",
				},
				Database: DatabaseConfig{Path: "/path/to/db"},
				UI:       UIConfig{Theme: "default", RefreshRate: 100},
				Logging:  LoggingConfig{Level: "info"},
			},
			wantError: false,
		},
		{
			name: "invalid config - no database path",
			config: &Config{
				APIKeys: map[string]string{"anthropic": "key1"},
				UI:      UIConfig{Theme: "default", RefreshRate: 100},
				Logging: LoggingConfig{Level: "info"},
			},
			wantError: true,
		},
		{
			name: "invalid config - invalid log level",
			config: &Config{
				APIKeys:  map[string]string{"anthropic": "key1"},
				Database: DatabaseConfig{Path: "/path/to/db"},
				UI:       UIConfig{Theme: "default", RefreshRate: 100},
				Logging:  LoggingConfig{Level: "invalid"},
			},
			wantError: true,
		},
		{
			name: "valid config without API keys (can be added later)",
			config: &Config{
				APIKeys:  map[string]string{},
				Database: DatabaseConfig{Path: "/path/to/db"},
				UI:       UIConfig{Theme: "default", RefreshRate: 100},
				Logging:  LoggingConfig{Level: "info"},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestConfigSave(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	cfg := &Config{
		APIKeys: map[string]string{
			"anthropic": "save-test-key",
		},
		Database: DatabaseConfig{Path: "/save/test/db"},
		UI:       UIConfig{Theme: "saved", RefreshRate: 150},
		Logging:  LoggingConfig{Level: "warn", File: "/save/test/log"},
	}

	err := cfg.SaveToFile(configFile)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	loaded, err := LoadFromFile(configFile)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loaded.APIKeys["anthropic"] != cfg.APIKeys["anthropic"] {
		t.Error("Saved config doesn't match original")
	}

	if loaded.Database.Path != cfg.Database.Path {
		t.Error("Saved database path doesn't match")
	}
}

func TestExpandHomePath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		input    string
		expected string
	}{
		{"~/test/path", filepath.Join(home, "test/path")},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	for _, tt := range tests {
		result := expandPath(tt.input)
		if result != tt.expected {
			t.Errorf("expandPath(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}
