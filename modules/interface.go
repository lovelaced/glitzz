package modules

import (
	"github.com/thoj/go-ircevent"
)

type Module interface {
	HandleEvent(event *irc.Event)
}

type Sender interface {
	Reply(e *irc.Event, text string)
}
