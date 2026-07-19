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
	code := handleInstallCommand([]string{"--dry-run"}, &out)
	if code != 0 {
		t.Fatalf("expected dry-run success, code=%d out=%s", code, out.String())
	}
	if !strings.Contains(out.String(), "Dry run only") {
		t.Fatalf("expected dry-run output, got %q", out.String())
	}
}
