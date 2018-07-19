package modules

import (
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
		log:      logging.New("modules/" + moduleName + "/base"),
	}
}

type Command func(e *irc.Event)

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

func (b *Base) HandleEvent(e *irc.Event) {
	if e.Code == "PRIVMSG" {
		if name, err := b.GetCommandName(e.Message()); err == nil {
			command, ok := b.commands[name]
			if ok {
				b.log.Debug("executing command", "name", name)
				command(e)
			}
		}
	}
}

func (b *Base) GetCommandName(msg string) (string, error) {
	return util.GetCommandName(msg, b.Config.CommandPrefix)
}
func (b *Base) GetCommandArguments(msg string) ([]string, error) {
	return util.GetCommandArguments(msg, b.Config.CommandPrefix)
}
