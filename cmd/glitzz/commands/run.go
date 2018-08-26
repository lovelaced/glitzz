package commands

import (
	"github.com/boreq/guinea"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/logging"
	"github.com/lovelaced/glitzz/modules"
	"github.com/pkg/errors"
	"github.com/thoj/go-ircevent"
	"strconv"
)

var runLog = logging.New("cmd/glitzz/commands/run")

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
		return errors.Wrap(err, "error loading config")
	}

	con := irc.IRC(conf.Nick, conf.User)
	con.UseTLS = conf.TLS
	if err = con.Connect(conf.Server); err != nil {
		return errors.Wrap(err, "connection failed")
	}

	sender := core.NewSender(con)
	loadedModules, err := modules.CreateModules(sender, conf)
	if err != nil {
		return errors.Wrap(err, "error creating modules")
	}

	con.AddCallback("001", func(e *irc.Event) {
		for _, room := range conf.Rooms {
			con.Join(room)
		}
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		go handleEvent(loadedModules, e)
		go runCommand(loadedModules, e, sender, conf.CommandPrefix)
	})
	con.AddCallback("*", func(e *irc.Event) {
		code, err := strconv.Atoi(e.Code)
		if err == nil {
			if code >= 400 {
				runLog.Error("server returned an error", "raw", e.Raw)
			}
		}
	})
	con.Loop()
	return nil
}

func handleEvent(loadedModules []core.Module, e *irc.Event) {
	for _, module := range loadedModules {
		go module.HandleEvent(e)
	}
}

func runCommand(loadedModules []core.Module, e *irc.Event, sender core.Sender, commandPrefix string) {
	command := core.Command{
		Text:   e.Message(),
		Nick:   e.Nick,
		Target: e.Arguments[0],
	}
	output, err := core.RunCommand(loadedModules, command, commandPrefix)
	if err != nil {
		if core.IsPipingError(err) {
			runLog.Debug("pipe is broken", "command", command, "err", err)
		} else {
			runLog.Error("error executing command", "command", command, "err", err)
			sender.Reply(e, "Internal error occured, check the logs!")
		}
	} else {
		for _, line := range output {
			sender.Reply(e, line)
		}
	}
}
