package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type folderBrowserModel struct {
	currentDir string
	folders    []string
	cursor     int
	selected   string
	quitted    bool
	err        error
}

func initialModel() (folderBrowserModel, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return folderBrowserModel{}, err
	}

	m := folderBrowserModel{
		currentDir: home,
	}
	if err := m.updateFolders(); err != nil {
		m.err = err
	}
	return m, nil
}

func (m *folderBrowserModel) updateFolders() error {
	m.err = nil
	entries, err := os.ReadDir(m.currentDir)
	if err != nil {
		m.folders = []string{".."}
		return err
	}

	var folders []string
	parent := filepath.Dir(m.currentDir)
	if parent != m.currentDir {
		folders = append(folders, "..")
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if entry.IsDir() {
			folders = append(folders, entry.Name())
		}
	}

	m.folders = folders
	if m.cursor >= len(m.folders) {
		m.cursor = 0
	}
	return nil
}

func (m folderBrowserModel) Init() tea.Cmd {
	return nil
}

func (m folderBrowserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitted = true
			return m, tea.Quit

		case "up", "k":
			if len(m.folders) == 0 {
				break
			}
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.folders) - 1
			}

		case "down", "j":
			if len(m.folders) == 0 {
				break
			}
			if m.cursor < len(m.folders)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}

		case "enter":
			if len(m.folders) == 0 {
				break
			}
			sel := m.folders[m.cursor]
			if sel == ".." {
				m.currentDir = filepath.Dir(m.currentDir)
				m.cursor = 0
				_ = m.updateFolders()
			} else {
				newDir := filepath.Join(m.currentDir, sel)
				testEntries, err := os.ReadDir(newDir)
				if err != nil {
					m.err = fmt.Errorf("cannot open folder: %w", err)
				} else {
					_ = testEntries
					m.currentDir = newDir
					m.cursor = 0
					_ = m.updateFolders()
				}
			}

		case "backspace":
			parent := filepath.Dir(m.currentDir)
			if parent != m.currentDir {
				m.currentDir = parent
				m.cursor = 0
				_ = m.updateFolders()
			}

		case " ", "s":
			if len(m.folders) == 0 {
				m.selected = m.currentDir
				return m, tea.Quit
			}
			sel := m.folders[m.cursor]
			if sel == ".." {
				m.selected = m.currentDir
			} else {
				m.selected = filepath.Join(m.currentDir, sel)
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m folderBrowserModel) View() string {
	if m.quitted {
		return "Setup cancelled. Exiting.\n"
	}
	if m.selected != "" {
		return fmt.Sprintf("Selected: %s\n", m.selected)
	}

	var s strings.Builder
	s.WriteString("==================================================\n")
	s.WriteString("  Setup Workbench: Choose Your Root Project Folder \n")
	s.WriteString("==================================================\n\n")
	s.WriteString(fmt.Sprintf("Current Directory: [ %s ]\n", m.currentDir))

	if m.err != nil {
		s.WriteString(fmt.Sprintf("\n⚠️ Error: %v\n\n", m.err))
	} else {
		s.WriteString("\n")
	}

	s.WriteString("Folders:\n")
	if len(m.folders) == 0 {
		s.WriteString("  (No subfolders found)\n")
	} else {
		for i, folder := range m.folders {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s.WriteString(fmt.Sprintf(" %s %s\n", cursor, folder))
		}
	}

	s.WriteString("\n--------------------------------------------------\n")
	s.WriteString("Controls:\n")
	s.WriteString("  [↑/↓] or [j/k]  Move cursor\n")
	s.WriteString("  [Enter]         Open highlighted folder\n")
	s.WriteString("  [Backspace]     Go back to parent folder\n")
	s.WriteString("  [Space] or [s]  Confirm and select highlighted folder\n")
	s.WriteString("  [q]             Cancel & Quit\n")
	s.WriteString("--------------------------------------------------\n")

	target := m.currentDir
	if len(m.folders) > 0 && m.folders[m.cursor] != ".." {
		target = filepath.Join(m.currentDir, m.folders[m.cursor])
	}
	s.WriteString(fmt.Sprintf("\nTarget selection: %s\n", target))

	return s.String()
}

// RunFolderBrowser launches the Bubble Tea folder browser.
// Returns the absolute path chosen, or empty string if aborted.
func RunFolderBrowser() (string, error) {
	m, err := initialModel()
	if err != nil {
		return "", err
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	m = finalModel.(folderBrowserModel)
	if m.quitted {
		return "", nil
	}
	return m.selected, nil
}
