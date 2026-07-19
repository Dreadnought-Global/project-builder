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
	for _, want := range []string{"Project Builder v2.0.0", appDescription, "Release: 2026-07-15", "Creator: dreadnought.studio", "Repo:    github.com/Dreadnought-Global/project-builder"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in banner: %q", want, out)
		}
	}
}

func TestBannerUsesPBMark(t *testing.T) {
	theme, _ := GetTheme("violet")
	out := RenderStartupBanner(theme, ReleaseMetadata{Version: "2.0.0", ReleaseDate: "2026-07-15"}, RenderOptions{UseColor: false, Width: 140, Height: 40})
	lines := strings.Split(out, "\n")
	if lines[0] != pbBanner[0]+"  Project Builder v2.0.0" {
		t.Fatalf("expected PB banner with title, got %q", lines[0])
	}
}

func TestBannerStacksWhenNarrow(t *testing.T) {
	theme, _ := GetTheme("violet")
	out := RenderStartupBanner(theme, ReleaseMetadata{Version: "2.0.0", ReleaseDate: "2026-07-15"}, RenderOptions{UseColor: false, Width: 50, Height: 40})
	lines := strings.Split(out, "\n")
	if lines[0] != pbBanner[0] {
		t.Fatalf("expected stacked PB icon first, got %q", lines[0])
	}
	if lines[len(pbBanner)] != "Project Builder v2.0.0" {
		t.Fatalf("expected title below PB icon, got %q", lines[len(pbBanner)])
	}
}

func TestBannerColorIncludesHyperlinks(t *testing.T) {
	theme, _ := GetTheme("violet")
	out := RenderStartupBanner(theme, ReleaseMetadata{Version: "2.0.0", ReleaseDate: "2026-07-15"}, RenderOptions{UseColor: true, Width: 140, Height: 40})
	for _, want := range []string{studioURL, repoURL, "\x1b["} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in color banner: %q", want, out)
		}
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
