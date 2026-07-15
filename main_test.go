package main

import (
	"os"
	"path/filepath"
	"testing"
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
