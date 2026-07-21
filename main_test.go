package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

func TestClientFlagSubfolderResolution(t *testing.T) {
	if got := ProjectCategorySubfolder(true); got != "00_Client_Projects" {
		t.Errorf("ProjectCategorySubfolder(true) = %q; expected %q", got, "00_Client_Projects")
	}
	if got := ProjectCategorySubfolder(false); got != "01_Passion_Projects" {
		t.Errorf("ProjectCategorySubfolder(false) = %q; expected %q", got, "01_Passion_Projects")
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
	expectedDesignClientCount := 14
	if len(designClient) != expectedDesignClientCount {
		t.Errorf("expected %d folders for Design with Client, got %d: %v", expectedDesignClientCount, len(designClient), designClient)
	}
}

func TestCreateFolderStructureClientVsPassion(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "project-builder-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	clientPath := filepath.Join(tempDir, "client_proj")
	err = CreateFolderStructure(clientPath, Design, true)
	if err != nil {
		t.Fatalf("failed to create client folder structure: %v", err)
	}

	passionPath := filepath.Join(tempDir, "passion_proj")
	err = CreateFolderStructure(passionPath, Design, false)
	if err != nil {
		t.Fatalf("failed to create passion folder structure: %v", err)
	}

	// Verify client folders exist
	if _, err := os.Stat(filepath.Join(clientPath, "00_Client_Docs")); os.IsNotExist(err) {
		t.Errorf("expected 00_Client_Docs to exist in client project")
	}

	// Verify passion folders do not have client folders
	if _, err := os.Stat(filepath.Join(passionPath, "00_Client_Docs")); err == nil {
		t.Errorf("expected 00_Client_Docs NOT to exist in passion project")
	}
}

func TestConfigLoadSaveAndMigration(t *testing.T) {
	tempDir := t.TempDir()
	tempConfigPath := filepath.Join(tempDir, "config.yaml")

	// Set override
	configFilePathOverride = tempConfigPath
	defer func() { configFilePathOverride = "" }()

	// Write old format config directly
	oldYaml := []byte(`workbench_path: /legacy/path`)
	if err := os.WriteFile(tempConfigPath, oldYaml, 0644); err != nil {
		t.Fatalf("failed to write legacy config: %v", err)
	}

	// Test migration loading
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed during migration test: %v", err)
	}
	if cfg.DefaultWorkbench != "/legacy/path" {
		t.Errorf("expected migrated DefaultWorkbench %q, got %q", "/legacy/path", cfg.DefaultWorkbench)
	}

	// Test saving new fields
	cfg.DefaultWorkbench = "/new/default/path"
	cfg.SetDisciplinePath(Design, "/design/path")
	cfg.SetDisciplinePath(VideoMotion, declinedDisciplinePath)
	err = SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Test loading saved config
	loadedCfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if loadedCfg.DefaultWorkbench != "/new/default/path" {
		t.Errorf("expected DefaultWorkbench %q, got %q", "/new/default/path", loadedCfg.DefaultWorkbench)
	}
	if loadedCfg.GetDisciplinePath(Design) != "/design/path" {
		t.Errorf("expected Design path %q, got %q", "/design/path", loadedCfg.GetDisciplinePath(Design))
	}

	// Test declined status
	if !loadedCfg.HasDeclinedDefault(VideoMotion) {
		t.Errorf("expected VideoMotion to be marked as declined")
	}
	if loadedCfg.GetDisciplinePath(VideoMotion) != "" {
		t.Errorf("expected VideoMotion path to be empty string due to 'declined', got %q", loadedCfg.GetDisciplinePath(VideoMotion))
	}

	loadedCfg.SetDisciplinePath(Audio, " /audio/path\n")
	if got := loadedCfg.GetDisciplinePath(Audio); got != "/audio/path" {
		t.Errorf("expected trimmed Audio path %q, got %q", "/audio/path", got)
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

	if len(m.folders) < 2 {
		t.Fatalf("expected at least 2 folders, got %v", m.folders)
	}

	idxA, idxB := -1, -1
	for idx, f := range m.folders {
		switch f {
		case "FolderA":
			idxA = idx
		case "FolderB":
			idxB = idx
		}
	}
	if idxA == -1 || idxB == -1 {
		t.Fatalf("folders not found in TUI: %v", m.folders)
	}

	m.cursor = idxA
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	m = newModel.(folderBrowserModel)

	if m.selected != subDir1 {
		t.Errorf("expected selected path to be %q, got %q", subDir1, m.selected)
	}
}

func TestResolveTargetPath(t *testing.T) {
	root := filepath.Join("tmp", "workbench", "design")

	clientPath := ResolveTargetPath(root, true, "My_Project")
	expectedClient := filepath.Join(root, "00_Client_Projects", "My_Project")
	if clientPath != expectedClient {
		t.Errorf("expected client path %q, got %q", expectedClient, clientPath)
	}

	passionPath := ResolveTargetPath(root, false, "My_Project")
	expectedPassion := filepath.Join(root, "01_Passion_Projects", "My_Project")
	if passionPath != expectedPassion {
		t.Errorf("expected passion path %q, got %q", expectedPassion, passionPath)
	}

	legacyRoot := filepath.Join(root, "00_Client_Projects")
	legacyPassionPath := ResolveTargetPath(legacyRoot, false, "My_Project")
	if legacyPassionPath != expectedPassion {
		t.Errorf("expected normalized legacy path %q, got %q", expectedPassion, legacyPassionPath)
	}
}

func TestNextAvailableProjectName(t *testing.T) {
	root := filepath.Join("tmp", "workbench", "design")
	existing := map[string]bool{
		ResolveTargetPath(root, false, "Demo_01"): true,
		ResolveTargetPath(root, false, "Demo_02"): true,
	}

	name, path := NextAvailableProjectName(root, false, "Demo", func(path string) bool {
		return existing[path]
	})

	if name != "Demo_03" {
		t.Errorf("expected Demo_03, got %q", name)
	}
	expectedPath := ResolveTargetPath(root, false, "Demo_03")
	if path != expectedPath {
		t.Errorf("expected path %q, got %q", expectedPath, path)
	}
}

func TestParseMenuChoice(t *testing.T) {
	if choice, ok := ParseMenuChoice(" 2\n", 1, 3); !ok || choice != 2 {
		t.Errorf("expected valid choice 2, got choice=%d ok=%t", choice, ok)
	}
	if _, ok := ParseMenuChoice("4", 1, 3); ok {
		t.Errorf("expected out-of-range choice to be invalid")
	}
	if _, ok := ParseMenuChoice("nope", 1, 3); ok {
		t.Errorf("expected non-number choice to be invalid")
	}
}

func TestParseYesNo(t *testing.T) {
	if value, ok := ParseYesNo(" yes\n"); !ok || !value {
		t.Errorf("expected yes to parse true, got value=%t ok=%t", value, ok)
	}
	if value, ok := ParseYesNo("N"); !ok || value {
		t.Errorf("expected N to parse false, got value=%t ok=%t", value, ok)
	}
	if _, ok := ParseYesNo("maybe"); ok {
		t.Errorf("expected maybe to be invalid")
	}
}

func TestInvalidDisciplineReturnsError(t *testing.T) {
	if _, err := GetFolderList(Discipline(99), false); err == nil {
		t.Errorf("expected invalid discipline to return error")
	}
}

func TestParseMenuChoiceWithDefault(t *testing.T) {
	choice, ok := ParseMenuChoiceWithDefault("\n", 1, 4, 1)
	if !ok || choice != 1 {
		t.Fatalf("expected blank input to select default 1, got choice=%d ok=%t", choice, ok)
	}

	choice, ok = ParseMenuChoiceWithDefault("3\n", 1, 4, 1)
	if !ok || choice != 3 {
		t.Fatalf("expected explicit input to select 3, got choice=%d ok=%t", choice, ok)
	}

	if _, ok := ParseMenuChoiceWithDefault("\n", 1, 4, 9); ok {
		t.Fatal("expected invalid default to fail")
	}
}

func TestHelpAndSettingsRenderTables(t *testing.T) {
	var help strings.Builder
	renderHelp(&help)
	helpText := stripANSI(help.String())
	for _, want := range []string{"| command", "install status", "theme set <name>", "project-builder"} {
		if !strings.Contains(helpText, want) {
			t.Fatalf("expected help output to contain %q, got:\n%s", want, helpText)
		}
	}

	var settings strings.Builder
	renderSettings(&settings, Config{DefaultWorkbench: "/workbench", Theme: "cyan"})
	settingsText := stripANSI(settings.String())
	for _, want := range []string{"Active theme", "cyan", "Global workbench", "/workbench", "project-builder install"} {
		if !strings.Contains(settingsText, want) {
			t.Fatalf("expected settings output to contain %q, got:\n%s", want, settingsText)
		}
	}
}

func TestFolderBrowserQuitConfirmation(t *testing.T) {
	m := folderBrowserModel{}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	m = updated.(folderBrowserModel)
	if !m.confirmQuit || m.quitted {
		t.Fatalf("first Ctrl+C should request confirmation, got confirmQuit=%t quitted=%t", m.confirmQuit, m.quitted)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m = updated.(folderBrowserModel)
	if m.confirmQuit || m.quitted {
		t.Fatalf("n should resume browser, got confirmQuit=%t quitted=%t", m.confirmQuit, m.quitted)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	m = updated.(folderBrowserModel)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	m = updated.(folderBrowserModel)
	if !m.quitted {
		t.Fatal("y should confirm folder browser cancellation")
	}
}

func TestShortenPath(t *testing.T) {
	if got := shortenPath("/very/long/path/to/project", 15); got != "/very/...roject" {
		t.Errorf("shortenPath() = %q", got)
	}
	if got := shortenPath("/very/long/path", 5); got != "..." {
		t.Errorf("shortenPath() for narrow width = %q", got)
	}
}

func TestConfigPathUsesXDGConfigHome(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG_CONFIG_HOME is not used on Windows")
	}
	configFilePathOverride = ""
	t.Cleanup(func() { configFilePathOverride = "" })
	t.Setenv("XDG_CONFIG_HOME", "/tmp/project-builder-xdg")

	got, err := GetConfigFilePath()
	if err != nil {
		t.Fatalf("GetConfigFilePath() error = %v", err)
	}
	want := filepath.Join("/tmp/project-builder-xdg", "project-builder", "config.yaml")
	if got != want {
		t.Errorf("GetConfigFilePath() = %q, want %q", got, want)
	}
}
