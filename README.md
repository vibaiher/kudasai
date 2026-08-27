# kudasai

`kudasai` is a tool for generating aliases for the most frequently used commands in your repository and sharing them with your team.

 ください (kudasai) it's a Japanese word. Fundamentally, is the polite form of the imperative form.

When you use ください with someone, you're fundamentally telling them to do something. Is close in meaning to the English "please".

Full documentation: **https://vibaiher.github.io/kudasai/**

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/vibaiher/kudasai/main/install.sh | sh
```

Or, if you have Go installed:

```bash
go install github.com/vibaiher/kudasai
```

## Getting Started

Scaffold a `.kudasai.json` in your project:

```bash
kudasai --init
```

This detects your project type (Go, Node.js, Ruby, PHP, Python) and generates appropriate commands. You can also create the file manually:

```json
{
  "commands": {
    "build": {
      "run": "go build -o myapp main.go",
      "description": "Build the binary"
    },
    "test": {
      "run": "go test -v ./tests/...",
      "description": "Run unit tests"
    },
    "deploy": "git push origin main"
  }
}
```

Then run your commands:

```bash
kudasai build
kudasai test
kudasai deploy
```

## Usage

### Command Format

Commands can be defined in two ways:

**String form** (just the shell command):

```json
{
  "commands": {
    "deploy": "git push origin main"
  }
}
```

**Object form** (with description):

```json
{
  "commands": {
    "deploy": {
      "run": "git push origin main",
      "description": "Deploy to production"
    }
  }
}
```

Both forms can be mixed in the same file. Descriptions appear in `--help` and `--json` output.

### Passing Arguments

You can pass extra arguments to commands. By default, arguments are appended at the end:

```bash
kudasai test -run TestFoo
# Executes: go test -v ./tests/... -run TestFoo
```

Use placeholders to control where arguments are inserted:

- `$1`, `$2`, ..., `$N` — positional arguments
- `$@` — all arguments

```json
{
  "commands": {
    "test": "go test $@ ./tests/...",
    "greet": "echo hello $1, welcome to $2"
  }
}
```

```bash
kudasai test -v -run Foo
# Executes: go test -v -run Foo ./tests/...

kudasai greet world earth
# Executes: echo hello world, welcome to earth
```

### Interactive Commands

Commands support stdin, so you can use interactive tools:

```json
{
  "commands": {
    "commit": "git add . && git commit",
    "search": "grep -r --color"
  }
}
```

### Built-in Flags

All built-in functionality is accessed via flags:

| Flag | Description |
|------|-------------|
| `--help` | List available commands |
| `--version` | Show kudasai version |
| `--init` | Scaffold a `.kudasai.json` with project detection |
| `--check` | Validate your `.kudasai.json` |
| `--json` | Output commands as structured JSON |

All other names are available for your custom commands.

## AI Agent Integration

AI coding agents (Claude Code, Cursor, Copilot, etc.) can discover your project commands automatically. Add this to your agent context file (`CLAUDE.md`, `.cursorrules`, etc.):

```markdown
Run `kudasai --json` to discover available commands in this repository. Use `kudasai <command>` to execute them.
```

The `--json` flag outputs structured data that includes descriptions when available:

```bash
kudasai --json
```

```json
{
  "commands": {
    "build": {
      "run": "go build -o kudasai main.go",
      "description": "Build the binary"
    },
    "test": {
      "run": "go test -v ./tests/...",
      "description": "Run unit tests"
    }
  }
}
```

## Contributing

Once you've cloned the repo and [set up the environment](DEVELOPMENT.md),
you can run the unit tests and acceptance suite, or submit a pull request.
