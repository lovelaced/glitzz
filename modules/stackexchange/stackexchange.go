package stackexchange

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	// "github.com/lovelaced/glitzz/logging"
)

const (
	defaultSite = "stackoverflow"
	seprefix    = "se"
	seanswer    = "a"
)

// se = seprefix
// seprefix: seq (question title)
// seq: question title

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &stackexchange{
		Base: core.NewBase("stackexchange", sender, conf),
	}
	rv.AddCommand(seprefix, rv.setitle)
	rv.AddCommand(seprefix+seanswer, rv.setitle)
	return rv, nil
}

type stackexchange struct {
	core.Base
}

func (s *stackexchange) getSite(arguments core.CommandArguments) string {
	if len(arguments.Arguments) > 0 {
		return arguments.Arguments[0]
	} else {
		return defaultSite
	}
}

func (s *stackexchange) setitle(arguments core.CommandArguments) ([]string, error) {
	site := s.getSite(arguments)
	return []string{site}, nil
}
