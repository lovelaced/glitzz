package core

import (
	"github.com/lovelaced/glitzz/util"
	"github.com/pkg/errors"
	"strings"
)

var commandNotExecutedError = errors.New("modules returned no response")
var commandReturnedNoOutputError = errors.New("command returned no output")
var commandPrefixMissingError = errors.New("command prefix missing")

const pipeSeparator = "|"

func RunCommand(loadedModules []Module, command Command, commandPrefix string) ([]string, error) {
	if !util.IsCommand(command.Text, commandPrefix) {
		return nil, commandPrefixMissingError
	}
	log.Debug("starting piping", "command", command)
	parts := strings.Split(command.Text, pipeSeparator)
	prevOutput := make([]string, 0)
	for _, part := range parts {
		text := assembleCommand(part, prevOutput, commandPrefix)
		command := Command{
			Text:   text,
			Nick:   command.Nick,
			Target: command.Target,
		}
		log.Debug("piping", "part", part, "command", command)
		output, err := findModuleResponse(loadedModules, command)
		log.Debug("find module response returned", "output", output, "err", err)
		if err != nil {
			return nil, err
		}
		if len(output) == 0 {
			return nil, commandReturnedNoOutputError
		}
		prevOutput = output
	}
	return prevOutput, nil
}

func IsPipingError(err error) bool {
	return err == commandNotExecutedError ||
		err == commandReturnedNoOutputError ||
		err == commandPrefixMissingError ||
		isMalformedCommandError(err)
}

func assembleCommand(part string, prevOutput []string, commandPrefix string) string {
	command := part
	if len(prevOutput) > 0 {
		command = command + " " + prevOutput[0]
	}
	command = strings.TrimSpace(command)
	if !strings.HasPrefix(command, commandPrefix) {
		command = commandPrefix + command
	}
	return command
}

func findModuleResponse(loadedModules []Module, command Command) ([]string, error) {
	for _, module := range loadedModules {
		output, err := module.RunCommand(command)
		if err == nil {
			return output, nil
		} else {
			if !isMalformedCommandError(err) {
				return nil, errors.Wrapf(err, "error executing command in module %T", module)
			}
		}
	}
	return nil, commandNotExecutedError
}
