package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
	"github.com/gen0cide/gscript"
	"github.com/gen0cide/gscript/compiler/printer"
	"github.com/gen0cide/gscript/generator"
	"github.com/gen0cide/gscript/logging"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	logger        *logrus.Logger
	defPath       string
	outputPath    string
	outputPackage string
	docPath       string
	outputSource  = false
)

func main() {
	cli.AppHelpTemplate = fmt.Sprintf("%s\n\n%s", logging.AsciiLogo(), cli.AppHelpTemplate)
	cli.CommandHelpTemplate = fmt.Sprintf("%s\n\n%s", logging.AsciiLogo(), cli.CommandHelpTemplate)
	app := cli.NewApp()
	app.Writer = color.Output
	app.ErrWriter = color.Output
	app.Name = "gsegen"
	app.Usage = "Generator for VM functions within the gscript SDK."
	app.Version = gscript.Version
	app.Authors = []cli.Author{
		{
			Name:  "Alex Levinson",
			Email: "gen0cide.threats@gmail.com",
		},
	}
	app.Copyright = "(c) 2018 Alex Levinson"

	logger = logrus.New()
	logger.Formatter = &logging.GSEFormatter{}
	logger.Out = logging.LogWriter{Name: "generator"}
	logger.Level = logrus.DebugLevel

	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate the functions off the function map.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "config",
					Value:       "",
					Usage:       "Path to the YAML function definitions",
					Destination: &defPath,
				},
				cli.StringFlag{
					Name:        "package",
					Value:       "engine",
					Usage:       "Golang package you want in the generated file.",
					Destination: &outputPackage,
				},
				cli.StringFlag{
					Name:        "out",
					Value:       "",
					Usage:       "Path to the file of the final golang source.",
					Destination: &outputPath,
				},
				cli.StringFlag{
					Name:        "docs",
					Value:       "",
					Usage:       "Path to the markdown docs.",
					Destination: &docPath,
				},
				cli.BoolFlag{
					Name:        "source",
					Usage:       "Do not write the generated code to a file. Output source instead.",
					Destination: &outputSource,
				},
			},
			Action: GenFunctions,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func GenFunctions(c *cli.Context) error {
	gen := &generator.Generator{
		Package: outputPackage,
		Funcs:   map[string]*generator.FunctionDef{},
		Logger:  logger,
		Libs:    map[string]*generator.Library{},
		Config:  defPath,
	}

	gen.ParseYAML(defPath)
	gen.ParseLibs()

	for name, l := range gen.Libs {
		gen.Logger.Infof("Modifying lib_%s.go with comments.", name)
		l.WriteModifiedSource()
	}

	buf := gen.BuildSource()

	if outputSource {
		printer.PrettyPrintSource(buf.String())
		return nil
	}

	err := ioutil.WriteFile(outputPath, buf.Bytes(), 0644)
	if err != nil {
		logger.Fatalf("Could not write final code: %s", err.Error())
	}

	if docPath != "" {
		docBuf := gen.BuildDocs()
		err := ioutil.WriteFile(docPath, docBuf.Bytes(), 0644)
		if err != nil {
			logger.Fatalf("Could not write final docs: %s", err.Error())
		}
	}

	funcFile, err := ioutil.ReadFile(defPath)
	if err != nil {
		logger.Fatalf("Could not read the function list: %s", err.Error())
	}
	absPath, err := filepath.Abs(filepath.Dir(defPath))
	if err != nil {
		logger.Fatalf("Could not get the absolute path to the function file!")
	}
	newFile := filepath.Join(absPath, "templates", "functions.yml")
	err = ioutil.WriteFile(newFile, funcFile, 0644)
	if err != nil {
		logger.Fatalf("Could not write function list: %s", err.Error())
	}
	return nil
}
