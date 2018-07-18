package logging

import (
	"log"
	"os"
)

func GetLogger(name string) *log.Logger {
	return log.New(os.Stdout, name+": ", 0)
}
