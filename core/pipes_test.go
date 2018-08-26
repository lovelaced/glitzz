package core

import (
	"github.com/lovelaced/glitzz/config"
	"strings"
	"testing"
)

var defaultCommandPrefix = config.Default().CommandPrefix

func newModuleMock(name string, command ModuleCommand) *moduleMock {
	rv := &moduleMock{
		Base: NewBase("moduleMock"+name, nil, config.Default()),
	}
	rv.AddCommand(name, func(arguments CommandArguments) ([]string, error) {
		rv.WasExecuted = true
		return command(arguments)
	})
	return rv
}

type moduleMock struct {
	Base
	WasExecuted bool
}

// RunCommand should return a piping error if there are no modules since no
// response could be produced but it was an input error.
func TestRunCommandNoModules(t *testing.T) {
	command := Command{}
	_, err := RunCommand(nil, command, defaultCommandPrefix)
	if err == nil {
		t.Error("error was nil")
	}
	if !IsPipingError(err) {
		t.Errorf("not a piping error: %s", err)
	}
}

// RunCommand should return a piping error if there is no command in the
// message since no response could be produced but it was an input error.
func TestRunCommandNoCommand(t *testing.T) {
	modules := []Module{
		newModuleMock("a", func(arguments CommandArguments) ([]string, error) {
			return nil, nil
		}),
		newModuleMock("b", func(arguments CommandArguments) ([]string, error) {
			return nil, nil
		}),
	}
	command := Command{Text: "no command"}

	_, err := RunCommand(modules, command, defaultCommandPrefix)
	if err == nil {
		t.Error("error was nil")
	}
	if !IsPipingError(err) {
		t.Errorf("not a piping error: %s", err)
	}
}

// RunCommand should return a piping error if there was a command that doesn't
// exist in the message since no response could be produced but it was an input
// error.
func TestRunCommandMissingCommand(t *testing.T) {
	modules := []Module{
		newModuleMock("a", func(arguments CommandArguments) ([]string, error) {
			return nil, nil
		}),
		newModuleMock("b", func(arguments CommandArguments) ([]string, error) {
			return nil, nil
		}),
	}
	command := Command{Text: ".no command"}

	_, err := RunCommand(modules, command, defaultCommandPrefix)
	if err == nil {
		t.Error("error was nil")
	}
	if !IsPipingError(err) {
		t.Errorf("not a piping error: %s", err)
	}
}

// RunCommand should return a piping error if the message contained a command
// that didn't produce any output.
func TestRunCommandCommandProducesNoOutput(t *testing.T) {
	mockA := newModuleMock("a", func(arguments CommandArguments) ([]string, error) {
		return nil, nil
	})
	mockB := newModuleMock("b", func(arguments CommandArguments) ([]string, error) {
		return nil, nil
	})
	modules := []Module{mockA, mockB}
	command := Command{Text: ".a"}

	_, err := RunCommand(modules, command, defaultCommandPrefix)
	if err == nil {
		t.Error("error was nil")
	}
	if !IsPipingError(err) {
		t.Errorf("not a piping error: %s", err)
	}
	if !mockA.WasExecuted {
		t.Error("mockA wasn't executed")
	}
}

// RunCommand should return a piping error if an existing command is being
// piped into a command that doesn't exist.
func TestRunCommandInvalidSecondCommand(t *testing.T) {
	mockA := newModuleMock("a", func(arguments CommandArguments) ([]string, error) {
		return []string{"a"}, nil
	})
	mockB := newModuleMock("b", func(arguments CommandArguments) ([]string, error) {
		return []string{"b"}, nil
	})
	modules := []Module{mockA, mockB}
	command := Command{Text: ".a | .no command"}

	_, err := RunCommand(modules, command, defaultCommandPrefix)
	if err == nil {
		t.Error("error was nil")
	}
	if !IsPipingError(err) {
		t.Errorf("not a piping error: %s", err)
	}
}

