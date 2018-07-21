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
	loadedModules, err := core.CreateModules(sender, conf)
	if err != nil {
		return err
	}

	con := irc.IRC(conf.Nick, conf.User)
	if err = con.Connect(conf.Server); err != nil {
		return err
	}
	con.AddCallback("001", func(e *irc.Event) {
		con.Join(conf.Room)
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		handleEvent(loadedModules, e)
		runCommand(loadedModules, e, sender)
	})
	con.Loop()
	return nil
}

func handleEvent(modules []modules.Module, e *irc.Event) {
	for _, module := range modules {
		go module.HandleEvent(e)
	}
}

func runCommand(loadedModules []modules.Module, e *irc.Event, sender modules.Sender) {
	output, err := createPipeOutput(loadedModules, modules.Command{
		Text: e.Message(),
		Nick: e.Nick,
	})
	if err != nil {
		sender.Reply(e, err.Error())
	} else {
		for _, line := range output {
			sender.Reply(e, line)
		}
	}
}

func createPipeOutput(loadedModules []modules.Module, command modules.Command) ([]string, error) {
	parts := strings.Split(command.Text, "|")
	prevOutput := make([]string, 0)
	for _, part := range parts {
		text := assembleCommand(part, prevOutput)
		runLog.Debug("piping", "part", part, "command", command)
		output, err := findModuleResponse(loadedModules, modules.Command{
			Text: text,
			Nick: command.Nick,
		})
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

func findModuleResponse(loadedModules []modules.Module, command modules.Command) ([]string, error) {
	runLog.Debug("findModuleResponse executing", "command", command)
	for _, module := range loadedModules {
		output, err := module.RunCommand(command)
		if err == nil {
			return output, nil
		}
	}
	return nil, errors.New("modules returned no response")
}
