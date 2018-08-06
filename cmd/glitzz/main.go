package main

import (
	"fmt"
	"github.com/boreq/guinea"
	"github.com/lovelaced/glitzz/cmd/glitzz/commands"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	err := guinea.Run(&commands.MainCmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
