package main

//type command interface {
//	name() string
//	run() string
//}

type Command func(chan<- string, string)

var commandMap = make(map[string]Command, 0)

func name(name string) string {
	return name
}

func run(msgs chan string, commandName string) string {
	commandMap[commandName] = echo
	echo(msgs, commandName)
}

func echo(msgs chan string, commandString string) {
	//	msgs<-commandString
	return commandString
}
