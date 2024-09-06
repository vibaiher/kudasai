package tests

import (
	"testing"

	"github.com/vibaiher/kudasai/pkg/kudasai"
)

func TestRun_NoArgs(t *testing.T) {
	err := kudasai.Run([]string{})
	if err == nil {
		t.Errorf("Expected an error when no arguments are provided")
	}
}

func TestRun_Help(t *testing.T) {
	err := kudasai.Run([]string{"help"})
	if err != nil {
		t.Errorf("Did not expect an error when running the help command")
	}
}

func TestRun_Start(t *testing.T) {
	err := kudasai.Run([]string{"start"})
	if err != nil {
		t.Errorf("Did not expect an error when running the help command")
	}
}

func ExampleStart() {
	kudasai.Run([]string{"start"})
	// Output: sushi, kudasai
}

func TestRun_InvalidCommand(t *testing.T) {
	err := kudasai.Run([]string{"invalid"})
	if err == nil {
		t.Errorf("Expected an error for an unrecognized command")
	}
}
