package pornhub

import (
	"github.com/PuerkitoBio/goquery"
)

func (p *Pornhub) setInfos() {
	p.URL = url(p.doc)
	p.Title = title(p.doc)
	p.Uploader = uploader(p.doc)
	p.Categories = categories(p.doc)
}

func url(doc *goquery.Document) (url string) {
	url, _ = doc.Find("link[rel=canonical]").First().Attr("href")
	return
}

func title(doc *goquery.Document) string {
	return doc.Find(".title-container .inlineFree").Text()
}

func uploader(doc *goquery.Document) string {
	return doc.Find(".video-info-row .usernameWrap").Find("a").Text()
}

func categories(doc *goquery.Document) (categories []string) {
	doc.Find(".categoriesWrapper a[href]").Each(func(i int, s *goquery.Selection) {
		categories = append(categories, s.Text())
	})
	return
}
