package modules

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/logging"
	"github.com/lovelaced/glitzz/modules/c3"
	"github.com/lovelaced/glitzz/modules/fourchan"
	"github.com/lovelaced/glitzz/modules/info"
	"github.com/lovelaced/glitzz/modules/pipes"
	"github.com/lovelaced/glitzz/modules/quotes"
	"github.com/lovelaced/glitzz/modules/tell"
	"github.com/lovelaced/glitzz/modules/untappd"
)

var log = logging.New("core")

func CreateModules(sender core.Sender, conf config.Config) ([]core.Module, error) {
	var rv []core.Module
	rv = append(rv, info.New(sender, conf))
	rv = append(rv, pipes.New(sender, conf))
	rv = append(rv, c3.New(sender, conf))
	rv = append(rv, fourchan.New(sender, conf))

	if m, err := quotes.New(sender, conf); err != nil {
		return nil, err
	} else {
		rv = append(rv, m)
	}

	if conf.Untappd != nil {
		if m, err := untappd.New(sender, conf); err != nil {
			return nil, err
		} else {
			rv = append(rv, m)
		}
	} else {
		log.Warn("Module untappd was not loaded! Missing config entry!")
	}

	if m, err := tell.New(sender, conf); err != nil {
		return nil, err
	} else {
		rv = append(rv, m)
	}

	return rv, nil
}
