package main

//type command interface {
//	name() string
//	run() string
//}

var commandMap = make(map[string]func(chan string, string) string,  0)

func name(name string)string {
	return name
}

func run(msgs chan string, commandName string) string {
	commandMap[commandName] = echo
	result := echo(msgs, commandName)
	return result
}

func echo(msgs chan string, commandString string) string {
//	msgs<-commandString
	return commandString
}
