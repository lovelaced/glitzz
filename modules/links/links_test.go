package links

import (
	"fmt"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/tests"
	"github.com/thoj/go-ircevent"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestIsLink(t *testing.T) {
	type isLinkTest struct {
		Link   string
		IsLink bool
	}

	var linkTests = []isLinkTest{
		{
			Link:   "http://example.com",
			IsLink: true,
		},
		{
			Link:   "http://example.com/a/b/c",
			IsLink: true,
		},
		{
			Link:   "https://example.com",
			IsLink: true,
		},
		{
			Link:   "https://example.com/a/b/c",
			IsLink: true,
		},
		{
			Link:   "not a link",
			IsLink: false,
		},
	}

	for _, test := range linkTests {
		output := isLink(test.Link)
		if output != test.IsLink {
			t.Errorf("invalid output for %#v: %t != %t", test.Link, test.IsLink, output)
		}
	}
}

func TestExtractLinks(t *testing.T) {
	args := []string{"http://example.com", "not a link", "https://example.com"}
	links := extractLinks(args)
	if len(links) != 2 || links[0] != args[0] || links[1] != args[2] {
		t.Fatalf("invalid output: %#v", links)
	}
}

func setup(f func(*url.URL) string, t *testing.T) (core.Module, *tests.SenderMock, string, func()) {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, f(req.URL))
	})
	server := httptest.NewServer(handler)

	senderMock := &tests.SenderMock{}
	m, err := New(senderMock, config.Default())
	if err != nil {
		t.Fatalf("could not create the module: %s", err)
	}
	return m, senderMock, server.URL, server.Close
}

func makeEvent(msg string) *irc.Event {
	return &irc.Event{Code: "PRIVMSG", Arguments: []string{"#channel", msg}}
}

const checkEvery = 100 * time.Millisecond
const checkN = 100

func TestLinks(t *testing.T) {
	const expectedLinks = 3

	f := func(u *url.URL) string {
		title := fmt.Sprintf("title for %s", u)
		return "<html><head><title>" + title + "</title></head></html>"
	}
	m, senderMock, url, cleanup := setup(f, t)
	defer cleanup()

	msg := fmt.Sprintf("message %s containing %s links %s", url, url+"/a", url+"/b")
	e := makeEvent(msg)

	m.HandleEvent(e)

	for i := 0; i < checkN; i++ {
		if len(senderMock.Replies) == expectedLinks {
			break
		}
		<-time.After(checkEvery)
	}

	if len(senderMock.Replies) != expectedLinks {
		t.Fatalf("invalid number of replies: %d", len(senderMock.Replies))
	}
}

func TestLinksTooLong(t *testing.T) {
	const expectedLinks = 1

	f := func(u *url.URL) string {
		title := strings.Repeat("a", 1000)
		return "<html><head><title>" + title + "</title></head></html>"
	}
	m, senderMock, url, cleanup := setup(f, t)
	defer cleanup()

	msg := fmt.Sprintf("message containing links %s", url)
	e := makeEvent(msg)

	m.HandleEvent(e)

	for i := 0; i < checkN; i++ {
		if len(senderMock.Replies) == expectedLinks {
			break
		}
		<-time.After(checkEvery)
	}

	if len(senderMock.Replies) != expectedLinks {
		t.Fatalf("invalid number of replies: %d", len(senderMock.Replies))
	}
	if len(senderMock.Replies[0].Text) > characterLimit+4 { // text + 2 brackets + 2 spaces
		t.Fatalf("reply too long: %d - %s", len(senderMock.Replies[0].Text), senderMock.Replies[0].Text)
	}
}
