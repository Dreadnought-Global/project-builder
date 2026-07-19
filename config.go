package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

// configFilePathOverride allows tests to redirect the config file path.
var configFilePathOverride string

// Config represents the application settings structure.
type Config struct {
	DefaultWorkbench string            `yaml:"default_workbench"`
	DisciplinePaths  map[string]string `yaml:"discipline_paths"`
	Theme            string            `yaml:"theme"`
}

func (c *Config) GetDisciplinePath(d Discipline) string {
	if c.DisciplinePaths == nil {
		return ""
	}
	val := strings.TrimSpace(c.DisciplinePaths[d.DisciplineKey()])
	if val == declinedDisciplinePath {
		return ""
	}
	return val
}

func (c *Config) HasDeclinedDefault(d Discipline) bool {
	if c.DisciplinePaths == nil {
		return false
	}
	return strings.TrimSpace(c.DisciplinePaths[d.DisciplineKey()]) == declinedDisciplinePath
}

func (c *Config) SetDisciplinePath(d Discipline, path string) {
	if c.DisciplinePaths == nil {
		c.DisciplinePaths = make(map[string]string)
	}
	c.DisciplinePaths[d.DisciplineKey()] = strings.TrimSpace(path)
}

// GetConfigFilePath returns the resolved OS-specific path for config.yaml.
// Linux/macOS: ~/.config/project-builder/config.yaml
// Windows: %APPDATA%/project-builder/config.yaml
func GetConfigFilePath() (string, error) {
	if configFilePathOverride != "" {
		return configFilePathOverride, nil
	}

	var baseDir string
	if runtime.GOOS == "windows" {
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			userConfig, err := os.UserConfigDir()
			if err != nil {
				return "", err
			}
			baseDir = userConfig
		}
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(home, ".config")
	}

	return filepath.Join(baseDir, "project-builder", "config.yaml"), nil
}

// LoadConfig reads the YAML configuration file.
// If the file does not exist, it returns an empty config and no error.
func LoadConfig() (Config, error) {
	var cfg Config
	path, err := GetConfigFilePath()
	if err != nil {
		return cfg, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	// Try migration first
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err == nil {
		if wp, ok := raw["workbench_path"].(string); ok && strings.TrimSpace(wp) != "" {
			cfg.DefaultWorkbench = strings.TrimSpace(wp)
			cfg.DisciplinePaths = make(map[string]string)
			_ = SaveConfig(cfg) // Auto-save migrated config
			return cfg, nil
		}
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse yaml: %w", err)
	}
	cfg.DefaultWorkbench = strings.TrimSpace(cfg.DefaultWorkbench)
	if cfg.DisciplinePaths == nil {
		cfg.DisciplinePaths = make(map[string]string)
	}
	for key, value := range cfg.DisciplinePaths {
		cfg.DisciplinePaths[key] = strings.TrimSpace(value)
	}
	cfg.Theme = normalizeThemeName(cfg.Theme)

	return cfg, nil
}

// SaveConfig writes the YAML configuration file, creating parent folders if needed.
func SaveConfig(cfg Config) error {
	path, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	cfg.Theme = normalizeThemeName(cfg.Theme)
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
