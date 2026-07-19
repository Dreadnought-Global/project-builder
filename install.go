package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const installCommandName = "project-builder"

type installOptions struct {
	DryRun bool
	Force  bool
	Status bool
	System bool
}

type installPlan struct {
	Source      string
	TargetDir   string
	TargetPath  string
	DisplayPath string
	System      bool
}

var errInstallNeedsElevation = errors.New("install needs elevated access")

func handleInstallCommand(args []string, out io.Writer) int {
	opts, err := parseInstallOptions(args)
	if err != nil {
		fmt.Fprintf(out, "%v\n", err)
		fmt.Fprintln(out, "Usage: project-builder install [--dry-run|--force|status]")
		return 1
	}
	if opts.Status {
		return printInstallStatus(out, opts)
	}
	if err := installProjectBuilder(out, opts); err != nil {
		if errors.Is(err, errInstallNeedsElevation) && !opts.System {
			return runElevatedInstall(out, opts)
		}
		fmt.Fprintf(out, "Install failed: %v\n", err)
		return 1
	}
	return 0
}

func parseInstallOptions(args []string) (installOptions, error) {
	opts := installOptions{}
	for _, arg := range args {
		switch arg {
		case "--dry-run":
			opts.DryRun = true
		case "--force":
			opts.Force = true
		case "--system":
			opts.System = true
		case "status":
			opts.Status = true
		default:
			return opts, fmt.Errorf("unknown install option: %s", arg)
		}
	}
	return opts, nil
}

func installProjectBuilder(out io.Writer, opts installOptions) error {
	plan, err := newInstallPlan(opts.System)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "Source: %s\n", plan.Source)
	fmt.Fprintf(out, "Target: %s\n", plan.TargetPath)
	if opts.DryRun {
		fmt.Fprintln(out, "Dry run only. No files changed.")
		return nil
	}
	if err := copyExecutable(plan.Source, plan.TargetPath, opts.Force); err != nil {
		if errors.Is(err, os.ErrPermission) {
			return errInstallNeedsElevation
		}
		return err
	}
	if err := ensureInstallPath(plan.TargetDir, plan.DisplayPath, plan.System); err != nil {
		if errors.Is(err, os.ErrPermission) {
			return errInstallNeedsElevation
		}
		return err
	}
	fmt.Fprintln(out, "Project Builder installed successfully.")
	fmt.Fprintln(out, "Open a new terminal session, then run: project-builder")
	return nil
}

func newInstallPlan(system bool) (installPlan, error) {
	source, err := os.Executable()
	if err != nil {
		return installPlan{}, err
	}
	source, err = filepath.EvalSymlinks(source)
	if err != nil {
		// ponytail: best-effort path cleanup only; hard-fail if real installs need stricter source validation.
		source, _ = os.Executable()
	}
	dir, display, err := installDir(system)
	if err != nil {
		return installPlan{}, err
	}
	return installPlan{Source: source, TargetDir: dir, TargetPath: filepath.Join(dir, installBinaryName()), DisplayPath: display, System: system}, nil
}

func installDir(system bool) (string, string, error) {
	if system {
		if runtime.GOOS == "windows" {
			base := os.Getenv("ProgramFiles")
			if base == "" {
				base = `C:\Program Files`
			}
			dir := filepath.Join(base, "Project Builder")
			return dir, dir, nil
		}
		return "/usr/local/bin", "/usr/local/bin", nil
	}
	switch runtime.GOOS {
	case "windows":
		base := os.Getenv("LOCALAPPDATA")
		if base == "" {
			if profile := os.Getenv("USERPROFILE"); profile != "" {
				base = filepath.Join(profile, "AppData", "Local")
			}
		}
		if base == "" {
			if drive, path := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); drive != "" && path != "" {
				base = filepath.Join(drive+path, "AppData", "Local")
			}
		}
		if base == "" {
			if home, err := os.UserHomeDir(); err == nil {
				base = filepath.Join(home, "AppData", "Local")
			}
		}
		if base == "" {
			base = os.TempDir()
		}
		dir := filepath.Join(base, "ProjectBuilder", "bin")
		return dir, dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "bin"), "$HOME/bin", nil
	default:
		return filepath.Join(home, ".local", "bin"), "$HOME/.local/bin", nil
	}
}

func installBinaryName() string {
	if runtime.GOOS == "windows" {
		return installCommandName + ".exe"
	}
	return installCommandName
}

func copyExecutable(source, target string, force bool) error {
	if samePath(source, target) {
		return nil
	}
	if _, err := os.Stat(target); err == nil && !force {
		return fmt.Errorf("target already exists: %s (use --force to replace it)", target)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func samePath(a, b string) bool {
	absA, errA := filepath.Abs(a)
	absB, errB := filepath.Abs(b)
	if errA != nil || errB != nil {
		return false
	}
	return strings.EqualFold(filepath.Clean(absA), filepath.Clean(absB))
}

func printInstallStatus(out io.Writer, opts installOptions) int {
	plan, err := newInstallPlan(opts.System)
	if err != nil {
		fmt.Fprintf(out, "Install status failed: %v\n", err)
		return 1
	}
	_, existsErr := os.Stat(plan.TargetPath)
	fmt.Fprintf(out, "Install path: %s\n", plan.TargetPath)
	if existsErr == nil {
		fmt.Fprintln(out, "Binary: installed")
	} else if errors.Is(existsErr, os.ErrNotExist) {
		fmt.Fprintln(out, "Binary: not installed")
	} else {
		fmt.Fprintf(out, "Binary: unknown (%v)\n", existsErr)
	}
	if pathContainsDir(os.Getenv("PATH"), plan.TargetDir) {
		fmt.Fprintln(out, "PATH: active in this terminal")
	} else {
		fmt.Fprintln(out, "PATH: not active in this terminal")
	}
	return 0
}

func pathContainsDir(pathValue, dir string) bool {
	for _, part := range filepath.SplitList(pathValue) {
		if samePath(part, dir) {
			return true
		}
	}
	return false
}
