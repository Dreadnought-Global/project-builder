package main

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
)

//go:embed CHANGELOG.md
var embeddedChangelog string

const (
	appDescription = "Project Scaffolding & Workspace Automation Tool"
	studioLabel    = "dreadnought.studio"
	studioURL      = "https://www.instagram.com/dreadnought.sc/"
)

type ReleaseMetadata struct {
	Version     string
	ReleaseDate string
}

func CurrentReleaseMetadata() ReleaseMetadata {
	metadata, err := ParseLatestReleaseMetadata(embeddedChangelog)
	if err != nil {
		return ReleaseMetadata{Version: "unknown", ReleaseDate: "release date unavailable"}
	}
	return metadata
}

func ParseLatestReleaseMetadata(changelog string) (ReleaseMetadata, error) {
	re := regexp.MustCompile(`(?m)^## \[([0-9]+\.[0-9]+\.[0-9]+)\] - ([0-9]{4}-[0-9]{2}-[0-9]{2})\s*$`)
	match := re.FindStringSubmatch(changelog)
	if match == nil {
		return ReleaseMetadata{}, fmt.Errorf("no valid release heading found")
	}
	return ReleaseMetadata{Version: match[1], ReleaseDate: match[2]}, nil
}

func (m ReleaseMetadata) DisplayVersion() string {
	version := strings.TrimSpace(m.Version)
	if version == "" {
		return "vunknown"
	}
	if strings.HasPrefix(version, "v") {
		return version
	}
	return "v" + version
}

func FormatMetadataLine(metadata ReleaseMetadata, opts RenderOptions) string {
	studio := studioLabel
	if opts.UseColor {
		studio = hyperlink(studioLabel, studioURL)
	}
	return fmt.Sprintf("Project Builder %s  |  %s  |  %s  |  %s", metadata.DisplayVersion(), appDescription, metadata.ReleaseDate, studio)
}
