package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("    Project Builder - Dreadnought Studio")
	fmt.Println("========================================")
	fmt.Println()

	// Parse --reconfigure flag
	reconfigure := false
	for _, arg := range os.Args[1:] {
		if arg == "--reconfigure" || arg == "-r" {
			reconfigure = true
			break
		}
	}

	// Load configuration
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// 1. Global Default Workbench setup
	if cfg.DefaultWorkbench == "" || reconfigure {
		if reconfigure {
			fmt.Println("Reconfiguring global default workbench path...")
		} else {
			fmt.Println("First-run setup: Global default workbench path not configured.")
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
		fmt.Printf("Saved global default workbench path: %s\n\n", cfg.DefaultWorkbench)
	}

	reader := bufio.NewReader(os.Stdin)

	// 2. Project Name
	var projectName string
	var sanitizedName string
	for {
		fmt.Print("Project name: ")
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

	// 3. Select Discipline
	fmt.Println("\nSelect discipline:")
	fmt.Println("[1] Design")
	fmt.Println("[2] Video & Motion")
	fmt.Println("[3] Audio")
	fmt.Println("[4] 3D & Animation")

	var disciplineChoice Discipline
	for {
		fmt.Print("Selection (1-4): ")
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
		fmt.Println("Invalid selection. Please enter a number between 1 and 4.")
	}

	// 4. Per-Discipline Destination Selection Flow
	var disciplineRoot string
	savedPath := cfg.GetDisciplinePath(disciplineChoice)
	declined := cfg.HasDeclinedDefault(disciplineChoice)

	if savedPath != "" {
		disciplineRoot = savedPath
	} else {
		// Show 3-option menu
		fmt.Printf("\nWhere should this project be created?\n")
		fmt.Printf("[1] Use Default Workbench (%s)\n", cfg.DefaultWorkbench)
		fmt.Println("[2] Select folder (native picker)")
		fmt.Println("[3] Select folder (terminal browser)")

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
			fmt.Println("Invalid selection. Please enter 1, 2, or 3.")
		}

		if folderChoice == 1 {
			disciplineRoot = cfg.DefaultWorkbench
		} else if folderChoice == 2 {
			fmt.Println("Opening native folder picker...")
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
				fmt.Println("No folder selected. Aborting.")
				os.Exit(0)
			}
			disciplineRoot = path
		} else if folderChoice == 3 {
			path, err := RunFolderBrowser()
			if err != nil {
				fmt.Printf("Error running folder browser: %v\n", err)
				os.Exit(1)
			}
			if path == "" {
				fmt.Println("No folder selected. Aborting.")
				os.Exit(0)
			}
			disciplineRoot = path
		}

		// Ask to set as default ONLY if not previously declined
		if !declined {
			for {
				fmt.Printf("\nSet this as the default workbench for %s? (y/n): ", disciplineChoice.String())
				input, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("Error reading input: %v\n", err)
					os.Exit(1)
				}
				isYes, ok := ParseYesNo(input)
				if ok && isYes {
					cfg.SetDisciplinePath(disciplineChoice, disciplineRoot)
					_ = SaveConfig(cfg)
					fmt.Println("Default saved.")
					break
				} else if ok {
					cfg.SetDisciplinePath(disciplineChoice, declinedDisciplinePath)
					_ = SaveConfig(cfg)
					break
				}
				fmt.Println("Invalid input. Please enter 'y' or 'n'.")
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
		fmt.Print("\nClient project? (y/n): ")
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
		fmt.Println("Invalid input. Please enter 'y' or 'n'.")
	}

	// 6. Target Path Resolution & Collision Check
	var targetPath string
	for {
		targetPath = ResolveTargetPath(disciplineRoot, isClient, sanitizedName)

		if _, err := os.Stat(targetPath); err == nil {
			fmt.Printf("\nWarning: A folder named '%s' already exists at:\n%s\n\n", sanitizedName, targetPath)
			fmt.Println("[1] Rename project")
			fmt.Println("[2] Append suffix (automatically append _01, _02, etc.)")
			fmt.Println("[3] Abort")

			var collisionChoice int
			for {
				fmt.Print("Selection (1-3): ")
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
				fmt.Println("Invalid selection. Please enter 1, 2, or 3.")
			}

			if collisionChoice == 1 {
				for {
					fmt.Print("\nNew project name: ")
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
				fmt.Printf("Using unique name: %s\n", sanitizedName)
				break
			} else {
				fmt.Println("Aborted by user.")
				os.Exit(0)
			}
		} else {
			// No collision
			break
		}
	}

	// 7. Show Summary & Confirm Generation
	fmt.Println("\n--- Project Summary ---")
	fmt.Printf("Project Name:    %s\n", sanitizedName)
	fmt.Printf("Discipline:      %s\n", disciplineChoice.String())
	if isClient {
		fmt.Println("Client Project:  Yes")
	} else {
		fmt.Println("Client Project:  No")
	}
	fmt.Printf("Target Path:     %s\n", targetPath)
	fmt.Println("-----------------------")

	for {
		fmt.Print("Create project? (y/n): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			os.Exit(1)
		}
		isYes, ok := ParseYesNo(input)
		if ok && isYes {
			break
		} else if ok {
			fmt.Println("Cancelled project creation.")
			os.Exit(0)
		}
		fmt.Println("Invalid input. Please enter 'y' or 'n'.")
	}

	// 8. Generate Folder Structure
	fmt.Printf("\nScaffolding folders in %s...\n", targetPath)
	err = CreateFolderStructure(targetPath, disciplineChoice, isClient)
	if err != nil {
		fmt.Printf("Error during folder generation: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Folders successfully created.")

	// 9. Post-creation actions
	fmt.Println("\n[1] Open project in file manager")
	fmt.Println("[2] Exit")
	for {
		fmt.Print("Selection (1-2): ")
		choiceStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			os.Exit(1)
		}
		choice, ok := ParseMenuChoice(choiceStr, 1, 2)
		if ok {
			if choice == 1 {
				fmt.Printf("Opening folder: %s\n", targetPath)
				if err := OpenFolder(targetPath); err != nil {
					fmt.Printf("Failed to open folder: %v\n", err)
				}
				fmt.Println("\nPress Enter to exit...")
				_, _ = reader.ReadString('\n')
			}
			break
		}
		fmt.Println("Invalid selection. Please enter 1 or 2.")
	}

	fmt.Println("Exiting Project Builder. Goodbye!")
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
