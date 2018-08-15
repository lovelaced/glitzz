package core

import (
	"github.com/lovelaced/glitzz/logging"
	"github.com/lovelaced/glitzz/util"
	"github.com/thoj/go-ircevent"
	"time"
)

const messageDelay = 2200 * time.Millisecond
const messageQueueLength = 100

func NewSender(conn *irc.Connection) Sender {
	rv := &sender{
		log:              logging.New("core/sender"),
		outgoingMessages: make(chan outgoingMessage, messageQueueLength),
		conn:             conn,
	}
	go rv.run()
	return rv
}

type outgoingMessage struct {
	target string
	text   string
}

type sender struct {
	log              logging.Logger
	outgoingMessages chan outgoingMessage
	conn             *irc.Connection
}

func (s *sender) run() {
	for {
		msg := <-s.outgoingMessages
		s.log.Debug("sending message", "target", msg.target, "text", msg.text)
		s.conn.Privmsg(msg.target, msg.text)
		<-time.After(messageDelay)
	}
}

func (s *sender) Message(target string, text string) {
	s.queueMessage(target, text)
}

func (s *sender) Reply(e *irc.Event, text string) {
	target := util.SelectReplyTarget(e)
	s.queueMessage(target, text)
}

func (s *sender) queueMessage(target string, text string) {
	s.log.Debug("queueing message", "queued_messages", len(s.outgoingMessages))
	s.outgoingMessages <- outgoingMessage{
		target: target,
		text:   text,
	}
}
