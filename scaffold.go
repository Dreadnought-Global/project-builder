package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Discipline int

const (
	Design Discipline = iota + 1
	VideoMotion
	Audio
	Animation3D
)

func (d Discipline) String() string {
	switch d {
	case Design:
		return "Design"
	case VideoMotion:
		return "Video & Motion"
	case Audio:
		return "Audio"
	case Animation3D:
		return "3D & Animation"
	default:
		return "Unknown"
	}
}

// GetFolderList returns the relative folder structure paths for a discipline and client overlay.
func GetFolderList(d Discipline, isClient bool) ([]string, error) {
	var folders []string

	switch d {
	case Design:
		folders = []string{
			"01_Client_Brief_&_Docs",
			"02_Raw_Assets",
			"03_Design_Files",
			filepath.Join("03_Design_Files", "Figma_Exports"),
			filepath.Join("03_Design_Files", "Adobe_Source"),
			filepath.Join("03_Design_Files", "Assets_Vector"),
			"04_Review_&_Deliverables",
			filepath.Join("04_Review_&_Deliverables", "Drafts_For_Review"),
			filepath.Join("04_Review_&_Deliverables", "Final_Print_Ready"),
		}
		if isClient {
			folders = append(folders, filepath.Join("04_Review_&_Deliverables", "Final_Print_Ready", "Client_Handoff"))
		}
	case VideoMotion:
		folders = []string{
			"01_Pre_Production",
			"02_Project_Files",
			"03_Footage_Local",
			filepath.Join("03_Footage_Local", "Raw_Camera_Proxy"),
			filepath.Join("03_Footage_Local", "Assets_Stock"),
			"04_Audio_Production",
			filepath.Join("04_Audio_Production", "Voiceover_&_Dialog"),
			filepath.Join("04_Audio_Production", "BGM"),
			filepath.Join("04_Audio_Production", "SFX"),
			"05_Graphics_&_VFX",
			"06_Render_Outputs",
			filepath.Join("06_Render_Outputs", "Drafts"),
			filepath.Join("06_Render_Outputs", "Final_Masters"),
		}
		if isClient {
			folders = append(folders, filepath.Join("06_Render_Outputs", "Final_Masters", "Client_Handoff"))
		}
	case Audio:
		folders = []string{
			"01_Scripts_&_Copy",
			"02_Session_Files",
			"03_Raw_Recordings",
			"04_Assets_&_SFX",
			"05_Mastered_Outputs",
		}
		if isClient {
			folders = append(folders, filepath.Join("05_Mastered_Outputs", "Client_Handoff"))
		}
	case Animation3D:
		folders = []string{
			"01_Pre_Production_&_References",
			"02_Project_Files",
			"03_Assets",
			filepath.Join("03_Assets", "Models"),
			filepath.Join("03_Assets", "Textures"),
			filepath.Join("03_Assets", "Rigs"),
			"04_Renders",
			filepath.Join("04_Renders", "Drafts"),
			filepath.Join("04_Renders", "Final_Masters"),
		}
		if isClient {
			folders = append(folders, filepath.Join("04_Renders", "Final_Masters", "Client_Handoff"))
		}
	default:
		return nil, fmt.Errorf("invalid discipline: %d", d)
	}

	if isClient {
		clientDocs := []string{
			"00_Client_Docs",
			filepath.Join("00_Client_Docs", "Agreements"),
			filepath.Join("00_Client_Docs", "Invoices"),
			filepath.Join("00_Client_Docs", "Briefs"),
		}
		// Prepend client docs to root
		folders = append(clientDocs, folders...)
	}

	return folders, nil
}

// CreateFolderStructure generates the directories at the target path.
// If it fails, it attempts to clean up the newly created root directory.
func CreateFolderStructure(root string, d Discipline, isClient bool) error {
	folders, err := GetFolderList(d, isClient)
	if err != nil {
		return err
	}

	// Create root directory first
	if err := os.MkdirAll(root, 0755); err != nil {
		return fmt.Errorf("failed to create project root folder: %w", err)
	}

	// Keep track of directories created in this run, in case we need to roll back.
	rollback := func() {
		// Clean up the entire root folder
		os.RemoveAll(root)
	}

	for _, folder := range folders {
		path := filepath.Join(root, folder)
		if err := os.MkdirAll(path, 0755); err != nil {
			rollback()
			return fmt.Errorf("failed to create subdirectory %s: %w", folder, err)
		}
	}

	return nil
}
