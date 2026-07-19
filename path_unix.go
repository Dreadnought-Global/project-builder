//go:build !windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const pathBlockMarker = "# Project Builder PATH"

func ensureInstallPath(dir, display string, system bool) error {
	if system {
		return nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	paths := []struct {
		path  string
		block string
	}{
		{filepath.Join(home, ".profile"), posixPathBlock(display)},
		{filepath.Join(home, ".bashrc"), posixPathBlock(display)},
		{filepath.Join(home, ".zshrc"), posixPathBlock(display)},
		{filepath.Join(home, ".config", "fish", "config.fish"), fishPathBlock(display)},
	}
	for _, item := range paths {
		if err := addPathBlock(item.path, display, item.block); err != nil {
			return err
		}
	}
	return nil
}

func posixPathBlock(display string) string {
	return fmt.Sprintf(`
%s
case ":$PATH:" in
  *":%s:"*) ;;
  *) export PATH="%s:$PATH" ;;
esac
`, pathBlockMarker, display, display)
}

func fishPathBlock(display string) string {
	return fmt.Sprintf(`
%s
if not contains "%s" $PATH
    set -gx PATH "%s" $PATH
end
`, pathBlockMarker, display, display)
}

func addPathBlock(path, display, block string) error {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	content := string(data)
	if strings.Contains(content, pathBlockMarker) || strings.Contains(content, display) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(block)
	return err
}
