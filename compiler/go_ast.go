package compiler

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"reflect"
	"regexp"
	"sync"

	"github.com/gen0cide/gscript/compiler/computil"
	"github.com/gen0cide/gscript/compiler/translator"
)

var (
	builtInGoTypes = map[string]bool{
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

	// TODO (gen0cide): fix this
	invalidGoTypes = map[string]bool{
		"complex128": true,
		"complex64":  true,
	}

	binaryImports = map[string]string{
		"bytes":           "bytes",
		"compress/gzip":   "gzip",
		"crypto/aes":      "aes",
		"crypto/cipher":   "cipher",
		"encoding/base64": "base64",
		"fmt":             "fmt",
		"io":              "io",
		"github.com/gen0cide/gscript/engine":          "engine",
		"github.com/robertkrimen/otto":                "otto",
		"github.com/gen0cide/gscript/debugger":        "debugger",
		"github.com/gen0cide/gscript/logger/standard": "standard",
	}

	funcRegexp  = regexp.MustCompile(`^func\({1}(?P<args>.*?)?\){1}\s*\(?(?P<rets>.*?)\)??$`)
	multipleRet = regexp.MustCompile(`,`)
)

// MaskedImport is used to separate import namespaces within the intermediate representation
type MaskedImport struct {
	// ImportPath of the masked Import
	ImportPath string

	// OldAlias represents the alias in the target package source
	OldAlias string

	// NewAlias represents the aliased package name in the intermediate representation
	NewAlias string
}

// GoPackage holds all the information about a Golang package that is being resolved to a given script
type GoPackage struct {
	sync.RWMutex

	// Dir is the local path where this package is found
	Dir string

	// MaskedName is the masked import representation of this gopackage
	MaskedName string

	// ImportPath is the golang import path used for this package
	ImportPath string

	// Name defines the go package's name
	Name string

	// VM references the script that is importing this package
	VM *GenesisVM

	// Namespace is the namespace aliased in the parent script
	Namespace string

	// ImportKey is the import path defined in the parent script
	ImportKey string

	// ScriptCallers maps the go package function names to genesis AST function calls
	ScriptCallers map[string]*FunctionCall

	// ImportsByFile is the map of each file within the package and what dependencies it imports
	ImportsByFile map[string][]*ast.ImportSpec

	// ImportsByAlias is a map of each import's alias and its dependency information
	ImportsByAlias map[string]*ast.ImportSpec

	// FuncToFileMap maps each public function within the package to it's corrasponding source file
	FuncToFileMap map[string]string

	// FuncTable is the map of each function name to it's Golang AST declaration
	FuncTable map[string]*ast.FuncDecl

	// LinkedFuncs defines references to the dynamically linked functions for this go package
	LinkedFuncs []*LinkedFunction

	// Reference to know if this go package is part of the standard library
	IsStandardLib bool

	// FileSet is used by the parser to interpret the current golang's file tokens
	FileSet *token.FileSet

	// GoTypes are struct type declarations that need to be accounted for in the engine
	GoTypes map[string]*GoStructDef

	// GoConsts are const declarations that need to be accounted for in the engine
	GoConsts map[string]*GoConst

	// GoVars are top level var declarations that need to be accounted for in the engine
	GoVars map[string]*GoVar
}

// GoParamDef defines a type to represent parameters found in a Golang function declaration (arguments or return types)
type GoParamDef struct {
	// Type denotes the type of parameter def (struct, function, etc.)
	Type string

	// OriginalName is to reference the original name of the field incase of problems
	OriginalName string

	// SigBuffer is used to create a definition to the actual type declaration in the genesis compiler's linker
	SigBuffer bytes.Buffer

	// NameBuffer is used to create a label that the genesis linker can use to create it's translations
	NameBuffer bytes.Buffer

	// ImportRefs holds a mapping of any golang dependencies required for this parameter's type declaration
	ImportRefs map[string]*ast.ImportSpec

	// VarName holds a representation of the logical representation of a parameter's label with it's offset appended
	VarName string

	// ParamIdx will be the relative position within the parameter declaration.
	// This is inclusive to multiple declarations of the same type. Example:
	//
	// func Foo(a, b string, c int) {}
	// GoParamDef objects for "a" and "b" would hold the same ParamIdx, but different
	// ArgOffset values
	ParamIdx int

	// ArgOffset defines the absolute position within the parameter declaration. It will increment
	// Regardless of multiple labels defined in the same type.
	// This is used by the linker to correctly translate arguments for golang functions
	ArgOffset int

	// ExtSig will hold the final result of the SigBuffer rendering as a string
	ExtSig string

	// GoLabel is used to represent the label name within Golang
	GoLabel string

	// NeedsMapping defines a hard coded type mapping between golang and javascript types
	NeedsMapping bool

	// MappedTypeAlias is a helper method to look up the type a method should be
	MappedTypeAlias string

	// LinkedFUnction is used to reference the parent LinkedFunction object
	LinkedFunction *LinkedFunction

	// IsInterfaceType defines ifthe GoParamDef an interface type?
	IsInterfaceType bool

	// SkipResolution defines whether the left recursive lexer should skip package resolution
	SkipResolution bool

	// Reference to the parent GoPackage
	GoPackage *GoPackage

	// False, unless the interpretation fails
	Errored bool

	// Holds the error message if the interpretation fails
	ErrorMessage string
}

// GoStructDef defines a struct type definition to be used within the genesis engine
type GoStructDef struct {
	sync.RWMutex

	// Package is a reference to the parent GoPackage
	Package *GoPackage

	// File is a reference to the parent AST file
	File *ast.File

	// TypeSpec references the typespec node within the file's AST
	TypeSpec *ast.TypeSpec

	// AST references the actual struct definition within the file's AST
	AST *ast.StructType

	// Name represents the name of the struct definition
	Name string

	// Key represents the masked name of the function declaration in the intermediate representation
	Key string

	// ImportRefs hold a mapping of any golang dependencies required for this type declaration
	ImportRefs map[string]*ast.ImportSpec

	// Fields are the fields that can be used within this struct (compatible with genesis type conversion)
	Fields map[string]*GoParamDef

	// Incompatible are fields that CANNOT be used with genesis
	Incompatible map[string]*GoParamDef

	// Embeds is a convenience reference showing any embedded types that might exist within this struct
	Embeds map[string]*GoParamDef
}

// GoConst defines the name of a golang constant declaration
type GoConst struct {
	// Name refers to the const name within the package
	Name string

	// Key refers to the unique identifier for a const in the intermediate representation
	Key string
}

// GoVar defines a top level var declaration within a Golang package
type GoVar struct {
	sync.RWMutex

	// Package is a reference to the parent GoPackage
	Package *GoPackage

	// File is a reference to the parent AST file
	File *ast.File

	// ValueSpec is the parent ValueSpec declaration in the AST
	ValueSpec *ast.ValueSpec

	// VSOffset is the offset within the ValueSpec's Ident this var is found
	VSOffset int

	// Ident is a reference to the actual Golang AST ident object
	Ident *ast.Ident

	// Expr holds the data we are checking out!
	Expr ast.Expr

	// Name is the name of the exported var
	Name string

	// Key is the swizzled unique identifier to the declaration
	Key string

	// Valid means it has passed swizzle sanity check for types we can't mess with
	Valid bool

	// Def is the swizzling point of reference for this particular package
	Def *GoParamDef
}

// NewGoPackage is a constructor for a gopackage that will be used in dynamically linking native code
func NewGoPackage(v *GenesisVM, ns, ikey string, stdlib bool) *GoPackage {
	return &GoPackage{
		VM:             v,
		Namespace:      ns,
		ImportKey:      ikey,
		ImportPath:     ikey,
		ScriptCallers:  map[string]*FunctionCall{},
		ImportsByFile:  map[string][]*ast.ImportSpec{},
		ImportsByAlias: map[string]*ast.ImportSpec{},
		FuncToFileMap:  map[string]string{},
		FuncTable:      map[string]*ast.FuncDecl{},
		LinkedFuncs:    []*LinkedFunction{},
		IsStandardLib:  stdlib,
		MaskedName:     computil.RandLowerAlphaString(6),
		GoTypes:        map[string]*GoStructDef{},
		GoConsts:       map[string]*GoConst{},
		GoVars:         map[string]*GoVar{},
	}
}

// NewGoParamDef creates a new definition object for a go parameter (either return or argument) and returns a pointer to itself.
func NewGoParamDef(l *LinkedFunction, idx int) *GoParamDef {
	gpd := &GoParamDef{
		LinkedFunction: l,
		ImportRefs:     map[string]*ast.ImportSpec{},
		ParamIdx:       idx,
		OriginalName:   l.GoDecl.Name.Name,
		Type:           "function",
		GoPackage:      l.GoPackage,
	}
	gpd.NameBuffer.WriteString("_")
	return gpd
}

// NewMaskedImport creates a new import mask based on an import path and old alias
func NewMaskedImport(ip, oa string) *MaskedImport {
	alias := computil.RandLowerAlphaString(6)
	if val, ok := binaryImports[ip]; ok {
		alias = val
	}
	return &MaskedImport{
		ImportPath: ip,
		OldAlias:   oa,
		NewAlias:   alias,
	}
}

// NewConst creates a new const reference within a golang package
func (gop *GoPackage) NewConst(name string) *GoConst {
	gop.Lock()
	defer gop.Unlock()
	val, ok := gop.GoConsts[name]
	if ok && val != nil {
		return val
	}
	gconst := &GoConst{
		Name: name,
		Key:  computil.RandUpperAlphaString(7),
	}
	gop.GoConsts[name] = gconst
	return gconst
}

// NewGoVar creates a new GoVar object if one doesn't already exist with that name in the GoPackage
func (gop *GoPackage) NewGoVar(goast *ast.File, vs *ast.ValueSpec, offset int, ident *ast.Ident, expr ast.Expr, name string) (*GoVar, error) {
	gop.Lock()
	defer gop.Unlock()
	ngv, ok := gop.GoVars[name]
	if ok {
		return ngv, errors.New("var declaration already exists in this package")
	}
	p := &GoParamDef{
		ParamIdx:       offset,
		SkipResolution: true,
		Type:           "var",
		OriginalName:   name,
		GoPackage:      gop,
	}
	gv := &GoVar{
		Package:   gop,
		File:      goast,
		ValueSpec: vs,
		VSOffset:  offset,
		Ident:     ident,
		Expr:      expr,
		Name:      name,
		Key:       computil.RandAlphaNumericString(9),
		Def:       p,
	}
	gop.GoVars[name] = gv
	return gv, nil
}

// IsDefaultImport tests a golang import path to determine if it is already defined in the intermediate representation
func IsDefaultImport(ip string) bool {
	return binaryImports[ip] != ""
}

// GetDefaultImportNamespace returns the corrasponding namespace to the import path provided that is used in the intermediate representation
func GetDefaultImportNamespace(ip string) string {
	return binaryImports[ip]
}

// import (
//   rekt "net/url"
// )
// EXAMPLE: func Foo(a0 map[*rekt.URL][]*ast.Field)
// 0 ""
// 1 "map["
// 2 "map[*"
// 3 "map[*url"
// 3.5 "map[*url."
// 3.7 "map[*url.URL"
// 4 "map[*url.URL]"
// 5 "map[*url.URL][]"
// 6 "map[*url.URL][]*"
// 7 "map[*url.URL][]*ast"
// 7.5 "map[*url.URL][]*ast."
// 7.7 "map[*url.URL][]*ast.Field"

// Interpret is a recursive walk function that is used to dispatch the next walker
// depending on the type of the provided interface (i). This is used to build up
// buffers of both names and golang type declarations to be used during linking.
func (p *GoParamDef) Interpret(i interface{}) error {
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
		return fmt.Errorf("%s %s includes an unsupported parameter type: %s", p.Type, p.OriginalName, "chan")
	case *ast.FuncType:
		return fmt.Errorf("%s %s includes an unsupported parameter type: %s", p.Type, p.OriginalName, "func")
	case *ast.InterfaceType:
		return p.ParseInterfaceType(t)
		//return fmt.Errorf("%s %s includes an unsupported parameter type: %s", p.Type, p.OriginalName, "interface{}")
	case *ast.StructType:
		return fmt.Errorf("%s %s includes an unsupported parameter type: %s", p.Type, p.OriginalName, "struct")
	case *ast.UnaryExpr:
		if p.Type == "var" {
			p.SigBuffer.Reset()
			printer.Fprint(&p.SigBuffer, p.GoPackage.FileSet, t)
			return nil
		}
		valType := reflect.ValueOf(t)
		return fmt.Errorf("could not determine the golang ast type of %s in %s %s.%s", valType.Type().String(), p.Type, p.GoPackage.Name, p.OriginalName)
	case *ast.CallExpr:
		if p.Type == "var" {
			p.SigBuffer.Reset()
			printer.Fprint(&p.SigBuffer, p.GoPackage.FileSet, t)
			return nil
		}
		valType := reflect.ValueOf(t)
		return fmt.Errorf("could not determine the golang ast type of %s in %s %s.%s", valType.Type().String(), p.Type, p.GoPackage.Name, p.OriginalName)
	case *ast.CompositeLit:
		if p.Type == "var" {
			p.SigBuffer.Reset()
			printer.Fprint(&p.SigBuffer, p.GoPackage.FileSet, t)
			return nil
		}
		valType := reflect.ValueOf(t)
		return fmt.Errorf("could not determine the golang ast type of %s in %s %s.%s", valType.Type().String(), p.Type, p.GoPackage.Name, p.OriginalName)
	default:
		if t == nil {
			return nil
		}
		valType := reflect.ValueOf(t)
		// spew.Dump(t)
		return fmt.Errorf("could not determine the golang ast type of %s in %s %s.%s", valType.Type().String(), p.Type, p.GoPackage.Name, p.OriginalName)
	}
}

