# Development

## Prerequisites

- **Go** 1.23 — [install](https://golang.org/doc/install)
- **clitest** — [install](https://github.com/aureliojargas/clitest)

## Setup

```bash
git clone https://github.com/vibaiher/kudasai.git
cd kudasai
go build -o kudasai main.go
```

## Running Tests

### Unit tests

Located in `tests/kudasai_test.go`. They test the `Run()` function and helpers directly:

```bash
go test -v ./tests/...
```

### Acceptance tests

Located in `examples/*.txt`. Each file is a [clitest](https://github.com/aureliojargas/clitest) script that runs kudasai end-to-end and checks its output:

```bash
bash scripts/acceptance.sh
```

The script builds a fresh binary, creates a temp directory for each test file, and runs it in isolation.

### Test coverage by feature

| Feature | Unit test | Acceptance test |
|---|---|---|
| Run custom command | `TestRun_Start` | `custom.txt` |
| Argument passing (`$1`, `$@`, append) | `TestInterpolateArgs_*` | `args.txt` |
| `--help` | `TestRun_Help` | `help.txt` |
| `--help` with descriptions | — | `descriptions.txt` |
| `--version` | `TestRun_Version` | `version.txt` |
| `--json` (with descriptions) | `TestRun_JSON` | `json.txt` |
| `--init` (project detection, confirm) | `TestRun_Init*` | `init.txt` |
| `--check` (validation) | `TestRun_Check*` | `check.txt` |
| Unrecognized command error | `TestRun_InvalidCommand` | `invalid.txt` |
| Exit code propagation | `TestExecute_PropagatesExitCode` | `exit-code.txt` |
| String and object command formats | — | `descriptions.txt` |
| Shell quoting | `TestShellQuote_*` | — |

### Adding a new acceptance test

Create a file in `examples/` following the clitest format:

```
$ echo '{"commands":{"hello": "echo world"}}' > .kudasai.json
$ kudasai hello
world
```

Lines starting with `$` are commands; the lines that follow are expected output. Each test file runs in its own temporary directory.

## Building

```bash
go build -o kudasai main.go
```

This generates a binary in the project root. Alternatively, install to your `$GOBIN`:

```bash
go install
```

## Project Structure

```
main.go                  Entry point, delegates to pkg/kudasai
pkg/kudasai/kudasai.go   All core logic (commands, config, help, init, etc.)
tests/kudasai_test.go    Unit tests
examples/*.txt           Acceptance tests (clitest format)
scripts/acceptance.sh    Acceptance test runner
.kudasai.json            This project's own commands (dogfooding)
```
