package kudasai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Commands map[string]string
type KudasaiConfig struct {
	Commands Commands `json:"commands"`
}

func KudasaiDefaultCommands() Commands {
	return Commands{
		"start": "echo \"sushi, kudasai\"",
		"help":  "echo \"Usage: kudasai [command]\"",
	}
}

func mergeMaps(map1, map2 map[string]string) map[string]string {
	for key, value := range map2 {
		map1[key] = value
	}
	return map1
}

func GetCommands(filepath string) Commands {
	configFile, err := os.Open(filepath)
	if err != nil {
		return KudasaiDefaultCommands()
	}
	defer configFile.Close()

	content, err := io.ReadAll(configFile)
	if err != nil {
		return KudasaiDefaultCommands()
	}
	var config KudasaiConfig
	json.Unmarshal(content, &config)

	return mergeMaps(KudasaiDefaultCommands(), config.Commands)
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

func Run(args []string) error {
	if len(args) == 0 {
		return errors.New("No arguments provided")
	}

	commands := GetCommands("./.kudasai.json")
	command, exists := commands[args[0]]
	if !exists {
		return fmt.Errorf("Unrecognized command: %s", args[0])
	}

	command = InterpolateArgs(command, args[1:])

	err := Execute(command)
	if err != nil {
		return fmt.Errorf("Unexpected error: %s", err)
	}
	return nil

}
