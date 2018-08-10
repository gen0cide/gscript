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
	cliLogger             = standard.NewStandardLogger(nil, "gscript", "cli", false, false)
	displayBefore         = true
	debugOutput           = false
)

func init() {
	cli.HelpFlag = cli.BoolFlag{Name: "help, h"}
	cli.VersionFlag = cli.BoolFlag{Name: "version"}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%s\n", gscript.Version)
	}
}

func main() {
	app := cli.NewApp()

	app.Writer = color.Output
	app.ErrWriter = color.Output

	cli.AppHelpTemplate = fmt.Sprintf("%s\n%s", standard.ASCIILogo(), cli.AppHelpTemplate)
	app.Name = "gscript"
	app.Usage = "Cross platform dropper framework"
	app.Description = "Framework to rapidly implement custom droppers for all three major operating systems."

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
		cli.Author{
			Name:  "Lucas Morris",
			Email: "emperorcow@gmail.com",
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
		if debugOutput {
			cliLogger.Logger.SetLevel(logrus.DebugLevel)
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
