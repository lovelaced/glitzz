package modules

import (
	"github.com/thoj/go-ircevent"
)

type EventHandler interface {
	HandleEvent(event *irc.Event)
}

type Module interface {
	EventHandler
}

type Sender interface {
	Reply(e *irc.Event, text string)
}
