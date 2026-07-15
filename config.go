package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// configFilePathOverride allows tests to redirect the config file path.
var configFilePathOverride string

// Config represents the application settings structure.
type Config struct {
	WorkbenchPath string `yaml:"workbench_path"`
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

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse yaml: %w", err)
	}

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

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
