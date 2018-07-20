package core

import (
	"github.com/lovelaced/glitzz/logging"
	"github.com/lovelaced/glitzz/modules"
	"github.com/thoj/go-ircevent"
)

func NewSender() modules.Sender {
	return &sender{
		log: logging.New("core/sender"),
	}
}

type sender struct {
	log logging.Logger
}

func (s *sender) Reply(e *irc.Event, text string) {
	s.log.Debug("reply", "target", e.Arguments[0], "text", text)
	e.Connection.Privmsg(e.Arguments[0], text)
}
