package main

import (
	"strings"
	"testing"
)

func TestBannerNoColorHasNoANSI(t *testing.T) {
	theme, _ := GetTheme("violet")
	out := RenderStartupBanner(theme, ReleaseMetadata{Version: "2.0.0", ReleaseDate: "2026-07-15"}, RenderOptions{UseColor: false, Width: 140})
	if strings.Contains(out, "\x1b[") || strings.Contains(out, "\x1b]") {
		t.Fatalf("expected no ANSI in no-color output: %q", out)
	}
	if !strings.Contains(out, "Project Builder v2.0.0  |  Project Scaffolding & Workspace Automation Tool  |  2026-07-15  |  dreadnought.studio") {
		t.Fatalf("metadata line missing: %q", out)
	}
}

func TestBannerNarrowKeepsAnsiHeading(t *testing.T) {
	theme, _ := GetTheme("violet")
	out := RenderStartupBanner(theme, ReleaseMetadata{Version: "2.0.0", ReleaseDate: "2026-07-15"}, RenderOptions{UseColor: false, Width: 40})
	lines := strings.Split(out, "\n")
	if lines[0] != projectBuilderBanner[0] {
		t.Fatalf("expected full banner heading, got %q", lines[0])
	}
}

func TestBannerHasNoBorders(t *testing.T) {
	theme, _ := GetTheme("violet")
	out := RenderStartupBanner(theme, ReleaseMetadata{Version: "2.0.0", ReleaseDate: "2026-07-15"}, RenderOptions{UseColor: false, Width: 140})
	for _, line := range strings.Split(out, "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 5 && (strings.Count(trimmed, "=") == len(trimmed) || strings.Count(trimmed, "-") == len(trimmed)) {
			t.Fatalf("unexpected border line: %q", line)
		}
	}
}
