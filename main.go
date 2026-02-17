package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/vibaiher/kudasai/pkg/kudasai"
)

func main() {
	args := os.Args[1:]  // Capture command-line arguments

	err := kudasai.Run(args)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Println(err)
		os.Exit(1)
	}
}
