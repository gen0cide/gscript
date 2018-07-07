package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	gast "github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/file"
	gparser "github.com/robertkrimen/otto/parser"
	"github.com/uudashr/gopkgs"
)

var (
	macros = []string{
		"priority",
		"os",
		"go_import",
		"timeout",
		"arch",
	}

	logger = logrus.New()

	basicComment    = regexp.MustCompile(`(?P<key>\S*?):(?P<value>\S*)$`)
	goImportComment = regexp.MustCompile(`(?P<key>\S*?):(?P<gopkg>\S*?) as (?P<namespace>\w*)`)

	BuiltInGoTypes = map[string]bool{
		"bool":       true,
		"byte":       true,
		"complex128": true,
		"complex64":  true,
		"error":      true,
		"float32":    true,
		"float64":    true,
		"int":        true,
		"int8":       true,
		"int16":      true,
		"int32":      true,
		"int64":      true,
		"rune":       true,
		"string":     true,
		"uint":       true,
		"uint8":      true,
		"uint16":     true,
		"uint32":     true,
		"uint64":     true,
		"uintptr":    true,
	}
)

//order of operations
//create script object
//parse script AST
//walk comments looking for macros
//set macros as appropriate
//walk ast for known go_import namespace
//Create GoPackage objects and populate with callers
//gather local golang packages
//for each package matching go_import macros
//	if package could not be found locally
//		go get -u $go_package_url
//rescan local golang packages (if needed)
//diff the golang import path with only the packages defined in namespace
//iterate the GoImports struct within Script
//for each go package
//	use golang parser to parse package dir
//  for each file in the package:
//		foreach function in the file:
//			if function name is a function in our Callers map
//				Create a swizzler and set Function name
//        retrieve file imports
//        set GoFunc in the gopackage
//  for each file with a compatible function:
//		populate the GoPackage FileImports map with the array of imports
//		sanity check that there will not be any collisions and no global imports
//	foreach swizzler
//		error if the swizzler call has a method receiver
//    error if the caller has a delta in # of arguments
//    parse golang function parameters
//    parse golang function returns
//

//parsing golang function arguments:
// return if argument list is of size 0
// attempt to build the ParamDef Buffers
// if a selectortype is found, check the file

type Script struct {
	Name   string
	Data   []byte
	AST    *gast.Program
	Macros []*Macro
	GoMap  map[string]*GoPackage // key = golang import path, value = namespace
	NSMap  map[string]*GoPackage // key = genesis namespace
}

type FunctionCall struct {
	Namespace    string
	FuncName     string
	ArgumentList []gast.Expression
}

type GoPackage struct {
	Script        *Script
	Namespace     string
	ImportKey     string
	Callers       map[string]*FunctionCall     // key = function name, val = gast decl
	FileImports   map[string][]*ast.ImportSpec // key = filename, val = file imports
	ImportMapping map[string]*ast.ImportSpec   // key = import name, val = ast.ImportSpec
	FuncToFileMap map[string]string            // key = function name, val = filename
	SwizzledFuncs map[string]*Swizzler         // key = function name, val = swizzler
	GoFuncs       map[string]*ast.FuncDecl     // key = function name, val = golang decl
	Dir           string
	ImportPath    string
	Name          string
}

type Swizzler struct {
	Function  string
	Caller    *FunctionCall
	File      *ast.File
	GoDecl    *ast.FuncDecl
	Imports   []*ast.ImportSpec
	GoPackage *GoPackage
	GoArgs    []*ParamDef
	GoReturns []*ParamDef
	VMShim    string
}

type ParamDef struct {
	SigBuffer  bytes.Buffer
	NameBuffer bytes.Buffer
	Swizzler   *Swizzler
	ImportRefs map[string]*ast.ImportSpec
	VarName    string
	ParamIdx   int
	ArgOffset  int
	ExtSig     string
	GoLabel    string
}

type Macro struct {
	Key    string
	Params map[string]string
}

type gWalker struct {
	script *Script
	source string
	shift  file.Idx
}

