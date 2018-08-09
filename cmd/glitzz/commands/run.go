package commands

import (
	"github.com/boreq/guinea"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/logging"
	"github.com/lovelaced/glitzz/modules"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "error loading config")
	}

	sender := core.NewSender()
	loadedModules, err := modules.CreateModules(sender, conf)
	if err != nil {
		return errors.Wrap(err, "error creating modules")
	}

	con := irc.IRC(conf.Nick, conf.User)
	if err = con.Connect(conf.Server); err != nil {
		return errors.Wrap(err, "connection failed")
	}
	con.AddCallback("001", func(e *irc.Event) {
		for _, room := range conf.Rooms {
			con.Join(room)
		}
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		handleEvent(loadedModules, e)
		runCommand(loadedModules, e, sender)
	})
	con.Loop()
	return nil
}

func handleEvent(loadedModules []core.Module, e *irc.Event) {
	for _, module := range loadedModules {
		go module.HandleEvent(e)
	}
}

func runCommand(loadedModules []core.Module, e *irc.Event, sender core.Sender) {
	command := core.Command{
		Text: e.Message(),
		Nick: e.Nick,
	}
	output, err := createPipeOutput(loadedModules, command)
	if err != nil {
		runLog.Error("error executing command", "command", command, "err", err)
		sender.Reply(e, "Internal error occured, check the logs!")
	} else {
		for _, line := range output {
			sender.Reply(e, line)
		}
	}
}

func createPipeOutput(loadedModules []core.Module, command core.Command) ([]string, error) {
	parts := strings.Split(command.Text, "|")
	prevOutput := make([]string, 0)
	for _, part := range parts {
		text := assembleCommand(part, prevOutput)
		runLog.Debug("piping", "part", part, "command", command)
		output, err := findModuleResponse(loadedModules, core.Command{
			Text: text,
			Nick: command.Nick,
		})
		if err != nil && !isPippingError(err) {
			return nil, err
		}
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

func isPippingError(err error) bool {
	return err == commandNotExecutedError || core.IsMalformedCommandError(err)
}

var commandNotExecutedError = errors.New("modules returned no response")

func findModuleResponse(loadedModules []core.Module, command core.Command) ([]string, error) {
	runLog.Debug("findModuleResponse executing", "command", command)
	for _, module := range loadedModules {
		output, err := module.RunCommand(command)
		if err == nil {
			return output, nil
		} else {
			if !core.IsMalformedCommandError(err) {
				return nil, errors.Wrapf(err, "error executing command in module %T", module)
			}
		}
	}
	return nil, commandNotExecutedError
}
