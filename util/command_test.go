package util

import (
	"fmt"
	"log"
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
	if err != nil {
		t.Errorf("err: %s", err)
	}
	if name != "command" {
		t.Errorf("name was: %s", name)
	}
}

func TestGetCommandNameArguments(t *testing.T) {
	name, err := GetCommandName(".command arg1 arg2", ".")
	if err != nil {
		t.Errorf("err: %s", err)
	}
	if name != "command" {
		t.Errorf("name was: %s", name)
	}
}

func ExampleGetCommandName() {
	name, err := GetCommandName(".interesting_command arg1 arg2", ".")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
	// Output: interesting_command
}

func TestGetCommandArgumentsNoCommand(t *testing.T) {
	_, err := GetCommandArguments("no command", ".")
	if err == nil {
		t.Errorf("err was nil")
	}
}

func TestGetCommandArgumentsNoArguments(t *testing.T) {
	args, err := GetCommandArguments(".command", ".")
	if err != nil {
		t.Errorf("err: %s", err)
	}
	if len(args) != 0 {
		t.Errorf("arg length was: %d", len(args))
	}
}

func TestGetCommandArguments(t *testing.T) {
	args, err := GetCommandArguments(".command arg1 arg2", ".")
	if err != nil {
		t.Errorf("err: %s", err)
	}
	if len(args) != 2 {
		t.Errorf("arg length was: %d", len(args))
	}
	if args[0] != "arg1" {
		t.Errorf("invalid first argument: %s", args[0])
	}
	if args[1] != "arg2" {
		t.Errorf("invalid second argument: %s", args[1])
	}
}

func ExampleGetCommandArguments() {
	args, err := GetCommandArguments(".interesting_command arg1 arg2", ".")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(args)
	// Output: [arg1 arg2]
}

func TestIsCommandYes(t *testing.T) {
	if !IsCommand(".command arg1 arg2", ".") {
		t.Fatalf("this should be a command")
	}
}

func TestIsCommandNo(t *testing.T) {
	if IsCommand("command arg1 arg2", ".") {
		t.Fatalf("this should not be a command")
	}
}
