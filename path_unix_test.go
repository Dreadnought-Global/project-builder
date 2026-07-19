//go:build !windows

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddPathBlockIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".profile")
	if err := addPathBlock(path, "$HOME/.local/bin"); err != nil {
		t.Fatalf("addPathBlock failed: %v", err)
	}
	if err := addPathBlock(path, "$HOME/.local/bin"); err != nil {
		t.Fatalf("second addPathBlock failed: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read profile: %v", err)
	}
	if strings.Count(string(data), pathBlockMarker) != 1 {
		t.Fatalf("expected one path block, got %q", data)
	}
}
