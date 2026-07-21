package main

import (
	"fmt"
	"io"
	"strings"
)

type menuRow struct {
	Command     string
	Description string
}

func clearTerminal(out io.Writer) {
	fmt.Fprint(out, "\033[H\033[2J")
}

func renderHomeDashboard(out io.Writer) {
	fmt.Fprintln(out, accentText("Home"))
	fmt.Fprintln(out, promptText("[1]")+" "+primaryText("Create project folder")+mutedText(" (default)"))
	fmt.Fprintln(out, promptText("[2]")+" "+primaryText("Help"))
	fmt.Fprintln(out, promptText("[3]")+" "+primaryText("Settings"))
	fmt.Fprintln(out, promptText("[4]")+" "+primaryText("Exit"))
	fmt.Fprint(out, promptText("Selection")+mutedText(": "))
}

func renderHelp(out io.Writer) {
	fmt.Fprintln(out, accentText("Usage:"), primaryText("project-builder"), mutedText("[command] [options]"))
	fmt.Fprintln(out)
	fmt.Fprintln(out, mutedText("Run without commands to open the interactive dashboard."))
	fmt.Fprintln(out)
	fmt.Fprintln(out, accentText("Commands:"))
	printMenuTable(out, []menuRow{
		{"install", "Install Project Builder and add it to PATH"},
		{"install status", "Show install and PATH status"},
		{"version, --version", "Show installed release version"},
		{"theme list", "List available themes / profiles"},
		{"theme set <name>", "Set theme: violet, cyan, emerald, amber, mono"},
		{"theme reset", "Reset theme to default"},
		{"--reconfigure", "Change global default workbench on next launch"},
		{"--no-color", "Disable ANSI colors"},
		{"help", "Show this help screen"},
	})
	fmt.Fprintln(out)
	fmt.Fprintln(out, accentText("Folder browser:"))
	printMenuTable(out, []menuRow{
		{"↑/↓ or j/k", "Move highlighted folder"},
		{"Enter", "Open highlighted folder"},
		{"Backspace", "Go to parent folder"},
		{"Space or s", "Select highlighted folder"},
		{"q or Ctrl+C", "Cancel selection; confirmation required"},
	})
}

func renderSettings(out io.Writer, cfg Config) {
	configPath, err := GetConfigFilePath()
	if err != nil {
		configPath = "unavailable: " + err.Error()
	}
	fmt.Fprintln(out, accentText("Settings"))
	printMenuTable(out, []menuRow{
		{"Active theme", activeThemeName(cfg)},
		{"Global workbench", valueOrUnset(cfg.DefaultWorkbench)},
		{"Config file", configPath},
		{"PATH install", "Run: project-builder install"},
	})
	fmt.Fprintln(out)
	fmt.Fprintln(out, promptText("[1]")+" "+primaryText("Change theme / profile"))
	fmt.Fprintln(out, promptText("[2]")+" "+primaryText("Change global default workbench"))
	fmt.Fprintln(out, promptText("[3]")+" "+primaryText("Back"))
	fmt.Fprint(out, promptText("Selection")+mutedText(": "))
}

func printMenuTable(out io.Writer, rows []menuRow) {
	commandWidth := 0
	descWidth := 0
	for _, row := range rows {
		if len(row.Command) > commandWidth {
			commandWidth = len(row.Command)
		}
		if len(row.Description) > descWidth {
			descWidth = len(row.Description)
		}
	}

	border := "+-" + strings.Repeat("-", commandWidth) + "-+-" + strings.Repeat("-", descWidth) + "-+"
	fmt.Fprintln(out, border)
	fmt.Fprintf(out, "| %-*s | %-*s |\n", commandWidth, "command", descWidth, "description")
	fmt.Fprintln(out, border)
	for _, row := range rows {
		fmt.Fprintf(out, "| %-*s | %-*s |\n", commandWidth, row.Command, descWidth, row.Description)
	}
	fmt.Fprintln(out, border)
}

func activeThemeName(cfg Config) string {
	if cfg.Theme == "" {
		return defaultThemeName
	}
	return cfg.Theme
}

func valueOrUnset(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "not set"
	}
	return trimmed
}
