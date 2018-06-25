package main

import (
	"fmt"
	"github.com/dedeibel/go-4chan-api/api"
	html2 "golang.org/x/net/html"
	"log"
	"math/rand"
	"strings"
)

type Command func(chan<- string, []string)

var commandMap = map[string]Command{
	"test":     echo,
	"shitpost": shitpost,
}

func htmlParser(html string) string {
	doc, err := html2.Parse(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}
	var textBody string
	var f func(*html2.Node)
	f = func(n *html2.Node) {
		if n.Type == html2.ElementNode && (n.Data == "br" || n.Data == "p") {
			textBody = textBody + " "
		}
		if n.Type != html2.ElementNode {
			textBody = textBody + n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return textBody
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

	post := htmlParser(posts[rand.Intn(len(posts))].Comment)
	msgs <- post
	//for i, c := range post {
	//	str , err := strconv.Atoi(post[i])
	//}
	if err != nil {
		msgs <- err.Error()
	}
}
