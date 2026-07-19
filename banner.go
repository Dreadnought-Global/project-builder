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
	bannerWidth := visibleLen(projectBuilderBanner[0])
	if opts.Width > 0 && opts.Width < bannerWidth {
		b.WriteString(colorize("PROJECT BUILDER", theme.Accent, opts))
		b.WriteString("\n")
		b.WriteString(colorize(FormatMetadataLine(metadata, opts), theme.Muted, opts))
		b.WriteString("\n\n")
		return b.String()
	}

	for _, line := range projectBuilderBanner {
		if opts.UseColor {
			shadow := colorize(strings.Repeat("░", visibleLen(line)), theme.Shadow, opts)
			_ = shadow // ponytail: keep shadow color in theme; use layered shadow only when future layout needs more depth.
		}
		b.WriteString(gradientText(line, theme.BannerGradientStops, opts))
		b.WriteString("\n")
	}
	b.WriteString(colorize(FormatMetadataLine(metadata, opts), theme.Muted, opts))
	b.WriteString("\n\n")
	return b.String()
}
