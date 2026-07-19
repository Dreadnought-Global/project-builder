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
	fmt.Fprint(out, promptText("Selection")+mutedText(" (1-4, Enter=1): "))
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
		{"theme list", "List available themes / profiles"},
		{"theme set <name>", "Set theme: violet, cyan, emerald, amber, mono"},
		{"theme reset", "Reset theme to default"},
		{"--reconfigure", "Change global default workbench on next launch"},
		{"--no-color", "Disable ANSI colors"},
		{"help", "Show this help screen"},
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
	fmt.Fprint(out, promptText("Selection")+mutedText(" (1-3): "))
}

func printMenuTable(out io.Writer, rows []menuRow) {
	commandWidth := len("command")
	for _, row := range rows {
		if len(row.Command) > commandWidth {
			commandWidth = len(row.Command)
		}
	}
	border := "+-" + strings.Repeat("-", commandWidth) + "-+-------------+"
	fmt.Fprintln(out, border)
	fmt.Fprintf(out, "| %-*s | description |\n", commandWidth, "command")
	fmt.Fprintln(out, border)
	for _, row := range rows {
		fmt.Fprintf(out, "| %-*s | %s |\n", commandWidth, row.Command, row.Description)
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
