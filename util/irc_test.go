package util

import (
	"github.com/thoj/go-ircevent"
	"testing"
)

func TestIsChannelNameChannel(t *testing.T) {
	if !IsChannelName("#channel") {
		t.Errorf("#channel should be a channel name")
	}
}

func TestIsChannelNameNick(t *testing.T) {
	if IsChannelName("nick") {
		t.Errorf("nick should not be a channel name")
	}
}

func TestSelectReplyTargetChannel(t *testing.T) {
	e := &irc.Event{
		Code:      "PRIVMSG",
		Nick:      "nick",
		Arguments: []string{"#channel", ".command"},
	}
	target := SelectReplyTarget(e)
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
	target := SelectReplyTarget(e)
	if target != "nick" {
		t.Errorf("target is %s", target)
	}
}
