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

func stripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

func visibleLen(s string) int {
	return utf8.RuneCountInString(stripANSI(s))
}

func gradientText(text string, stops []RGB, opts RenderOptions) string {
	if !opts.UseColor || len(stops) == 0 {
		return text
	}
	visible := []rune(text)
	count := 0
	for _, r := range visible {
		if r != ' ' {
			count++
		}
	}
	if count == 0 {
		return text
	}
	var b strings.Builder
	idx := 0
	for _, r := range visible {
		if r == ' ' {
			b.WriteRune(r)
			continue
		}
		ratio := 0.0
		if count > 1 {
			ratio = float64(idx) / float64(count-1)
		}
		b.WriteString(colorize(string(r), interpolateStops(stops, ratio), opts))
		idx++
	}
	return b.String()
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
