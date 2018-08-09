package quotes

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/util"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &quotes{
		Base: core.NewBase("quotes", sender, conf),
	}
	if err := rv.initializeQuotes(conf.Quotes.QuotesDirectory); err != nil {
		return nil, errors.Wrap(err, "could not read quotes")
	}
	return rv, nil
}

type quotes struct {
	core.Base
}

func (q *quotes) initializeQuotes(directory string) error {
	return filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			q.AddCommand(f.Name(), func(arguments core.CommandArguments) ([]string, error) {
				line, err := getRandomLine(path)
				if err != nil {
					return nil, err
				}
				return []string{line}, nil
			})
		}
		return nil
	})
}

func getRandomLine(filename string) (string, error) {
	lines, err := getAllQuotes(filename)
	if err != nil {
		return "", err
	}
	line, err := util.GetRandomArrayElement(lines)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func getAllQuotes(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimRight(string(content), "\n"), "\n"), nil
}
