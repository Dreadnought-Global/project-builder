# Project Builder

Project Builder is a cross-platform command line utility written in Go that scaffolds standardized project directory structures for creative disciplines.

## Overview

The tool automates the creation of directory trees for four creative disciplines:
- Design
- Video & Motion
- Audio
- 3D & Animation

It provides a conditional client project overlay (`00_Client_Docs` at the root and `Client_Handoff` within the discipline's final export folders).

Upon the first run, the tool launches an interactive, terminal-based folder browser to select a root workbench directory. This path is stored in a configuration file (`config.yaml`) and used for all subsequent project initializations.

## Installation

Pre-built binaries for Windows, macOS, and Linux are available on the GitHub Releases page.

1. Download the executable appropriate for your operating system.
2. Place the binary in your system path or execute it directly.

### Linux Installation & Desktop Launcher

Running the binary directly from a terminal (`./project-builder`) is the standard, fully-supported way to run the tool on Linux.

For double-click launching via a file manager, a `project-builder.desktop` file is provided in the release. To install:
1. Copy `project-builder.desktop` to `~/.local/share/applications/`.
2. Edit the `Exec` line in the copied file to point to the absolute path of your compiled binary (e.g. `Exec=/home/user/bin/project-builder`).
3. Make the `.desktop` file executable:
   ```bash
   chmod +x ~/.local/share/applications/project-builder.desktop
   ```

### Configuration Path

The tool stores its configuration file at the following locations:
- **Linux/macOS**: `~/.config/project-builder/config.yaml`
- **Windows**: `%APPDATA%\project-builder\config.yaml`

## Usage

Run the compiled executable to start the interactive initialization flow:

```bash
./project-builder
```

### Command Line Flags

To reset and reconfigure the root workbench path, pass the reconfigure flag:

```bash
./project-builder --reconfigure
```

## Building from Source

### Prerequisites

- Go 1.26 or higher

### Build Instructions

1. Clone the repository:

```bash
git clone https://github.com/Dreadnought-Global/project-builder.git
cd project-builder
```

2. Compile the binary:

```bash
go build -o project-builder .
```

3. Run the test suite:

```bash
go test -v ./...
```

## License

Proprietary. Copyright (c) Dreadnought Studio. All rights reserved.
