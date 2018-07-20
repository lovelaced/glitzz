package info

import (
	"fmt"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/modules"
	"github.com/thoj/go-ircevent"
)

const repoUrl = "https://github.com/lovelaced/glitzz"

func New(sender modules.Sender, conf config.Config) modules.Module {
	rv := &info{
		Base: modules.NewBase("info", sender, conf),
	}
	rv.AddCommand("git", rv.git)
	return rv
}

type info struct {
	modules.Base
}

func (i *info) git(arguments []string) ([]string, error) {
	return []string{repoUrl}, nil
}

func (i *info) HandleEvent(event *irc.Event) {
	if event.Code == "PRIVMSG" {
		// The IRC Bot Identification Protocol Standard
		if event.Message() == ".bots" {
			text := fmt.Sprintf("Reporting in! [Go] %s", repoUrl)
			i.Sender.Reply(event, text)
		}
	}
}
