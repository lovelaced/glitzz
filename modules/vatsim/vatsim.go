package vatsim

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/pkg/errors"
)

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	return newVatsim(sender, conf, newApi())
}

func newVatsim(sender core.Sender, conf config.Config, api vatsimApi) (core.Module, error) {
	rv := &vatsim{
		Base: core.NewBase("vatsim", sender, conf),
		api:  api,
	}
	rv.AddCommand("metar", rv.metar)
	return rv, nil
}

type vatsim struct {
	core.Base
	api vatsimApi
}

func (v *vatsim) metar(arguments core.CommandArguments) ([]string, error) {
	if len(arguments.Arguments) == 1 {
		metar, err := v.api.GetMetar(arguments.Arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "could not get metar")
		}
		return []string{metar}, nil
	}
	return nil, nil
}
