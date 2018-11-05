package decide

import (
	"fmt"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/util"
	"regexp"
	"strings"
)

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &decide{
		Base: core.NewBase("decide", sender, conf),
	}
	rv.AddCommand("decide", rv.decide)
	return rv, nil
}

type decide struct {
	core.Base
}

func (d *decide) decide(arguments core.CommandArguments) ([]string, error) {
	if len(arguments.Arguments) > 0 {
		s := strings.Join(arguments.Arguments, " ")
		elements := regexp.MustCompile(", | or ").Split(s, -1)
		if len(elements) < 2 {
			return nil, nil
		}
		element, err := util.GetRandomArrayElement(elements)
		if err != nil {
			d.Log.Debug("get random array element failed", "err", err, "arguments", arguments)
			return nil, nil
		}
		text := fmt.Sprintf("%s: %s", arguments.Nick, element)
		return []string{text}, nil
	}
	return nil, nil
}
