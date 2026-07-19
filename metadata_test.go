package main

import "testing"

func TestParseLatestReleaseMetadata(t *testing.T) {
	changelog := `# Changelog

## [1.2.0] - 2026-07-19
### Added
- Theme system

## [1.1.0] - 2026-07-15
`
	metadata, err := ParseLatestReleaseMetadata(changelog)
	if err != nil {
		t.Fatalf("ParseLatestReleaseMetadata failed: %v", err)
	}
	if metadata.Version != "1.2.0" || metadata.ReleaseDate != "2026-07-19" {
		t.Fatalf("unexpected metadata: %+v", metadata)
	}
	if metadata.DisplayVersion() != "v1.2.0" {
		t.Fatalf("expected v1.2.0, got %q", metadata.DisplayVersion())
	}
}

func TestParseLatestReleaseMetadataSkipsMalformed(t *testing.T) {
	changelog := `## [unreleased]
## [1.2] - yesterday
## [2.0.0] - 2026-07-15
`
	metadata, err := ParseLatestReleaseMetadata(changelog)
	if err != nil {
		t.Fatalf("ParseLatestReleaseMetadata failed: %v", err)
	}
	if metadata.Version != "2.0.0" || metadata.ReleaseDate != "2026-07-15" {
		t.Fatalf("unexpected metadata: %+v", metadata)
	}
}

func TestParseLatestReleaseMetadataNoValidRelease(t *testing.T) {
	if _, err := ParseLatestReleaseMetadata("## [bad] - nope"); err == nil {
		t.Fatalf("expected error for missing valid release")
	}
}
