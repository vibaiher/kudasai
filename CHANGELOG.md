# Changelog

## 0.2.1

### Fixed
- Exit codes for commands terminated by a signal now follow the shell convention (`128 + N`): `130` for SIGINT, `137` for SIGKILL. They previously collapsed to `255`.

### Added
- Documentation site under `docs/`, published with GitHub Pages, including an `llms.txt` reference for AI agents.

## 0.2.0

### Added
- Command descriptions: commands can now use an object form with `run` and `description` fields. Descriptions appear in `--help` and `--json` output.
- `kudasai --init` scaffolds `.kudasai.json` with project detection (Go, Node.js, Ruby, PHP, Python).
- `kudasai --check` validates `.kudasai.json` syntax.
- `kudasai --version` shows the current version.
- `kudasai --json` outputs commands as structured JSON for AI agent integration.
- `AGENTS.md` symlink to `CLAUDE.md`.

### Changed
- Built-in functionality is now accessed exclusively via flags (`--help`, `--init`, etc.), freeing all other names for custom commands.
- Custom commands can override the built-in `start` command.
- Exit codes from executed commands are now propagated correctly.

## 0.1.0

### Added
- Support for passing extra arguments to commands.
- Placeholder interpolation: `$1`, `$2`, ..., `$N` for positional arguments and `$@` for all arguments.
- When no placeholders are present, extra arguments are appended to the end of the command (backwards compatible).

## 0.0.1

### Added
- CLI tool that reads command definitions from `.kudasai.json` and executes them via `sh -c`.
- Built-in `help` and `start` default commands.
- Support for custom commands that override defaults.
- Interactive command support with stdin/stdout/stderr piped to terminal.
- `curl | sh` installation support.