func (s *Script) WalkGASTForGoCalls() {
	w := &gWalker{source: string(s.Data), script: s}
	gast.Walk(w, s.AST)
}

func (w *gWalker) Enter(n gast.Node) gast.Visitor {
	switch v1 := n.(type) {
	case *gast.CallExpression:
		switch v2 := v1.Callee.(type) {
		case *gast.DotExpression:
			namespace, ok := v2.Left.(*gast.Identifier)
			if !ok {
				return w
			}
			if _, ok := w.script.NSMap[namespace.Name]; !ok {
				return w
			}
			funcName := v2.Identifier.Name
			gop := w.script.NSMap[namespace.Name]

			gop.Callers[funcName] = &FunctionCall{
				Namespace:    namespace.Name,
				FuncName:     funcName,
				ArgumentList: v1.ArgumentList,
			}
		}
	}
	return w
}

func (w *gWalker) Exit(n gast.Node) {
	return
}

func (s *Script) ScanForMacros() {
	for _, comments := range s.AST.Comments {
		for _, comment := range comments {
			if basicComment.MatchString(comment.Text) {
				n1 := basicComment.SubexpNames()
				r2 := basicComment.FindAllStringSubmatch(comment.Text, -1)[0]
				md := map[string]string{}
				m1 := &Macro{}
				for i, n := range r2 {
					if i == 0 {
						continue
					} else if i == 1 {
						m1.Key = n
						continue
					}
					md[n1[i]] = n
				}
				m1.Params = md
				s.Macros = append(s.Macros, m1)
			} else if goImportComment.MatchString(comment.Text) {
				n1 := goImportComment.SubexpNames()
				r2 := goImportComment.FindAllStringSubmatch(comment.Text, -1)[0]
				md := map[string]string{}
				m1 := &Macro{}
				for i, n := range r2 {
					if i == 0 {
						continue
					} else if i == 1 {
						m1.Key = n
						continue
					}
					md[n1[i]] = n
				}
				m1.Params = md
				s.Macros = append(s.Macros, m1)
			}
		}
	}
}

func (s *Script) InitializeGoImports() {
	for _, m := range s.Macros {
		if m.Key != "go_import" {
			continue
		}
		gop := &GoPackage{
			Script:        s,
			Namespace:     m.Params["namespace"],
			ImportKey:     m.Params["gopkg"],
			Callers:       map[string]*FunctionCall{},
			FileImports:   map[string][]*ast.ImportSpec{},
			ImportMapping: map[string]*ast.ImportSpec{},
			FuncToFileMap: map[string]string{},
			SwizzledFuncs: map[string]*Swizzler{},
			GoFuncs:       map[string]*ast.FuncDecl{},
		}
		s.GoMap[m.Params["gopkg"]] = gop
		s.NSMap[m.Params["namespace"]] = gop
		logger.Infof("Found Go Import Macro: %s as %s", m.Params["gopkg"], m.Params["namespace"])
	}
}

func (s *Script) LocateGolangDependencies() {
	gopks, err := gopkgs.Packages(gopkgs.Options{NoVendor: true})
	if err != nil {
		panic(err)
	}

	for _, gopkg := range gopks {
		if gop, ok := s.GoMap[gopkg.ImportPath]; ok {
			logger.Infof("Located golang directory for import %s", gopkg.Name)
			gop.Dir = gopkg.Dir
			gop.ImportPath = gopkg.ImportPath
			gop.Name = gopkg.Name
		}
	}
}

