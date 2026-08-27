package kudasai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"syscall"
)

var Version = "dev"

type Command struct {
	Run         string `json:"run"`
	Description string `json:"description,omitempty"`
}

func (c *Command) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		c.Run = s
		return nil
	}

	type commandAlias Command
	var alias commandAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*c = Command(alias)
	return nil
}

type Commands map[string]Command
type KudasaiConfig struct {
	Commands Commands `json:"commands"`
}

func KudasaiDefaultCommands() Commands {
	return Commands{
		"start": {Run: "echo \"sushi, kudasai\""},
	}
}

func mergeMaps(map1, map2 Commands) Commands {
	for key, value := range map2 {
		map1[key] = value
	}
	return map1
}

func GetCommands(filepath string) (Commands, bool, error) {
	configFile, err := os.Open(filepath)
	if err != nil {
		return KudasaiDefaultCommands(), false, nil
	}
	defer configFile.Close()

	content, err := io.ReadAll(configFile)
	if err != nil {
		return KudasaiDefaultCommands(), false, nil
	}
	var config KudasaiConfig
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, true, fmt.Errorf("failed to parse .kudasai.json: %s", err)
	}

	return mergeMaps(KudasaiDefaultCommands(), config.Commands), true, nil
}

func Check(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf(".kudasai.json not found in current directory. Run `kudasai init` to create one")
	}

	commands, _, err := GetCommands(filepath)
	if err != nil {
		return fmt.Errorf("%s. Fix the JSON syntax and run `kudasai check` again", err)
	}

	count := 0
	defaults := KudasaiDefaultCommands()
	for name := range commands {
		if _, isBuiltin := defaults[name]; !isBuiltin {
			count++
		}
	}

	fmt.Printf(".kudasai.json is valid (%d commands defined)\n", count)
	return nil
}

func ShellQuote(arg string) string {
	return "'" + strings.ReplaceAll(arg, "'", "'\\''") + "'"
}

func InterpolateArgs(command string, args []string) string {
	hasPlaceholder := regexp.MustCompile(`\$(\d+|@)`).MatchString(command)

	if !hasPlaceholder {
		if len(args) == 0 {
			return command
		}
		quoted := make([]string, len(args))
		for i, a := range args {
			quoted[i] = ShellQuote(a)
		}
		return command + " " + strings.Join(quoted, " ")
	}

	// Replace $@ with all args
	if strings.Contains(command, "$@") {
		quoted := make([]string, len(args))
		for i, a := range args {
			quoted[i] = ShellQuote(a)
		}
		command = strings.ReplaceAll(command, "$@", strings.Join(quoted, " "))
	}

	// Replace $N positional placeholders (highest first to avoid $1 matching in $10)
	re := regexp.MustCompile(`\$(\d+)`)
	command = re.ReplaceAllStringFunc(command, func(match string) string {
		// Extract the number after $
		numStr := match[1:]
		n := 0
		for _, c := range numStr {
			n = n*10 + int(c-'0')
		}
		idx := n - 1 // $1 -> args[0]
		if idx >= 0 && idx < len(args) {
			return ShellQuote(args[idx])
		}
		return ""
	})

	return command
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return 1
	}

	if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
		return 128 + int(status.Signal())
	}

	return exitErr.ExitCode()
}

func Execute(command string) error {
	cmd := Prepare(command)
	err := cmd.Run()

	return err
}

func Prepare(command string) *exec.Cmd {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir, _ = os.Getwd()

	return cmd
}

func Help(commands Commands, configFound bool) {
	fmt.Println("Usage: kudasai [command]")
	fmt.Println()

	if !configFound {
		fmt.Println("Warning: no .kudasai.json found in current directory")
		fmt.Println()
		fmt.Println("Built-in commands:")
	} else {
		fmt.Println("Available commands:")
	}

	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)

	maxLen := 0
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	for _, name := range names {
		desc := commands[name].Description
		if desc != "" {
			fmt.Printf("  %-*s  %s\n", maxLen, name, desc)
		} else {
			fmt.Printf("  %s\n", name)
		}
	}
}

