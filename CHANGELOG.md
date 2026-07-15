# Changelog

All notable changes to Project Builder will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) as detailed in [AGENTS.md](AGENTS.md).

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