// RunCommand should not return any errors if the commands are correct.
func TestRunCommandPipe(t *testing.T) {
	mockA := newModuleMock("a", func(arguments CommandArguments) ([]string, error) {
		return []string{strings.Join(arguments.Arguments, "") + "a"}, nil
	})
	mockB := newModuleMock("b", func(arguments CommandArguments) ([]string, error) {
		return []string{strings.Join(arguments.Arguments, "") + "b"}, nil
	})
	mockC := newModuleMock("c", func(arguments CommandArguments) ([]string, error) {
		return []string{strings.Join(arguments.Arguments, "") + "c"}, nil
	})
	modules := []Module{mockA, mockB, mockC}
	command := Command{Text: ".a argA | .b argB | .c argC"}

	output, err := RunCommand(modules, command, defaultCommandPrefix)
	if err != nil {
		t.Errorf("error was not nil: %s", err)
	}
	if len(output) != 1 || output[0] != "argCargBargAabc" {
		t.Errorf("invalid output: %+v", output)
	}
}

// RunCommand should fix the missing command prefixes in all parts of the pipe
// but the first one.
func TestRunCommandPipeMissingPrefixes(t *testing.T) {
	mockA := newModuleMock("a", func(arguments CommandArguments) ([]string, error) {
		return []string{strings.Join(arguments.Arguments, "") + "a"}, nil
	})
	mockB := newModuleMock("b", func(arguments CommandArguments) ([]string, error) {
		return []string{strings.Join(arguments.Arguments, "") + "b"}, nil
	})
	mockC := newModuleMock("c", func(arguments CommandArguments) ([]string, error) {
		return []string{strings.Join(arguments.Arguments, "") + "c"}, nil
	})
	modules := []Module{mockA, mockB, mockC}
	command := Command{Text: ".a argA | b argB | c argC"}

	output, err := RunCommand(modules, command, defaultCommandPrefix)
	if err != nil {
		t.Errorf("error was not nil: %s", err)
	}
	if len(output) != 1 || output[0] != "argCargBargAabc" {
		t.Errorf("invalid output: %+v", output)
	}
}

// RunCommand should fix the missing command prefixes in all parts of the pipe
// but the first one.
func TestRunCommandPipeMissingAllPrefixes(t *testing.T) {
	mockA := newModuleMock("a", func(arguments CommandArguments) ([]string, error) {
		return []string{strings.Join(arguments.Arguments, "") + "a"}, nil
	})
	mockB := newModuleMock("b", func(arguments CommandArguments) ([]string, error) {
		return []string{strings.Join(arguments.Arguments, "") + "b"}, nil
	})
	mockC := newModuleMock("c", func(arguments CommandArguments) ([]string, error) {
		return []string{strings.Join(arguments.Arguments, "") + "c"}, nil
	})
	modules := []Module{mockA, mockB, mockC}
	command := Command{Text: "a argA | b argB | c argC"}

	_, err := RunCommand(modules, command, defaultCommandPrefix)
	if err == nil {
		t.Error("error was nil")
	}
	if !IsPipingError(err) {
		t.Errorf("not a piping error: %s", err)
	}
}

// RunCommand should return a piping error if the message contains a pipe
// separator but is not a command.
func TestRunCommandNotAPipe(t *testing.T) {
	mockA := newModuleMock("a", func(arguments CommandArguments) ([]string, error) {
		return nil, nil
	})
	mockB := newModuleMock("b", func(arguments CommandArguments) ([]string, error) {
		return nil, nil
	})
	mockC := newModuleMock("this", func(arguments CommandArguments) ([]string, error) {
		return nil, nil
	})
	modules := []Module{mockA, mockB, mockC}
	command := Command{Text: "this is not | a | command"}

	_, err := RunCommand(modules, command, defaultCommandPrefix)
	if err == nil {
		t.Error("error was nil")
	}
	if !IsPipingError(err) {
		t.Errorf("not a piping error: %s", err)
	}
}
