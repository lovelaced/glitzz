package c3

import (
	"fmt"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"time"
)

const firstCongressYear = 1984
const congressMonth = time.December
const congressDay = 27

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &c3{
		Base: core.NewBase("c3", sender, conf),
	}
	rv.AddCommand("c3", rv.c3)
	return rv, nil
}

type c3 struct {
	core.Base
}

func (i *c3) c3(arguments core.CommandArguments) ([]string, error) {
	days, number := getDaysToNextCongressAndCongressNumber(time.Now())
	text := fmt.Sprintf("Time to %dC3: %d days", number, days)
	return []string{text}, nil
}

func getDaysToNextCongressAndCongressNumber(now time.Time) (int, int) {
	congress := createDate(now.Year(), congressMonth, congressDay)
	if congress.Before(now) {
		congress = createDate(now.Year()+1, congressMonth, congressDay)
	}
	number := congress.Year() - firstCongressYear + 1
	days := secondsToDays(congress.Sub(now).Seconds())
	return days, number
}

func secondsToDays(seconds float64) int {
	return int(seconds / 60 / 60 / 24)
}

func createDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
