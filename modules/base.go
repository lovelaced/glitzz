package modules

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/logging"
	"github.com/thoj/go-ircevent"
	"log"
)

func NewBase(moduleName string, sender Sender, conf config.Config) Base {
	return Base{
		Config:   conf,
		Sender:   sender,
		commands: make(map[string]Command),
		log:      logging.GetLogger(moduleName + "/base"),
	}
}

type Command func(e *irc.Event)

type Base struct {
	Config   config.Config
	Sender   Sender
	commands map[string]Command
	log      *log.Logger
}

func (b *Base) AddCommand(name string, command Command) {
	b.log.Printf("adding command %s", name)
	b.commands[name] = command
}

func (b *Base) HandleEvent(e *irc.Event) {
	if e.Code == "PRIVMSG" {
		if name, err := b.GetCommandName(e.Message()); err == nil {
			command, ok := b.commands[name]
			if ok {
				b.log.Printf("executing command %s", name)
				command(e)
			}
		}
	}
}

func (b *Base) GetCommandName(msg string) (string, error) {
	return GetCommandName(msg, b.Config.CommandPrefix)
}
func (b *Base) GetCommandArguments(msg string) ([]string, error) {
	return GetCommandArguments(msg, b.Config.CommandPrefix)
}
