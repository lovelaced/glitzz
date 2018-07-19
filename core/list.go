package core

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/modules"
	"github.com/lovelaced/glitzz/modules/info"
)

func CreateModules(sender modules.Sender, conf config.Config) []modules.Module {
	var rv []modules.Module
	rv = append(rv, info.New(sender, conf))
	return rv
}
