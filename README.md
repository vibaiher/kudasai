# kudasai

`kudasai` is a tool for generating aliases for the most frequently used commands in your repository and sharing them with your team.

ください (kudasai) it's a Japanese word. Fundamentally, is the polite form of the imperative form.

When you use ください with someone, you’re fundamentally telling them to do something. Is close in meaning to the English "please".

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/vibaiher/kudasai/main/install.sh | sh
```

Or, if you have Go installed:

```bash
go install github.com/vibaiher/kudasai
```

## Usage

### Default Commands

```bash
kudasai help  # Prints usage
kudasai start # Prints a message
```

### Custom Commands

Create a `.kudasai.json` file in your project root to define custom commands:

```json
{
  "commands": {
    "build": "go build -o kudasai main.go",
    "test": "go test -v ./tests/...",
    "deploy": "git push origin main"
  }
}
```

Then run your custom commands:

```bash
kudasai build
kudasai test
kudasai deploy
```

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

Custom commands take precedence over default commands if they share the same name.

## AI Agent Integration

AI coding agents (Claude Code, Cursor, Copilot, etc.) can discover your project commands automatically. Add this to your agent context file (`CLAUDE.md`, `.cursorrules`, etc.):

```markdown
Run `kudasai --json` to discover available commands in this repository. Use `kudasai <command>` to execute them.
```

You can also use `kudasai --json` programmatically to get structured output:

```bash
kudasai --json
```

## Contributing

Once you’ve cloned the repo and [set up the environment](DEVELOPMENT.md),
you can run the unit tests and acceptance suite, or submit a pull request.
