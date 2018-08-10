package core

import (
	"github.com/thoj/go-ircevent"
)

type Command struct {

	// Text of the command for example ".command argument1 argument2".
	Text string

	// Nick of the person that sent this command for example "nick".
	Nick string

	// Target of the message. A channel name or a nick if the bot is
	// messaged directly. Example: "#channel" or "bot_nick".
	Target string
}

type Module interface {
	HandleEvent(event *irc.Event)

	RunCommand(command Command) ([]string, error)
}

type Sender interface {
	Reply(e *irc.Event, text string)
}
