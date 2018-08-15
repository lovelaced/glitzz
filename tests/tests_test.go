package tests

import (
	"github.com/thoj/go-ircevent"
	"testing"
)

func TestSenderMockReplyChannel(t *testing.T) {
	const channel = "#channel"
	const text = "text"
	m := &SenderMock{}
	m.Reply(&irc.Event{Arguments: []string{channel, text}}, text)
	if len(m.Replies) != 1 {
		t.Fatalf("invalid length: %d", len(m.Replies))
	}
	if m.Replies[0].Target != channel {
		t.Errorf("invalid target: %s", m.Replies[0].Target)
	}
	if m.Replies[0].Text != text {
		t.Errorf("invalid text: %s", m.Replies[0].Text)
	}
}

func TestSenderMockReplyNick(t *testing.T) {
	const nick = "nick"
	const text = "text"
	m := &SenderMock{}
	m.Reply(&irc.Event{Nick: nick, Arguments: []string{"glitzz", text}}, text)
	if len(m.Replies) != 1 {
		t.Fatalf("invalid length: %d", len(m.Replies))
	}
	if m.Replies[0].Target != nick {
		t.Errorf("invalid target: %s", m.Replies[0].Target)
	}
	if m.Replies[0].Text != text {
		t.Errorf("invalid text: %s", m.Replies[0].Text)
	}
}

func TestSenderMockMessage(t *testing.T) {
	const target = "target"
	const text = "text"
	m := &SenderMock{}
	m.Message(target, text)
	if len(m.Replies) != 1 {
		t.Fatalf("invalid length: %d", len(m.Replies))
	}
	if m.Replies[0].Target != target {
		t.Errorf("invalid target: %s", m.Replies[0].Target)
	}
	if m.Replies[0].Text != text {
		t.Errorf("invalid text: %s", m.Replies[0].Text)
	}
}
