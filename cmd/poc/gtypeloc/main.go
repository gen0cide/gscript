package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/gen0cide/waiter"
	"github.com/uudashr/gopkgs"

	"github.com/gen0cide/gscript/compiler"

	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"

	"github.com/gen0cide/gscript/compiler/computil"
)

// GoPirate defines a type searching a given go package for type declarations
type GoPirate struct {
	sync.RWMutex
	Dir         string
	Pkg         *gopkgs.Pkg
	FilesByName map[string]*GoFile
	AST         *ast.Package
	Waiter      *waiter.Waiter
	Log         *logrus.Entry
}

// GoStructDef defines a struct declaration within Golang source code
type GoStructDef struct {
	sync.RWMutex
	File         *ast.File
	TypeSpec     *ast.TypeSpec
	AST          *ast.StructType
	Name         string
	Fields       map[string]*compiler.GoParamDef
	Incompatible map[string]*compiler.GoParamDef
	Embeds       map[string]*compiler.GoParamDef
}

// GoInterfaceDef defines an interface declaration within Golang source code
type GoInterfaceDef struct {
	sync.RWMutex
	File     *GoFile
	TypeSpec *ast.TypeSpec
	AST      *ast.InterfaceType
	Name     string
	Methods  []*GoInterfaceFunc
}

// GoInterfaceFunc defines a method within a golang interface declaration
type GoInterfaceFunc struct {
	InterfaceDef *GoInterfaceDef
	Name         string
	Args         []*compiler.GoParamDef
	Returns      []*compiler.GoParamDef
}

// GoFile describes a golang source file
type GoFile struct {
	sync.RWMutex
	Filename   string
	Waiter     *waiter.Waiter
	Log        *logrus.Entry
	AST        *ast.File
	Parent     *GoPirate
	Structs    map[string]*GoStructDef
	Interfaces map[string]*GoInterfaceDef
}

// Walk walks a golang AST file looking for discoverable types
func (g *GoFile) Walk() {
	defer g.Waiter.Done()
	ast.Inspect(g.AST, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.TypeSpec:
			if !v.Name.IsExported() {
				return true
			}
			if v.Name.Name == "Header" {
				spew.Dump(v)
			}
			g.Waiter.Add(1)
			go g.ParseTypeSpec(v)
		}
		return true
	})
	return
}

// ParseTypeSpec attempts to parse a type's specification against a given AST object
func (g *GoFile) ParseTypeSpec(v *ast.TypeSpec) {
	defer g.Waiter.Done()
	switch t := v.Type.(type) {
	case *ast.StructType:
		g.Waiter.Add(1)
		go g.ParseStructDef(v, t)
	case *ast.InterfaceType:
		g.Waiter.Add(1)
		go g.ParseInterfaceDef(v, t)
	}
	return
}

// WalkInterface attempts to walk the interface methods out of a golang interface declaration
func (g *GoFile) WalkInterface(i *GoInterfaceDef) {
	defer g.Waiter.Done()
	if i.AST.Methods == nil {
		return
	}

}

// WalkStruct attempts to walk the struct fields out of a golang struct declaration
func (g *GoFile) WalkStruct(s *GoStructDef) {
	defer g.Waiter.Done()
	if s.AST.Fields == nil {
		return
	}
	for fidx, tf := range s.AST.Fields.List {
		if len(tf.Names) == 0 {
			// struct embed
			g.Waiter.Add(1)
			go g.ParseStructField(s, tf, "_EMBED", fidx)
		}
		for noff, name := range tf.Names {
			if !name.IsExported() {
				continue
			}
			g.Waiter.Add(1)
			go g.ParseStructField(s, tf, name.Name, (noff + fidx))
		}
	}
}

// ParseStructField attempts to resolve the type of a given struct
func (g *GoFile) ParseStructField(s *GoStructDef, f *ast.Field, name string, offset int) {
	defer g.Waiter.Done()
	p := &compiler.GoParamDef{
		ParamIdx:       offset,
		SkipResolution: true,
		Type:           "struct_field",
		OriginalName:   name,
	}
	err := p.Interpret(f.Type)
	if err != nil {
		p.NameBuffer.Truncate(0)
		p.NameBuffer.WriteString(err.Error())
		if name == "_EMBED" {
			name = fmt.Sprintf("EMBED_%d", offset)
		}
		s.Lock()
		s.Incompatible[name] = p
		s.Unlock()
		return
	}
	if name == "_EMBED" {
		s.Lock()
		s.Embeds[p.SigBuffer.String()] = p
		s.Unlock()
		return
	}
	s.Lock()
	s.Fields[name] = p
	s.Unlock()
	return
}

