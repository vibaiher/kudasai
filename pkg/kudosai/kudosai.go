package kudosai

import (
	"errors"
	"fmt"
)

func Run(args []string) error {
	if len(args) == 0 {
		return errors.New("No arguments provided")
	}

	switch args[0] {
	case "help":
		fmt.Println("Usage: kudosai [command]")
	default:
		fmt.Printf("Unrecognized command: %s\n", args[0])
		return errors.New("unrecognized command")
	}

	return nil
}
