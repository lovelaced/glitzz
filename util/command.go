package util

import (
	"errors"
	"strings"
)

func parseCommand(msg string, sep string) (string, []string, error) {
	if strings.HasPrefix(msg, sep) {
		msg = strings.TrimPrefix(msg, sep)
		fields := strings.Fields(msg)
		if len(fields) > 0 {
			return fields[0], fields[1:], nil
		}

	}
	return "", nil, errors.New("this message is not a command")
}

func GetCommandName(msg string, sep string) (string, error) {
	name, _, err := parseCommand(msg, sep)
	return name, err
}

func GetCommandArguments(msg string, sep string) ([]string, error) {
	_, arguments, err := parseCommand(msg, sep)
	return arguments, err
}
