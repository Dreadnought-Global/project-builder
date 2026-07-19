//go:build !windows

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func runElevatedInstall(out io.Writer, opts installOptions) int {
	fmt.Fprintln(out, "Per-user install was blocked by permissions.")
	fmt.Fprint(out, "Request administrator access now? [y/N]: ")
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if strings.TrimSpace(strings.ToLower(line)) != "y" {
		fmt.Fprintln(out, "Admin install cancelled.")
		return 1
	}
	fmt.Fprintln(out, "Requesting administrator access with sudo...")
	source, err := os.Executable()
	if err != nil {
		fmt.Fprintf(out, "Admin install failed: %v\n", err)
		return 1
	}
	args := []string{source, "install", "--system", "--force"}
	if opts.DryRun {
		args = append(args, "--dry-run")
	}
	cmd := exec.Command("sudo", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(out, "Admin install failed: %v\n", err)
		return 1
	}
	fmt.Fprintln(out, "Admin install completed. Open a new terminal session, then run: project-builder")
	return 0
}
