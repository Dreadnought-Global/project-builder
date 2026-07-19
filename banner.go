package main

import "strings"

var fullProjectBuilderBanner = []string{
	"██████╗ ██████╗  ██████╗      ██╗███████╗ ██████╗████████╗    ██████╗ ██╗   ██╗██╗██╗     ██████╗ ███████╗██████╗ ",
	"██╔══██╗██╔══██╗██╔═══██╗     ██║██╔════╝██╔════╝╚══██╔══╝    ██╔══██╗██║   ██║██║██║     ██╔══██╗██╔════╝██╔══██╗",
	"██████╔╝██████╔╝██║   ██║     ██║█████╗  ██║        ██║       ██████╔╝██║   ██║██║██║     ██║  ██║█████╗  ██████╔╝",
	"██╔═══╝ ██╔══██╗██║   ██║██   ██║██╔══╝  ██║        ██║       ██╔══██╗██║   ██║██║██║     ██║  ██║██╔══╝  ██╔══██╗",
	"██║     ██║  ██║╚██████╔╝╚█████╔╝███████╗╚██████╗   ██║       ██████╔╝╚██████╔╝██║███████╗██████╔╝███████╗██║  ██║",
	"╚═╝     ╚═╝  ╚═╝ ╚═════╝  ╚════╝ ╚══════╝ ╚═════╝   ╚═╝       ╚═════╝  ╚═════╝ ╚═╝╚══════╝╚═════╝ ╚══════╝╚═╝  ╚═╝",
}

var mediumProjectBuilderBanner = []string{
	"█▀█ █▀█ █▀█ ░█ █▀▀ █▀▀ ▀█▀   █▀▄ █ █ █ █ █   █▀▄ █▀▀ █▀█",
	"█▀▀ █▀▄ █▄█ ░█ █▀▀ █   ░█░   █▀▄ █ █ █ █ █   █▄▀ █▀▀ █▀▄",
	"▀░░ ▀░▀ ▀░▀ █▄ █▄▄ █▄▄ ░▀░   ▀▀░ ▀▀▀ ▀ ▀ ▀▀▀ ▀▀░ ▀▀▀ ▀░▀",
}

var compactProjectBuilderBanner = []string{"PROJECT BUILDER"}

func RenderStartupBanner(theme Theme, metadata ReleaseMetadata, opts RenderOptions) string {
	banner := bannerForWidth(opts.Width)
	var b strings.Builder

	for i, line := range banner {
		b.WriteString(gradientLine(line, theme.BannerGradientStops, i, len(banner), opts))
		b.WriteString("\n")
	}
	b.WriteString(colorize(FormatMetadataLine(metadata, opts), theme.Muted, opts))
	b.WriteString("\n\n")
	return b.String()
}

func bannerForWidth(width int) []string {
	if width <= 0 || width >= bannerWidth(fullProjectBuilderBanner) {
		return fullProjectBuilderBanner
	}
	if width >= bannerWidth(mediumProjectBuilderBanner) {
		return mediumProjectBuilderBanner
	}
	return compactProjectBuilderBanner
}

func bannerWidth(lines []string) int {
	width := 0
	for _, line := range lines {
		if lineWidth := visibleLen(line); lineWidth > width {
			width = lineWidth
		}
	}
	return width
}
