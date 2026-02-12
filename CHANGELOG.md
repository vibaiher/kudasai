# Changelog

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
