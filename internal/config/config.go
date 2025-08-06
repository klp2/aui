package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	APIKeys  map[string]string `yaml:"api_keys"`
	Database DatabaseConfig    `yaml:"database"`
	UI       UIConfig          `yaml:"ui"`
	Logging  LoggingConfig     `yaml:"logging"`
}

// DatabaseConfig contains database-related settings
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// UIConfig contains UI-related settings
type UIConfig struct {
	Theme       string `yaml:"theme"`
	RefreshRate int    `yaml:"refresh_rate"`
}

// LoggingConfig contains logging-related settings
type LoggingConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file,omitempty"`
}

// NewDefault creates a new configuration with default values
func NewDefault() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		APIKeys: make(map[string]string),
		Database: DatabaseConfig{
			Path: filepath.Join(home, ".config", "aui", "aui.db"),
		},
		UI: UIConfig{
			Theme:       "default",
			RefreshRate: 100,
		},
		Logging: LoggingConfig{
			Level: "info",
			File:  filepath.Join(home, ".config", "aui", "aui.log"),
		},
	}
}

// LoadFromFile loads configuration from a YAML file
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := NewDefault()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Expand paths
	cfg.Database.Path = expandPath(cfg.Database.Path)
	if cfg.Logging.File != "" {
		cfg.Logging.File = expandPath(cfg.Logging.File)
	}

	return cfg, nil
}

// Load loads configuration from the default location or creates default
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(home, ".config", "aui", "config.yaml")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		cfg := NewDefault()
		cfg.LoadFromEnv()
		return cfg, nil
	}

	cfg, err := LoadFromFile(configPath)
	if err != nil {
		return nil, err
	}

	// Apply environment overrides
	cfg.LoadFromEnv()

	return cfg, nil
}

// LoadFromEnv loads configuration overrides from environment variables
func (c *Config) LoadFromEnv() {
	// API Keys
	if key := os.Getenv("AUI_API_KEY_ANTHROPIC"); key != "" {
		if c.APIKeys == nil {
			c.APIKeys = make(map[string]string)
		}
		c.APIKeys["anthropic"] = key
	}

	if key := os.Getenv("AUI_API_KEY_OPENAI"); key != "" {
		if c.APIKeys == nil {
			c.APIKeys = make(map[string]string)
		}
		c.APIKeys["openai"] = key
	}

	if key := os.Getenv("AUI_API_KEY_GOOGLE"); key != "" {
		if c.APIKeys == nil {
			c.APIKeys = make(map[string]string)
		}
		c.APIKeys["google"] = key
	}

	// Database
	if path := os.Getenv("AUI_DATABASE_PATH"); path != "" {
		c.Database.Path = expandPath(path)
	}

	// UI
	if theme := os.Getenv("AUI_UI_THEME"); theme != "" {
		c.UI.Theme = theme
	}

	if rate := os.Getenv("AUI_UI_REFRESH_RATE"); rate != "" {
		// Parse int, ignore errors and keep default if invalid
		var refreshRate int
		fmt.Sscanf(rate, "%d", &refreshRate)
		if refreshRate > 0 {
			c.UI.RefreshRate = refreshRate
		}
	}

	// Logging
	if level := os.Getenv("AUI_LOGGING_LEVEL"); level != "" {
		c.Logging.Level = level
	}

	if file := os.Getenv("AUI_LOGGING_FILE"); file != "" {
		c.Logging.File = expandPath(file)
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Database path is required
	if c.Database.Path == "" {
		return fmt.Errorf("database path is required")
	}

	// Validate log level
	validLevels := []string{"debug", "info", "warn", "error"}
	levelValid := false
	for _, level := range validLevels {
		if c.Logging.Level == level {
			levelValid = true
			break
		}
	}
	if !levelValid {
		return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", c.Logging.Level)
	}

	// Validate refresh rate
	if c.UI.RefreshRate <= 0 {
		return fmt.Errorf("refresh rate must be positive")
	}

	return nil
}

// SaveToFile saves the configuration to a YAML file
func (c *Config) SaveToFile(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Save saves the configuration to the default location
func (c *Config) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(home, ".config", "aui", "config.yaml")
	return c.SaveToFile(configPath)
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
