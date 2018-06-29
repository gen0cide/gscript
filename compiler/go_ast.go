package compiler

import (
	"bytes"
	"fmt"
	"go/ast"
	"reflect"
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
)

// GoPackage holds all the information about a Golang package that is being resolved to a given script
type GoPackage struct {
	// Dir is the local path where this package is found
	Dir string

	// ImportPath is the golang import path used for this package
	ImportPath string

	// Name defines the go package's name
	Name string

	// Script references the script that is importing this package
	Script *GenesisVM

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
}

// GoParamDef defines a type to represent parameters found in a Golang function declaration (arguments or return types)
type GoParamDef struct {
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

	// LinkedFUnction is used to reference the parent LinkedFunction object
	LinkedFunction *LinkedFunction
}

// NewGoParamDef creates a new definition object for a go parameter (either return or argument) and returns a pointer to itself.
func NewGoParamDef(l *LinkedFunction, idx int) *GoParamDef {
	gpd := &GoParamDef{
		LinkedFunction: l,
		ImportRefs:     map[string]*ast.ImportSpec{},
		ParamIdx:       idx,
	}
	gpd.NameBuffer.WriteString("_")
	return gpd
}

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
		return fmt.Errorf("function %s includes an unsupported parameter type: %s", p.LinkedFunction.GoDecl.Name.Name, "chan")
	case *ast.FuncType:
		return fmt.Errorf("function %s includes an unsupported parameter type: %s", p.LinkedFunction.GoDecl.Name.Name, "func")
	case *ast.InterfaceType:
		return fmt.Errorf("function %s includes an unsupported parameter type: %s", p.LinkedFunction.GoDecl.Name.Name, "interface{}")
	case *ast.StructType:
		return fmt.Errorf("function %s includes an unsupported parameter type: %s", p.LinkedFunction.GoDecl.Name.Name, "struct")
	default:
		valType := reflect.ValueOf(t)
		return fmt.Errorf("could not determine the golang ast type of %s in func %s", valType.Type().String(), p.LinkedFunction.Function)
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
	canResolve, err := p.LinkedFunction.CanResolveImportDep(x.Name)
	if err != nil {
		return err
	}
	if canResolve != true {
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

// ParseStarExpr interprets a golang pointer into the appropriate GoParamDef structure
func (p *GoParamDef) ParseStarExpr(a *ast.StarExpr) error {
	p.SigBuffer.WriteString("*")
	p.NameBuffer.WriteString("PointerTo")
	return p.Interpret(a.X)
}

// ParseIdent interprets a golang identifier into the appropriate GoParamDef structure
func (p *GoParamDef) ParseIdent(a *ast.Ident) error {
	if ok := builtInGoTypes[a.Name]; !ok {
		p.SigBuffer.WriteString(p.LinkedFunction.GoPackage.Name)
		p.SigBuffer.WriteString(".")
		p.NameBuffer.WriteString(p.LinkedFunction.GoPackage.Name)
		p.NameBuffer.WriteString("_")
	}
	p.SigBuffer.WriteString(a.Name)
	p.NameBuffer.WriteString(a.Name)
	return nil
}

// IsBuiltInGoType takes a string argument and determines if is a valid built-in type
// in golang
func IsBuiltInGoType(s string) bool {
	return builtInGoTypes[s]
}
