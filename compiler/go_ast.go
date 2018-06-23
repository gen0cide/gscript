package compiler

import (
	"bytes"
	"go/ast"
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
}
