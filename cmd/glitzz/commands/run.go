package commands

import (
	"fmt"
	"github.com/boreq/guinea"
	"github.com/lovelaced/glitzz/commands"
	"github.com/lovelaced/glitzz/config"
	"github.com/thoj/go-ircevent"
	"strings"
)

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
	config, err := config.Load(c.Arguments[0])
	if err != nil {
		return err
	}

	msgs := make(chan string)

	con := irc.IRC(config.Nick, config.User)
	err = con.Connect(config.Server)
	if err != nil {
		fmt.Println("Connection failed")
		return err
	}
	go botResponse(msgs, *con, config.Room)

	con.AddCallback("001", func(e *irc.Event) {
		con.Join(config.Room)
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		go parseMsg(msgs, e.Message(), config.CommandPrefix)
	})
	con.Loop()
	return nil
}

func parseMsg(msgs chan<- string, msg string, sep string) {
	msgSlice := strings.Split(msg, " ")
	println(msgSlice[:])
	prop := msgSlice[0]

	if prop[:1] == sep {
		msgSlice[0] = msgSlice[0][1:]
		commandString := msgSlice
		println(commandString)
		commands.Run(msgs, commandString)
	}
	return
}

func botResponse(msgs chan string, con irc.Connection, room string) {
	for {
		select {
		case message := <-msgs:
			con.Privmsg(room, message)
		}
	}
}
