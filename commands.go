package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type CLIOptions struct {
	Reconfigure bool
	NoColor     bool
}

func ParseGlobalOptions(args []string) (CLIOptions, []string) {
	opts := CLIOptions{}
	var rest []string
	for _, arg := range args {
		switch arg {
		case "--reconfigure", "-r":
			opts.Reconfigure = true
		case "--no-color":
			opts.NoColor = true
		default:
			rest = append(rest, arg)
		}
	}
	return opts, rest
}

func HandleCommand(args []string, cfg Config, in io.Reader, out io.Writer) (bool, Config, int) {
	if len(args) == 0 {
		return false, cfg, 0
	}
	switch args[0] {
	case "help", "--help", "-h":
		renderHelp(out)
		return true, cfg, 0
	case "install":
		return true, cfg, handleInstallCommand(args[1:], in, out)
	case "version", "--version", "-v":
		metadata := CurrentReleaseMetadata()
		fmt.Fprintf(out, "Project Builder %s (%s)\n", metadata.DisplayVersion(), metadata.ReleaseDate)
		return true, cfg, 0
	case "theme":
		cfg, code := handleThemeCommand(args[1:], cfg, in, out)
		return true, cfg, code
	default:
		fmt.Fprintf(out, "Unknown command: %s\n", args[0])
		fmt.Fprintln(out, "Run `project-builder help` for commands.")
		return true, cfg, 1
	}
}

func handleThemeCommand(args []string, cfg Config, in io.Reader, out io.Writer) (Config, int) {
	active := cfg.Theme
	if active == "" {
		active = defaultThemeName
	}
	if len(args) == 0 {
		fmt.Fprintf(out, "Available themes (active: %s):\n", active)
		printThemeList(out, active)
		return cfg, 0
	}
	switch args[0] {
	case "list":
		fmt.Fprintf(out, "Available themes (active: %s):\n", active)
		printThemeList(out, active)
		return cfg, 0
	case "set":
		if len(args) < 2 {
			fmt.Fprintf(out, "Usage: project-builder theme set <%s>\n", strings.Join(ThemeNames(), "|"))
			return cfg, 1
		}
		name := strings.ToLower(strings.TrimSpace(args[1]))
		if _, err := GetTheme(name); err != nil {
			fmt.Fprintf(out, "%v\n", err)
			return cfg, 1
		}
		cfg.Theme = name
		if err := SaveConfig(cfg); err != nil {
			fmt.Fprintf(out, "Error saving theme: %v\n", err)
			return cfg, 1
		}
		fmt.Fprintf(out, "Theme set to %s.\n", name)
		return cfg, 0
	case "reset":
		cfg.Theme = defaultThemeName
		if err := SaveConfig(cfg); err != nil {
			fmt.Fprintf(out, "Error saving theme: %v\n", err)
			return cfg, 1
		}
		fmt.Fprintf(out, "Theme reset to %s.\n", defaultThemeName)
		return cfg, 0
	case "select":
		return promptThemeSelector(cfg, in, out)
	default:
		fmt.Fprintf(out, "Unknown theme command: %s\n", args[0])
		fmt.Fprintln(out, "Usage: project-builder theme [list|set <name>|reset]")
		return cfg, 1
	}
}

func printThemeList(out io.Writer, active string) {
	for _, name := range ThemeNames() {
		marker := " "
		if name == active {
			marker = "*"
		}
		fmt.Fprintf(out, "%s %s\n", marker, name)
	}
}

func promptThemeSelector(cfg Config, in io.Reader, out io.Writer) (Config, int) {
	names := ThemeNames()
	fmt.Fprintln(out, "Select theme:")
	for i, name := range names {
		fmt.Fprintf(out, "[%d] %s\n", i+1, name)
	}
	fmt.Fprintf(out, "Selection (1-%d): ", len(names))
	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil {
		fmt.Fprintf(out, "Error reading input: %v\n", err)
		return cfg, 1
	}
	choice, ok := ParseMenuChoice(line, 1, len(names))
	if !ok {
		fmt.Fprintln(out, "Invalid selection.")
		return cfg, 1
	}
	cfg.Theme = names[choice-1]
	if err := SaveConfig(cfg); err != nil {
		fmt.Fprintf(out, "Error saving theme: %v\n", err)
		return cfg, 1
	}
	fmt.Fprintf(out, "Theme set to %s.\n", cfg.Theme)
	return cfg, 0
}
