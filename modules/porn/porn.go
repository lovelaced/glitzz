package porn

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/modules/porn/pornhub"
	"github.com/lovelaced/glitzz/modules/porn/pornmd"
)

// New registers the porn module
func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &porn{
		Base: core.NewBase("porn", sender, conf),
	}

	rv.AddCommand("pornmd", rv.pornMD)
	rv.AddCommand("porn", rv.pornHubTitle)
	rv.AddCommand("porntitle", rv.pornHubTitle)
	rv.AddCommand("pornlast", rv.pornHubLast)
	return rv, nil
}

type porn struct {
	core.Base
	hub pornhub.Pornhub
}

func (p *porn) pornMD(arguments core.CommandArguments) ([]string, error) {
	return pornmd.ReturnRandSearch()
}

func (p *porn) pornHubTitle(arguments core.CommandArguments) ([]string, error) {
	err := p.hub.GetPage()
	if err != nil {
		return nil, err
	}
	err = p.hub.SetTitle()
	if err != nil {
		return nil, err
	}
	return []string{p.hub.Title}, nil
}

func (p *porn) pornHubLast(arguments core.CommandArguments) ([]string, error) {
	return []string{p.hub.URL}, nil
}