func JSONOutput(commands Commands) error {
	defaults := KudasaiDefaultCommands()

	output := make(map[string]interface{})
	cmds := make(map[string]interface{})

	for name, cmd := range commands {
		if _, isBuiltin := defaults[name]; isBuiltin {
			cmds[name] = map[string]interface{}{"builtin": true}
		} else {
			entry := map[string]interface{}{"run": cmd.Run}
			if cmd.Description != "" {
				entry["description"] = cmd.Description
			}
			cmds[name] = entry
		}
	}

	output["commands"] = cmds

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

type ProjectDetection struct {
	Name     string
	File     string
	Commands Commands
}

func DetectProject() *ProjectDetection {
	detectors := []ProjectDetection{
		{"Go", "go.mod", Commands{
			"build": {Run: "go build ./...", Description: "Build the project"},
			"start": {Run: "go run .", Description: "Start the application"},
			"test":  {Run: "go test ./...", Description: "Run tests"},
			"lint":  {Run: "go vet ./...", Description: "Run linter"},
		}},
		{"Node.js", "package.json", Commands{
			"build": {Run: "npm run build", Description: "Build the project"},
			"start": {Run: "npm start", Description: "Start the application"},
			"test":  {Run: "npm test", Description: "Run tests"},
			"lint":  {Run: "npm run lint", Description: "Run linter"},
		}},
		{"Ruby", "Gemfile", Commands{
			"build": {Run: "bundle exec rake build", Description: "Build the project"},
			"start": {Run: "bundle exec ruby app.rb", Description: "Start the application"},
			"test":  {Run: "bundle exec rspec", Description: "Run tests"},
			"lint":  {Run: "bundle exec rubocop", Description: "Run linter"},
		}},
		{"PHP", "composer.json", Commands{
			"build": {Run: "composer install", Description: "Build the project"},
			"start": {Run: "php -S localhost:8000", Description: "Start the application"},
			"test":  {Run: "./vendor/bin/phpunit", Description: "Run tests"},
			"lint":  {Run: "./vendor/bin/phpcs", Description: "Run linter"},
		}},
		{"Python", "pyproject.toml", Commands{
			"build": {Run: "pip install -e .", Description: "Build the project"},
			"start": {Run: "python -m app", Description: "Start the application"},
			"test":  {Run: "pytest", Description: "Run tests"},
			"lint":  {Run: "ruff check .", Description: "Run linter"},
		}},
		{"Python", "requirements.txt", Commands{
			"build": {Run: "pip install -r requirements.txt", Description: "Build the project"},
			"start": {Run: "python -m app", Description: "Start the application"},
			"test":  {Run: "pytest", Description: "Run tests"},
			"lint":  {Run: "ruff check .", Description: "Run linter"},
		}},
	}

	for _, d := range detectors {
		if _, err := os.Stat(d.File); err == nil {
			return &d
		}
	}

	return nil
}

var Stdin io.Reader = os.Stdin

func Confirm(prompt string) bool {
	fmt.Print(prompt)
	var response string
	fmt.Fscanln(Stdin, &response)
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "" || response == "y" || response == "yes"
}

func Init(force bool) error {
	configPath := ".kudasai.json"

	if !force {
		if _, err := os.Stat(configPath); err == nil {
			return fmt.Errorf(".kudasai.json already exists (use --force to overwrite)")
		}
	}

	detection := DetectProject()

	var commands Commands
	if detection != nil {
		fmt.Printf("Detected %s project (%s)\n", detection.Name, detection.File)
		commands = detection.Commands
	} else {
		fmt.Println("No project type detected")
		commands = Commands{
			"build": {Run: "echo 'Replace with your build command'", Description: "Build the project"},
			"start": {Run: "echo 'Replace with your start command'", Description: "Start the application"},
			"test":  {Run: "echo 'Replace with your test command'", Description: "Run tests"},
			"lint":  {Run: "echo 'Replace with your lint command'", Description: "Run linter"},
		}
	}

	fmt.Println()
	fmt.Println("Commands:")

	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)

	maxLen := 0
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	for _, name := range names {
		cmd := commands[name]
		if cmd.Description != "" {
			fmt.Printf("  %-*s  %s (%s)\n", maxLen, name, cmd.Description, cmd.Run)
		} else {
			fmt.Printf("  %s: %s\n", name, cmd.Run)
		}
	}

	fmt.Println()
	if !Confirm("Create .kudasai.json? [Y/n] ") {
		fmt.Println("Aborted")
		return nil
	}

	config := KudasaiConfig{Commands: commands}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to generate .kudasai.json: %s", err)
	}

	if err := os.WriteFile(configPath, append(data, '\n'), 0644); err != nil {
		return fmt.Errorf("failed to create .kudasai.json: %s", err)
	}

	fmt.Println("Created .kudasai.json")
	return nil
}

func Run(args []string) error {
	if len(args) == 0 {
		return errors.New("No arguments provided")
	}

	// Flags (--) always run built-in behavior
	switch args[0] {
	case "--version":
		fmt.Printf("kudasai %s\n", Version)
		return nil
	case "--json":
		commands, _, err := GetCommands("./.kudasai.json")
		if err != nil {
			return err
		}
		return JSONOutput(commands)
	case "--help":
		commands, configFound, err := GetCommands("./.kudasai.json")
		if err != nil {
			return err
		}
		Help(commands, configFound)
		return nil
	case "--init":
		force := len(args) > 1 && args[1] == "--force"
		return Init(force)
	case "--check":
		return Check("./.kudasai.json")
	}

	commands, _, err := GetCommands("./.kudasai.json")
	if err != nil {
		return err
	}
	command, exists := commands[args[0]]
	if !exists {
		return fmt.Errorf("Unrecognized command: %s. Run `kudasai --help` to see available commands", args[0])
	}

	run := InterpolateArgs(command.Run, args[1:])
	return Execute(run)

}
