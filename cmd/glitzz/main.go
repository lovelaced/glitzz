package main

import (
	"fmt"
	"github.com/boreq/guinea"
	"github.com/lovelaced/glitzz/cmd/glitzz/commands"
	"os"
)

func main() {
	err := guinea.Run(&commands.MainCmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