// ParseMapType interprets a golang map type into the appropriate GoParamDef structure
func (p *GoParamDef) ParseMapType(a *ast.MapType) error {
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

// ParseInterfaceType is used to parse an interface type that is being passed into a linked function
func (p *GoParamDef) ParseInterfaceType(i *ast.InterfaceType) error {
	p.SigBuffer.WriteString("interface{}")
	p.NameBuffer.WriteString("InterfaceType")
	return nil
}

// ParseArrayType interprets a golang array/slice type into the appropriate GoParamDef structure
func (p *GoParamDef) ParseArrayType(a *ast.ArrayType) error {
	p.SigBuffer.WriteString("[]")
	p.NameBuffer.WriteString("ArrayOf")
	return p.Interpret(a.Elt)
}

// ParseSelectorExpr interprets a golang namespace external to the function declarations package
// and maps it into the appropriate GoParamDef structure
func (p *GoParamDef) ParseSelectorExpr(a *ast.SelectorExpr) error {
	x, ok := a.X.(*ast.Ident)
	if !ok {
		return fmt.Errorf("could not parse selector namespace in func %s", p.LinkedFunction.Function)
	}

	if p.SkipResolution {
		p.SigBuffer.WriteString(x.Name)
		p.NameBuffer.WriteString(x.Name)
		p.SigBuffer.WriteString(".")
		p.NameBuffer.WriteString("_")
		p.SigBuffer.WriteString(a.Sel.Name)
		p.NameBuffer.WriteString(a.Sel.Name)
		return nil
	}

	resolved, err := p.LinkedFunction.CanResolveImportDep(x.Name)
	if err != nil {
		return err
	}

	mappedType := p.MappedType(x.Name, a.Sel.Name)
	if mappedType != "" {
		p.MappedTypeAlias = mappedType
		p.NeedsMapping = true
	}

	p.SigBuffer.WriteString(resolved)
	p.NameBuffer.WriteString(resolved)
	p.SigBuffer.WriteString(".")
	p.NameBuffer.WriteString("_")
	p.SigBuffer.WriteString(a.Sel.Name)
	p.NameBuffer.WriteString(a.Sel.Name)
	return nil
}

// ParseStarExpr interprets a golang pointer into the appropriate GoParamDef structure
func (p *GoParamDef) ParseStarExpr(a *ast.StarExpr) error {
	p.SigBuffer.WriteString("*")
	p.NameBuffer.WriteString("PointerTo")
	return p.Interpret(a.X)
}

// ParseIdent interprets a golang identifier into the appropriate GoParamDef structure
func (p *GoParamDef) ParseIdent(a *ast.Ident) error {
	if ok := builtInGoTypes[a.Name]; !ok {

		if p.SkipResolution {
			p.SigBuffer.WriteString(a.Name)
			p.NameBuffer.WriteString(a.Name)
			return nil
		}

		if IsDefaultImport(p.LinkedFunction.GoPackage.ImportPath) {
			p.SigBuffer.WriteString(GetDefaultImportNamespace(p.LinkedFunction.GoPackage.ImportPath))
			p.SigBuffer.WriteString(".")
			p.NameBuffer.WriteString(GetDefaultImportNamespace(p.LinkedFunction.GoPackage.ImportPath))
			p.NameBuffer.WriteString("_")
		} else {
			p.SigBuffer.WriteString(p.LinkedFunction.GoPackage.MaskedName)
			p.SigBuffer.WriteString(".")
			p.NameBuffer.WriteString(p.LinkedFunction.GoPackage.MaskedName)
			p.NameBuffer.WriteString("_")
		}
	}
	p.SigBuffer.WriteString(a.Name)
	p.NameBuffer.WriteString(a.Name)
	return nil
}

// Signature is a helper function that returns the type representation of a GoVar
func (gv *GoVar) Signature() string {
	return gv.Def.SigBuffer.String()
}

// MappedType looks at whether there is a type mapping that needs to be honored
func (p *GoParamDef) MappedType(pkg, sel string) string {
	if val, ok := translator.TypeAliasMap[pkg]; ok {
		if t, ok := val[sel]; ok {
			return t
		}
	}
	return ""
}

// IsBuiltInGoType takes a string argument and determines if is a valid built-in type
// in golang
func IsBuiltInGoType(s string) bool {
	return builtInGoTypes[s]
}

// ParseTypeSpec attempts to filter type spec definitions within the AST to only struct types
func (gop *GoPackage) ParseTypeSpec(goast *ast.File, v *ast.TypeSpec, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	switch t := v.Type.(type) {
	case *ast.StructType:
		wg.Add(1)
		go gop.ParseStructDef(goast, v, t, wg, errChan)
	}
	return
}

// ParseStructDef creates a new mapping between a go package and a struct type definition
func (gop *GoPackage) ParseStructDef(goast *ast.File, v *ast.TypeSpec, t *ast.StructType, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	gop.Lock()
	defer gop.Unlock()
	if _, ok := gop.GoTypes[v.Name.Name]; ok {
		errChan <- fmt.Errorf("struct def already mapped for %s - file=%s pkg=%s", v.Name.Name, goast.Name.Name, gop.ImportPath)
		return
	}
	gsd := &GoStructDef{
		Package:      gop,
		File:         goast,
		TypeSpec:     v,
		AST:          t,
		Name:         v.Name.Name,
		ImportRefs:   map[string]*ast.ImportSpec{},
		Fields:       map[string]*GoParamDef{},
		Incompatible: map[string]*GoParamDef{},
		Embeds:       map[string]*GoParamDef{},
		Key:          computil.RandUpperAlphaString(13),
	}
	gop.GoTypes[v.Name.Name] = gsd
	wg.Add(1)
	go gsd.WalkStruct(wg, errChan)
	return
}

// WalkStruct walks the type definition AST for exported struct fields
func (gsd *GoStructDef) WalkStruct(wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	if gsd.AST.Fields == nil {
		return
	}
	for fidx, tf := range gsd.AST.Fields.List {
		if len(tf.Names) == 0 {
			// struct embed
			wg.Add(1)
			go gsd.ParseStructField(tf, "_EMBED", fidx, wg, errChan)
		}
		for noff, name := range tf.Names {
			if !name.IsExported() {
				continue
			}
			wg.Add(1)
			go gsd.ParseStructField(tf, name.Name, (noff + fidx), wg, errChan)
		}
	}
}

// ParseStructField attempts to interpret the struct fields within exported go types
func (gsd *GoStructDef) ParseStructField(f *ast.Field, name string, offset int, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	p := &GoParamDef{
		ParamIdx:       offset,
		SkipResolution: true,
		Type:           "struct_field",
		OriginalName:   name,
		GoPackage:      gsd.Package,
	}

	err := p.Interpret(f.Type)
	if err != nil {
		p.Errored = true
		p.ErrorMessage = err.Error()
		if name == "_EMBED" {
			name = fmt.Sprintf("EMBED_%d", offset)
		}
		gsd.Lock()
		gsd.Incompatible[name] = p
		gsd.Unlock()
		return
	}
	if name == "_EMBED" {
		gsd.Lock()
		gsd.Embeds[p.SigBuffer.String()] = p
		gsd.Unlock()
		return
	}
	gsd.Lock()
	gsd.Fields[name] = p
	gsd.Unlock()
	return
}

// ParseDeclaration attempts to parse a top level var declaration and account for any incompatible types
func (gv *GoVar) ParseDeclaration(wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	gv.Lock()
	defer gv.Unlock()
	err := gv.Def.Interpret(gv.Expr)
	if err == nil {
		gv.Valid = true
	} else {
		gv.Package.VM.Logger.Errorf("Could not parse declaration %s.%s: %v", gv.Package.Name, gv.Name, err)
	}
	return
}

func (gop *GoPackage) printResults() {
	structs := map[string]*GoStructDef{}
	for sname, s := range gop.GoTypes {
		structs[sname] = s
	}
	gop.VM.Logger.Infof("===  %s STRUCTS  ===", gop.ImportKey)
	for n, s := range structs {
		gop.VM.Logger.Infof("  %s (fields=%d, incompatibles=%d, embeds=%d)", n, len(s.Fields), len(s.Incompatible), len(s.Embeds))
		for fn, f := range s.Fields {
			gop.VM.Logger.Infof("    (F) %s (%s)", fn, f.SigBuffer.String())
		}
		for fn, f := range s.Incompatible {
			gop.VM.Logger.Infof("    (I) %s (%s)", fn, f.NameBuffer.String())
		}
		for _, f := range s.Embeds {
			gop.VM.Logger.Infof("    (E) %s", f.SigBuffer.String())
		}
	}
}

// ParseVarSpec walks a var declaration inside of a gopackage to make sure it can sanely be accounted for by the compiler
func (gop *GoPackage) ParseVarSpec(goast *ast.File, vardecl *ast.ValueSpec, offset int, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	gv, err := gop.NewGoVar(goast, vardecl, offset, vardecl.Names[offset], vardecl.Type, vardecl.Names[offset].Name)
	if err != nil {
		return
	}
	wg.Add(1)
	gv.ParseDeclaration(wg, errChan)
	return
}

// ParseFuncDecl walks a func declaration inside of a gopackage to result it's method signature
func (gop *GoPackage) ParseFuncDecl(goast *ast.File, funcdecl *ast.FuncDecl, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	funcName := funcdecl.Name.Name
	if funcdecl.Name.IsExported() && funcdecl.Recv == nil {
		gop.Lock()
		caller := gop.ScriptCallers[funcName]
		if caller == nil && !gop.IsStandardLib && !gop.VM.Options.ImportAllNativeFuncs {
			gop.Unlock()
			return
		}
		sig := new(bytes.Buffer)
		printer.Fprint(sig, gop.FileSet, funcdecl.Type)
		lf, err := gop.VM.Linker.NewLinkedFunction(
			funcName,
			caller,
			goast,
			funcdecl,
			goast.Imports,
			gop,
		)
		if err != nil {
			gop.Unlock()
			errChan <- err
			return
		}
		match := funcRegexp.FindStringSubmatch(sig.String())
		result := make(map[string]string)
		for i, name := range funcRegexp.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		newSigBuf := new(bytes.Buffer)
		if result["rets"] != "" {
			if multipleRet.MatchString(result["rets"]) {
				newSigBuf.WriteString("[")
			}
			newSigBuf.WriteString(result["rets"])
			if multipleRet.MatchString(result["rets"]) {
				newSigBuf.WriteString("]")
			}
			newSigBuf.WriteString(" = ")
		}
		if gop.IsStandardLib {
			newSigBuf.WriteString("G.")
			newSigBuf.WriteString(gop.Name)
		} else {
			newSigBuf.WriteString(gop.Namespace)
		}
		newSigBuf.WriteString(".")
		newSigBuf.WriteString(funcName)
		newSigBuf.WriteString("(")
		newSigBuf.WriteString(result["args"])
		newSigBuf.WriteString(")")
		lf.Signature = newSigBuf.String()
		if len(gop.ImportsByFile[goast.Name.Name]) == 0 {
			gop.ImportsByFile[goast.Name.Name] = goast.Imports
		}
		gop.FuncTable[funcName] = funcdecl
		gop.FuncToFileMap[funcName] = goast.Name.Name
		gop.LinkedFuncs = append(gop.LinkedFuncs, lf)
		gop.Unlock()
	}
}

// WalkGoFileAST walks the AST of a golang file and determines if it should be included as a linked
// function based on one of the following statements being true:
// Parent GoPackage is a member of the standard library
// OR
// Compiler option ImportAllNativeFunc is set to true
// OR
// VM Script calls this function explicitly
func (gop *GoPackage) WalkGoFileAST(goast *ast.File, wg *sync.WaitGroup, errChan chan error) {
	for _, n := range goast.Decls {
		switch v := n.(type) {
		case *ast.GenDecl:
			switch v.Tok {
			case token.CONST:
				for _, cdecl := range v.Specs {
					if valdecl, ok := cdecl.(*ast.ValueSpec); ok {
						if len(valdecl.Names) > 0 {
							for _, name := range valdecl.Names {
								if !name.IsExported() {
									continue
								}
								gop.NewConst(name.Name)
							}
						}
					}
				}
			case token.VAR:
				for _, vdecl := range v.Specs {
					if valdecl, ok := vdecl.(*ast.ValueSpec); ok {
						for idx, name := range valdecl.Names {
							if !name.IsExported() {
								continue
							}
							wg.Add(1)
							go gop.ParseVarSpec(goast, valdecl, idx, wg, errChan)
						}
					}
				}
			case token.TYPE:
				for _, tdecl := range v.Specs {
					if typedecl, ok := tdecl.(*ast.TypeSpec); ok {
						if !typedecl.Name.IsExported() {
							continue
						}
						wg.Add(1)
						go gop.ParseTypeSpec(goast, typedecl, wg, errChan)
					}
				}
			}
		case *ast.FuncDecl:
			wg.Add(1)
			go gop.ParseFuncDecl(goast, v, wg, errChan)
		}
	}

	wg.Done()
	return
}

// SanityCheckScriptCallers enumerates all of the parent gopackage's script callers looking for any
// callers who do not have an entry in this go package's symbol table
func (gop *GoPackage) SanityCheckScriptCallers() error {
	for fnName := range gop.ScriptCallers {
		if gop.FuncTable[fnName] == nil {
			return fmt.Errorf("function %s is not a valid function in package %s", fnName, gop.Name)
		}
	}
	return nil
}

// SuccessfullyLinkedFuncs is used during rendering to make sure that only linked functions that successfully
// swizzled are going to be built into the source
func (gop *GoPackage) SuccessfullyLinkedFuncs() []*LinkedFunction {
	lf := []*LinkedFunction{}
	for _, l := range gop.LinkedFuncs {
		if l.SwizzleSuccessful {
			lf = append(lf, l)
		}
	}

	return lf
}

// ValidVars is used during rendering to ensure that only valid var declarations that have been
// processed by the swizzler are included in the intermediate representation
func (gop *GoPackage) ValidVars() []*GoVar {
	gvs := []*GoVar{}
	for _, g := range gop.GoVars {
		if g.Valid {
			gvs = append(gvs, g)
		}
	}
	return gvs
}

// BuiltInTranslationRequired is a compiler helper to determine if the param definition requires a built in translation
func (p *GoParamDef) BuiltInTranslationRequired() bool {
	return translator.BuiltInMap[p.ExtSig] != ""
}

// BuiltInJSType returns the type that JS will return that we need to convert the ParamDef to golang
func (p *GoParamDef) BuiltInJSType() string {
	return translator.BuiltInMap[p.ExtSig]
}
