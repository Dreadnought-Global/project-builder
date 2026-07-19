package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	opts, args := ParseGlobalOptions(os.Args[1:])

	// Load configuration
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	if handled, _, code := HandleCommand(args, cfg, os.Stdin, os.Stdout); handled {
		os.Exit(code)
	}

	renderOpts := DetectRenderOptions(opts.NoColor)
	theme := ActiveTheme(cfg)
	SetActiveStyle(theme, renderOpts)
	fmt.Print(RenderStartupBanner(theme, CurrentReleaseMetadata(), renderOpts))

	// 1. Global Default Workbench setup
	if cfg.DefaultWorkbench == "" || opts.Reconfigure {
		if opts.Reconfigure {
			fmt.Println(accentText("Reconfiguring global default workbench path..."))
		} else {
			fmt.Println(accentText("First-run setup: Global default workbench path not configured."))
		}
		selectedPath, err := RunFolderBrowser()
		if err != nil {
			fmt.Printf("Error running folder browser: %v\n", err)
			os.Exit(1)
		}
		if selectedPath == "" {
			fmt.Println("No folder selected. Configuration aborted.")
			os.Exit(0)
		}
		cfg.DefaultWorkbench = selectedPath
		if err := SaveConfig(cfg); err != nil {
			fmt.Printf("Error saving configuration: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s %s\n\n", successText("Saved global default workbench path:"), cfg.DefaultWorkbench)
	}

	reader := bufio.NewReader(os.Stdin)

	// 2. Project Name
	var projectName string
	var sanitizedName string
	for {
		fmt.Print(promptText("Project name:") + " ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			os.Exit(1)
		}
		projectName = strings.TrimSpace(input)
		if projectName == "" {
			fmt.Println("Project name cannot be empty. Please try again.")
			continue
		}

		sanitizedName = sanitizeProjectName(projectName)
		if sanitizedName == "" {
			fmt.Println("Project name contains only illegal characters. Please enter a valid name.")
			continue
		}
		break
	}

	fmt.Println(accentText("\nSelect discipline:"))
	fmt.Println(promptText("[1]") + " " + primaryText("Design"))
	fmt.Println(promptText("[2]") + " " + primaryText("Video & Motion"))
	fmt.Println(promptText("[3]") + " " + primaryText("Audio"))
	fmt.Println(promptText("[4]") + " " + primaryText("3D & Animation"))

	var disciplineChoice Discipline
	for {
		fmt.Print(promptText("Selection") + mutedText(" (1-4): "))
		choiceStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			os.Exit(1)
		}
		choice, ok := ParseMenuChoice(choiceStr, 1, 4)
		if ok {
			disciplineChoice = Discipline(choice)
			break
		}
		fmt.Println(errorText("Invalid selection. Please enter a number between 1 and 4."))
	}

	// 4. Per-Discipline Destination Selection Flow
	var disciplineRoot string
	savedPath := cfg.GetDisciplinePath(disciplineChoice)
	declined := cfg.HasDeclinedDefault(disciplineChoice)

	if savedPath != "" {
		disciplineRoot = savedPath
	} else {
		// Show 3-option menu
		fmt.Printf("\n%s\n", accentText("Where should this project be created?"))
		fmt.Println(promptText("[1]") + " " + primaryText("Use Default Workbench") + mutedText(" ("+cfg.DefaultWorkbench+")"))
		fmt.Println(promptText("[2]") + " " + primaryText("Select folder") + mutedText(" (native picker)"))
		fmt.Println(promptText("[3]") + " " + primaryText("Select folder") + mutedText(" (terminal browser)"))

		var folderChoice int
		for {
			fmt.Print("Selection (1-3): ")
			choiceStr, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				os.Exit(1)
			}
			choice, ok := ParseMenuChoice(choiceStr, 1, 3)
			if ok {
				folderChoice = choice
				break
			}
			fmt.Println(errorText("Invalid selection. Please enter 1, 2, or 3."))
		}

		switch folderChoice {
		case 1:
			disciplineRoot = cfg.DefaultWorkbench
		case 2:
			fmt.Println(accentText("Opening native folder picker..."))
			path, err := RunNativeFolderPicker()
			if err != nil {
				fmt.Printf("No GUI folder picker found (%v) — using terminal browser instead.\n", err)
				path, err = RunFolderBrowser()
				if err != nil {
					fmt.Printf("Error running folder browser: %v\n", err)
					os.Exit(1)
				}
			}
			if path == "" {
				fmt.Println(warningText("No folder selected. Aborting."))
				os.Exit(0)
			}
			disciplineRoot = path
		case 3:
			path, err := RunFolderBrowser()
			if err != nil {
				fmt.Printf("Error running folder browser: %v\n", err)
				os.Exit(1)
			}
			if path == "" {
				fmt.Println(warningText("No folder selected. Aborting."))
				os.Exit(0)
			}
			disciplineRoot = path
		}

		// Ask to set as default ONLY if not previously declined
		if !declined {
			for {
				fmt.Printf("\n%s %s%s", promptText("Set this as the default workbench for"), primaryText(disciplineChoice.String()), mutedText("? (y/n): "))
				input, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("Error reading input: %v\n", err)
					os.Exit(1)
				}
				isYes, ok := ParseYesNo(input)
				if ok && isYes {
					cfg.SetDisciplinePath(disciplineChoice, disciplineRoot)
					_ = SaveConfig(cfg)
					fmt.Println(successText("Default saved."))
					break
				} else if ok {
					cfg.SetDisciplinePath(disciplineChoice, declinedDisciplinePath)
					_ = SaveConfig(cfg)
					break
				}
				fmt.Println(errorText("Invalid input. Please enter 'y' or 'n'."))
			}
		}
	}

	// Verify discipline root is valid
	if info, err := os.Stat(disciplineRoot); err != nil || !info.IsDir() {
		fmt.Printf("Error: The selected path '%s' is not accessible or is not a folder.\n", disciplineRoot)
		os.Exit(1)
	}

	// 5. Client Project?
	var isClient bool
	for {
		fmt.Print("\n" + promptText("Client project?") + mutedText(" (y/n): "))
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			os.Exit(1)
		}
		isYes, ok := ParseYesNo(input)
		if ok {
			isClient = isYes
			break
		}
		fmt.Println(errorText("Invalid input. Please enter 'y' or 'n'."))
	}

	// 6. Target Path Resolution & Collision Check
	var targetPath string
	for {
		targetPath = ResolveTargetPath(disciplineRoot, isClient, sanitizedName)

		if _, err := os.Stat(targetPath); err == nil {
			fmt.Printf("\n%s %s\n%s\n\n", warningText("Warning: A folder named"), warningText("'"+sanitizedName+"' already exists at:"), targetPath)
			fmt.Println(promptText("[1]") + " " + primaryText("Rename project"))
			fmt.Println(promptText("[2]") + " " + primaryText("Append suffix") + mutedText(" (automatically append _01, _02, etc.)"))
			fmt.Println(promptText("[3]") + " " + primaryText("Abort"))

			var collisionChoice int
			for {
				fmt.Print(promptText("Selection") + mutedText(" (1-3): "))
				choiceStr, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("Error reading input: %v\n", err)
					os.Exit(1)
				}
				choice, ok := ParseMenuChoice(choiceStr, 1, 3)
				if ok {
					collisionChoice = choice
					break
				}
				fmt.Println(errorText("Invalid selection. Please enter 1, 2, or 3."))
			}

			if collisionChoice == 1 {
				for {
					fmt.Print("\n" + promptText("New project name:") + " ")
					input, err := reader.ReadString('\n')
					if err != nil {
						fmt.Printf("Error reading input: %v\n", err)
						os.Exit(1)
					}
					newName := strings.TrimSpace(input)
					if newName == "" {
						continue
					}
					newSanitized := sanitizeProjectName(newName)
					if newSanitized == "" {
						fmt.Println("Project name contains only illegal characters.")
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
				fmt.Printf("%s %s\n", successText("Using unique name:"), sanitizedName)
				break
			} else {
				fmt.Println(warningText("Aborted by user."))
				os.Exit(0)
			}
		} else {
			// No collision
			break
		}
	}

	// 7. Show Summary & Confirm Generation
	fmt.Println("\n" + accentText("Project Summary"))
	fmt.Printf("%s    %s\n", mutedText("Project Name:"), primaryText(sanitizedName))
	fmt.Printf("%s      %s\n", mutedText("Discipline:"), primaryText(disciplineChoice.String()))
	if isClient {
		fmt.Printf("%s  %s\n", mutedText("Client Project:"), successText("Yes"))
	} else {
		fmt.Printf("%s  %s\n", mutedText("Client Project:"), primaryText("No"))
	}
	fmt.Printf("%s     %s\n", mutedText("Target Path:"), primaryText(targetPath))

	for {
		fmt.Print(promptText("Create project?") + mutedText(" (y/n): "))
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			os.Exit(1)
		}
		isYes, ok := ParseYesNo(input)
		if ok && isYes {
			break
		} else if ok {
			fmt.Println(warningText("Cancelled project creation."))
			os.Exit(0)
		}
		fmt.Println(errorText("Invalid input. Please enter 'y' or 'n'."))
	}

	// 8. Generate Folder Structure
	fmt.Printf("\n%s %s...\n", accentText("Scaffolding folders in"), primaryText(targetPath))
	err = CreateFolderStructure(targetPath, disciplineChoice, isClient)
	if err != nil {
		fmt.Printf("Error during folder generation: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(successText("Folders successfully created."))

	// 9. Post-creation actions
	fmt.Println("\n" + promptText("[1]") + " " + primaryText("Open project in file manager"))
	fmt.Println(promptText("[2]") + " " + primaryText("Exit"))
	for {
		fmt.Print(promptText("Selection") + mutedText(" (1-2): "))
		choiceStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			os.Exit(1)
		}
		choice, ok := ParseMenuChoice(choiceStr, 1, 2)
		if ok {
			if choice == 1 {
				fmt.Printf("%s %s\n", accentText("Opening folder:"), targetPath)
				if err := OpenFolder(targetPath); err != nil {
					fmt.Printf("Failed to open folder: %v\n", err)
				}
				fmt.Println("\n" + mutedText("Press Enter to exit..."))
				_, _ = reader.ReadString('\n')
			}
			break
		}
		fmt.Println(errorText("Invalid selection. Please enter 1 or 2."))
	}

	fmt.Println(mutedText("Exiting Project Builder. Goodbye!"))
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
