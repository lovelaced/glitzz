package main

import (
	"fmt"
	"github.com/thoj/go-ircevent"
	"strings"
)

func parseMsg(msgs chan string, msg string, sep string) {
	msgSlice := strings.Split(msg, " ")
	println(msgSlice[:])
	prop := msgSlice[0]
	if prop[:1] == sep {
		msgSlice[0] = msgSlice[0][1:]
		commandString := msgSlice
		println(commandString)
		go run(msgs, commandString)
	}
}

func main() {

	var room = "#test"
	var separator = "."

	msgs := make(chan string)

	con := irc.IRC("glitz-test", "glitz-test")
	err := con.Connect("irc.rizon.net:6667")
	if err != nil {
		fmt.Println("Connection failed")
		return
	}
	con.AddCallback("001", func(e *irc.Event) {
		con.Join(room)
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		parseMsg(msgs, e.Message(), separator)
		message := <-msgs
		if msgs != nil {
			con.Privmsg(room, message)
		}
	})
	con.Loop()
}
