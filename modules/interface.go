package modules

import (
	"github.com/thoj/go-ircevent"
)

type Module interface {
	HandleEvent(event *irc.Event)
	RunCommand(text string) ([]string, error)
}

type Sender interface {
	Reply(e *irc.Event, text string)
}
