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
			fmt.Printf("Unexpected error: %s\n", err)
			return errors.New("command error")
		}
		fmt.Print(string(out))
	case "help":
		fmt.Println("Usage: kudasai [command]")
	default:
		fmt.Printf("Unrecognized command: %s\n", args[0])
		return errors.New("unrecognized command")
	}

	return nil
}
