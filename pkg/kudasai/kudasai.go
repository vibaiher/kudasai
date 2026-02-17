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
)

var Version = "dev"

type Commands map[string]string
type KudasaiConfig struct {
	Commands Commands `json:"commands"`
}

func KudasaiDefaultCommands() Commands {
	return Commands{
		"start": "echo \"sushi, kudasai\"",
		"help":  "",
		"init":  "",
	}
}

func mergeMaps(map1, map2 map[string]string) map[string]string {
	for key, value := range map2 {
		map1[key] = value
	}
	return map1
}

func GetCommands(filepath string) (Commands, bool) {
	configFile, err := os.Open(filepath)
	if err != nil {
		return KudasaiDefaultCommands(), false
	}
	defer configFile.Close()

	content, err := io.ReadAll(configFile)
	if err != nil {
		return KudasaiDefaultCommands(), false
	}
	var config KudasaiConfig
	json.Unmarshal(content, &config)

	return mergeMaps(KudasaiDefaultCommands(), config.Commands), true
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

	for _, name := range names {
		fmt.Printf("  %s\n", name)
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
			cmds[name] = map[string]interface{}{"run": cmd}
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
		{"Go", "go.mod", Commands{"build": "go build ./...", "start": "go run .", "test": "go test ./...", "lint": "go vet ./..."}},
		{"Node.js", "package.json", Commands{"build": "npm run build", "start": "npm start", "test": "npm test", "lint": "npm run lint"}},
		{"Ruby", "Gemfile", Commands{"build": "bundle exec rake build", "start": "bundle exec ruby app.rb", "test": "bundle exec rspec", "lint": "bundle exec rubocop"}},
		{"PHP", "composer.json", Commands{"build": "composer install", "start": "php -S localhost:8000", "test": "./vendor/bin/phpunit", "lint": "./vendor/bin/phpcs"}},
		{"Python", "pyproject.toml", Commands{"build": "pip install -e .", "start": "python -m app", "test": "pytest", "lint": "ruff check ."}},
		{"Python", "requirements.txt", Commands{"build": "pip install -r requirements.txt", "start": "python -m app", "test": "pytest", "lint": "ruff check ."}},
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
			"build": "echo 'Replace with your build command'",
			"start": "echo 'Replace with your start command'",
			"test":  "echo 'Replace with your test command'",
			"lint":  "echo 'Replace with your lint command'",
		}
	}

	fmt.Println()
	fmt.Println("Commands:")

	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("  %s: %s\n", name, commands[name])
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

	switch args[0] {
	case "--version", "version":
		fmt.Printf("kudasai %s\n", Version)
		return nil
	case "--json":
		commands, _ := GetCommands("./.kudasai.json")
		return JSONOutput(commands)
	case "help", "--help":
		commands, configFound := GetCommands("./.kudasai.json")
		Help(commands, configFound)
		return nil
	case "init":
		force := len(args) > 1 && args[1] == "--force"
		return Init(force)
	}

	commands, _ := GetCommands("./.kudasai.json")
	command, exists := commands[args[0]]
	if !exists {
		return fmt.Errorf("Unrecognized command: %s", args[0])
	}

	command = InterpolateArgs(command, args[1:])

	return Execute(command)

}
