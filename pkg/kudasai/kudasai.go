package kudasai

import (
	"os/exec"
	"errors"
	"fmt"
)

func Run(args []string) error {
	if len(args) == 0 {
		return errors.New("No arguments provided")
	}

	switch args[0] {
	case "start":
		cmd := exec.Command("echo", "sushi, kudasai")
		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("Unexpected error: %s", err)
		}
		fmt.Print(string(out))
	case "help":
		fmt.Println("Usage: kudasai [command]")
	default:
		return fmt.Errorf("Unrecognized command: %s", args[0])
	}

	return nil
}
