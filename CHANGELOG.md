# Changelog

All notable changes to Project Builder will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) as detailed in [AGENTS.md](AGENTS.md).

## [2.3.0] - 2026-07-21

### Added
- Added confirmation before cancelling the terminal folder browser with `Ctrl+C` or `q`, helping prevent accidental loss of the current selection.
- Added folder browser controls to the Help screen.

### Changed
- Redesigned terminal folder browser with a clean themed header, full-screen redraws, and compact folder selection view.
- Shortened deep folder paths with `~` and `...` so narrow terminal windows stay readable.
- Cleared terminal between major project-creation steps and now show the selected destination before later project questions.
- Linux and macOS configuration now honor `XDG_CONFIG_HOME`, allowing isolated first-run testing without reading existing settings.

## [2.2.0] - 2026-07-20

### Added
- Added a home dashboard on interactive launch with Create Project, Help, Settings, and Exit options.
- Added a table-style help screen and `project-builder help` command.
- Added settings screen entries for theme/profile selection, default workbench changes, config path, and PATH install guidance.

### Changed
- Interactive launch now clears the terminal before rendering the banner and dashboard for a clean focused start.
- Blank Enter on the dashboard selects the default Create Project flow.
- Cancelled folder selections and declined project creation now return to the previous menu/dashboard instead of immediately closing the app.

## [2.1.0] - 2026-07-19

### Added
- Added `project-builder install` to copy the app into an OS-specific install directory and add it to PATH automatically.
- Added installer options for `--dry-run`, `--force`, and `install status`.
- Added administrator fallback when per-user install is blocked by permissions.

### Fixed
- Corrected the repository link shown in the startup banner to `https://github.com/Dreadnought-Global/project-builder`.

## [2.0.0] - 2026-07-15

### Added
- Native OS folder picker dialogs for discipline destination selection (Windows FolderBrowserDialog, macOS `osascript`, Linux `zenity`/`kdialog` with terminal TUI fallback).
- Per-discipline destination selection flow with options to use the Global Default Workbench, Native OS Picker, or Terminal Browser.
- Option to save or decline a default workbench path on a per-discipline basis.

### Changed
- **BREAKING:** Changed configuration schema in `config.yaml` from a single `workbench_path` to a `default_workbench` and `discipline_paths` map. Old configs are automatically migrated on next run.
- Redesigned the CLI prompt flow to request project name and discipline before asking for the destination and client flags.

### Fixed
- Fixed bug where the "Client Project" flag did not correctly resolve the root target folder to `00_Client_Projects` vs `01_Passion_Projects`.
- Added automated tests to ensure client flag structures and folder resolution logic operate correctly.

## [1.1.3] - 2026-07-15
### Fixed
- Updated release workflow to automatically populate GitHub Release notes from the matching `CHANGELOG.md` entry for the tagged version.
- Workflow now fails loudly if no matching `CHANGELOG.md` entry is found for the pushed tag, enforcing the changelog-before-tagging rule.
- Added GitHub Release Notes rule to `AGENTS.md`.

## [1.1.2] - 2026-07-15

### Fixed
- Fixed folder opening behavior where the terminal window closed prematurely on option 1. The terminal now waits for the user to press Enter before exiting.

## [1.1.1] - 2026-07-15

### Fixed
- Fixed Linux desktop double-click execution behavior by introducing `project-builder.desktop` to auto-open in a terminal.
- Updated release pipeline to include the Linux `.desktop` launcher asset.
- Updated `README.md` to document Linux installation and `.desktop` configuration.
- Added Local Development Build Workflow standing rule in `AGENTS.md`.

## [1.1.0] - 2026-07-15

### Added
- Interactive terminal folder browser (TUI) powered by Bubble Tea to select the root workbench directory on first run.
- Configuration persistence using `config.yaml` stored in OS-standard directories.
- Command-line flag `--reconfigure` to reset and change the saved root path.
- Permanent project versioning rules in `AGENTS.md` and this `CHANGELOG.md` file.

## [1.0.0] - 2026-07-15

### Added
- Initial standalone release of Project Builder CLI.
- Directory structure scaffolding templates for four disciplines: Design, Video & Motion, Audio, and 3D & Animation.
- Client project overlay (adding `00_Client_Docs` at root and `Client_Handoff` in final output directories).
- Interactive CLI prompt flow with input sanitization (replacing spaces with underscores and stripping illegal characters).
- Project name collision handling (options to rename, append suffix, or abort).
- Cross-platform support to open the newly generated project folder in the OS native file explorer.
