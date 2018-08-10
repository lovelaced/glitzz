package seen

import (
	"encoding/json"
	"fmt"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/logging"
	"github.com/lovelaced/glitzz/util"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type messageStore interface {
	Save(nick string, channel string, date time.Time) error
	Get(nick string, channel string) (*time.Time, error)
}

func newFileMessageStore(filepath string, log logging.Logger) (messageStore, error) {
	rv := &fileMessageStore{
		filepath: filepath,
		messages: make(map[string]map[string]time.Time),
		log:      log,
	}
	if err := rv.readMessages(); err != nil {
		return nil, err
	}
	return rv, nil
}

type fileMessageStore struct {
	filepath      string
	messages      map[string]map[string]time.Time
	messagesMutex sync.Mutex
	log           logging.Logger
}

func (fms *fileMessageStore) readMessages() error {
	content, err := ioutil.ReadFile(fms.filepath)
	if os.IsNotExist(err) {
		fms.log.Warn("File containing saved messages does not exist! Is this the first run?", "error", err)
		return nil
	}
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &fms.messages)
	if err != nil {
		return err
	}
	return nil
}

func (fms *fileMessageStore) saveMessages() error {
	data, err := json.Marshal(fms.messages)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fms.filepath, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (fms *fileMessageStore) Save(nick string, channel string, date time.Time) error {
	fms.messagesMutex.Lock()
	defer fms.messagesMutex.Unlock()

	nick = strings.ToLower(nick)
	channel = strings.ToLower(channel)

	channelData, ok := fms.messages[channel]
	if !ok {
		channelData = make(map[string]time.Time)
	}
	channelData[nick] = date
	fms.messages[channel] = channelData
	if err := fms.saveMessages(); err != nil {
		return err
	}
	return nil
}

func (fms *fileMessageStore) Get(nick string, channel string) (*time.Time, error) {
	fms.messagesMutex.Lock()
	defer fms.messagesMutex.Unlock()

	nick = strings.ToLower(nick)
	channel = strings.ToLower(channel)

	channelData, ok := fms.messages[channel]
	if ok {
		date, ok := channelData[nick]
		if ok {
			return &date, nil
		}
	}
	return nil, nil
}

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	base := core.NewBase("seen", sender, conf)
	ms, err := newFileMessageStore(conf.Seen.SeenFile, base.Log)
	if err != nil {
		return nil, err
	}
	rv := &seen{
		Base: base,
		ms:   ms,
	}
	rv.AddCommand("seen", rv.seen)
	return rv, nil
}

type seen struct {
	core.Base
	ms messageStore
}

func (s *seen) seen(arguments core.CommandArguments) ([]string, error) {
	if len(arguments.Arguments) == 1 {
		target := arguments.Target
		nick := arguments.Arguments[0]
		if util.IsChannelName(target) {
			date, err := s.ms.Get(nick, target)
			if err != nil {
				s.Log.Error("could not retrieve data", "err", err)
			} else {
				if date != nil {
					text := formatMessage(nick, target, *date)
					return []string{text}, nil
				} else {
					text := formatMessageNotFound(nick, target)
					return []string{text}, nil
				}
			}
		}
	}
	return nil, nil
}

func (s *seen) HandleEvent(event *irc.Event) {
	if event.Code == "PRIVMSG" {
		target := event.Arguments[0]
		if util.IsChannelName(target) {
			if err := s.ms.Save(event.Nick, target, time.Now().UTC()); err != nil {
				s.Log.Error("could not record a message", "err", err)
			}
		}
	}
}

const dateFormat = "2006-01-02 15:04:05 MST"

func formatMessage(nick string, channel string, date time.Time) string {
	dateText := date.Format(dateFormat)
	return fmt.Sprintf("%s was last seen in %s on %s", nick, channel, dateText)
}

func formatMessageNotFound(nick string, channel string) string {
	return fmt.Sprintf("I have never seen %s in %s", nick, channel)
}
