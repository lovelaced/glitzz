package core

import (
	ircutil "github.com/lovelaced/glitzz/irc"
	"github.com/lovelaced/glitzz/logging"
	"github.com/thoj/go-ircevent"
	"time"
)

const messageDelay = 2200 * time.Millisecond
const messageQueueLength = 100

func NewSender() Sender {
	rv := &sender{
		log:              logging.New("core/sender"),
		outgoingMessages: make(chan outgoingMessage, messageQueueLength),
	}
	go rv.run()
	return rv
}

type outgoingMessage struct {
	e    *irc.Event
	text string
}

type sender struct {
	log              logging.Logger
	outgoingMessages chan outgoingMessage
}

func (s *sender) run() {
	for {
		msg := <-s.outgoingMessages
		target := selectReplyTarget(msg.e)
		s.log.Debug("sending message", "target", target, "text", msg.text)
		msg.e.Connection.Privmsg(target, msg.text)
		<-time.After(messageDelay)
	}
}

func (s *sender) Reply(e *irc.Event, text string) {
	s.log.Debug("queueing message", "queued_messages", len(s.outgoingMessages))
	s.outgoingMessages <- outgoingMessage{e: e, text: text}
}

func selectReplyTarget(e *irc.Event) string {
	if ircutil.IsChannelName(e.Arguments[0]) {
		return e.Arguments[0]
	} else {
		return e.Nick
	}
}
