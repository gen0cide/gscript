package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli"
)

var (
	templatesSubcommands = []cli.Command{
		{
			Name:   "list",
			Usage:  "show a list of all templates and their descriptions",
			Action: templatesListSubcommand,
		},
		{
			Name:   "show",
			Usage:  "print the named template to standard output",
			Action: templatesPrintSubcommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "template, t",
					Usage: "Show code for `TEMPLATE`",
				},
			},
		},
	}
	templatesCommand = cli.Command{
		Name:        "templates",
		Usage:       "access a library of pre-templated gscripts",
		Subcommands: templatesSubcommands,
	}
)

func templatesListSubcommand(c *cli.Context) error {
	return commandNotImplemented(c)
}

func templatesPrintSubcommand(c *cli.Context) error {
	if c.String("template") == "" {
		return errors.New("must provide a --template/-t NAME argument")
	}
	fmt.Fprintf(c.App.Writer, "Template to print: %s\n", c.String("template"))
	return commandNotImplemented(c)
}
