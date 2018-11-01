package sed

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/logging"
	"github.com/lovelaced/glitzz/util"
	"github.com/thoj/go-ircevent"
	"regexp"
	"strings"
	"sync"
)

var sedLog = logging.New("sed")

const historyLimit = 100

var sedPattern = regexp.MustCompile(`^.*[a-zA-Z]?/.*/.*/?[a-zA-Z]?$`)

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &sed{
		Base: core.NewBase("sed", sender, conf),
	}
	return rv, nil
}

type sed struct {
	core.Base
	sync.RWMutex
	items []*irc.Event
}

func (r *sed) HandleEvent(event *irc.Event) {
	if event.Code == "PRIVMSG" {
		var replaced []string
		r.items = append(r.items, event)
		if len(r.items) >= historyLimit {
			r.Lock()
			defer r.Unlock()
			r.items = r.items[1:]
		}
		replaced = r.sedReplace(strings.Fields(event.Message()), event.Nick)
		go r.processReplace(replaced, event)
	}
}

func (r *sed) processReplace(replaced []string, e *irc.Event) {
	text := strings.Join(replaced, " ")
	if text != "" {
		r.Sender.Reply(e, text)
	}
}

func (r *sed) sedReplace(arguments []string, selfnick string) []string {
	var replaced string
	var rpl []string
	myself, other := r.isSed(arguments)
	if myself || other {
		if sedPattern.MatchString(strings.Join(arguments, " ")) {
			sed := strings.Split(strings.Join(arguments, " "), "/")
			first, err := regexp.Compile(sed[1])
			if err != nil {
				r.Log.Debug("error compiling regexp", "sed", replaced, "err", err)
			}
			second, err := regexp.Compile(sed[2])
			if err != nil {
				r.Log.Debug("error compiling regexp", "sed", replaced, "err", err)
			}
			rep_string := util.Returntonormal(util.Boldtext(second.String()))
			for i := range r.items {
				if i == 0 {
					i = 1
					continue
				} else if i == len(r.items) {
					break
				}
				var err error
				if myself {
					if r.items[len(r.items)-i-1].Nick == selfnick && first.MatchString(r.items[len(r.items)-i-1].Message()) {
						replaced = first.ReplaceAllLiteralString(r.items[len(r.items)-i-1].Message(), rep_string)
						replaced = selfnick + util.Returntonormal(util.Boldtext(" meant")) + " to say: " + replaced
						if err != nil {
							r.Log.Debug("error running sed", "sed", replaced, "err", err)
						}
						rpl = append(rpl, replaced)
						return rpl
					}
				} else if other {
					oNick, err := regexp.MatchString(r.items[len(r.items)-i-1].Nick, arguments[0])
					if err != nil {
						sedLog.Error("something's fucked with your regex")
					}
					if oNick {
						replaced = first.ReplaceAllLiteralString(r.items[len(r.items)-i-1].Message(), rep_string)
						replaced = selfnick + " thinks " + r.items[len(r.items)-i-1].Nick + util.Returntonormal(util.Boldtext(" meant")) +
							" to say: " + replaced
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
	return rpl
}

func (r *sed) isSed(s []string) (bool, bool) {
	var other bool
	var tmp bool
	myself, _ := regexp.MatchString("^s/.+", s[0])
	for i := range r.items {
		nick := r.items[i].Nick
		isNick, err := regexp.MatchString(nick, s[0])
		if err != nil {
			sedLog.Error("somethin's fucked with your regex")
		}
		if isNick {
			tmp, _ = regexp.MatchString("^s/.+", s[1])
			other = other || tmp
			if other {
				return myself, other
			}
		}
	}
	return myself, other
}
