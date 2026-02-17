# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Kudasai is a Go-based CLI tool for generating aliases for frequently used commands in a repository. It reads command definitions from a `.kudasai.json` file in the current directory and executes shell commands via `sh -c`.

## Development Commands

Run `kudasai --json` to discover available commands in this repository. Use `kudasai <command>` to execute them.

### Prerequisites
- Go 1.23
- clitest (for acceptance tests): https://github.com/aureliojargas/clitest

## Architecture

### Core Components

**main.go**: Entry point that captures CLI args and delegates to the kudasai package.

**pkg/kudasai/kudasai.go**: Contains all core logic:
- `Run(args)` - Main entry point that loads commands and executes the requested command
- `GetCommands(filepath)` - Loads `.kudasai.json` config and merges with default commands
- `Execute(command)` / `Prepare(command)` - Executes shell commands via `sh -c` with stdin/stdout/stderr piped
- `KudasaiDefaultCommands()` - Returns built-in commands (`help`, `start`)

### Command Resolution Flow

1. User runs `kudasai <command>`
2. `Run()` reads `.kudasai.json` from current directory
3. Custom commands are merged with default commands (custom takes precedence)
4. Command string is executed via `sh -c` in current working directory
5. Stdin/stdout/stderr are piped directly to terminal (supports interactive commands)

### Configuration Format

`.kudasai.json` structure:
```json
{
  "commands": {
    "command-name": "shell command to execute"
  }
}
```

## Testing Strategy

- **Unit tests** (tests/kudasai_test.go): Test the `Run()` function with various inputs
- **Acceptance tests** (examples/*.txt): Use clitest to verify CLI behavior end-to-end

The project's own `.kudasai.json` defines the development commands. Run `kudasai --json` to see them.
