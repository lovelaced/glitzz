package logging

import (
	"github.com/inconshreveable/log15"
)

type Logger = log15.Logger

func New(name string) Logger {
	return log15.New("source", name)
}
