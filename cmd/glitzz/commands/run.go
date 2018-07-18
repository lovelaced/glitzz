package commands

import (
	"github.com/boreq/guinea"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/logging"
	"github.com/lovelaced/glitzz/modules"
	"github.com/thoj/go-ircevent"
)

var runLog = logging.GetLogger("commands/run")

var runCmd = guinea.Command{
	Run: runRun,
	Arguments: []guinea.Argument{
		guinea.Argument{
			Name:        "config",
			Multiple:    false,
			Description: "Config file",
		},
	},
	ShortDescription: "runs the bot",
}

func runRun(c guinea.Context) error {
	conf, err := config.Load(c.Arguments[0])
	if err != nil {
		return err
	}

	sender := core.NewSender()
	modules := core.CreateModules(sender, conf)

	con := irc.IRC(conf.Nick, conf.User)
	if err = con.Connect(conf.Server); err != nil {
		return err
	}
	con.AddCallback("001", func(e *irc.Event) {
		con.Join(conf.Room)
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		runLog.Printf("PRIVMSG received: %s", e.Message())
		runModules(modules, e)
	})
	con.Loop()
	return nil
}

func runModules(modules []modules.Module, e *irc.Event) {
	for _, module := range modules {
		go module.HandleEvent(e)
	}
}
