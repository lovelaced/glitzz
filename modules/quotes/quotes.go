package quotes

import (
	"errors"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/modules"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

func New(sender modules.Sender, conf config.Config) (modules.Module, error) {
	rv := &quotes{
		Base: modules.NewBase("quotes", sender, conf),
	}
	if err := rv.initializeQuotes(conf.QuotesDirectory); err != nil {
		return nil, err
	}
	return rv, nil
}

type quotes struct {
	modules.Base
}

func (q *quotes) initializeQuotes(directory string) error {
	return filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			q.AddCommand(f.Name(), func(arguments []string) ([]string, error) {
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
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(content), "\n")
	line, err := randomArrayEntry(lines)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func randomArrayEntry(array []string) (string, error) {
	if len(array) == 0 {
		return "", errors.New("array length is zero")
	}
	return array[rand.Intn(len(array))], nil
}
