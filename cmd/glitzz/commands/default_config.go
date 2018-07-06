package commands

import (
	"encoding/json"
	"fmt"
	"github.com/boreq/guinea"
	"github.com/lovelaced/glitzz/config"
)

var defaultConfigCmd = guinea.Command{
	Run:              runDefaultConfig,
	ShortDescription: "prints the default configuration to stdout",
	Description: `
This command prints out the default config in the json format to stdout. The
available config keys are:

Debug
	Specifies if the program should run in the debug mode. The program
	running in the debug mode prints more log messages.
	Allowed values: true or false
	Default: false

Room
	Specifies which channel the bot is going to join.
	Allowed values: a string containing a valid channel name

Nick
	Nick used by the bot.
	Allowed values: a string containing a valid nick

User
	Username which appears in the full name such as nick!user@host.
	Allowed values: a string containing a valid username

Server
	Server name in the form host:port.
	Allowed values: a string containing a server address with a port

CommandPrefix
	A string used to prefix commands.
	Allowed values: a string`,
}

func runDefaultConfig(c guinea.Context) error {
	defaultConfig := config.Default()
	j, err := json.MarshalIndent(defaultConfig, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(j))
	return nil
}
