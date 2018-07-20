package commands

import (
	"errors"
	"github.com/boreq/guinea"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/logging"
	"github.com/lovelaced/glitzz/modules"
	"github.com/thoj/go-ircevent"
	"strings"
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
		handleEvent(modules, e)
		runCommand(modules, e, sender)
	})
	con.Loop()
	return nil
}

func handleEvent(modules []modules.Module, e *irc.Event) {
	for _, module := range modules {
		go module.HandleEvent(e)
	}
}

func runCommand(modules []modules.Module, e *irc.Event, sender modules.Sender) {
	output, err := createPipeOutput(modules, e.Message())
	if err != nil {
		sender.Reply(e, err.Error())
	} else {
		for _, line := range output {
			sender.Reply(e, line)
		}
	}
}

func createPipeOutput(modules []modules.Module, text string) ([]string, error) {
	parts := strings.Split(text, "|")
	prevOutput := make([]string, 0)
	for _, part := range parts {
		command := assembleCommand(part, prevOutput)
		runLog.Debug("piping", "part", part, "command", command)
		output, err := findModuleResponse(modules, command)
		if (err != nil || len(output) == 0) && len(parts) > 1 {
			return nil, errors.New("malformed pipe")
		}
		prevOutput = output
	}
	return prevOutput, nil

}

func assembleCommand(part string, prevOutput []string) string {
	command := part
	if len(prevOutput) > 0 {
		command = command + " " + prevOutput[0]
	}
	return strings.TrimSpace(command)
}

func findModuleResponse(modules []modules.Module, text string) ([]string, error) {
	runLog.Debug("findModuleResponse executing", "text", text)
	for _, module := range modules {
		output, err := module.RunCommand(text)
		if err == nil {
			return output, nil
		}
	}
	return nil, errors.New("modules returned no response")
}
