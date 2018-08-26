package util

import (
	"errors"
	"strings"
)

func parseCommand(msg string, prefix string) (string, []string, error) {
	if strings.HasPrefix(msg, prefix) {
		msg = strings.TrimPrefix(msg, prefix)
		fields := strings.Fields(msg)
		if len(fields) > 0 {
			return fields[0], fields[1:], nil
		}
	}
	return "", nil, errors.New("this message is not a command")
}

// GetCommandName extracts the command name from the command string.
func GetCommandName(msg string, prefix string) (string, error) {
	name, _, err := parseCommand(msg, prefix)
	return name, err
}

// GetCommandArguments extracts the command arguments from the command string.
func GetCommandArguments(msg string, prefix string) ([]string, error) {
	_, arguments, err := parseCommand(msg, prefix)
	return arguments, err
}

func IsCommand(msg string, prefix string) bool {
	return strings.HasPrefix(msg, prefix)
}
