package tests

import (
	"github.com/lovelaced/glitzz/util"
	"github.com/thoj/go-ircevent"
)

type Reply struct {
	Target string
	Text   string
}

type SenderMock struct {
	Replies []Reply
}

func (s *SenderMock) Reply(e *irc.Event, text string) {
	target := util.SelectReplyTarget(e)
	s.Replies = append(s.Replies, Reply{Target: target, Text: text})
}

func (s *SenderMock) Message(target string, text string) {
	s.Replies = append(s.Replies, Reply{Target: target, Text: text})
}
