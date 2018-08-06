package info

import (
	"fmt"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/thoj/go-ircevent"
)

const repoUrl = "https://github.com/lovelaced/glitzz"

func New(sender core.Sender, conf config.Config) core.Module {
	rv := &info{
		Base: core.NewBase("info", sender, conf),
	}
	rv.AddCommand("git", rv.git)
	return rv
}

type info struct {
	core.Base
}

func (i *info) git(arguments core.CommandArguments) ([]string, error) {
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
