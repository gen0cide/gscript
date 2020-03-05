package main

import (
	"errors"

	"github.com/urfave/cli/v2"
)

var (
	linkerCommand = &cli.Command{
		Name:      "linker",
		Usage:     "provide information about how compatible functions in a native golang package from a genesis script",
		UsageText: "gscript linker GO_IMPORT_PATH",
		Action:    linkerPackageLookupCommand,
	}
)

func linkerPackageLookupCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return errors.New("must supply a golang package import path as an argument")
	}
	return commandNotImplemented(c)
}
