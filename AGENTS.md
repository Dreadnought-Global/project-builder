# Project Rules — Versioning & Release Guidelines

This document outlines the rules governing version numbers and releases for Project Builder. All future contributors and agent sessions must adhere to these rules.

## Semantic Versioning (SemVer)

Versions must follow the standard `MAJOR.MINOR.PATCH` format. Version bumps are determined by the impact on the user rather than the volume of internal code changes:

- **PATCH (e.g. v1.0.0 → v1.0.1)**:
  - Bug fixes only.
  - No new functionality.
  - No user-visible behavior changes.
  - Ensures something that was broken now works correctly.
  
- **MINOR (e.g. v1.0.0 → v1.1.0)**:
  - New, fully backward-compatible functionality.
  - Existing projects, configs, and command line workflows continue to work exactly as before without modification.
  - Examples: Adding a new discipline template, adding an interactive directory browser, adding custom icon bindings.

- **MAJOR (e.g. v1.0.0 → v2.0.0)**:
  - Breaking changes.
  - Changes to the output folder structure.
  - Changes to the configuration file format that invalidate existing configurations.
  - CLI behavior changes that require user action or break existing scripting integrations.

## Changelog Requirement

A `CHANGELOG.md` file must be maintained at the repository root and updated concurrently with every version bump.
- Each entry must specify the version number, release date, and a clear, plain-language description of what changed and why.
- Do not dump raw commit messages. Entries must be reader-friendly and summarize user-facing impact.

## Local Development Build Workflow

Once any task, fix, or feature is completed and verified via automated tests, always run the following command to produce an immediately runnable local binary:

```bash
go build -o project-builder && ./project-builder
```

This compiles the current source into a binary named `project-builder` (overwriting any previous local build) and immediately runs it. This ensures that the binary is left in an immediately testable/runnable state for direct manual validation.

The release pipeline (version tagging, cross-platform builds, GitHub Releases) is reserved for actual versioned releases meant for distribution. Do not tag/release for local iteration or testing.

