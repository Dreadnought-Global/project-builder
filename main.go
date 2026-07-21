package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	opts, args := ParseGlobalOptions(os.Args[1:])

	cfg, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	renderOpts := DetectRenderOptions(opts.NoColor)
	theme := ActiveTheme(cfg)
	SetActiveStyle(theme, renderOpts)

	if handled, _, code := HandleCommand(args, cfg, os.Stdin, os.Stdout); handled {
		os.Exit(code)
	}

	reader := bufio.NewReader(os.Stdin)
	if err := runInteractiveApp(opts, cfg, reader, os.Stdout, renderOpts); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runInteractiveApp(opts CLIOptions, cfg Config, reader *bufio.Reader, out io.Writer, renderOpts RenderOptions) error {
	for {
		theme := ActiveTheme(cfg)
		SetActiveStyle(theme, renderOpts)
		clearTerminal(out)
		fmt.Fprint(out, RenderStartupBanner(theme, CurrentReleaseMetadata(), renderOpts))
		renderHomeDashboard(out)

		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading dashboard selection: %w", err)
		}
		choice, ok := ParseMenuChoiceWithDefault(input, 1, 4, 1)
		if !ok {
			fmt.Fprintln(out, errorText("Invalid selection. Please enter 1, 2, 3, or 4."))
			waitForEnter(reader, out)
			continue
		}

		switch choice {
		case 1:
			updated, err := ensureDefaultWorkbench(&cfg, opts.Reconfigure, reader, out)
			if err != nil {
				return err
			}
			if !updated && cfg.DefaultWorkbench == "" {
				continue
			}
			opts.Reconfigure = false
			if completed, err := runCreateProjectFlow(&cfg, reader, out); err != nil {
				return err
			} else if !completed {
				continue
			}
			return nil
		case 2:
			renderHelp(out)
			waitForEnter(reader, out)
		case 3:
			if err := runSettings(&cfg, reader, out); err != nil {
				return err
			}
		case 4:
			fmt.Fprintln(out, mutedText("Exiting Project Builder. Goodbye!"))
			return nil
		}
	}
}

func ensureDefaultWorkbench(cfg *Config, reconfigure bool, reader *bufio.Reader, out io.Writer) (bool, error) {
	if cfg.DefaultWorkbench != "" && !reconfigure {
		return true, nil
	}
	if reconfigure {
		fmt.Fprintln(out, accentText("Reconfiguring global default workbench path..."))
	} else {
		fmt.Fprintln(out, accentText("First-run setup: Global default workbench path not configured."))
	}
	clearTerminal(out)
	selectedPath, err := RunFolderBrowser()
	if err != nil {
		return false, fmt.Errorf("running folder browser: %w", err)
	}
	if selectedPath == "" {
		fmt.Fprintln(out, warningText("No folder selected. Returning to dashboard."))
		waitForEnter(reader, out)
		return false, nil
	}
	cfg.DefaultWorkbench = selectedPath
	if err := SaveConfig(*cfg); err != nil {
		return false, fmt.Errorf("saving configuration: %w", err)
	}
	fmt.Fprintf(out, "%s %s\n\n", successText("Saved global default workbench path:"), cfg.DefaultWorkbench)
	return true, nil
}

func runSettings(cfg *Config, reader *bufio.Reader, out io.Writer) error {
	for {
		renderSettings(out, *cfg)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading settings selection: %w", err)
		}
		choice, ok := ParseMenuChoice(input, 1, 3)
		if !ok {
			fmt.Fprintln(out, errorText("Invalid selection. Please enter 1, 2, or 3."))
			continue
		}
		switch choice {
		case 1:
			updated, code := promptThemeSelector(*cfg, reader, out)
			*cfg = updated
			if code != 0 {
				fmt.Fprintln(out, warningText("Theme unchanged."))
			}
			waitForEnter(reader, out)
		case 2:
			if _, err := ensureDefaultWorkbench(cfg, true, reader, out); err != nil {
				return err
			}
		case 3:
			return nil
		}
	}
}

