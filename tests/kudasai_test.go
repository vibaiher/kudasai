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

func TestRun_HelpFlag(t *testing.T) {
	err := kudasai.Run([]string{"--help"})
	if err != nil {
		t.Errorf("Did not expect an error when running --help")
	}
}

func TestRun_Start(t *testing.T) {
	err := kudasai.Run([]string{"start"})
	if err != nil {
		t.Errorf("Did not expect an error when running the start command")
	}
}

func TestRun_InvalidCommand(t *testing.T) {
	err := kudasai.Run([]string{"invalid"})
	if err == nil {
		t.Errorf("Expected an error for an unrecognized command")
	}
}

func TestRun_Version(t *testing.T) {
	err := kudasai.Run([]string{"--version"})
	if err != nil {
		t.Errorf("Did not expect an error when running --version")
	}
}

func TestRun_VersionSubcommand(t *testing.T) {
	err := kudasai.Run([]string{"version"})
	if err != nil {
		t.Errorf("Did not expect an error when running version")
	}
}

func TestRun_JSON(t *testing.T) {
	err := kudasai.Run([]string{"--json"})
	if err != nil {
		t.Errorf("Did not expect an error when running --json")
	}
}

func TestShellQuote_SimpleArg(t *testing.T) {
	result := kudasai.ShellQuote("hello")
	expected := "'hello'"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestShellQuote_ArgWithSpaces(t *testing.T) {
	result := kudasai.ShellQuote("hello world")
	expected := "'hello world'"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestShellQuote_ArgWithSingleQuotes(t *testing.T) {
	result := kudasai.ShellQuote("it's")
	expected := "'it'\\''s'"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestInterpolateArgs_AllArgsPlaceholder(t *testing.T) {
	result := kudasai.InterpolateArgs("go test $@ ./tests/...", []string{"-v", "-run", "Foo"})
	expected := "go test '-v' '-run' 'Foo' ./tests/..."
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestInterpolateArgs_PositionalPlaceholders(t *testing.T) {
	result := kudasai.InterpolateArgs("echo hello $1, welcome to $2", []string{"world", "earth"})
	expected := "echo hello 'world', welcome to 'earth'"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestInterpolateArgs_NoPlaceholder_AppendsArgs(t *testing.T) {
	result := kudasai.InterpolateArgs("go test -v ./tests/...", []string{"-run", "TestFoo"})
	expected := "go test -v ./tests/... '-run' 'TestFoo'"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestInterpolateArgs_NoPlaceholder_NoArgs(t *testing.T) {
	result := kudasai.InterpolateArgs("go test -v ./tests/...", []string{})
	expected := "go test -v ./tests/..."
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestInterpolateArgs_MissingPositionalArg(t *testing.T) {
	result := kudasai.InterpolateArgs("echo $1 and $2", []string{"hello"})
	expected := "echo 'hello' and "
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
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
