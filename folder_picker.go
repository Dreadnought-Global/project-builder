package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var linuxPickerCommand string
var linuxPickerChecked bool

// RunNativeFolderPicker opens an OS-native folder selection dialog.
// Returns the selected path, an empty string if cancelled, and an error if fallback is required.
func RunNativeFolderPicker() (string, error) {
	switch runtime.GOOS {
	case "windows":
		// PowerShell script to invoke FolderBrowserDialog
		script := `
[System.Reflection.Assembly]::LoadWithPartialName("System.windows.forms") | Out-Null
$browser = New-Object System.Windows.Forms.FolderBrowserDialog
$browser.Description = "Select a folder"
$browser.ShowNewFolderButton = $true
$result = $browser.ShowDialog((New-Object System.Windows.Forms.Form -Property @{TopMost = $true}))
if ($result -eq [System.Windows.Forms.DialogResult]::OK) {
    Write-Output $browser.SelectedPath
}
`
		cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
		out, err := cmd.Output()
		if err != nil {
			return "", nil // Treat as cancelled or failed silently
		}
		path := strings.TrimSpace(string(out))
		return path, nil

	case "darwin":
		cmd := exec.Command("osascript", "-e", `POSIX path of (choose folder)`)
		out, err := cmd.Output()
		if err != nil {
			// user cancelled or error
			return "", nil
		}
		return strings.TrimSpace(string(out)), nil

	case "linux":
		if !linuxPickerChecked {
			linuxPickerChecked = true
			if _, err := exec.LookPath("zenity"); err == nil {
				linuxPickerCommand = "zenity"
			} else if _, err := exec.LookPath("kdialog"); err == nil {
				linuxPickerCommand = "kdialog"
			}
		}

		switch linuxPickerCommand {
		case "zenity":
			cmd := exec.Command("zenity", "--file-selection", "--directory", "--title=Select a folder")
			out, err := cmd.Output()
			if err != nil {
				return "", nil
			}
			return strings.TrimSpace(string(out)), nil
		case "kdialog":
			cmd := exec.Command("kdialog", "--getexistingdirectory", ".", "--title", "Select a folder")
			out, err := cmd.Output()
			if err != nil {
				return "", nil
			}
			return strings.TrimSpace(string(out)), nil
		}

		return "", fmt.Errorf("no gui folder picker found")

	default:
		return "", fmt.Errorf("unsupported platform for native picker")
	}
}
