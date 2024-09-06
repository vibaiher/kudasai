package main

import (
	"fmt"
	"os"

	"github.com/vibaiher/kudasai/pkg/kudasai"
)

func main() {
	args := os.Args[1:]  // Capture command-line arguments

	err := kudasai.Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
