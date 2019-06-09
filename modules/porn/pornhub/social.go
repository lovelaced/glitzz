package pornhub

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

type votes struct {
	UpVotes   int
	DownVotes int
}

type comment struct {
	Author   string
	Verified bool
	Score    int
	Message  string
}

func (p *Pornhub) setSocial() {
	p.Views = views(p.doc)

	// Using dumb function names to avoid namespace conflicts
	p.Votes = upboats(p.doc)
	p.Comments = shitposts(p.doc)
}

func views(doc *goquery.Document) (views int) {
	views, _ = strconv.Atoi(strings.ReplaceAll(doc.Find(".rating-info-container .views .count").Text(), ",", ""))
	return
}

func upboats(doc *goquery.Document) (upboats votes) {
	voteContainer := doc.Find(".votes-count-container")

	upboats.UpVotes, _ = strconv.Atoi(voteContainer.Find(".votesUp").Text())
	upboats.DownVotes, _ = strconv.Atoi(voteContainer.Find(".votesDown").Text())
	return
}

func shitposts(doc *goquery.Document) (shitposts []comment) {
	doc.Find("#cmtContent").First().Find(".topCommentBlock").Each(func(i int, s *goquery.Selection) {
		var shitpost comment
		shitpost.Author = s.Find(".usernameLink").First().Text()
		shitpost.Score, _ = strconv.Atoi(s.Find(".voteTotal").First().Text())
		shitpost.Message = s.Find(".commentMessage span").First().Text()
		shitpost.Verified = s.Find(".verified-user").Length() != 0

		shitposts = append(shitposts, shitpost)
	})
	return
}
