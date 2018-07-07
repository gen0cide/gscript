package main

import (
	"github.com/urfave/cli"
)

var (
	shellCommand = cli.Command{
		Name:      "shell",
		Usage:     "drop into an interactive REPL within the genesis runtime",
		UsageText: "gscript shell [--macro MACRO] [--macro MACRO] ...",
		Action:    interactiveShellCommand,
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "macro, m",
				Usage: "apply a compiler macro to the interactive shell",
			},
		},
	}
)

func interactiveShellCommand(c *cli.Context) error {
	if len(c.StringSlice("macro")) > 0 {
		for _, m := range c.StringSlice("macro") {
			cliLogger.Infof("Found macro: %s", m)
		}
	}
	return commandNotImplemented(c)
}
