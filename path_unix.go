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
	if system || pathContainsDir(os.Getenv("PATH"), dir) {
		return nil
	}
	rc, err := shellStartupFile()
	if err != nil {
		return err
	}
	return addPathBlock(rc, display)
}

func shellStartupFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	shell := filepath.Base(os.Getenv("SHELL"))
	switch shell {
	case "zsh":
		return filepath.Join(home, ".zshrc"), nil
	case "bash":
		return filepath.Join(home, ".bashrc"), nil
	default:
		return filepath.Join(home, ".profile"), nil
	}
}

func addPathBlock(path, display string) error {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	content := string(data)
	if strings.Contains(content, pathBlockMarker) || strings.Contains(content, display) {
		return nil
	}
	block := fmt.Sprintf(`
%s
case ":$PATH:" in
  *":%s:"*) ;;
  *) export PATH="%s:$PATH" ;;
esac
`, pathBlockMarker, display, display)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(block)
	return err
}
