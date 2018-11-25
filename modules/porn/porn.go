package porn

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	// "github.com/lovelaced/glitzz/modules/porn/pornhub"
	"github.com/lovelaced/glitzz/modules/porn/pornmd"
)

// New registers the porn module
func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &porn{
		Base: core.NewBase("porn", sender, conf),
	}

	rv.AddCommand("pornmd", rv.pornMD)
	rv.AddCommand("porn", rv.pornMD)
	rv.AddCommand("porntitle", rv.pornMD)
	rv.AddCommand("pornlast", rv.pornMD)
	return rv, nil
}

type porn struct {
	core.Base
}

func (p *porn) pornMD(arguments core.CommandArguments) ([]string, error) {
	return pornmd.ReturnRandSearch()
}

// func (p *porn) pornhub(arguments core.CommandArguments) ([]string, error) {
// 	pornhub.GetRandPage()
// 	return nil, nil
// }
