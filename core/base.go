package core

import (
	"errors"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/logging"
	"github.com/lovelaced/glitzz/util"
	"github.com/thoj/go-ircevent"
)

func NewBase(moduleName string, sender Sender, conf config.Config) Base {
	return Base{
		Config:   conf,
		Sender:   sender,
		commands: make(map[string]ModuleCommand),
		Log:      logging.New("modules/" + moduleName),
	}
}

type CommandArguments struct {
	Arguments []string
	Nick      string
	Target    string
}

type ModuleCommand func(arguments CommandArguments) ([]string, error)

type Base struct {
	Config   config.Config
	Sender   Sender
	commands map[string]ModuleCommand
	Log      logging.Logger
}

func (b *Base) AddCommand(name string, moduleCommand ModuleCommand) {
	b.Log.Debug("adding command", "name", name)
	b.commands[name] = moduleCommand
}

var argumentParsingError = errors.New("Argument parsing failed")
var commandNotFoundError = errors.New("Command not found")
var commandParsingError = errors.New("Command parsing failed")

func IsMalformedCommandError(err error) bool {
	return err == argumentParsingError ||
		err == commandNotFoundError ||
		err == commandParsingError
}

func (b *Base) RunCommand(command Command) ([]string, error) {
	if name, err := b.GetCommandName(command.Text); err == nil {
		moduleCommand, ok := b.commands[name]
		if ok {
			if arguments, err := b.GetCommandArguments(command.Text); err == nil {
				b.Log.Debug("executing command", "name", name)
				commandArguments := CommandArguments{
					Arguments: arguments,
					Nick:      command.Nick,
					Target:    command.Target,
				}
				return moduleCommand(commandArguments)
			}
			return nil, argumentParsingError
		}
		return nil, commandNotFoundError
	}
	return nil, commandParsingError
}

func (b *Base) HandleEvent(event *irc.Event) {
}

func (b *Base) GetCommandName(msg string) (string, error) {
	return util.GetCommandName(msg, b.Config.CommandPrefix)
}
func (b *Base) GetCommandArguments(msg string) ([]string, error) {
	return util.GetCommandArguments(msg, b.Config.CommandPrefix)
}
