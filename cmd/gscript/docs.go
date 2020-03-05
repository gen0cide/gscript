package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	docsSubcommands = []*cli.Command{
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
	docsCommand = &cli.Command{
		Name:        "docs",
		Usage:       "Shows documentation on a variety of gscript topics",
		Subcommands: docsSubcommands,
	}

	macroHelp = fmt.Sprintf(
		"For more information on the available macros, please check the README at:\n  %s",
		color.HiGreenString("https://github.com/gen0cide/gscript"),
	)
	scriptHelp = fmt.Sprintf(
		"For more information on the scripts, please refer to:\n  %s: %s\n  %s: %s",
		color.HiWhiteString("README"),
		color.HiGreenString("https://github.com/gen0cide/gscript"),
		color.HiWhiteString("EXAMPLES"),
		color.HiGreenString("https://github.com/ahhh/gscripts"),
	)
	stdlibHelp = fmt.Sprintf(
		"For more information on the standard library, please refer to:\n  %s: %s\n  %s: %s\n  %s: %s",
		color.HiWhiteString("README"),
		color.HiGreenString("https://github.com/gen0cide/gscript"),
		color.HiWhiteString("EXAMPLES"),
		color.HiGreenString("https://github.com/ahhh/gscripts"),
		color.HiWhiteString("INTERACTIVE"),
		color.HiCyanString("Run \"gscript shell\" and enter the command \"SymbolTable()\""),
	)
)

func docsMacrosSubcommand(c *cli.Context) error {
	cliLogger.Infoln(macroHelp)
	return nil
}

func docsScriptsSubcommand(c *cli.Context) error {
	cliLogger.Infoln(scriptHelp)
	return nil
}

func docsStandardLibSubcommand(c *cli.Context) error {
	cliLogger.Infoln(stdlibHelp)
	return nil
}
