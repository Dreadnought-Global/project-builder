//go:build windows

package main

import "golang.org/x/sys/windows"

func detectTerminalSize(fd uintptr) (int, int, bool) {
	handle := windows.Handle(fd)
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(handle, &info); err != nil {
		return 0, 0, false
	}
	width := int(info.Window.Right - info.Window.Left + 1)
	height := int(info.Window.Bottom - info.Window.Top + 1)
	if width <= 0 || height <= 0 {
		return 0, 0, false
	}
	return width, height, true
}
