package main

import "strings"

var projectBuilderBanner = []string{
	"██████╗ ██████╗  ██████╗      ██╗███████╗ ██████╗████████╗    ██████╗ ██╗   ██╗██╗██╗     ██████╗ ███████╗██████╗ ",
	"██╔══██╗██╔══██╗██╔═══██╗     ██║██╔════╝██╔════╝╚══██╔══╝    ██╔══██╗██║   ██║██║██║     ██╔══██╗██╔════╝██╔══██╗",
	"██████╔╝██████╔╝██║   ██║     ██║█████╗  ██║        ██║       ██████╔╝██║   ██║██║██║     ██║  ██║█████╗  ██████╔╝",
	"██╔═══╝ ██╔══██╗██║   ██║██   ██║██╔══╝  ██║        ██║       ██╔══██╗██║   ██║██║██║     ██║  ██║██╔══╝  ██╔══██╗",
	"██║     ██║  ██║╚██████╔╝╚█████╔╝███████╗╚██████╗   ██║       ██████╔╝╚██████╔╝██║███████╗██████╔╝███████╗██║  ██║",
	"╚═╝     ╚═╝  ╚═╝ ╚═════╝  ╚════╝ ╚══════╝ ╚═════╝   ╚═╝       ╚═════╝  ╚═════╝ ╚═╝╚══════╝╚═════╝ ╚══════╝╚═╝  ╚═╝",
}

func RenderStartupBanner(theme Theme, metadata ReleaseMetadata, opts RenderOptions) string {
	var b strings.Builder

	for i, line := range projectBuilderBanner {
		b.WriteString(gradientLine(line, theme.BannerGradientStops, i, len(projectBuilderBanner), opts))
		b.WriteString("\n")
	}
	b.WriteString(colorize(FormatMetadataLine(metadata, opts), theme.Muted, opts))
	b.WriteString("\n\n")
	return b.String()
}
