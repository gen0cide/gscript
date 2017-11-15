package main

import (
	"os"
	"sort"

	"github.com/gen0cide/gscript"
	"github.com/urfave/cli"
)

// func main() {

// 	a := gscript.New()
// 	a.EnableLogging()
// 	a.CreateVM()
// 	a.VM.Run(gscript.DefaultScript)

// }

func main() {
	app := cli.NewApp()
	app.Name = "gscript"
	app.Usage = "Interact with the Genesis Scripting Engine (GSE)"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Alex Levinson",
			Email: "gen0cide.threats@gmail.com",
		},
	}
	app.Copyright = "(c) 2017 Alex Levinson"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Run gscript in debug mode.",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Suppress all logging output.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Check a GSE script for syntax errors.",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "shell",
			Aliases: []string{"s"},
			Usage:   "Run an interactive GSE console session.",
			Action: func(c *cli.Context) error {
				a := gscript.New()
				a.CreateVM()
				a.InteractiveSession()
				return nil
			},
		},
		{
			Name:    "compile",
			Aliases: []string{"c"},
			Usage:   "Compile a Genesis script into a stand alone binary.",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "Bundle multiple Genesis scripts and files into a single package.",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run a Genesis script (Careful, don't infect yourself!).",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}