func (s *Script) ResolveGolangASTCalls() {
	for _, gop := range s.GoMap {
		fs := token.NewFileSet()
		pkgs, err := parser.ParseDir(fs, gop.Dir, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		if _, ok := pkgs[gop.Name]; !ok {
			panic(fmt.Errorf("should have found golang package %s but didnt", gop.ImportPath))
		}
		for _, file := range pkgs[gop.Name].Files {
			ast.Inspect(file, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if ok {
					if fn.Name.IsExported() && gop.Callers[fn.Name.Name] != nil {
						logger.Infof("Resolved Golang and Genesis caller for function %s", fn.Name.Name)
						if gop.FileImports[file.Name.Name] == nil {
							gop.FileImports[file.Name.Name] = file.Imports
						}
						gop.GoFuncs[fn.Name.Name] = fn
						gop.FuncToFileMap[fn.Name.Name] = file.Name.Name
						gop.SwizzledFuncs[fn.Name.Name] = &Swizzler{
							Function:  fn.Name.Name,
							Caller:    gop.Callers[fn.Name.Name],
							File:      file,
							GoDecl:    fn,
							Imports:   file.Imports,
							GoPackage: gop,
							GoArgs:    []*ParamDef{},
							GoReturns: []*ParamDef{},
						}
					}
				}
				return true
			})
		}
	}
}

func (s *Script) LinkGolangASTWithGAST() {
	for _, gop := range s.GoMap {
		gop.LinkASTFunctions()
	}
}

func (g *GoPackage) LinkASTFunctions() {
	for _, swiz := range g.SwizzledFuncs {
		logger.Infof("Swizzling for function %s", swiz.Function)
		err := swiz.SwizzleAllTheThings()
		if err != nil {
			panic(err)
		}
	}
}

func (s *Swizzler) CanResolveImportDep(pkg string) bool {
	if pkg == "." {
		panic(fmt.Errorf("should not attempt to import anonymously in package %s", s.File.Name.Name))
	}
	for _, i := range s.Imports {
		if i.Name != nil {
			if i.Name.Name == pkg {
				return true
			}
		} else {
			pkgParts := strings.Split(i.Path.Value, "/")
			packageAlias := pkgParts[len(pkgParts)-1]
			newAlias := strings.Replace(packageAlias, `"`, ``, -1)
			if newAlias == pkg {
				return true
			}
		}
	}
	spew.Dump(s.Imports)
	logger.Errorf("could not resolve package %s", pkg)
	return false
}

func (s *Swizzler) SwizzleAllTheThings() error {
	if s.GoDecl.Recv != nil {
		return fmt.Errorf("golang method %s in package %s cannot have a receiver", s.Function, s.GoPackage.ImportPath)
	}
	err := s.SwizzleToTheLeft()
	if err != nil {
		return err
	}
	err = s.SwizzleToTheRight()
	if err != nil {
		return err
	}
	return nil
}

func (s *Swizzler) SwizzleToTheLeft() error {
	var aOff = 0
	for idx, p := range s.GoDecl.Type.Params.List {
		masterP := &ParamDef{
			Swizzler:   s,
			ImportRefs: map[string]*ast.ImportSpec{},
			ParamIdx:   idx,
		}
		masterP.NameBuffer.WriteString("__")
		err := masterP.Interpret(p.Type)
		if err != nil {
			panic(err)
		}
		masterP.VarName = masterP.NameBuffer.String()
		masterP.ExtSig = masterP.SigBuffer.String()
		for i := 0; i < len(p.Names); i++ {
			newP := &ParamDef{
				Swizzler:   s,
				ImportRefs: map[string]*ast.ImportSpec{},
				VarName:    fmt.Sprintf("%s%d", masterP.VarName, aOff),
				ParamIdx:   idx,
				ArgOffset:  aOff,
				ExtSig:     masterP.ExtSig,
				GoLabel:    p.Names[i].Name,
			}
			aOff++
			s.GoArgs = append(s.GoArgs, newP)
			logger.Infof("In function %s, parameter %s was swizzled into %s and %s", s.Function, newP.GoLabel, newP.ExtSig, newP.VarName)
		}
	}
	return nil
}

