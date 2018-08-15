package reminders

import (
	"encoding/json"
	"github.com/lovelaced/glitzz/logging"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type reminder struct {
	Nick    string
	ReplyTo string
	Message string
	Created time.Time
	Expires time.Time
}

type remindersStore interface {
	Set(r reminder) error
	Get(expireBefore time.Time) ([]reminder, error)
}

func newFileRemindersStore(filepath string, log logging.Logger) (remindersStore, error) {
	rv := &fileRemindersStore{
		filepath: filepath,
		log:      log,
	}
	if err := rv.readReminders(); err != nil {
		return nil, err
	}
	return rv, nil
}

type fileRemindersStore struct {
	filepath       string
	reminders      []reminder
	remindersMutex sync.Mutex
	log            logging.Logger
}

func (frs *fileRemindersStore) readReminders() error {
	content, err := ioutil.ReadFile(frs.filepath)
	if os.IsNotExist(err) {
		frs.log.Warn("File containing saved reminders does not exist! Is this the first run?", "error", err)
		return nil
	}
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &frs.reminders)
	if err != nil {
		return err
	}
	return nil
}

func (frs *fileRemindersStore) saveReminders() error {
	data, err := json.Marshal(frs.reminders)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(frs.filepath, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (frs *fileRemindersStore) Set(r reminder) error {
	frs.remindersMutex.Lock()
	defer frs.remindersMutex.Unlock()

	frs.reminders = append(frs.reminders, r)
	if err := frs.saveReminders(); err != nil {
		return err
	}
	return nil
}

func (frs *fileRemindersStore) Get(expireBefore time.Time) ([]reminder, error) {
	frs.remindersMutex.Lock()
	defer frs.remindersMutex.Unlock()

	var reminders []reminder
	for i := len(frs.reminders) - 1; i >= 0; i-- {
		r := frs.reminders[i]
		if r.Expires.Before(expireBefore) {
			reminders = append(reminders, r)
			frs.reminders = append(frs.reminders[:i], frs.reminders[i+1:]...)
		}
	}
	if err := frs.saveReminders(); err != nil {
		frs.reminders = append(frs.reminders, reminders...)
		return nil, err
	}
	return reminders, nil
}