func waitForEnter(reader *bufio.Reader, out io.Writer) {
	fmt.Fprintln(out, "\n"+mutedText("Press Enter to continue..."))
	_, _ = reader.ReadString('\n')
}

func runCreateProjectFlow(cfg *Config, reader *bufio.Reader, out io.Writer) (bool, error) {
	var projectName string
	var sanitizedName string
	for {
		fmt.Fprint(out, promptText("Project name:")+" ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("reading project name: %w", err)
		}
		projectName = strings.TrimSpace(input)
		if projectName == "" {
			fmt.Fprintln(out, "Project name cannot be empty. Please try again.")
			continue
		}

		sanitizedName = sanitizeProjectName(projectName)
		if sanitizedName == "" {
			fmt.Fprintln(out, "Project name contains only illegal characters. Please enter a valid name.")
			continue
		}
		break
	}

	clearTerminal(out)
	fmt.Fprintln(out, accentText("Select discipline:"))
	fmt.Fprintln(out, promptText("[1]")+" "+primaryText("Design"))
	fmt.Fprintln(out, promptText("[2]")+" "+primaryText("Video & Motion"))
	fmt.Fprintln(out, promptText("[3]")+" "+primaryText("Audio"))
	fmt.Fprintln(out, promptText("[4]")+" "+primaryText("3D & Animation"))
	fmt.Fprintln(out, promptText("[5]")+" "+primaryText("Back"))

	var disciplineChoice Discipline
	for {
		fmt.Fprint(out, promptText("Selection")+mutedText(": "))
		choiceStr, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("reading discipline selection: %w", err)
		}
		choice, ok := ParseMenuChoice(choiceStr, 1, 5)
		if ok {
			if choice == 5 {
				return false, nil
			}
			disciplineChoice = Discipline(choice)
			break
		}
		fmt.Fprintln(out, errorText("Invalid selection. Please enter a number between 1 and 5."))
	}

	var disciplineRoot string
	savedPath := cfg.GetDisciplinePath(disciplineChoice)
	declined := cfg.HasDeclinedDefault(disciplineChoice)

	if savedPath != "" {
		disciplineRoot = savedPath
	} else {
		for disciplineRoot == "" {
			clearTerminal(out)
			fmt.Fprintln(out, accentText("Where should this project be created?"))
			fmt.Fprintln(out, promptText("[1]")+" "+primaryText("Use Default Workbench")+mutedText(" ("+cfg.DefaultWorkbench+")"))
			fmt.Fprintln(out, promptText("[2]")+" "+primaryText("Select folder")+mutedText(" (native picker)"))
			fmt.Fprintln(out, promptText("[3]")+" "+primaryText("Select folder")+mutedText(" (terminal browser)"))
			fmt.Fprintln(out, promptText("[4]")+" "+primaryText("Back"))

			var folderChoice int
			for {
				fmt.Fprint(out, promptText("Selection")+mutedText(": "))
				choiceStr, err := reader.ReadString('\n')
				if err != nil {
					return false, fmt.Errorf("reading folder selection: %w", err)
				}
				choice, ok := ParseMenuChoice(choiceStr, 1, 4)
				if ok {
					folderChoice = choice
					break
				}
				fmt.Fprintln(out, errorText("Invalid selection. Please enter a number between 1 and 4."))
			}

			switch folderChoice {
			case 1:
				disciplineRoot = cfg.DefaultWorkbench
			case 2:
				fmt.Fprintln(out, accentText("Opening native folder picker..."))
				path, err := RunNativeFolderPicker()
				if err != nil {
					fmt.Fprintf(out, "No GUI folder picker found (%v) — using terminal browser instead.\n", err)
					path, err = RunFolderBrowser()
					if err != nil {
						return false, fmt.Errorf("running folder browser: %w", err)
					}
				}
				if path == "" {
					fmt.Fprintln(out, warningText("No folder selected. Returning to folder choices."))
					continue
				}
				disciplineRoot = path
			case 3:
				path, err := RunFolderBrowser()
				if err != nil {
					return false, fmt.Errorf("running folder browser: %w", err)
				}
				if path == "" {
					fmt.Fprintln(out, warningText("No folder selected. Returning to folder choices."))
					continue
				}
				disciplineRoot = path
			case 4:
				return false, nil
			}
		}

		if !declined {
			clearTerminal(out)
			fmt.Fprintf(out, "%s\n%s %s\n", accentText("Project location selected"), mutedText("Folder:"), primaryText(disciplineRoot))
			for {
				fmt.Fprintf(out, "\n%s %s%s", promptText("Set this as the default workbench for"), primaryText(disciplineChoice.String()), mutedText("? (y/n): "))
				input, err := reader.ReadString('\n')
				if err != nil {
					return false, fmt.Errorf("reading default save selection: %w", err)
				}
				isYes, ok := ParseYesNo(input)
				if ok && isYes {
					cfg.SetDisciplinePath(disciplineChoice, disciplineRoot)
					_ = SaveConfig(*cfg)
					fmt.Fprintln(out, successText("Default saved."))
					break
				} else if ok {
					cfg.SetDisciplinePath(disciplineChoice, declinedDisciplinePath)
					_ = SaveConfig(*cfg)
					break
				}
				fmt.Fprintln(out, errorText("Invalid input. Please enter 'y' or 'n'."))
			}
		}
	}

	if info, err := os.Stat(disciplineRoot); err != nil || !info.IsDir() {
		return false, fmt.Errorf("selected path %q is not accessible or is not a folder", disciplineRoot)
	}

	var isClient bool
	for {
		fmt.Fprint(out, "\n"+promptText("Client project?")+mutedText(" (y/n): "))
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("reading client project selection: %w", err)
		}
		isYes, ok := ParseYesNo(input)
		if ok {
			isClient = isYes
			break
		}
		fmt.Fprintln(out, errorText("Invalid input. Please enter 'y' or 'n'."))
	}

	var targetPath string
	for {
		targetPath = ResolveTargetPath(disciplineRoot, isClient, sanitizedName)

		if _, err := os.Stat(targetPath); err == nil {
			fmt.Fprintf(out, "\n%s %s\n%s\n\n", warningText("Warning: A folder named"), warningText("'"+sanitizedName+"' already exists at:"), targetPath)
			fmt.Fprintln(out, promptText("[1]")+" "+primaryText("Rename project"))
			fmt.Fprintln(out, promptText("[2]")+" "+primaryText("Append suffix")+mutedText(" (automatically append _01, _02, etc.)"))
			fmt.Fprintln(out, promptText("[3]")+" "+primaryText("Abort"))

			var collisionChoice int
			for {
				fmt.Fprint(out, promptText("Selection")+mutedText(": "))
				choiceStr, err := reader.ReadString('\n')
				if err != nil {
					return false, fmt.Errorf("reading collision selection: %w", err)
				}
				choice, ok := ParseMenuChoice(choiceStr, 1, 3)
				if ok {
					collisionChoice = choice
					break
				}
				fmt.Fprintln(out, errorText("Invalid selection. Please enter 1, 2, or 3."))
			}

			if collisionChoice == 1 {
				for {
					fmt.Fprint(out, "\n"+promptText("New project name:")+" ")
					input, err := reader.ReadString('\n')
					if err != nil {
						return false, fmt.Errorf("reading new project name: %w", err)
					}
					newName := strings.TrimSpace(input)
					if newName == "" {
						continue
					}
					newSanitized := sanitizeProjectName(newName)
					if newSanitized == "" {
						fmt.Fprintln(out, "Project name contains only illegal characters.")
						continue
					}
					sanitizedName = newSanitized
					break
				}
				continue
			} else if collisionChoice == 2 {
				pathExists := func(path string) bool {
					_, err := os.Stat(path)
					return err == nil
				}
				sanitizedName, targetPath = NextAvailableProjectName(disciplineRoot, isClient, sanitizedName, pathExists)
				fmt.Fprintf(out, "%s %s\n", successText("Using unique name:"), sanitizedName)
				break
			} else {
				fmt.Fprintln(out, warningText("Project creation cancelled. Returning to dashboard."))
				waitForEnter(reader, out)
				return false, nil
			}
		} else {
			break
		}
	}

	clearTerminal(out)
	fmt.Fprintln(out, accentText("Project Summary"))
	fmt.Fprintf(out, "%s    %s\n", mutedText("Project Name:"), primaryText(sanitizedName))
	fmt.Fprintf(out, "%s      %s\n", mutedText("Discipline:"), primaryText(disciplineChoice.String()))
	if isClient {
		fmt.Fprintf(out, "%s  %s\n", mutedText("Client Project:"), successText("Yes"))
	} else {
		fmt.Fprintf(out, "%s  %s\n", mutedText("Client Project:"), primaryText("No"))
	}
	fmt.Fprintf(out, "%s     %s\n", mutedText("Target Path:"), primaryText(targetPath))

	for {
		fmt.Fprint(out, promptText("Create project?")+mutedText(" (y/n): "))
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("reading create confirmation: %w", err)
		}
		isYes, ok := ParseYesNo(input)
		if ok && isYes {
			break
		} else if ok {
			fmt.Fprintln(out, warningText("Project creation cancelled. Returning to dashboard."))
			waitForEnter(reader, out)
			return false, nil
		}
		fmt.Fprintln(out, errorText("Invalid input. Please enter 'y' or 'n'."))
	}

	fmt.Fprintf(out, "\n%s %s...\n", accentText("Scaffolding folders in"), primaryText(targetPath))
	err := CreateFolderStructure(targetPath, disciplineChoice, isClient)
	if err != nil {
		return false, fmt.Errorf("folder generation failed: %w", err)
	}

	fmt.Fprintln(out, successText("Folders successfully created."))

	fmt.Fprintln(out, "\n"+promptText("[1]")+" "+primaryText("Open project in file manager"))
	fmt.Fprintln(out, promptText("[2]")+" "+primaryText("Exit"))
	for {
		fmt.Fprint(out, promptText("Selection")+mutedText(": "))
		choiceStr, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("reading post-creation selection: %w", err)
		}
		choice, ok := ParseMenuChoice(choiceStr, 1, 2)
		if ok {
			if choice == 1 {
				fmt.Fprintf(out, "%s %s\n", accentText("Opening folder:"), targetPath)
				if err := OpenFolder(targetPath); err != nil {
					fmt.Fprintf(out, "Failed to open folder: %v\n", err)
				}
				waitForEnter(reader, out)
			}
			break
		}
		fmt.Fprintln(out, errorText("Invalid selection. Please enter 1 or 2."))
	}

	fmt.Fprintln(out, mutedText("Exiting Project Builder. Goodbye!"))
	return true, nil
}

// sanitizeProjectName cleans up a user's input name to be safe across filesystems.
func sanitizeProjectName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "\t", "_")

	var sb strings.Builder
	for _, r := range name {
		if r < 32 || r == 127 {
			continue
		}
		switch r {
		case '<', '>', ':', '"', '/', '\\', '|', '?', '*', ';', '&', '(', ')', '[', ']', '{', '}':
			continue
		default:
			sb.WriteRune(r)
		}
	}

	sanitized := sb.String()

	for strings.Contains(sanitized, "__") {
		sanitized = strings.ReplaceAll(sanitized, "__", "_")
	}

	sanitized = strings.Trim(sanitized, "_.")
	return sanitized
}
