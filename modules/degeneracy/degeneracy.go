package degeneracy

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/thoj/go-ircevent"
	"time"
)

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &degeneracy{
		Base: core.NewBase("degeneracy", sender, conf),
	}
	return rv, nil
}

type degeneracy struct {
	core.Base
}

const shotsDelay = time.Second

func (d *degeneracy) HandleEvent(event *irc.Event) {
	if event.Code == "PRIVMSG" {
		commandName, err := d.GetCommandName(event.Message())
		if err == nil {
			if commandName == "shots" {
				for _, elem := range []string{"3", "2", "1", "GO!"} {
					d.Sender.Reply(event, elem)
					<-time.After(shotsDelay)
				}
			}
		}
	}
}
