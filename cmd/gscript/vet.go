package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/robertkrimen/otto/parser"

	"github.com/urfave/cli"
)

var (
	vetCommand = cli.Command{
		Name:      "vet",
		Usage:     "verifies the syntax of a supplied gscripts",
		UsageText: "gscript vet LOCATION",
		Action:    vetScriptCommand,
	}
)

func vetScriptCommand(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return errors.New("must supply a script to this command")
	}
	for _, a := range c.Args() {
		if _, err := os.Stat(a); err != nil {
			cliLogger.Errorf("Error locating file %s:\n  %v", a, err)
			continue
		}
		basename := filepath.Base(a)
		fi, err := os.Open(a)
		if err != nil {
			cliLogger.Errorf("Error loading file %s:\n  %v", basename, err)
			continue
		}
		_, err = parser.ParseFile(nil, a, fi, 0)
		if err != nil {
			fmt.Fprintf(
				color.Output,
				"%s%s%s\n  %s: %s\n  %s: %s\n  %s: %s\n",
				color.HiWhiteString("["),
				color.HiRedString("SYNTAX ERROR"),
				color.HiWhiteString("]"),
				color.HiWhiteString("FILE"),
				color.HiYellowString(a),
				color.HiWhiteString("REASON"),
				color.HiYellowString(err.Error()),
				color.HiWhiteString("STATUS"),
				color.HiRedString("failed"),
			)
		} else {
			fmt.Fprintf(
				color.Output,
				"%s%s%s\n  %s: %s\n  %s: %s\n",
				color.HiWhiteString("["),
				color.HiGreenString("SYNTAX OK"),
				color.HiWhiteString("]"),
				color.HiWhiteString("FILE"),
				color.HiGreenString(a),
				color.HiWhiteString("STATUS"),
				color.HiGreenString("passed"),
			)
		}
	}
	return nil
}
