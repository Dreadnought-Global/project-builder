package main

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSanitizeProjectName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple Name", "Simple_Name"},
		{"  Trim spaces  ", "Trim_spaces"},
		{"Illegal /\\:*?\"<>| Chars", "Illegal_Chars"},
		{"Multiple   Spaces   And___Underscores", "Multiple_Spaces_And_Underscores"},
		{"....Leading.Trailing....", "Leading.Trailing"},
		{"Control\u0000Characters", "ControlCharacters"},
		{"&;Special()[]{}", "Special"},
		{"&;()[]{}", ""},
		{"Complex @!# Project $%^", "Complex_@!#_Project_$%^"},
	}

	for _, test := range tests {
		result := sanitizeProjectName(test.input)
		if result != test.expected {
			t.Errorf("sanitizeProjectName(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

func TestGetFolderList(t *testing.T) {
	// Design
	designNoClient, err := GetFolderList(Design, false)
	if err != nil {
		t.Fatalf("failed to get design folder list: %v", err)
	}
	expectedDesignCount := 9
	if len(designNoClient) != expectedDesignCount {
		t.Errorf("expected %d folders for Design, got %d: %v", expectedDesignCount, len(designNoClient), designNoClient)
	}

	designClient, err := GetFolderList(Design, true)
	if err != nil {
		t.Fatalf("failed to get design client folder list: %v", err)
	}
	// Expected: 9 (Design) + 4 (Client Docs) + 1 (Client_Handoff inside print ready) = 14 folders total
	// Wait, Client_Handoff is appended in designNoClient if client is true, which is included in the case.
	// Let's verify list length:
	expectedDesignClientCount := 14
	if len(designClient) != expectedDesignClientCount {
		t.Errorf("expected %d folders for Design with Client, got %d: %v", expectedDesignClientCount, len(designClient), designClient)
	}
}

func TestCreateFolderStructure(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "project-builder-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	targetProjPath := filepath.Join(tempDir, "test_project")

	err = CreateFolderStructure(targetProjPath, Design, true)
	if err != nil {
		t.Fatalf("failed to create folder structure: %v", err)
	}

	// Verify folders
	expectedFolders := []string{
		"00_Client_Docs",
		filepath.Join("00_Client_Docs", "Agreements"),
		filepath.Join("00_Client_Docs", "Invoices"),
		filepath.Join("00_Client_Docs", "Briefs"),
		"01_Client_Brief_&_Docs",
		"02_Raw_Assets",
		"03_Design_Files",
		filepath.Join("03_Design_Files", "Figma_Exports"),
		filepath.Join("03_Design_Files", "Adobe_Source"),
		filepath.Join("03_Design_Files", "Assets_Vector"),
		"04_Review_&_Deliverables",
		filepath.Join("04_Review_&_Deliverables", "Drafts_For_Review"),
		filepath.Join("04_Review_&_Deliverables", "Final_Print_Ready"),
		filepath.Join("04_Review_&_Deliverables", "Final_Print_Ready", "Client_Handoff"),
	}

	for _, folder := range expectedFolders {
		path := filepath.Join(targetProjPath, folder)
		if info, err := os.Stat(path); err != nil {
			t.Errorf("expected directory %s does not exist", folder)
		} else if !info.IsDir() {
			t.Errorf("expected path %s is not a directory", folder)
		}
	}
}

func TestConfigLoadSave(t *testing.T) {
	tempDir := t.TempDir()
	tempConfigPath := filepath.Join(tempDir, "config.yaml")

	// Set override
	configFilePathOverride = tempConfigPath
	defer func() { configFilePathOverride = "" }()

	// Test loading when file does not exist
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed when file does not exist: %v", err)
	}
	if cfg.WorkbenchPath != "" {
		t.Errorf("expected empty WorkbenchPath, got %q", cfg.WorkbenchPath)
	}

	// Test saving
	expectedPath := "/some/mock/workbench/path"
	cfg.WorkbenchPath = expectedPath
	err = SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Test loading saved config
	loadedCfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if loadedCfg.WorkbenchPath != expectedPath {
		t.Errorf("expected WorkbenchPath %q, got %q", expectedPath, loadedCfg.WorkbenchPath)
	}
}

func TestTUIFolderBrowser(t *testing.T) {
	tempDir := t.TempDir()
	subDir1 := filepath.Join(tempDir, "FolderA")
	subDir2 := filepath.Join(tempDir, "FolderB")
	_ = os.Mkdir(subDir1, 0755)
	_ = os.Mkdir(subDir2, 0755)

	m := folderBrowserModel{
		currentDir: tempDir,
	}
	err := m.updateFolders()
	if err != nil {
		t.Fatalf("failed to update folders: %v", err)
	}

	// Verify both subfolders are loaded
	if len(m.folders) < 2 {
		t.Fatalf("expected at least 2 folders, got %v", m.folders)
	}

	// Find index of FolderA and FolderB
	idxA, idxB := -1, -1
	for idx, f := range m.folders {
		if f == "FolderA" {
			idxA = idx
		} else if f == "FolderB" {
			idxB = idx
		}
	}
	if idxA == -1 || idxB == -1 {
		t.Fatalf("folders not found in TUI: %v", m.folders)
	}

	// Test moving cursor to FolderA
	m.cursor = idxA

	// Press Space/s to select
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	m = newModel.(folderBrowserModel)

	if m.selected != subDir1 {
		t.Errorf("expected selected path to be %q, got %q", subDir1, m.selected)
	}
}
