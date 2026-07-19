//go:build !windows

package main

import "golang.org/x/sys/unix"

func detectTerminalSize(fd uintptr) (int, int, bool) {
	ws, err := unix.IoctlGetWinsize(int(fd), unix.TIOCGWINSZ)
	if err != nil || ws == nil || ws.Col == 0 {
		return 0, 0, false
	}
	return int(ws.Col), int(ws.Row), true
}
