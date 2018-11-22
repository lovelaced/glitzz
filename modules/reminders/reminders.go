package reminders

import (
	"fmt"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/util"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	base := core.NewBase("reminders", sender, conf)
	rs, err := newFileRemindersStore(conf.Reminders.RemindersFile, base.Log)
	if err != nil {
		return nil, err
	}
	rv := &reminders{
		Base: core.NewBase("reminders", sender, conf),
		rs:   rs,
	}
	rv.AddCommand("in", rv.in)
	go rv.run()
	return rv, nil
}

type reminders struct {
	core.Base
	rs remindersStore
}

func (r *reminders) run() {
	for range time.Tick(10 * time.Second) {
		err := r.processReminders()
		if err != nil {
			r.Log.Error("could not process the reminders", "err", err)
		}
	}
}

func (r *reminders) processReminders() error {
	reminders, err := r.rs.Get(time.Now())
	if err != nil {
		return err
	}
	for _, rmd := range reminders {
		msg := formatReminder(rmd)
		r.Sender.Message(rmd.ReplyTo, msg)
	}
	return nil
}

func (r *reminders) in(arguments core.CommandArguments) ([]string, error) {
	if len(arguments.Arguments) >= 2 {
		duration, message, err := parseMessage(arguments.Arguments)
		if err != nil {
			return nil, nil
		}
		rmd := reminder{
			Nick:    arguments.Nick,
			Message: message,
			Created: time.Now().UTC(),
			Expires: time.Now().UTC().Add(duration),
		}
		if util.IsChannelName(arguments.Target) {
			rmd.ReplyTo = arguments.Target
		} else {
			rmd.ReplyTo = arguments.Nick
		}
		if err := r.rs.Set(rmd); err != nil {
			return nil, err
		} else {
			return []string{"Reminder set!"}, nil
		}
	}
	return nil, nil
}

const timeFormat = "2006-01-02 15:04:05 MST"

func formatReminder(rmd reminder) string {
	t := rmd.Created.Format(timeFormat)
	return fmt.Sprintf("%s: %s (set %s)", rmd.Nick, rmd.Message, t)
}

var unitDurations = map[time.Duration][]string{
	time.Second:          []string{"s", "sec", "second", "seconds"},
	time.Minute:          []string{"m", "min", "minute", "minutes"},
	time.Hour:            []string{"h", "hour", "hours"},
	24 * time.Hour:       []string{"d", "day", "days"},
	7 * 24 * time.Hour:   []string{"w", "week", "weeks"},
	365 * 24 * time.Hour: []string{"y", "year", "years"},
}

func getUnitDuration(u string) (time.Duration, error) {
	for duration, units := range unitDurations {
		for _, unit := range units {
			if unit == u {
				return duration, nil
			}
		}
	}
	return time.Second, errors.New("invalid unit")
}

func parseMessage(arguments []string) (time.Duration, string, error) {
	re := regexp.MustCompile("^([0-9.]+)([a-z]*)$")
	result := re.FindAllStringSubmatch(arguments[0], -1)
	if len(result) == 0 {
		return time.Second, "", errors.New("first argument doesn't match the regex")
	}
	amountString := result[0][1]
	unitString := result[0][2]
	messageSlice := arguments[1:]
	if unitString == "" {
		unitString = arguments[1]
		messageSlice = arguments[2:]
	}
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return time.Second, "", err
	}
	unitDuration, err := getUnitDuration(unitString)
	if err != nil {
		return time.Second, "", err
	}
	message := strings.Join(messageSlice, " ")
	duration := time.Duration(amount * float64(unitDuration))
	return duration, message, nil
}
