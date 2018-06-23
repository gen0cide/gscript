package compiler

import "go/ast"

// LinkedFunction is the type that represents the gscript <=> golang native binding
// so proper interfaces can be generated at compile time for calling native go from
// the genesis VM.
type LinkedFunction struct {
	// string representation of the function basename
	Function string

	// reference to the caller in the genesis script AST
	Caller *FunctionCall

	// reference to the golang AST tree of the file this function is declared in
	File *ast.File

	// reference to the declaration of this function in the golang AST
	GoDecl *ast.FuncDecl

	// list of references to any imports needed in declaring argument and return parameters
	// for this linked function
	Imports []*ast.ImportSpec

	// reference to the compiler's go package object
	GoPackage *GoPackage

	// a slice of all the go parameters that make up the argument signature
	GoArgs []*GoParamDef

	// a slice of all the go parameters that make up the return signature
	GoReturns []*GoParamDef

	// a reference to the parent genesis VM
	GenesisVM *GenesisVM
}

// Linker holds the maps between functions called from the genesis script and
// their associated golang equivalent, including package references. The linker
// will use this mapping to generate import shims for each golang public golang
// function called.
type Linker struct {
	// a reference to the parent genesis VM
	VM *GenesisVM

	// mapping of function name to the linked function object used during generation
	Funcs map[string]*LinkedFunction
}
