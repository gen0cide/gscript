package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/gen0cide/gscript"
	"github.com/gen0cide/gscript/compiler/printer"
	"github.com/gen0cide/gscript/logging"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
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

type Generator struct {
	Funcs   []FunctionDef
	Package string
	Logger  *logrus.Logger
}

type FunctionDef struct {
	Name                string   `yaml:"name"`
	Description         string   `yaml:"description"`
	Author              string   `yaml:"author"`
	Package             string   `yaml:"package"`
	ExpectedArgTypes    []ArgDef `yaml:"args"`
	ExpectedReturnTypes []RetDef `yaml:"returns"`
}

type ArgDef struct {
	Name   string `yaml:"label"`
	GoType string `yaml:"gotype"`
}

type RetDef struct {
	Name       string `yaml:"label"`
	GoType     string `yaml:"gotype"`
	ReturnToVM bool   `yaml:"return,omitempty"`
}

func (f *FunctionDef) ReceiverString() string {
	var buf bytes.Buffer
	buf.WriteString(f.Name)
	buf.WriteString("(")
	argNames := []string{}
	for _, a := range f.ExpectedArgTypes {
		argNames = append(argNames, a.Name)
	}
	buf.WriteString(strings.Join(argNames, ", "))
	buf.WriteString(")")
	return buf.String()
}

func (f *FunctionDef) ReturnString() string {
	var buf bytes.Buffer
	argNames := []string{}
	for _, a := range f.ExpectedReturnTypes {
		argNames = append(argNames, a.Name)
	}
	buf.WriteString(strings.Join(argNames, ", "))
	return buf.String()
}

func (g *Generator) ParseYAML(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		g.Logger.Fatalf("Error reading function config: %s", err.Error())
	}

	err = yaml.Unmarshal(file, &g.Funcs)
	if err != nil {
		g.Logger.Fatalf("Error parsing YAML: %s", err.Error())
	}
}

func (g *Generator) BuildSource() bytes.Buffer {
	var buf bytes.Buffer
	tmpl := template.New("generator")
	newTmpl, err := tmpl.Parse(string(MustAsset("templates/vm_functions.go.tmpl")))
	if err != nil {
		g.Logger.Fatalf("Error generating source: %s", err.Error())
	}
	err = newTmpl.Execute(&buf, &g)
	if err != nil {
		g.Logger.Fatalf("Error generating source: %s", err.Error())
	}

	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		g.Logger.Fatalf("Could not parse final Go source: %s", err.Error())
	}

	var finalBuf bytes.Buffer
	finalBuf.Write(formattedCode)
	return finalBuf
}

func (g *Generator) BuildDocs() bytes.Buffer {
	var buf bytes.Buffer
	tmpl := template.New("generator")
	newTmpl, err := tmpl.Parse(string(MustAsset("templates/docs.md.tmpl")))
	if err != nil {
		g.Logger.Fatalf("Error generating source: %s", err.Error())
	}
	err = newTmpl.Execute(&buf, &g)
	if err != nil {
		g.Logger.Fatalf("Error generating source: %s", err.Error())
	}

	return buf
}

func GenFunctions(c *cli.Context) error {
	gen := &Generator{
		Package: outputPackage,
		Funcs:   []FunctionDef{},
		Logger:  logger,
	}

	gen.ParseYAML(defPath)

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

	return nil
}