func (s *Swizzler) SwizzleToTheRight() error {
	return nil
}

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Println("~~~ Genesis Scripting Engine AST Compiler ~~~")
	fmt.Println("         \"...here be dragons...\"           \n")

	fmt.Println("===== INPUT SCRIPT =====")
	fmt.Println(string(data))
	fmt.Println("========================\n")

	program, err := gparser.ParseFile(nil, os.Args[1], data, 2)

	script := &Script{
		Name:   os.Args[1],
		AST:    program,
		Macros: []*Macro{},
		Data:   data,
		GoMap:  map[string]*GoPackage{},
		NSMap:  map[string]*GoPackage{},
	}

	script.ScanForMacros()
	script.InitializeGoImports()
	script.WalkGASTForGoCalls()
	script.LocateGolangDependencies()
	script.ResolveGolangASTCalls()
	script.LinkGolangASTWithGAST()
}

func (p *ParamDef) Interpret(i interface{}) error {
	switch t := i.(type) {
	case *ast.StarExpr:
		return p.ParseStarExpr(t)
	case *ast.SelectorExpr:
		return p.ParseSelectorExpr(t)
	case *ast.Ident:
		return p.ParseIdent(t)
	case *ast.ArrayType:
		return p.ParseArrayType(t)
	case *ast.MapType:
		return p.ParseMapType(t)
	case *ast.ChanType:
		return fmt.Errorf("function %s includes an unsupported parameter type: %s", p.Swizzler.GoDecl.Name.Name, "chan")
	case *ast.FuncType:
		return fmt.Errorf("function %s includes an unsupported parameter type: %s", p.Swizzler.GoDecl.Name.Name, "func")
	case *ast.InterfaceType:
		return fmt.Errorf("function %s includes an unsupported parameter type: %s", p.Swizzler.GoDecl.Name.Name, "interface{}")
	case *ast.StructType:
		return fmt.Errorf("function %s includes an unsupported parameter type: %s", p.Swizzler.GoDecl.Name.Name, "struct")
	default:
		valType := reflect.ValueOf(t)
		return fmt.Errorf("could not determine the golang ast type of %s in func %s", valType.Type().String, p.Swizzler.Function)
	}
}

func (p *ParamDef) ParseMapType(a *ast.MapType) error {
	p.SigBuffer.WriteString("map[")
	p.NameBuffer.WriteString("MapOf")
	err := p.Interpret(a.Key)
	if err != nil {
		return err
	}
	p.SigBuffer.WriteString("]")
	p.NameBuffer.WriteString("WithValType")
	err = p.Interpret(a.Value)
	return err
}

func (p *ParamDef) ParseArrayType(a *ast.ArrayType) error {
	p.SigBuffer.WriteString("[]")
	p.NameBuffer.WriteString("ArrayOf")
	return p.Interpret(a.Elt)
}

func (p *ParamDef) ParseSelectorExpr(a *ast.SelectorExpr) error {
	x, ok := a.X.(*ast.Ident)
	if !ok {
		return fmt.Errorf("could not parse selector namespace in func %s", p.Swizzler.Function)
	}
	if p.Swizzler.CanResolveImportDep(x.Name) == false {
		return fmt.Errorf("the package %s was not found in the import map", x.Name)
	}
	p.SigBuffer.WriteString(x.Name)
	p.NameBuffer.WriteString(x.Name)
	p.SigBuffer.WriteString(".")
	p.NameBuffer.WriteString("_")
	p.SigBuffer.WriteString(a.Sel.Name)
	p.NameBuffer.WriteString(a.Sel.Name)
	return nil
}

func (p *ParamDef) ParseStarExpr(a *ast.StarExpr) error {
	p.SigBuffer.WriteString("*")
	p.NameBuffer.WriteString("PointerTo")
	return p.Interpret(a.X)
}

func (p *ParamDef) ParseIdent(a *ast.Ident) error {
	if ok := BuiltInGoTypes[a.Name]; !ok {
		p.SigBuffer.WriteString(p.Swizzler.GoPackage.Name)
		p.SigBuffer.WriteString(".")
		p.NameBuffer.WriteString(p.Swizzler.GoPackage.Name)
		p.NameBuffer.WriteString("_")
	}
	p.SigBuffer.WriteString(a.Name)
	p.NameBuffer.WriteString(a.Name)
	return nil
}

func IsBuiltInGoType(s string) bool {
	return BuiltInGoTypes[s]
}
