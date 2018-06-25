package main

import (
	"github.com/dedeibel/go-4chan-api/api"
	"glitzz/util"
	html2 "golang.org/x/net/html"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
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
			textBody = textBody + "\n"
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
	rand.Seed(time.Now().Unix())
	var argBoard string
	if len(args) > 0 {
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
		println("Error getting index of %s", argBoard)
		println("Error: %s", err.Error())
		shitpost(msgs, args)
		return
	}
	if len(threads) < 1 {
		msgs <- "No threads found on board " + argBoard
		shitpost(msgs, args)
		return
	}
	print(threads)
	if len(threads) < 2 {
		msgs <- "No threads found for some reason..."
		shitpost(msgs, args)
		return
	}
	rnum := rand.Intn(len(threads) - 1)
	text := "Got some threads: " + strconv.Itoa(rnum)
	println(text)
	//	msgs <- text
	posts := threads[rnum].Posts
	var post string
	if len(posts) < 2 {
		println("Not enough posts")
		shitpost(msgs, args)
		return
	}
	text = "Got some posts: " + strconv.Itoa(len(posts))
	println(text)
	post = posts[rand.Intn(len(posts)-1)].Comment
	if len(post) > 0 {
		post = htmlParser(post)
		textSlice := strings.Split(post, "\n")
		println(textSlice)
		if len(textSlice) > 1 {
			if strings.HasPrefix(textSlice[0], ">") && len(textSlice) > 1 {
				textSlice[0] = util.Greentext(textSlice[0])
				textSlice[1] = util.Normaltext(textSlice[1])
				post = strings.Join(textSlice, " ")
			} else {
				post = util.Normaltext(post)
			}
		}
	} else {
		shitpost(msgs, args)
	}
	msgs <- post
}
