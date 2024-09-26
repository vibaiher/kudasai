package kudasai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type KudasaiConfig struct {
	Commands map[string]string `json:"commands"`
}

func GetCommands(filepath string) (map[string]string, error) {
	configFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Warn! .kudasai.json file not found")
		return nil, err
	}

	defer configFile.Close()

	content, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}
	var config KudasaiConfig
	json.Unmarshal(content, &config)

	return config.Commands, nil
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

	commands, err := GetCommands("./.kudasai.json")
	if err == nil {
		command, exists := commands[args[0]]
		if exists {
			err := Execute(command)
			if err != nil {
				return fmt.Errorf("Unexpected error: %s", err)
			}
			return nil
		}
	}

	switch args[0] {
	case "start":
		err := Execute("echo \"sushi, kudasai\"")
		if err != nil {
			return fmt.Errorf("Unexpected error: %s", err)
		}
	case "help":
		fmt.Println("Usage: kudasai [command]")
	default:
		return fmt.Errorf("Unrecognized command: %s", args[0])
	}

	return nil
}
