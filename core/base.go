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

func (b *Base) RunCommand(command Command) ([]string, error) {
	if name, err := b.GetCommandName(command.Text); err == nil {
		moduleCommand, ok := b.commands[name]
		if ok {
			if arguments, err := b.GetCommandArguments(command.Text); err == nil {
				b.Log.Debug("executing command", "name", name)
				commandArguments := CommandArguments{
					Arguments: arguments,
					Nick:      command.Nick,
				}
				return moduleCommand(commandArguments)
			}
			return nil, errors.New("Argument parsing failed")
		}
		return nil, errors.New("Command not found")
	}
	return nil, errors.New("Command parsing failed")
}

func (b *Base) HandleEvent(event *irc.Event) {
}

func (b *Base) GetCommandName(msg string) (string, error) {
	return util.GetCommandName(msg, b.Config.CommandPrefix)
}
func (b *Base) GetCommandArguments(msg string) ([]string, error) {
	return util.GetCommandArguments(msg, b.Config.CommandPrefix)
}
