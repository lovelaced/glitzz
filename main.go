package main

import (
	"strings"
	"github.com/thoj/go-ircevent"
	"fmt"
)

func parseMsg(msgs chan string, msg string, sep string) string {
	msgSlice := strings.Split(msg, " ")
	println(msgSlice[:])
	prop := msgSlice[0]
	var result = ""
	if prop[:1] == sep {
		commandname := prop[1:]
		println(commandname)
		result = run(msgs, commandname)
	//	go commandMap[commandname](msgs, msg)
	} else {
		result = ""
	}
	msgs <- result
	return "uhhhh"
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
	con.AddCallback("001", func (e *irc.Event) {
		con.Join(room)
		})
	con.AddCallback("PRIVMSG", func (e *irc.Event) {
		go parseMsg(msgs, e.Message(), separator)
		message := <-msgs
		if msgs != nil {
			con.Privmsg(room, message)
		}
	})
	con.Loop()
}

