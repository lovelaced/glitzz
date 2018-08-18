package fourchan

import (
	"errors"
	"fmt"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/util"
	"github.com/moshee/go-4chan-api/api"
	html2 "golang.org/x/net/html"
	"log"
	"math/rand"
	"strings"
)

// numberOfRetries specifies how many threads will be downloaded at most to
// select a random post or image.
const numberOfRetries = 10

// minNumberOfPosts specifies the minimum amount of replies in a thread from
// which random posts or images will be selected. Threads with a small number
// of replies can be the stickies or complete trash.
const minNumberOfPosts = 3

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &fourchan{
		Base: core.NewBase("fourchan", sender, conf),
	}
	rv.AddCommand("shitpost", rv.shitpost)
	rv.AddCommand("pic", rv.pic)
	return rv, nil
}

type fourchan struct {
	core.Base
}

func (f *fourchan) shitpost(arguments core.CommandArguments) ([]string, error) {
	board, err := f.getBoardOrSelectRandomBoard(arguments)
	if err != nil {
		return nil, err
	}
	for i := 0; i < numberOfRetries; i++ {
		post, err := f.getRandomPost(board)
		if err == nil {
			return []string{post}, nil
		}
	}
	return nil, errors.New("Could not find a random post")
}

func (f *fourchan) pic(arguments core.CommandArguments) ([]string, error) {
	board, err := f.getBoardOrSelectRandomBoard(arguments)
	if err != nil {
		return nil, err
	}
	for i := 0; i < numberOfRetries; i++ {
		board = strings.Replace(board, "/", "", -1)
		url, err := f.getRandomImage(board)
		if err == nil {
			return []string{url}, nil
		}
	}
	return nil, errors.New("Could not find a random image")
}

func (f *fourchan) getBoardOrSelectRandomBoard(arguments core.CommandArguments) (string, error) {
	if len(arguments.Arguments) > 0 {
		return arguments.Arguments[0], nil
	} else {
		return f.getRandomBoard()
	}
}

func (f *fourchan) getRandomPost(board string) (string, error) {
	thread, err := f.getRandomThread(board)
	if err != nil {
		return "", err
	}
	if len(thread.Posts) < minNumberOfPosts {
		return "", fmt.Errorf("Found only %d posts in the thread, rejecting", len(thread.Posts))
	}
	randomPostIndex := rand.Intn(len(thread.Posts))
	post := thread.Posts[randomPostIndex].Comment
	if len(post) > 0 {
		return formatPost(post), nil
	} else {
		return "", errors.New("Post was empty")
	}
}

func (f *fourchan) getRandomImage(board string) (string, error) {
	thread, err := f.getRandomThread(board)
	if err != nil {
		return "", err
	}
	if len(thread.Posts) < minNumberOfPosts {
		return "", fmt.Errorf("Found only %d posts in the thread, rejecting", len(thread.Posts))
	}
	var postsWithFiles []*api.Post
	for _, post := range thread.Posts {
		if post.File != nil {
			postsWithFiles = append(postsWithFiles, post)
		}
	}
	if len(postsWithFiles) == 0 {
		return "", errors.New("No posts with files in this thread")
	}
	randomPostIndex := rand.Intn(len(postsWithFiles))
	randomPostWithFile := postsWithFiles[randomPostIndex]
	return randomPostWithFile.ImageURL(), nil
}

func (f *fourchan) getRandomThread(board string) (*api.Thread, error) {
	f.Log.Debug("getting threads", "board", board)
	threads, err := api.GetIndex(board, 0)
	if err != nil {
		return nil, fmt.Errorf("Error getting index of %s: %s", board, err)
	}
	f.Log.Debug("got threads", "board", board, "amount", len(threads))
	if len(threads) < 2 {
		return nil, fmt.Errorf("Found only %d threads on board %s, this is most likely a bug", len(threads), board)
	}
	randomThreadIndex := rand.Intn(len(threads))
	return threads[randomThreadIndex], nil
}

func (f *fourchan) getRandomBoard() (string, error) {
	boardList, err := api.GetBoards()
	if err != nil {
		return "", err
	} else {
		if len(boardList) == 0 {
			return "", errors.New("no boards")
		}
		return boardList[rand.Intn(len(boardList))].Board, nil
	}
}

func formatPost(post string) string {
	post = htmlParser(post)
	textSlice := strings.Split(post, "\n")
	if len(textSlice) >= 2 {
		if strings.HasPrefix(textSlice[0], ">") {
			textSlice[0] = util.Greentext(textSlice[0])
			textSlice[1] = util.Normaltext(textSlice[1])
			return strings.Join(textSlice, " ")
		}
	}
	return post
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
