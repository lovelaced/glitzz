package core

import (
	"github.com/thoj/go-ircevent"
	"testing"
)

func TestSelectReplyTargetChannel(t *testing.T) {
	e := &irc.Event{
		Code:      "PRIVMSG",
		Nick:      "nick",
		Arguments: []string{"#channel", ".command"},
	}
	target := selectReplyTarget(e)
	if target != "#channel" {
		t.Errorf("target is %s", target)
	}
}

func TestSelectReplyTargetNick(t *testing.T) {
	e := &irc.Event{
		Code:      "PRIVMSG",
		Nick:      "nick",
		Arguments: []string{"glitz", ".command"},
	}
	target := selectReplyTarget(e)
	if target != "nick" {
		t.Errorf("target is %s", target)
	}
}
