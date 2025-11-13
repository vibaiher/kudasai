# kudasai

`kudasai` is a tool for generating aliases for the most frequently used commands in your repository and sharing them with your team.

ください (kudasai) it's a Japanese word. Fundamentally, is the polite form of the imperative form.

When you use ください with someone, you’re fundamentally telling them to do something. Is close in meaning to the English "please".

## Installation

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

## Contributing

Once you’ve cloned the repo and [set up the environment](DEVELOPMENT.md),
you can run the unit tests and acceptance suite, or submit a pull request.
