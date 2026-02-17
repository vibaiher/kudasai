package tests

import (
	"errors"
	"os"
	"os/exec"
	"strings"
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

func withFakeStdin(input string, fn func()) {
	original := kudasai.Stdin
	kudasai.Stdin = strings.NewReader(input)
	defer func() { kudasai.Stdin = original }()
	fn()
}

func TestRun_Init(t *testing.T) {
	dir, _ := os.MkdirTemp("", "kudasai-test")
	defer os.RemoveAll(dir)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	withFakeStdin("y\n", func() {
		err := kudasai.Run([]string{"init"})
		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}
	})

	content, err := os.ReadFile(".kudasai.json")
	if err != nil {
		t.Fatalf("Expected .kudasai.json to exist")
	}

	if !strings.Contains(string(content), "build") {
		t.Errorf("Expected template to contain 'build' command")
	}
}

func TestRun_InitDetectsGo(t *testing.T) {
	dir, _ := os.MkdirTemp("", "kudasai-test")
	defer os.RemoveAll(dir)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)

	withFakeStdin("y\n", func() {
		err := kudasai.Run([]string{"init"})
		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}
	})

	content, _ := os.ReadFile(".kudasai.json")
	if !strings.Contains(string(content), "go build") {
		t.Errorf("Expected Go commands, got %s", string(content))
	}
}

func TestRun_InitDetectsNode(t *testing.T) {
	dir, _ := os.MkdirTemp("", "kudasai-test")
	defer os.RemoveAll(dir)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.WriteFile("package.json", []byte("{}"), 0644)

	withFakeStdin("y\n", func() {
		err := kudasai.Run([]string{"init"})
		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}
	})

	content, _ := os.ReadFile(".kudasai.json")
	if !strings.Contains(string(content), "npm test") {
		t.Errorf("Expected Node commands, got %s", string(content))
	}
}

func TestRun_InitAborted(t *testing.T) {
	dir, _ := os.MkdirTemp("", "kudasai-test")
	defer os.RemoveAll(dir)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	withFakeStdin("n\n", func() {
		err := kudasai.Run([]string{"init"})
		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}
	})

	if _, err := os.Stat(".kudasai.json"); err == nil {
		t.Error("Expected .kudasai.json to NOT be created when user says no")
	}
}

func TestRun_InitAlreadyExists(t *testing.T) {
	dir, _ := os.MkdirTemp("", "kudasai-test")
	defer os.RemoveAll(dir)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.WriteFile(".kudasai.json", []byte("{}"), 0644)

	err := kudasai.Run([]string{"init"})
	if err == nil {
		t.Fatal("Expected error when .kudasai.json already exists")
	}
}

func TestRun_InitForce(t *testing.T) {
	dir, _ := os.MkdirTemp("", "kudasai-test")
	defer os.RemoveAll(dir)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.WriteFile(".kudasai.json", []byte("{}"), 0644)

	withFakeStdin("y\n", func() {
		err := kudasai.Run([]string{"init", "--force"})
		if err != nil {
			t.Fatalf("Expected no error with --force, got %s", err)
		}
	})
}

func TestRun_Check(t *testing.T) {
	dir, _ := os.MkdirTemp("", "kudasai-test")
	defer os.RemoveAll(dir)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.WriteFile(".kudasai.json", []byte(`{"commands":{"build":"make","test":"make test"}}`), 0644)

	err := kudasai.Run([]string{"check"})
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}
}

func TestRun_CheckMissing(t *testing.T) {
	dir, _ := os.MkdirTemp("", "kudasai-test")
	defer os.RemoveAll(dir)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	err := kudasai.Run([]string{"check"})
	if err == nil {
		t.Fatal("Expected error when .kudasai.json missing")
	}
	if !strings.Contains(err.Error(), "kudasai init") {
		t.Errorf("Expected suggestion to run init, got: %s", err)
	}
}

func TestRun_CheckInvalidJSON(t *testing.T) {
	dir, _ := os.MkdirTemp("", "kudasai-test")
	defer os.RemoveAll(dir)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.WriteFile(".kudasai.json", []byte(`{invalid`), 0644)

	err := kudasai.Run([]string{"check"})
	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "failed to parse") || !strings.Contains(err.Error(), "kudasai check") {
		t.Errorf("Expected parse error with suggestion, got: %s", err)
	}
}

func TestExecute_PropagatesExitCode(t *testing.T) {
	err := kudasai.Execute("exit 42")
	if err == nil {
		t.Fatal("Expected an error for non-zero exit code")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("Expected exec.ExitError, got %T", err)
	}

	if exitErr.ExitCode() != 42 {
		t.Errorf("Expected exit code 42, got %d", exitErr.ExitCode())
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
