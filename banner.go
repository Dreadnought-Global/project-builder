package main

import "strings"

var pbBanner = []string{
	"██████╗ ██████╗ ",
	"██╔══██╗██╔══██╗",
	"██████╔╝██████╔╝",
	"██╔═══╝ ██╔══██╗",
	"██║     ██████╔╝",
	"╚═╝     ╚═════╝ ",
}

func RenderStartupBanner(theme Theme, metadata ReleaseMetadata, opts RenderOptions) string {
	if opts.Width >= 72 {
		return renderSideBySideBanner(theme, metadata, opts)
	}
	return renderStackedBanner(theme, metadata, opts)
}

func renderSideBySideBanner(theme Theme, metadata ReleaseMetadata, opts RenderOptions) string {
	info := bannerInfoLines(theme, metadata, opts)
	var b strings.Builder
	for i, iconLine := range pbBanner {
		b.WriteString(gradientLine(iconLine, theme.BannerGradientStops, i, len(pbBanner), opts))
		b.WriteString("  ")
		if i < len(info) {
			b.WriteString(info[i])
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.String()
}

func renderStackedBanner(theme Theme, metadata ReleaseMetadata, opts RenderOptions) string {
	var b strings.Builder
	for i, iconLine := range pbBanner {
		b.WriteString(gradientLine(iconLine, theme.BannerGradientStops, i, len(pbBanner), opts))
		b.WriteString("\n")
	}
	for _, line := range bannerInfoLines(theme, metadata, opts) {
		b.WriteString(line)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.String()
}

func bannerInfoLines(theme Theme, metadata ReleaseMetadata, opts RenderOptions) []string {
	return []string{
		colorize("Project Builder "+metadata.DisplayVersion(), theme.Primary, opts),
		colorize(appDescription, theme.Muted, opts),
		"",
		bannerInfoLine("Release", metadata.ReleaseDate, theme, opts),
		bannerInfoLine("Creator", linkIfColor(studioLabel, studioURL, opts), theme, opts),
		bannerInfoLine("Repo", linkIfColor(repoLabel, repoURL, opts), theme, opts),
	}
}

func bannerInfoLine(label, value string, theme Theme, opts RenderOptions) string {
	return colorize(padRight(label+":", 9), theme.Accent, opts) + colorize(value, theme.Primary, opts)
}

func linkIfColor(label, url string, opts RenderOptions) string {
	if !opts.UseColor {
		return label
	}
	return hyperlink(label, url)
}

func padRight(text string, width int) string {
	for visibleLen(text) < width {
		text += " "
	}
	return text
}
