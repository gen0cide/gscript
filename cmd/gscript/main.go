package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/fatih/color"
	"github.com/gen0cide/gscript"
	"github.com/gen0cide/gscript/compiler/computil"
	"github.com/gen0cide/gscript/logger/standard"
	"github.com/urfave/cli"
)

var (
	defaultCompileOptions = computil.DefaultOptions()
	cliLogger             = standard.NewStandardLogger(nil, "cli", false, false)
	displayBefore         = true
	debugOutput           = false
)

func init() {
	cli.HelpFlag = cli.BoolFlag{Name: "help, h"}
	cli.VersionFlag = cli.BoolFlag{Name: "version, v"}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "Genesis Engine Version: %s\n", gscript.Version)
	}
}

func main() {
	app := cli.NewApp()

	app.Writer = color.Output
	app.ErrWriter = color.Output
	app.Name = "gscript"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "enables verbose debug output",
			Destination: &debugOutput,
		},
	}
	app.Version = gscript.Version
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Alex Levinson",
			Email: "gen0cide.threats@gmail.com",
		},
		cli.Author{
			Name:  "Dan Borges",
			Email: "ahhh.db@gmail.com",
		},
		cli.Author{
			Name:  "Vyrus",
			Email: "vyrus@dc949.org",
		},
	}
	app.Copyright = "(c) 2018 Alex Levinson"
	app.Commands = []cli.Command{
		docsCommand,
		templatesCommand,
		shellCommand,
		linkerCommand,
		vetCommand,
		tidyCommand,
		compileCommand,
	}

	app.Before = func(c *cli.Context) error {
		subcmd := ""
		if debugOutput {
			cliLogger.Logger.SetLevel(logrus.DebugLevel)
		}
		if len(c.Args()) > 0 {
			subcmd = c.Args().Get(0)
		}
		if subcmd != "shell" {
			fmt.Fprintf(c.App.Writer, "%s\n", standard.ASCIILogo())
		}
		return nil
	}

	// ignore error so we don't exit non-zero and break gfmrun README example tests
	err := app.Run(os.Args)
	if err != nil {
		cliLogger.Fatalf("Error Encountered: %v", err)
	}
}

func commandNotImplemented(c *cli.Context) error {
	return fmt.Errorf("%s command not implemented", c.Command.FullName())
}
