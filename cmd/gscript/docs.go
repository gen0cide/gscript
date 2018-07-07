package main

import "github.com/urfave/cli"

var (
	packageName     = ""
	docsSubcommands = []cli.Command{
		{
			Name:   "macros",
			Usage:  "shows macros available to the gscript compiler",
			Action: docsMacrosSubcommand,
		},
		{
			Name:   "scripts",
			Usage:  "shows an overview of genesis scripting and entry points to use",
			Action: docsScriptsSubcommand,
		},
		{
			Name:   "stdlib",
			Usage:  "prints a reference to the genesis standard library and each library's commands",
			Action: docsStandardLibSubcommand,
		},
	}
	docsCommand = cli.Command{
		Name:        "docs",
		Usage:       "Shows documentation on a variety of gscript topics",
		Subcommands: docsSubcommands,
	}
)

func docsMacrosSubcommand(c *cli.Context) error {
	return commandNotImplemented(c)
}

func docsScriptsSubcommand(c *cli.Context) error {
	return commandNotImplemented(c)
}

func docsStandardLibSubcommand(c *cli.Context) error {
	return commandNotImplemented(c)
}
