package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

const declinedDisciplinePath = "declined"

func ResolveTargetPath(disciplineRoot string, isClient bool, sanitizedName string) string {
	return filepath.Join(NormalizeDisciplineRoot(disciplineRoot), ProjectCategorySubfolder(isClient), sanitizedName)
}

func NormalizeDisciplineRoot(path string) string {
	cleaned := filepath.Clean(strings.TrimSpace(path))
	switch filepath.Base(cleaned) {
	case "00_Client_Projects", "01_Passion_Projects":
		return filepath.Dir(cleaned)
	default:
		return cleaned
	}
}

func NextAvailableProjectName(disciplineRoot string, isClient bool, sanitizedName string, exists func(string) bool) (string, string) {
	for suffix := 1; ; suffix++ {
		candidateName := fmt.Sprintf("%s_%02d", sanitizedName, suffix)
		candidatePath := ResolveTargetPath(disciplineRoot, isClient, candidateName)
		if !exists(candidatePath) {
			return candidateName, candidatePath
		}
	}
}

func ParseMenuChoice(input string, min int, max int) (int, bool) {
	choice, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || choice < min || choice > max {
		return 0, false
	}
	return choice, true
}

func ParseMenuChoiceWithDefault(input string, min int, max int, defaultChoice int) (int, bool) {
	if strings.TrimSpace(input) == "" {
		if defaultChoice < min || defaultChoice > max {
			return 0, false
		}
		return defaultChoice, true
	}
	return ParseMenuChoice(input, min, max)
}

func ParseYesNo(input string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "y", "yes":
		return true, true
	case "n", "no":
		return false, true
	default:
		return false, false
	}
}
