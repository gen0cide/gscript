package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/gen0cide/gscript"
	"github.com/urfave/cli"
)

func main() {

	var outputFile string
	var compilerOS string
	var compilerArch string
	var outputSource = false

	app := cli.NewApp()
	app.Name = "gscript"
	app.Usage = "Interact with the Genesis Scripting Engine (GSE)"
	app.Version = "0.0.3"
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
				gse := gscript.New("main")
				gse.EnableLogging()
				filename := c.Args().Get(0)
				if len(filename) == 0 {
					gse.LogCritf("You did not supply a filename!")
				}
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					gse.LogCritf("File does not exist: %s", filename)
				}
				_, err := exec.LookPath("jshint")
				if err != nil {
					gse.LogCritf("You do not have jshint in your path. Run: npm install -g jshint")
				}

				jshCmd := exec.Command("jshint", filename)
				jshOutput, err := jshCmd.CombinedOutput()
				if err != nil {
					gse.LogCritf("File Not Valid Javascript!\n -- JSHint Output:\n%s", jshOutput)
				}
				data, err := ioutil.ReadFile(filename)
				gse.SetName(filename)
				gse.CreateVM()
				err = gse.ValidateAST(data)
				if err != nil {
					gse.LogErrorf("Invalid Script Error: %s", err.Error())
				} else {
					gse.LogInfof("Script Valid: %s", filename)
				}
				return nil
			},
		},
		{
			Name:    "shell",
			Aliases: []string{"s"},
			Usage:   "Run an interactive GSE console session.",
			Action: func(c *cli.Context) error {
				gse := gscript.New("shell")
				gse.EnableLogging()
				gse.CreateVM()
				gse.InteractiveSession()
				return nil
			},
		},
		{
			Name:    "compile",
			Aliases: []string{"c"},
			Usage:   "Compile genesis scripts into a stand alone binary.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "outfile",
					Value:       "-",
					Usage:       "Location of the compiled binary (STDOUT if none specified)",
					Destination: &outputFile,
				},
				cli.StringFlag{
					Name:        "os",
					Value:       runtime.GOOS,
					Usage:       "The GOOS you wish to use for your compiled binary.",
					Destination: &compilerOS,
				},
				cli.StringFlag{
					Name:        "arch",
					Value:       runtime.GOARCH,
					Usage:       "The GOARCH you wish to use for your compiled binary.",
					Destination: &compilerArch,
				},
				cli.BoolFlag{
					Name:        "source",
					Usage:       "Do not compile the generated code. Output source instead.",
					Destination: &outputSource,
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					gse := gscript.NewCompiler([]string{}, "", "", "", false)
					gse.Logger.Critf("You did not specify a genesis script!")
				}
				scriptFiles := c.Args()
				if !outputSource && outputFile == "-" {
					outputFile = filepath.Join(os.TempDir(), fmt.Sprintf("%d_genesis.bin", time.Now().Unix()))
				}
				compiler := gscript.NewCompiler(scriptFiles, outputFile, compilerOS, compilerArch, outputSource)
				compiler.Do()
				if !outputSource {
					compiler.Logger.Logf("Your binary is located at: %s", outputFile)
				}
				return nil
			},
		},
		// {
		// 	Name:    "build",
		// 	Aliases: []string{"b"},
		// 	Usage:   "Bundle multiple Genesis scripts and files into a single package.",
		// 	Action: func(c *cli.Context) error {
		// 		return nil
		// 	},
		// },
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run a Genesis script (Careful, don't infect yourself!).",
			Action: func(c *cli.Context) error {
				gse := gscript.New("main")
				gse.EnableLogging()
				filename := c.Args().Get(0)
				if len(filename) == 0 {
					gse.LogCritf("You did not supply a filename!")
				}
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					gse.LogCritf("File does not exist: %s", filename)
				}
				data, err := ioutil.ReadFile(filename)
				gse.SetName(filename)
				gse.CreateVM()
				err = gse.LoadScript(data)
				if err != nil {
					gse.LogErrorf("Script Error: %s", err.Error())
				} else {
					gse.LogInfof("Script loaded successfully")
				}
				err = gse.ExecutePlan()
				if err != nil {
					gse.LogCritf("Hooks Failure: %s", err.Error())
				}
				gse.LogInfof("Hooks executed successfully")
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}
