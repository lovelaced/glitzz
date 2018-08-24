package core

import (
	"github.com/pkg/errors"
	"strings"
)

func RunCommand(loadedModules []Module, command Command) ([]string, error) {
	parts := strings.Split(command.Text, "|")
	prevOutput := make([]string, 0)
	for _, part := range parts {
		text := assembleCommand(part, prevOutput)
		log.Debug("piping", "part", part, "command", command)
		output, err := findModuleResponse(loadedModules, Command{
			Text:   text,
			Nick:   command.Nick,
			Target: command.Target,
		})
		if err != nil && !isPippingError(err) {
			return nil, err
		}
		if (err != nil || len(output) == 0) && len(parts) > 1 {
			return nil, errors.New("malformed pipe")
		}
		prevOutput = output
	}
	return prevOutput, nil

}

func isPippingError(err error) bool {
	return err == commandNotExecutedError || IsMalformedCommandError(err)
}

func assembleCommand(part string, prevOutput []string) string {
	command := part
	if len(prevOutput) > 0 {
		command = command + " " + prevOutput[0]
	}
	return strings.TrimSpace(command)
}

var commandNotExecutedError = errors.New("modules returned no response")

func findModuleResponse(loadedModules []Module, command Command) ([]string, error) {
	log.Debug("findModuleResponse executing", "command", command)
	for _, module := range loadedModules {
		output, err := module.RunCommand(command)
		if err == nil {
			return output, nil
		} else {
			if !IsMalformedCommandError(err) {
				return nil, errors.Wrapf(err, "error executing command in module %T", module)
			}
		}
	}
	return nil, commandNotExecutedError
}
