package main

import (
	"github.com/boreq/guinea"
	"github.com/lovelaced/glitzz/cmd/glitzz/commands"
	"github.com/lovelaced/glitzz/logging"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"time"
)

var log = logging.New("cmd/glitzz")

func main() {
	rand.Seed(time.Now().UnixNano())
	injectGlobalBehaviour(&commands.MainCmd)
	err := guinea.Run(&commands.MainCmd)
	if err != nil {
		log.Crit(err.Error())
		os.Exit(1)
	}
}

func injectGlobalBehaviour(cmd *guinea.Command) {
	cmd.Options = append(cmd.Options, guinea.Option{
		Name:        "loglevel",
		Type:        guinea.String,
		Default:     "info",
		Description: "One of: debug, info, warn, error or crit. Default: info",
	})
	oldRun := cmd.Run
	cmd.Run = func(c guinea.Context) error {
		level, err := logging.LevelFromString(c.Options["loglevel"].Str())
		if err != nil {
			return errors.Wrap(err, "Could not select a log level")
		}
		logging.SetLoggingLevel(level)
		if oldRun != nil {
			return oldRun(c)
		}
		return nil
	}
	for _, subCmd := range cmd.Subcommands {
		injectGlobalBehaviour(subCmd)
	}
}
