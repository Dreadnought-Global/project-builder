package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseInstallOptions(t *testing.T) {
	opts, err := parseInstallOptions([]string{"--dry-run", "--force", "status"})
	if err != nil {
		t.Fatalf("parseInstallOptions failed: %v", err)
	}
	if !opts.DryRun || !opts.Force || !opts.Status {
		t.Fatalf("expected all options enabled, got %+v", opts)
	}
	if _, err := parseInstallOptions([]string{"--wat"}); err == nil {
		t.Fatalf("expected unknown option error")
	}
}

func TestInstallBinaryName(t *testing.T) {
	if installBinaryName() == "" {
		t.Fatalf("expected install binary name")
	}
}

func TestCopyExecutableRefusesOverwriteWithoutForce(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source")
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(source, []byte("new"), 0o755); err != nil {
		t.Fatalf("write source: %v", err)
	}
	if err := os.WriteFile(target, []byte("old"), 0o755); err != nil {
		t.Fatalf("write target: %v", err)
	}
	if err := copyExecutable(source, target, false); err == nil || !strings.Contains(err.Error(), "--force") {
		t.Fatalf("expected force error, got %v", err)
	}
}

func TestCopyExecutableForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source")
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(source, []byte("new"), 0o755); err != nil {
		t.Fatalf("write source: %v", err)
	}
	if err := os.WriteFile(target, []byte("old"), 0o755); err != nil {
		t.Fatalf("write target: %v", err)
	}
	if err := copyExecutable(source, target, true); err != nil {
		t.Fatalf("copyExecutable failed: %v", err)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read target: %v", err)
	}
	if string(data) != "new" {
		t.Fatalf("expected overwritten target, got %q", data)
	}
}

func TestInstallCommandDryRun(t *testing.T) {
	var out bytes.Buffer
	code := handleInstallCommand([]string{"--dry-run"}, strings.NewReader(""), &out)
	if code != 0 {
		t.Fatalf("expected dry-run success, code=%d out=%s", code, out.String())
	}
	if !strings.Contains(out.String(), "Dry run only") {
		t.Fatalf("expected dry-run output, got %q", out.String())
	}
	if strings.Contains(out.String(), "Folder browser controls") {
		t.Fatalf("dry-run must not show onboarding, got %q", out.String())
	}
}

func TestAcknowledgeFolderBrowserControls(t *testing.T) {
	configFilePathOverride = filepath.Join(t.TempDir(), "config.yaml")
	t.Cleanup(func() { configFilePathOverride = "" })

	var out bytes.Buffer
	if err := acknowledgeFolderBrowserControls(strings.NewReader("no\ny\n"), &out); err != nil {
		t.Fatalf("acknowledgeFolderBrowserControls() error = %v", err)
	}
	if !strings.Contains(out.String(), "Enter y after reviewing the controls.") {
		t.Fatalf("expected reprompt, got %q", out.String())
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if !cfg.FolderBrowserControlsAcknowledged {
		t.Fatal("expected acknowledgement to persist")
	}
}

func TestVersionCommand(t *testing.T) {
	var out bytes.Buffer
	handled, _, code := HandleCommand([]string{"--version"}, Config{}, strings.NewReader(""), &out)
	if !handled || code != 0 {
		t.Fatalf("expected version command success, handled=%t code=%d", handled, code)
	}
	if !strings.Contains(out.String(), "Project Builder v2.4.0 (2026-07-21)") {
		t.Fatalf("unexpected version output: %q", out.String())
	}
}
