package fourchan

import (
	"errors"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/modules"
	"github.com/lovelaced/glitzz/util"
	"github.com/moshee/go-4chan-api/api"
	html2 "golang.org/x/net/html"
	"log"
	"math/rand"
	"strings"
	"time"
)

func New(sender modules.Sender, conf config.Config) modules.Module {
	rv := &fourchan{
		Base: modules.NewBase("fourchan", sender, conf),
		rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
	rv.AddCommand("shitpost", rv.shitpost)
	return rv
}

type fourchan struct {
	modules.Base
	rand *rand.Rand
}

func (f *fourchan) shitpost(arguments modules.CommandArguments) ([]string, error) {
	var argBoard string
	if len(arguments.Arguments) > 0 {
		argBoard = arguments.Arguments[0]
	} else {
		var err error
		argBoard, err = f.getRandomBoard()
		if err != nil {
			return nil, err
		}
	}
	pageNo := f.rand.Intn(10)
	threads, err := api.GetIndex(argBoard, pageNo)
	if err != nil {
		f.Log.Error("Error getting index\n", "board", argBoard, "error", err.Error())
		return f.shitpost(arguments)
	}
	if len(threads) < 2 {
		f.Log.Error("no threads found for some reason")
		return f.shitpost(arguments)
	}
	rnum := f.rand.Intn(len(threads) - 1)
	f.Log.Debug("selected thread", "number", rnum)
	posts := threads[rnum].Posts
	if len(posts) < 2 {
		f.Log.Warn("not enough posts")
		return f.shitpost(arguments)
	}
	f.Log.Debug("got some posts", "number", len(posts))
	post := posts[f.rand.Intn(len(posts)-1)].Comment
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
		return f.shitpost(arguments)
	}
	return []string{post}, nil
}

func (f *fourchan) getRandomBoard() (string, error) {
	boardList, err := api.GetBoards()
	if err != nil {
		return "", err
	} else {
		if len(boardList) == 0 {
			return "", errors.New("no boards")
		}
		return boardList[f.rand.Intn(len(boardList))].Board, nil
	}
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
