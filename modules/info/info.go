package info

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/modules"
)

func New(sender modules.Sender, conf config.Config) modules.Module {
	rv := &info{
		Base: modules.NewBase("info", sender, conf),
	}
	rv.AddCommand("git", rv.git)
	return rv
}

type info struct {
	modules.Base
}

func (i *info) git(arguments []string) ([]string, error) {
	return []string{"https://github.com/lovelaced/glitzz"}, nil
}
