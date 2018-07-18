package info

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/modules"
	"github.com/thoj/go-ircevent"
)

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

func (i *info) git(e *irc.Event) {
	i.Sender.Reply(e, "https://github.com/lovelaced/glitzz")
}
