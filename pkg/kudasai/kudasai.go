package kudasai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
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

func Execute(command string) error {
	cmd := Prepare(command)
	err := cmd.Run()

	return err
}

func Prepare(command string) *exec.Cmd {
	cmd := exec.Command("sh", "-c", command)
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

	err := Execute(command)
	if err != nil {
		return fmt.Errorf("Unexpected error: %s", err)
	}
	return nil

}
