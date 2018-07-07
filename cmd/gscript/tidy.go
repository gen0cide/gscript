package main

import (
	"errors"

	"github.com/urfave/cli"
)

var (
	tidyCommand = cli.Command{
		Name:      "tidy",
		Usage:     "tidy the syntax of a provided gscript",
		UsageText: "gscript tidy GSCRIPT",
		Action:    tidyScriptCommand,
	}
)

func tidyScriptCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return errors.New("must supply a gscript to this command")
	}
	cliLogger.Infof("Script: %s", c.Args().First())
	return commandNotImplemented(c)
}
