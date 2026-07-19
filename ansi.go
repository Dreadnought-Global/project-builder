package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]|\x1b\]8;;.*?\x07|\x1b\]8;;\x07`)

type RenderOptions struct {
	UseColor bool
	Width    int
}

var activeRenderOptions = RenderOptions{UseColor: false, Width: 140}
var activeTheme, _ = GetTheme(defaultThemeName)

func SetActiveStyle(theme Theme, opts RenderOptions) {
	activeTheme = theme
	activeRenderOptions = opts
}

func DetectRenderOptions(noColor bool) RenderOptions {
	useColor := !noColor && os.Getenv("NO_COLOR") == "" && os.Getenv("PROJECT_BUILDER_NO_COLOR") == ""
	width := 140
	if columns := strings.TrimSpace(os.Getenv("COLUMNS")); columns != "" {
		if parsed, err := strconv.Atoi(columns); err == nil && parsed > 0 {
			width = parsed
		}
	}
	return RenderOptions{UseColor: useColor, Width: width}
}

func colorize(text string, color RGB, opts RenderOptions) string {
	if !opts.UseColor || text == "" {
		return text
	}
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", clamp(color.R), clamp(color.G), clamp(color.B), text)
}

func hyperlink(label, url string) string {
	return fmt.Sprintf("\x1b]8;;%s\x07%s\x1b]8;;\x07", url, label)
}

func primaryText(text string) string { return colorize(text, activeTheme.Primary, activeRenderOptions) }
func accentText(text string) string  { return colorize(text, activeTheme.Accent, activeRenderOptions) }
func mutedText(text string) string   { return colorize(text, activeTheme.Muted, activeRenderOptions) }
func successText(text string) string { return colorize(text, activeTheme.Success, activeRenderOptions) }
func warningText(text string) string { return colorize(text, activeTheme.Warning, activeRenderOptions) }
func errorText(text string) string   { return colorize(text, activeTheme.Error, activeRenderOptions) }
func promptText(text string) string {
	return colorize(text, activeTheme.SelectionIndicator, activeRenderOptions)
}

func stripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

func visibleLen(s string) int {
	return utf8.RuneCountInString(stripANSI(s))
}

func gradientLine(text string, stops []RGB, row, rows int, opts RenderOptions) string {
	if !opts.UseColor || len(stops) == 0 {
		return text
	}
	ratio := 0.0
	if rows > 1 {
		ratio = float64(row) / float64(rows-1)
	}
	return colorize(text, interpolateStops(stops, ratio), opts)
}

func interpolateStops(stops []RGB, ratio float64) RGB {
	if len(stops) == 1 {
		return stops[0]
	}
	if ratio <= 0 {
		return stops[0]
	}
	if ratio >= 1 {
		return stops[len(stops)-1]
	}
	scaled := ratio * float64(len(stops)-1)
	idx := int(math.Floor(scaled))
	local := scaled - float64(idx)
	return interpolateRGB(stops[idx], stops[idx+1], local)
}

func interpolateRGB(a, b RGB, ratio float64) RGB {
	return RGB{
		R: int(math.Round(float64(a.R) + (float64(b.R)-float64(a.R))*ratio)),
		G: int(math.Round(float64(a.G) + (float64(b.G)-float64(a.G))*ratio)),
		B: int(math.Round(float64(a.B) + (float64(b.B)-float64(a.B))*ratio)),
	}
}

func clamp(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}
