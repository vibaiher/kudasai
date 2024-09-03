package main

import (
	"fmt"
	"os"

	"github.com/vibaiher/kudosai/pkg/kudosai"
)

func main() {
	args := os.Args[1:]  // Capture command-line arguments

	err := kudosai.Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
