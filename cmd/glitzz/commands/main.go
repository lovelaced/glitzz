package commands

import (
	"github.com/boreq/guinea"
)

var MainCmd = guinea.Command{
	Run: runMain,
	Subcommands: map[string]*guinea.Command{
		"run": &runCmd,
	},
	ShortDescription: "an IRC bot",
	Description:      "An IRC bot written in Go.",
}

func runMain(c guinea.Context) error {
	return guinea.ErrInvalidParms
}