// ParseInterfaceDef attempts to parse the interface definition from source
func (g *GoFile) ParseInterfaceDef(s *ast.TypeSpec, v *ast.InterfaceType) {
	defer g.Waiter.Done()
	g.Lock()
	defer g.Unlock()
	if _, ok := g.Interfaces[s.Name.Name]; ok {
		g.Log.Errorf("Interface already defined for %s", s.Name.Name)
		return
	}
	gid := &GoInterfaceDef{
		File:     g,
		TypeSpec: s,
		AST:      v,
		Name:     s.Name.Name,
		Methods:  []*GoInterfaceFunc{},
	}
	g.Interfaces[s.Name.Name] = gid
	g.Waiter.Add(1)
	go g.WalkInterface(gid)
	return
}

// ParseStructDef attempts to parse a type definition from source
func (g *GoFile) ParseStructDef(s *ast.TypeSpec, v *ast.StructType) {
	defer g.Waiter.Done()
	g.Lock()
	defer g.Unlock()
	if _, ok := g.Structs[s.Name.Name]; ok {
		g.Log.Errorf("Struct already defined for %s", s.Name.Name)
		return
	}
	gsd := &GoStructDef{
		TypeSpec:     s,
		AST:          v,
		Name:         s.Name.Name,
		Fields:       map[string]*compiler.GoParamDef{},
		Incompatible: map[string]*compiler.GoParamDef{},
		Embeds:       map[string]*compiler.GoParamDef{},
	}
	g.Structs[s.Name.Name] = gsd
	g.Waiter.Add(1)
	go g.WalkStruct(gsd)
	return
}

// NewFile adds a new file to the index
func (g *GoPirate) NewFile(name string, f *ast.File) (*GoFile, error) {
	g.Lock()
	defer g.Unlock()
	if file, found := g.FilesByName[name]; found {
		g.Log.Errorf("Already indexed file %s", name)
		return file, errors.New("already exists")
	}
	file := &GoFile{
		Filename:   name,
		Waiter:     g.Waiter,
		Log:        g.Log,
		AST:        f,
		Parent:     g,
		Structs:    map[string]*GoStructDef{},
		Interfaces: map[string]*GoInterfaceDef{},
	}
	g.FilesByName[name] = file
	return file, nil
}

// PrintResults prints the results of the search to the console
func (g *GoPirate) PrintResults() {
	structs := map[string]*GoStructDef{}
	interfaces := map[string]*GoInterfaceDef{}
	for _, f := range g.FilesByName {
		for sname, s := range f.Structs {
			teststruct, ok := structs[sname]
			_ = teststruct
			if ok {
				//g.Log.Errorf("Duplicate Struct Def: %s is declared in %s, ignoring declaration in %s", sname, teststruct.File.Filename, s.File.Filename)
			}
			structs[sname] = s
		}
		for iname, i := range f.Interfaces {
			testinterface, ok := interfaces[iname]
			_ = testinterface
			if ok {
				//g.Log.Errorf("Duplicate Interface Def: %s is declared in %s, ignoring declaration in %s", iname, testinterface.File.Filename, i.File.Filename)
			}
			interfaces[iname] = i
		}
	}
	// g.Log.Infof("=== INTERFACES ===")
	// for n, i := range interfaces {
	// 	g.Log.Infof("  %s (len=%d)", n, len(i.Methods))
	// }
	g.Log.Infof("===  STRUCTS  ===")
	for n, s := range structs {
		g.Log.Infof("  %s (fields=%d, incompatibles=%d, embeds=%d)", n, len(s.Fields), len(s.Incompatible), len(s.Embeds))
		for fn, f := range s.Fields {
			g.Log.Infof("    (F) %s (%s)", fn, f.SigBuffer.String())
		}
		for fn, f := range s.Incompatible {
			g.Log.Infof("    (I) %s (%s)", fn, f.NameBuffer.String())
		}
		for _, f := range s.Embeds {
			g.Log.Infof("    (E) %s", f.SigBuffer.String())
		}
	}
}

func main() {
	log := logrus.New()
	wtr := waiter.New("typesearch", log.Out)
	log.Out = wtr
	formatter := &prefixed.TextFormatter{
		ForceFormatting: true,
		ForceColors:     true,
	}
	log.Formatter = formatter
	log.SetLevel(logrus.DebugLevel)
	fs := token.NewFileSet()
	pkg, err := computil.ResolveGlobalImport(os.Args[1])
	if err != nil {
		log.Error(err)
		log.Fatal("could not find mentioned package")
	}

	pkgs, err := parser.ParseDir(fs, pkg.Dir, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	var astpkg *ast.Package
	astpkg, ok := pkgs[pkg.Name]
	if !ok {
		log.Fatal("could not find mentioned package in package AST")
	}

	pirate := &GoPirate{
		Dir:         pkg.Dir,
		Pkg:         pkg,
		FilesByName: map[string]*GoFile{},
		AST:         astpkg,
		Waiter:      wtr,
		Log:         log.WithField("prefix", "pirate"),
	}

	for filename, file := range astpkg.Files {
		gf, err := pirate.NewFile(filename, file)
		if err != nil {
			continue
		}
		pirate.Waiter.Add(1)
		go gf.Walk()
	}

	pirate.Waiter.Wait()

	pirate.PrintResults()
}
