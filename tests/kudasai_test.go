package tests

import (
	"os"
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

func TestRun_InvalidCommand(t *testing.T) {
	err := kudasai.Run([]string{"invalid"})
	if err == nil {
		t.Errorf("Expected an error for an unrecognized command")
	}
}

func TestPrepare_ConfiguresStdin(t *testing.T) {
	cmd := kudasai.Prepare("echo test")

	if cmd.Stdin == nil {
		t.Errorf("Expected cmd.Stdin to be configured, but it was nil")
	}

	if cmd.Stdin != os.Stdin {
		t.Errorf("Expected cmd.Stdin to be os.Stdin")
	}
}
