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

func (p *pipes) upper(arguments modules.CommandArguments) ([]string, error) {
	return []string{strings.ToUpper(strings.Join(arguments.Arguments, " "))}, nil
}

func (p *pipes) lower(arguments modules.CommandArguments) ([]string, error) {
	return []string{strings.ToLower(strings.Join(arguments.Arguments, " "))}, nil
}

func (p *pipes) echo(arguments modules.CommandArguments) ([]string, error) {
	return []string{strings.Join(arguments.Arguments, " ")}, nil
}
