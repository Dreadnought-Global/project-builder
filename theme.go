package main

import (
	"fmt"
	"sort"
	"strings"
)

const defaultThemeName = "violet"

type RGB struct {
	R int
	G int
	B int
}

type Theme struct {
	Name                string
	BannerGradientStops []RGB
	Primary             RGB
	Accent              RGB
	Muted               RGB
	Success             RGB
	Warning             RGB
	Error               RGB
	SelectionIndicator  RGB
	Shadow              RGB
}

var themes = map[string]Theme{
	"violet": {
		Name:                "violet",
		BannerGradientStops: []RGB{{91, 39, 160}, {191, 64, 255}, {255, 185, 245}},
		Primary:             RGB{245, 230, 255}, Accent: RGB{205, 90, 255}, Muted: RGB{168, 145, 184},
		Success: RGB{86, 230, 166}, Warning: RGB{255, 199, 95}, Error: RGB{255, 106, 132},
		SelectionIndicator: RGB{218, 96, 255}, Shadow: RGB{58, 20, 86},
	},
	"cyan": {
		Name:                "cyan",
		BannerGradientStops: []RGB{{16, 46, 105}, {0, 172, 210}, {158, 242, 255}},
		Primary:             RGB{225, 248, 255}, Accent: RGB{34, 211, 238}, Muted: RGB{136, 165, 178},
		Success: RGB{94, 234, 212}, Warning: RGB{250, 204, 21}, Error: RGB{251, 113, 133},
		SelectionIndicator: RGB{34, 211, 238}, Shadow: RGB{7, 32, 56},
	},
	"emerald": {
		Name:                "emerald",
		BannerGradientStops: []RGB{{10, 80, 55}, {16, 185, 129}, {187, 247, 208}},
		Primary:             RGB{230, 255, 242}, Accent: RGB{52, 211, 153}, Muted: RGB{139, 170, 151},
		Success: RGB{74, 222, 128}, Warning: RGB{251, 191, 36}, Error: RGB{248, 113, 113},
		SelectionIndicator: RGB{52, 211, 153}, Shadow: RGB{5, 46, 34},
	},
	"amber": {
		Name:                "amber",
		BannerGradientStops: []RGB{{146, 64, 14}, {245, 158, 11}, {254, 240, 138}},
		Primary:             RGB{255, 248, 232}, Accent: RGB{251, 191, 36}, Muted: RGB{176, 148, 104},
		Success: RGB{132, 204, 22}, Warning: RGB{253, 186, 116}, Error: RGB{248, 113, 113},
		SelectionIndicator: RGB{251, 191, 36}, Shadow: RGB{70, 31, 8},
	},
	"mono": {
		Name:                "mono",
		BannerGradientStops: []RGB{{160, 160, 160}, {214, 214, 214}, {245, 245, 245}},
		Primary:             RGB{235, 235, 235}, Accent: RGB{210, 210, 210}, Muted: RGB{150, 150, 150},
		Success: RGB{220, 220, 220}, Warning: RGB{200, 200, 200}, Error: RGB{180, 180, 180},
		SelectionIndicator: RGB{235, 235, 235}, Shadow: RGB{82, 82, 82},
	},
}

func ThemeNames() []string {
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func ValidThemeChoices() string {
	return strings.Join(ThemeNames(), ", ")
}

func GetTheme(name string) (Theme, error) {
	cleaned := strings.ToLower(strings.TrimSpace(name))
	if cleaned == "" {
		cleaned = defaultThemeName
	}
	theme, ok := themes[cleaned]
	if !ok {
		return Theme{}, fmt.Errorf("invalid theme %q; valid themes: %s", name, ValidThemeChoices())
	}
	return theme, nil
}

func ActiveTheme(cfg Config) Theme {
	theme, err := GetTheme(cfg.Theme)
	if err != nil {
		theme, _ = GetTheme(defaultThemeName)
	}
	return theme
}

func normalizeThemeName(name string) string {
	cleaned := strings.ToLower(strings.TrimSpace(name))
	if _, err := GetTheme(cleaned); err != nil {
		return defaultThemeName
	}
	return cleaned
}
