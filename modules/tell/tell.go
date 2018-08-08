package tell

import (
	"encoding/json"
	"fmt"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/logging"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type savedMessage struct {
	Author  string
	Target  string
	Message string
	Time    time.Time
}

type messageStore interface {
	SaveMessage(msg savedMessage) error
	GetMessagesForNick(target string) ([]savedMessage, error)
}

func newFileMessageStore(filepath string, log logging.Logger) (messageStore, error) {
	rv := &fileMessageStore{
		filepath: filepath,
		messages: make(map[string][]savedMessage),
		log:      log,
	}
	if err := rv.readMessages(); err != nil {
		return nil, err
	}
	return rv, nil
}

type fileMessageStore struct {
	filepath      string
	messages      map[string][]savedMessage
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

func (fms *fileMessageStore) SaveMessage(msg savedMessage) error {
	fms.messagesMutex.Lock()
	defer fms.messagesMutex.Unlock()

	msg.Target = strings.ToLower(msg.Target)
	messages, ok := fms.messages[msg.Target]
	if !ok {
		messages = make([]savedMessage, 0)
	}
	messages = append(messages, msg)
	fms.messages[msg.Target] = messages
	err := fms.saveMessages()
	if err != nil {
		return err
	}
	return nil
}

func (fms *fileMessageStore) GetMessagesForNick(target string) ([]savedMessage, error) {
	fms.messagesMutex.Lock()
	defer fms.messagesMutex.Unlock()

	target = strings.ToLower(target)
	messages, ok := fms.messages[target]
	if ok {
		delete(fms.messages, target)
		err := fms.saveMessages()
		if err != nil {
			fms.log.Error("could not save messages", "error", err)
		}
		return messages, nil
	}
	return nil, nil
}

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	base := core.NewBase("tell", sender, conf)
	ms, err := newFileMessageStore(conf.Tell.TellFile, base.Log)
	if err != nil {
		return nil, err
	}
	rv := &tell{
		Base: base,
		ms:   ms,
	}
	rv.AddCommand("tell", rv.tell)
	return rv, nil
}

type tell struct {
	core.Base
	ms messageStore
}

func (t *tell) tell(arguments core.CommandArguments) ([]string, error) {
	if len(arguments.Arguments) >= 2 {
		msg := savedMessage{
			Author:  arguments.Nick,
			Target:  arguments.Arguments[0],
			Message: strings.Join(arguments.Arguments[1:], " "),
			Time:    time.Now().UTC(),
		}
		err := t.ms.SaveMessage(msg)
		if err != nil {
			t.Log.Error("could not save message",
				"msg", msg, "err", err)
		} else {
			return []string{"Will do!"}, nil
		}
	}
	return nil, nil
}

func (t *tell) HandleEvent(event *irc.Event) {
	if event.Code == "PRIVMSG" {
		messages, err := t.ms.GetMessagesForNick(event.Nick)
		if err != nil {
			t.Log.Error("could not retrieve saved messages",
				"nick", event.Nick, "err", err)
		} else {
			for _, msg := range messages {
				t.Sender.Reply(event, formatMessage(msg))
			}
		}
	}
}

func formatMessage(msg savedMessage) string {
	time := msg.Time.Format("15:04")
	return fmt.Sprintf("%s: %s <%s> %s", msg.Target, time, msg.Author, msg.Message)
}
