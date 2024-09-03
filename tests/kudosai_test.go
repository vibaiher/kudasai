package tests

import (
	"testing"

	"github.com/vibaiher/kudosai/pkg/kudosai"
)

func TestRun_NoArgs(t *testing.T) {
	err := kudosai.Run([]string{})
	if err == nil {
		t.Errorf("Expected an error when no arguments are provided")
	}
}

func TestRun_Help(t *testing.T) {
	err := kudosai.Run([]string{"help"})
	if err != nil {
		t.Errorf("Did not expect an error when running the help command")
	}
}

func TestRun_InvalidCommand(t *testing.T) {
	err := kudosai.Run([]string{"invalid"})
	if err == nil {
		t.Errorf("Expected an error for an unrecognized command")
	}
}
