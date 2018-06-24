package main

import "strings"

//type command interface {
//	name() string
//	run() string
//}

type Command func(chan<- string, []string)

var commandMap = map[string]Command{
	"test": echo,
}

func run(msgs chan<- string, commandString []string) {
	println(len(commandString))
	args := commandString[1:]
	commandName := commandString[0]
	if command, ok := commandMap[commandName]; ok {
		command(msgs, args)
	}
	return
	//	echo(msgs, commandName)
}

func echo(msgs chan<- string, args []string) {
	msgs <- strings.Join(args, " ")
	return
}
