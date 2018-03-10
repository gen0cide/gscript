package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
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
	Package string
	Logger  *logrus.Logger
	Libs    map[string]*Library
	Funcs   map[string]*FunctionDef
	Config  string
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

type Library struct {
	Name   string                  `yaml:"package"`
	Path   string                  `yaml:"path"`
	FSet   *token.FileSet          `yaml:"-"`
	AST    *ast.File               `yaml:"-"`
	Logger *logrus.Logger          `yaml:"-"`
	Funcs  map[string]*FunctionDef `yaml:"-"`
	Source []byte                  `yaml:"-"`
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

func (g *Generator) ResolveLibs() []string {
	libs := []string{}
	keys := make(map[string]bool)
	for _, f := range g.Funcs {
		if _, value := keys[f.Package]; !value {
			keys[f.Package] = true
			libs = append(libs, f.Package)
		}
	}
	return libs
}

func (g *Generator) ParseLibs() {
	libs := g.ResolveLibs()
	baseDir := filepath.Dir(g.Config)
	for _, name := range libs {
		libPath := filepath.Join(baseDir, "engine", fmt.Sprintf("lib_%s.go", name))
		if _, err := os.Stat(libPath); os.IsNotExist(err) {
			g.Logger.Fatalf("Package %s could not be found! Expecting to find %s", name, libPath)
		}
		lib := &Library{
			Name:   name,
			FSet:   token.NewFileSet(),
			AST:    nil,
			Path:   libPath,
			Logger: g.Logger,
			Funcs:  map[string]*FunctionDef{},
		}
		if _, value := g.Libs[name]; !value {
			g.Libs[name] = lib
			g.Logger.Infof("Discovered Library: %s", libPath)
			lib.ParseLibrary(g)
		}
	}
}

func (l *Library) ParseLibrary(g *Generator) {
	a, err := parser.ParseFile(l.FSet, l.Path, nil, parser.DeclarationErrors|parser.ParseComments)
	if err != nil {
		l.Logger.Fatalf("Error parsing %s library: %s", l.Name, err.Error())
	}
	l.AST = a
	cm := ast.NewCommentMap(l.FSet, l.AST, l.AST.Comments)

	if cm == nil {
		cm = ast.CommentMap{}
	}
	for _, d := range l.AST.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			if _, exists := g.Funcs[d.Name.Name]; exists {
				l.Logger.Infof("Found Function: %s", d.Name)
				l.Funcs[d.Name.Name] = g.Funcs[d.Name.Name]
				cmLines, err := g.Funcs[d.Name.Name].BuildGodoc()
				if err != nil {
					l.Logger.Fatalf("Error generating comments: %s", err.Error())
				}
				var newDoc []*ast.Comment
				for _, line := range cmLines {
					if line == "" {
						continue
					}
					cmLine := &ast.Comment{
						Text:  line,
						Slash: d.Pos() - 1,
					}
					newDoc = append(newDoc, cmLine)
					l.FSet.File(cmLine.End()).AddLine(int(cmLine.End()))
				}
				cmGroup := &ast.CommentGroup{
					List: newDoc,
				}
				d.Doc = cmGroup
				cm[d] = []*ast.CommentGroup{cmGroup}
			}
		}
	}
	l.AST.Comments = cm.Comments()
	var b bytes.Buffer
	format.Node(&b, l.FSet, l.AST)
	src, err := format.Source(b.Bytes())
	if err != nil {
		l.Logger.Fatalf("Error formatting library: %s", err.Error())
	}
	l.Source = src
}

func (l *Library) WriteModifiedSource() {
	f, err := os.OpenFile(l.Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		l.Logger.Fatalf("Failed to write modified source: %s", err.Error())
	}
	defer f.Close()
	f.Write(l.Source)
}

func (g *Generator) ParseYAML(path string) {
	fns := []*FunctionDef{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		g.Logger.Fatalf("Error reading function config: %s", err.Error())
	}

	err = yaml.Unmarshal(file, &fns)
	if err != nil {
		g.Logger.Fatalf("Error parsing YAML: %s", err.Error())
	}
	for _, f := range fns {
		if _, value := g.Funcs[f.Name]; value {
			g.Logger.Fatalf("Function declared twice in config: %s", f.Name)
		}
		g.Funcs[f.Name] = f
	}
}

func (f *FunctionDef) BuildGodoc() ([]string, error) {
	var buf bytes.Buffer
	tmpl := template.New("commentgenerator")
	newTmpl, err := tmpl.Parse(string(MustAsset("templates/comment.go.tmpl")))
	if err != nil {
		return nil, errors.New("could not parse comment template")
	}
	err = newTmpl.Execute(&buf, &f)
	if err != nil {
		return nil, fmt.Errorf("could not generate comment template for %s", f.Name)
	}
	lines := strings.Split(buf.String(), "\n")
	return lines, nil
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
		Funcs:   map[string]*FunctionDef{},
		Logger:  logger,
		Libs:    map[string]*Library{},
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

	return nil
}
