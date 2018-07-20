package modules

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
		commands: make(map[string]Command),
		log:      logging.New("modules/" + moduleName),
	}
}

type Command func(arguments []string) ([]string, error)

type Base struct {
	Config   config.Config
	Sender   Sender
	commands map[string]Command
	log      logging.Logger
}

func (b *Base) AddCommand(name string, command Command) {
	b.log.Debug("adding command", "name", name)
	b.commands[name] = command
}

func (b *Base) RunCommand(text string) ([]string, error) {
	if name, err := b.GetCommandName(text); err == nil {
		command, ok := b.commands[name]
		if ok {
			if arguments, err := b.GetCommandArguments(text); err == nil {
				b.log.Debug("executing command", "name", name)
				return command(arguments)
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
