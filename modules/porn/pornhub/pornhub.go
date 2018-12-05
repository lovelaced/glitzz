package pornhub

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	pornhubURL    = "https://www.pornhub.com"
	pornhubVidFMT = pornhubURL + "/view_video.php?viewkey="
	//userAgent  = "Mozilla/5.0 (X11; Linux x86_64; rv:63.0) Gecko/20100101 Firefox/63.0"
	userAgent = "glitzz, a cute irc bot written in go"
)

// Random video URL path, global variable to facilitate testing
var vidURL = pornhubURL + "/random"

// Pornhub structure, this is where the magic happens
type Pornhub struct {
	body  []byte
	URL   string
	Title string
}

func (p *Pornhub) getLinkURL(body []byte) (err error) {
	if body == nil {
		return fmt.Errorf("Empty page")
	}
	reader := bytes.NewReader(body)
	doc, err := html.Parse(reader)
	if err != nil {
		return
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			for _, a := range n.Attr {
				if a.Key == "href" && strings.Contains(a.Val, pornhubVidFMT) {
					p.URL = a.Val
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if p.URL == "" {
		return fmt.Errorf("Link not found")
	}
	return
}

// GetPage gets a random page and put it in the pornhub struct
func (p *Pornhub) GetPage() (err error) {
	client := &http.Client{Timeout: 10 * time.Second}

	request, err := http.NewRequest("GET", vidURL, nil)
	if err != nil {
		return
	}
	request.Header.Set("Accept-Language", "en-US,en;q=0.5")
	request.Header.Set("User-Agent", userAgent)

	response, err := client.Do(request)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Something went wrong: %d code", response.StatusCode)
	}
	defer response.Body.Close()

	p.body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	return p.getLinkURL(p.body)
}

// SetTitle returns a random pornhub video title
func (p *Pornhub) SetTitle() (err error) {
	reader := bytes.NewReader(p.body)
	doc, err := html.Parse(reader)
	if err != nil {
		return
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "data-video-title" {
					p.Title = a.Val
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if p.Title == "" {
		return fmt.Errorf("Title not found")
	}
	return
}

// GetRandComment returns a random pornhub comment
func GetRandComment() ([]string, error) {
	// GO's HTML PARSER IS FUCKING BROKEN SO I'M NOT DOING THIS ANY TIME SOON
	return nil, nil
}
