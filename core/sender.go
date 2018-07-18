package core

import (
	"github.com/lovelaced/glitzz/logging"
	"github.com/lovelaced/glitzz/modules"
	"github.com/thoj/go-ircevent"
	"log"
)

func NewSender() modules.Sender {
	return &sender{
		log: logging.GetLogger("sender"),
	}
}

type sender struct {
	log *log.Logger
}

func (s *sender) Reply(e *irc.Event, text string) {
	s.log.Printf("reply <%s, %s>", e.Arguments[0], text)
	e.Connection.Privmsg(e.Arguments[0], text)
}
