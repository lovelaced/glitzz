package util

import (
	"testing"
)

func TestGetCommandNameNoCommand(t *testing.T) {
	_, err := GetCommandName("no command", ".")
	if err == nil {
		t.Error("err was nil")
	}
}

func TestGetCommandNameEmptyString(t *testing.T) {
	_, err := GetCommandName("", ".")
	if err == nil {
		t.Error("err was nil")
	}
}

func TestGetCommandNameNoArguments(t *testing.T) {
	name, err := GetCommandName(".command", ".")

	if name != "command" {
		t.Errorf("name was: %s", name)
	}
	if err != nil {
		t.Errorf("err: %s", err)
	}
}

func TestGetCommandNameArguments(t *testing.T) {
	name, err := GetCommandName(".command arg1 arg2", ".")

	if name != "command" {
		t.Errorf("name was: %s", name)
	}
	if err != nil {
		t.Errorf("err: %s", err)
	}
}

func TestGetCommandNoCommand(t *testing.T) {
	_, err := GetCommandArguments("command", ".")

	if err == nil {
		t.Errorf("err was nil")
	}
}

func TestGetCommandNoArguments(t *testing.T) {
	args, err := GetCommandArguments(".command", ".")

	if len(args) != 0 {
		t.Errorf("arg length was: %d", len(args))
	}
	if err != nil {
		t.Errorf("err: %s", err)
	}
}

func TestGetCommandArguments(t *testing.T) {
	args, err := GetCommandArguments(".command arg1 arg2", ".")

	if len(args) != 2 {
		t.Errorf("arg length was: %d", len(args))
	}
	if err != nil {
		t.Errorf("err: %s", err)
	}
}
