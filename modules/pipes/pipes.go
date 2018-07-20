package pipes

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/modules"
	"strings"
)

func New(sender modules.Sender, conf config.Config) modules.Module {
	rv := &pipes{
		Base: modules.NewBase("pipes", sender, conf),
	}
	rv.AddCommand("upper", rv.upper)
	rv.AddCommand("lower", rv.lower)
	rv.AddCommand("echo", rv.echo)
	return rv
}

type pipes struct {
	modules.Base
}

func (p *pipes) upper(arguments []string) ([]string, error) {
	return []string{strings.ToUpper(strings.Join(arguments, " "))}, nil
}

func (p *pipes) lower(arguments []string) ([]string, error) {
	return []string{strings.ToLower(strings.Join(arguments, " "))}, nil
}

func (p *pipes) echo(arguments []string) ([]string, error) {
	return []string{strings.Join(arguments, " ")}, nil
}
