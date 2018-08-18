package sed

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/thoj/go-ircevent"
	"regexp"
	"strings"
	"sync"
)

var historyLimit = 100
var sedPattern = regexp.MustCompile(`^ *[a-zA-Z]?/.*/.*/?[a-zA-Z]?$`)
var hq HistoryQueue

type HistoryQueue struct {
	sync.RWMutex
	items []*irc.Event
}

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &sed{
		Base: core.NewBase("sed", sender, conf),
	}
	return rv, nil
}

type sed struct {
	core.Base
}

func (hq *HistoryQueue) Append(item *irc.Event) {
	hq.Lock()
	defer hq.Unlock()
	hq.items = append(hq.items, item)
}

func (r *sed) HandleEvent(event *irc.Event) {
	if event.Code == "PRIVMSG" {
		var nicks []string
		var repl []string
		if len(hq.items) >= historyLimit {
			hq.Lock()
			defer hq.Unlock()
			hq.items = hq.items[1:]
		}
		hq.Append(event)
		for _, msg := range hq.items {
			nicks = append(nicks, msg.Nick)
			repl = r.sedReplace(hq, nicks, strings.Fields(event.Message()), event.Nick)
		}
		go r.processReplace(repl, event)
	}
}

func (r *sed) processReplace(repl []string, e *irc.Event) {
	text := strings.Join(repl, " ")
	if text != "" {
		r.Sender.Reply(e, text)
	}
}

func (r *sed) sedReplace(hq HistoryQueue, nicks []string, arguments []string, selfnick string) []string {
	var replaced string
	var rpl []string
	for _, argument := range arguments {
		myself, other := isSed(nicks, argument)
		if myself || other {
			println("yep it's sed--->", strings.Join(arguments, " "))
			println(len(hq.items))
			if sedPattern.MatchString(strings.Join(arguments, " ")) {
				sed := strings.Split(strings.Join(arguments, " "), "/")
				println(sed[0], sed[1], sed[2])
				first, err := regexp.Compile(sed[1])
				if err != nil {
					r.Log.Debug("error compiling regexp", "sed", replaced, "err", err)
				}
				second, err := regexp.Compile(sed[2])
				if err != nil {
					r.Log.Debug("error compiling regexp", "sed", replaced, "err", err)
				}
				println("True")
				for i := range hq.items {
					if i == 0 {
						i = 1
						continue
					} else if i == len(hq.items) {
						break
					}
					println(i)
					var err error
					hq.Lock()
					defer hq.Unlock()
					if myself {
						println("Myself")
						if hq.items[len(hq.items)-i-1].Nick == selfnick && first.MatchString(hq.items[len(hq.items)-i-1].Message()) {
							println("Matched")
							println(hq.items[len(hq.items)-i-1].Message())
							replaced = first.ReplaceAllLiteralString(hq.items[len(hq.items)-i-1].Message(), second.String())
							if err != nil {
								r.Log.Debug("error running sed", "sed", replaced, "err", err)
							}
							rpl = append(rpl, replaced)
							return rpl
						}
					} else if other {
						if hq.items[len(hq.items)-i-1].Nick == arguments[0] {
							println("other")
							replaced = first.ReplaceAllLiteralString(hq.items[len(hq.items)-i-1].Message(), second.String())
							if err != nil {
								r.Log.Debug("error running sed", "sed", replaced, "err", err)
							}
							rpl = append(rpl, replaced)
							return rpl
						}
					}
				}
			}
		}
	}
	return rpl
}

func isSed(nicks []string, s string) (bool, bool) {
	var other bool
	var tmp bool
	myself, _ := regexp.MatchString("^s/.+", s)
	for _, nick := range nicks {
		tmp, _ = regexp.MatchString(nick+"([.+]s/)", s)
		other = other || tmp
	}
	return myself, other
}
