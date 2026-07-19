//go:build windows

package main

import (
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func ensureInstallPath(dir, display string, system bool) error {
	_ = display
	if pathContainsDir(os.Getenv("PATH"), dir) {
		return nil
	}
	root := registry.CURRENT_USER
	keyPath := `Environment`
	if system {
		root = registry.LOCAL_MACHINE
		keyPath = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
	}
	key, err := registry.OpenKey(root, keyPath, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()
	current, _, err := key.GetStringValue("Path")
	if err != nil && err != registry.ErrNotExist {
		return err
	}
	if windowsPathContains(current, dir) {
		return nil
	}
	sep := ""
	if strings.TrimSpace(current) != "" && !strings.HasSuffix(current, ";") {
		sep = ";"
	}
	return key.SetExpandStringValue("Path", current+sep+dir)
}

func windowsPathContains(value, dir string) bool {
	for _, part := range strings.Split(value, ";") {
		if samePath(strings.TrimSpace(part), dir) {
			return true
		}
	}
	return false
}
