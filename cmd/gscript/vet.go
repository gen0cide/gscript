package main

import (
	"errors"

	"github.com/urfave/cli"
)

var (
	vetCommand = cli.Command{
		Name:      "vet",
		Usage:     "verifies the syntax of a supplied gscript",
		UsageText: "gscript vet GSCRIPT",
		Action:    vetScriptCommand,
	}
)

func vetScriptCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return errors.New("must supply a gscript to this command")
	}
	cliLogger.Infof("Script: %s", c.Args().First())
	return commandNotImplemented(c)
}
