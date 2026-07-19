package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestThemeValidation(t *testing.T) {
	if theme, err := GetTheme("violet"); err != nil || theme.Name != "violet" {
		t.Fatalf("expected violet theme, got theme=%+v err=%v", theme, err)
	}
	if _, err := GetTheme("nope"); err == nil || !strings.Contains(err.Error(), "valid themes") {
		t.Fatalf("expected invalid theme error with choices, got %v", err)
	}
}

func TestThemePersistence(t *testing.T) {
	tempDir := t.TempDir()
	configFilePathOverride = filepath.Join(tempDir, "config.yaml")
	defer func() { configFilePathOverride = "" }()

	cfg := Config{Theme: "cyan"}
	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}
	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if loaded.Theme != "cyan" {
		t.Fatalf("expected cyan theme, got %q", loaded.Theme)
	}
}

func TestThemeCommands(t *testing.T) {
	tempDir := t.TempDir()
	configFilePathOverride = filepath.Join(tempDir, "config.yaml")
	defer func() { configFilePathOverride = "" }()

	cfg := Config{}
	var out bytes.Buffer
	_, cfg, code := HandleCommand([]string{"theme", "set", "emerald"}, cfg, strings.NewReader(""), &out)
	if code != 0 || cfg.Theme != "emerald" {
		t.Fatalf("expected emerald set, code=%d cfg=%+v out=%s", code, cfg, out.String())
	}
	out.Reset()
	_, cfg, code = HandleCommand([]string{"theme", "reset"}, cfg, strings.NewReader(""), &out)
	if code != 0 || cfg.Theme != defaultThemeName {
		t.Fatalf("expected reset theme, code=%d cfg=%+v out=%s", code, cfg, out.String())
	}
}
