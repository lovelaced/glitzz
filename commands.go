package main

import (
	"fmt"
	"github.com/dedeibel/go-4chan-api/api"
	"glitzz/strip"
	"html"
	"math/rand"
	"strings"
)

type Command func(chan<- string, []string)

var commandMap = map[string]Command{
	"test":     echo,
	"shitpost": shitpost,
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

func shitpost(msgs chan<- string, args []string) {
	var argBoard string
	if len(args) > 1 {
		argBoard = args[0]
	} else {
		boardList, err := api.GetBoards()
		if err != nil {
			msgs <- err.Error()
		} else {
			argBoard = boardList[rand.Intn(len(boardList))].Board
		}
	}
	pageNo := rand.Intn(10)
	threads, err := api.GetIndex(argBoard, pageNo)
	if err != nil {
		fmt.Printf("Error getting index of %s", argBoard)
		msgs <- err.Error()
	}
	if len(threads) < 1 {
		msgs <- "No threads found on board " + argBoard
	}
	posts := threads[rand.Intn(len(threads))].Posts
	msgs <- strip.StripTags(html.UnescapeString(posts[rand.Intn(len(posts))].Comment))
	if err != nil {
		msgs <- err.Error()
	}
}
