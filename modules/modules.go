package modules

import (
	"fmt"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/logging"
	"github.com/lovelaced/glitzz/modules/c3"
	"github.com/lovelaced/glitzz/modules/fourchan"
	"github.com/lovelaced/glitzz/modules/info"
	"github.com/lovelaced/glitzz/modules/links"
	"github.com/lovelaced/glitzz/modules/pipes"
	"github.com/lovelaced/glitzz/modules/quotes"
	"github.com/lovelaced/glitzz/modules/reactions"
	"github.com/lovelaced/glitzz/modules/reminders"
	"github.com/lovelaced/glitzz/modules/seen"
	"github.com/lovelaced/glitzz/modules/tell"
	"github.com/lovelaced/glitzz/modules/untappd"
	"github.com/lovelaced/glitzz/modules/vatsim"
	"github.com/pkg/errors"
)

var log = logging.New("modules")

type moduleConstructor func(sender core.Sender, conf config.Config) (core.Module, error)

func CreateModules(sender core.Sender, conf config.Config) ([]core.Module, error) {
	return createModules(getModuleConstructors(), conf.EnabledModules, sender, conf)
}

func createModules(moduleConstructors map[string]moduleConstructor, moduleNames []string, sender core.Sender, conf config.Config) ([]core.Module, error) {
	var modules []core.Module
	for _, moduleName := range conf.EnabledModules {
		moduleConstructor, ok := moduleConstructors[moduleName]
		if !ok {
			return nil, fmt.Errorf("module %s does not exist and could not be loaded", moduleName)
		}
		module, err := moduleConstructor(sender, conf)
		if err != nil {
			return nil, errors.Wrapf(err, "module %s could not be created", moduleName)
		}
		log.Info("created module", "name", moduleName)
		modules = append(modules, module)

	}
	return modules, nil
}

func getModuleConstructors() map[string]moduleConstructor {
	modules := map[string]moduleConstructor{
		"c3":        c3.New,
		"fourchan":  fourchan.New,
		"info":      info.New,
		"links":     links.New,
		"pipes":     pipes.New,
		"quotes":    quotes.New,
		"reactions": reactions.New,
		"reminders": reminders.New,
		"seen":      seen.New,
		"tell":      tell.New,
		"untappd":   untappd.New,
		"vatsim":    vatsim.New,
	}
	return modules
}
