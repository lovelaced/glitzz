package main

import "strings"

type Command func(chan<- string, []string)

var commandMap = map[string]Command{
	"test": echo,
}

func run(msgs chan<- string, commandString []string) {
	println(len(commandString))
	args := commandString[1:]
	println("run")
	commandName := commandString[0]
	command, ok := commandMap[commandName]

	if ok {
		println("ok")
		command(msgs, args)
	}
	return
}

func echo(msgs chan<- string, args []string) {
	msgs <- strings.Join(args, " ")
	return
}
