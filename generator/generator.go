package generator

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
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

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
	file, err := ioutil.ReadFile(path)
	if err != nil {
		g.Logger.Fatalf("Error reading function config: %s", err.Error())
	}

	g.Funcs = g.ExtractFunctionList(file)
}

func (g *Generator) ExtractFunctionList(yData []byte) map[string]*FunctionDef {
	fns := []*FunctionDef{}
	ret := map[string]*FunctionDef{}
	err := yaml.Unmarshal(yData, &fns)
	if err != nil {
		g.Logger.Fatalf("Error parsing YAML: %s", err.Error())
	}
	for _, f := range fns {
		if _, value := ret[f.Name]; value {
			g.Logger.Fatalf("Function declared twice in config: %s", f.Name)
		}
		ret[f.Name] = f
	}
	return ret
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
		tmpfile, _ := ioutil.TempFile("", "")
		tmpfile.Write(buf.Bytes())
		tmpfile.Close()
		g.Logger.Fatalf("Error generating source. Dumped render to %s. Error message: %s", tmpfile.Name(), err.Error())
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
