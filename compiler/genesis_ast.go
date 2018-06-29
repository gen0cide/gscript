package compiler

import (
	"go/ast"
	"go/token"

	gast "github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/file"
)

var (

	// NamespaceConversionMap is a temporary mapping to allow genesis scripts
	// written in the v0 specification to be backwards compatible with v1 scripts.
	// The major difference is that v1 did not have namespaces for the standard library
	// where as v1 implements them with a deprecation and v2 will remove them entirely
	NamespaceConversionMap = map[string]*LegacyFunctionCall{}
)

// LegacyLibrary defines one of the original standard library golang packages
// for backward compatibility during linking to associate global namespace
// function calls with the v1 package style
type LegacyLibrary struct {
	// name of the owning package in the legacy standard library
	Name string `yaml:"package"`

	// path to the golang file located in which this library is implemented
	Path string `yaml:"path"`

	// file set for the AST walk and re-write
	FSet *token.FileSet `yaml:"-"`

	// Golang AST representation of the library's source file
	AST *ast.File `yaml:"-"`

	// map of the function names to their LegacyFunctionCall objects
	Funcs map[string]*LegacyFunctionCall `yaml:"-"`

	// the raw file data for the legacy library source
	Source []byte `yaml:"-"`
}

// LegacyFunctionCall uses the old generator to represent v0 standard library
// functions so they can be deprecated in subsequent versions without
// forcing users to convert all v0 scripts to v1 at this time
type LegacyFunctionCall struct {
	// name of the legacy function
	Name string `yaml:"name"`

	// description of the legacy function
	Description string `yaml:"description"`

	// owner of the legacy function
	Author string `yaml:"author"`

	// package name the legacy function is in
	Package string `yaml:"package"`

	// the expected arguments to the legacy function call
	ExpectedArgTypes []LegacyArgDef `yaml:"args"`

	// the expected returns to the legacy function call
	ExpectedReturnTypes []LegacyRetDef `yaml:"returns"`
}

// LegacyArgDef is used by LegacyFunctionCall objects to create mappings for
// legacy function argument's for the linker to inject into the build
type LegacyArgDef struct {
	// the name of the argument passed into the function
	Name string `yaml:"label"`

	// the golang type argument is expected to be
	GoType string `yaml:"gotype"`
}

// LegacyRetDef is used by LegacyFunctionCall objects to create mappings for
// legacy return values for the linker to inject into the build
type LegacyRetDef struct {
	// name of the rerturn value parameter
	Name string `yaml:"label"`

	// the golang type return parameter is expected to be
	GoType string `yaml:"gotype"`

	// optional value to determine if you wish to return this value to the VM
	ReturnToVM bool `yaml:"return,omitempty"`
}

// FunctionCall contains information relating to a Golang native function call
// that is being used within a genesis script
type FunctionCall struct {
	// genesis script namespace the function call corrasponds to
	Namespace string

	// name of the function call as used in the genesis script
	FuncName string

	// list of arguments passed to the genesis script function caller
	ArgumentList []gast.Expression
}

// genesisWalker is a type used to recursively walk the genesis script AST
type genesisWalker struct {
	// reference to the parent genesis VM object
	script *GenesisVM

	// source as represented as text
	source string

	// offset record used during the AST talk
	shift file.Idx
}

func (g *genesisWalker) Enter(n gast.Node) gast.Visitor {
	switch a := n.(type) {
	case *gast.CallExpression:
		switch b := a.Callee.(type) {
		case *gast.DotExpression:
			namespace, ok := b.Left.(*gast.Identifier)
			if !ok {
				// caller's left side was not an identifier (probably a function)
				// move on
				return g
			}
			if _, ok := g.script.GoPackageByNamespace[namespace.Name]; !ok {
				// given caller was not a known golang namespace
				// move on
				return g
			}
			// we got em :) adding to the gopackge script caller table
			fnName := b.Identifier.Name
			gop := g.script.GoPackageByNamespace[namespace.Name]
			gop.ScriptCallers[fnName] = &FunctionCall{
				Namespace:    namespace.Name,
				FuncName:     fnName,
				ArgumentList: a.ArgumentList,
			}
		}
	}
	return g
}

func (g *genesisWalker) Exit(n gast.Node) {
	// we done here - bye!
	return
}
