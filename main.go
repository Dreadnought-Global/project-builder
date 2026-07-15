package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

	// Check if we need to run first-run setup TUI
	if cfg.WorkbenchPath == "" || reconfigure {
		if reconfigure {
			fmt.Println("Reconfiguring root workbench path...")
		} else {
			fmt.Println("First-run setup: Root workbench path not configured.")
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
		cfg.WorkbenchPath = selectedPath
		if err := SaveConfig(cfg); err != nil {
			fmt.Printf("Error saving configuration: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Saved root workbench path: %s\n\n", cfg.WorkbenchPath)
	}

	rootDir := cfg.WorkbenchPath

	// Verify that rootDir exists and is a directory
	if info, err := os.Stat(rootDir); err != nil || !info.IsDir() {
		fmt.Printf("Error: The saved workbench path '%s' is not accessible or is not a folder.\n", rootDir)
		fmt.Println("Please run project-builder with the '--reconfigure' flag to set a new workbench path.")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	var projectName string
	var sanitizedName string
	var targetPath string

	// 1. Project Name and Collision Handling Loop
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

		targetPath = filepath.Join(rootDir, sanitizedName)

		// Check collision
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
				choiceStr = strings.TrimSpace(choiceStr)
				choice, err := strconv.Atoi(choiceStr)
				if err == nil && choice >= 1 && choice <= 3 {
					collisionChoice = choice
					break
				}
				fmt.Println("Invalid selection. Please enter 1, 2, or 3.")
			}

			if collisionChoice == 1 {
				// Loop back to prompt for a new name
				continue
			} else if collisionChoice == 2 {
				// Append suffix
				suffix := 1
				var testPath string
				var testName string
				for {
					testName = fmt.Sprintf("%s_%02d", sanitizedName, suffix)
					testPath = filepath.Join(rootDir, testName)
					if _, err := os.Stat(testPath); os.IsNotExist(err) {
						sanitizedName = testName
						targetPath = testPath
						break
					}
					suffix++
				}
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

	// 2. Select Discipline
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
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err == nil && choice >= 1 && choice <= 4 {
			disciplineChoice = Discipline(choice)
			break
		}
		fmt.Println("Invalid selection. Please enter a number between 1 and 4.")
	}

	// 3. Client Project?
	var isClient bool
	for {
		fmt.Print("Client project? (y/n): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			os.Exit(1)
		}
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "y" || input == "yes" {
			isClient = true
			break
		} else if input == "n" || input == "no" {
			isClient = false
			break
		}
		fmt.Println("Invalid input. Please enter 'y' or 'n'.")
	}

	// 4. Show Summary & Confirm Generation
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
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "y" || input == "yes" {
			break
		} else if input == "n" || input == "no" {
			fmt.Println("Cancelled project creation.")
			os.Exit(0)
		}
		fmt.Println("Invalid input. Please enter 'y' or 'n'.")
	}

	// 5. Generate Folder Structure
	fmt.Printf("\nScaffolding folders in %s...\n", targetPath)
	err = CreateFolderStructure(targetPath, disciplineChoice, isClient)
	if err != nil {
		fmt.Printf("Error during folder generation: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Folders successfully created.")

	// 6. Post-creation actions
	fmt.Println("\n[1] Open project in file manager")
	fmt.Println("[2] Exit")
	for {
		fmt.Print("Selection (1-2): ")
		choiceStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			os.Exit(1)
		}
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err == nil && choice >= 1 && choice <= 2 {
			if choice == 1 {
				fmt.Printf("Opening folder: %s\n", targetPath)
				if err := OpenFolder(targetPath); err != nil {
					fmt.Printf("Failed to open folder: %v\n", err)
				}
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
