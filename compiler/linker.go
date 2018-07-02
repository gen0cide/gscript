package compiler

import (
	"fmt"
	"go/ast"
	"strings"
)

// LinkedFunction is the type that represents the gscript <=> golang native binding
// so proper interfaces can be generated at compile time for calling native go from
// the genesis VM.
type LinkedFunction struct {
	// ID of the function that is linked
	ID string

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

// NewLinker creates a new linker for the given genesis VM
func NewLinker(vm *GenesisVM) *Linker {
	return &Linker{
		VM:    vm,
		Funcs: map[string]*LinkedFunction{},
	}
}

// NewLinkedFunction creates a function mapping in the VMs linker between golang AST function declearations and genesis AST function calls
// so the compiler can build the function interfaces between the virtual machine and the native golang package
func (l *Linker) NewLinkedFunction(caller *FunctionCall, file *ast.File, godecl *ast.FuncDecl, imports []*ast.ImportSpec, gopkg *GoPackage) (*LinkedFunction, error) {
	if l.Funcs[caller.FuncName] != nil {
		return nil, fmt.Errorf("vm %s already has a linker for go func %s under package %s - new function is in package %s", l.VM.Name, caller.FuncName, l.Funcs[caller.FuncName].GoPackage.ImportPath, gopkg.ImportPath)
	}
	lf := &LinkedFunction{
		ID:        RandLowerAlphaString(16),
		Function:  caller.FuncName,
		Caller:    caller,
		File:      file,
		GoDecl:    godecl,
		Imports:   imports,
		GoPackage: gopkg,
		GenesisVM: l.VM,
		GoArgs:    []*GoParamDef{},
		GoReturns: []*GoParamDef{},
	}
	l.Funcs[caller.FuncName] = lf
	gopkg.LinkedFuncs = append(gopkg.LinkedFuncs, lf)
	return lf, nil
}

// SwizzleToTheLeft enumerates the function arguments of both the caller and the native function
// to build a structured list of parameters and their types. It also compares the caller argument
// signature and throws an error if the caller is providing an incompatible number of arguments.
func (l *LinkedFunction) SwizzleToTheLeft() error {
	aOff := 0
	for idx, p := range l.GoDecl.Type.Params.List {
		masterP := NewGoParamDef(l, idx)
		err := masterP.Interpret(p.Type)
		if err != nil {
			return err
		}
		masterP.VarName = masterP.NameBuffer.String()
		masterP.ExtSig = masterP.SigBuffer.String()
		for i := 0; i < len(p.Names); i++ {
			newP := NewGoParamDef(l, idx)
			newP.VarName = fmt.Sprintf("%s%d", masterP.VarName, aOff)
			newP.ArgOffset = aOff
			newP.ExtSig = masterP.ExtSig
			newP.GoLabel = p.Names[i].Name
			aOff++
			l.GoArgs = append(l.GoArgs, newP)
		}
	}
	return nil
}

// SwizzleToTheRight enumerates the function returns of the native function to build a structured
// list of the return value types. This is then used by the linker to generate a special wrapper
// to allow multiple return values to be returned in single value context (required by javascript)
func (l *LinkedFunction) SwizzleToTheRight() error {
	aOff := 0
	for idx, p := range l.GoDecl.Type.Results.List {
		masterP := NewGoParamDef(l, idx)
		err := masterP.Interpret(p.Type)
		if err != nil {
			return err
		}
		masterP.VarName = masterP.NameBuffer.String()
		masterP.ExtSig = masterP.SigBuffer.String()
		if len(p.Names) > 0 {
			for i := 0; i < len(p.Names); i++ {
				newP := NewGoParamDef(l, idx)
				newP.ArgOffset = aOff
				newP.ExtSig = masterP.ExtSig
				aOff++
				l.GoReturns = append(l.GoReturns, newP)
			}
		} else {
			masterP.ArgOffset = aOff
			aOff++
			l.GoReturns = append(l.GoReturns, masterP)
		}
	}
	return nil
}

// CanResolveImportDep takes a package string and compares it against the linked functions known import
// table to determine if the referenced namespace is declared in the golang AST as a referenced sub-type
func (l *LinkedFunction) CanResolveImportDep(pkg string) (bool, error) {
	if pkg == "." {
		return false, fmt.Errorf("should not attempt to import anonymously in package %s", l.File.Name.Name)
	}
	for _, i := range l.Imports {
		if i.Name != nil {
			if i.Name.Name == pkg {
				return true, nil
			}
		} else {
			pkgParts := strings.Split(i.Path.Value, "/")
			packageAlias := pkgParts[len(pkgParts)-1]
			newAlias := strings.Replace(packageAlias, `"`, ``, -1)
			if newAlias == pkg {
				return true, nil
			}
		}
	}
	return false, fmt.Errorf("could not resolve package %s used in function %s inside package %s", pkg, l.Function, l.GoPackage.ImportPath)
}

// GenerateReturnString generates a golang return signature to use in the interface code
func (l *LinkedFunction) GenerateReturnString(prefix string) string {
	rets := []string{}
	for idx := range l.GoReturns {
		rets = append(rets, fmt.Sprintf("%s%d", prefix, idx))
	}
	return strings.Join(rets, ", ")
}

// GenerateArgString generates a golang argument string to use in the interface code based on the number of
// arguments required for the supplied function
func (l *LinkedFunction) GenerateArgString(prefix string) string {
	args := []string{}
	for idx := range l.GoArgs {
		args = append(args, fmt.Sprintf("%s%d", prefix, idx))
	}
	return strings.Join(args, ", ")
}
