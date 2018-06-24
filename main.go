package main

import (
	"fmt"
	"github.com/thoj/go-ircevent"
	"strings"
	"time"
)

func parseMsg(msgs chan string, msg string, sep string) bool {
	msgSlice := strings.Split(msg, " ")
	println(msgSlice[:])
	prop := msgSlice[0]
	if prop[:1] == sep {
		msgSlice[0] = msgSlice[0][1:]
		commandString := msgSlice
		println(commandString)
		go run(msgs, commandString)
	}
	return true
}

func botResponse(msgs chan string, con irc.Connection, room string) {
	for {
		select {
		case message := <-msgs:
			con.Privmsg(room, message)
		default:
			time.Sleep(200 * time.Millisecond)
		}
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
		go parseMsg(msgs, e.Message(), separator)
		go botResponse(msgs, *con, room)
	})
	con.Loop()
}
