package info

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/tests"
	"github.com/thoj/go-ircevent"
	"strings"
	"testing"
)

func TestGit(t *testing.T) {
	p, err := New(nil, config.Default())
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	output, err := p.RunCommand(core.Command{Text: ".git", Nick: "nick"})
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
	if !strings.HasPrefix(output[0], "http") {
		t.Errorf("invalid output %s", output[0])
	}
}

func TestIbip(t *testing.T) {
	s := &tests.SenderMock{}
	p, err := New(s, config.Default())
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	e := irc.Event{Arguments: []string{".bots"}, Code: "PRIVMSG"}
	p.HandleEvent(&e)
	if len(s.Replies) != 1 {
		t.Errorf("invalid output length %d", len(s.Replies))
	}
	if !strings.HasPrefix(s.Replies[0].Text, "Reporting in") {
		t.Errorf("invalid output %s", s.Replies[0].Text)
	}
}
