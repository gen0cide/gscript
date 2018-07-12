package engine

import "github.com/robertkrimen/otto"

// NativePackage defines a golang library that is being imported into the genesis VM
type NativePackage struct {
	ImportPath  string
	Name        string
	SymbolTable map[string]*NativeFunc
}

// NativeFunc defines a golang library function that is callable from within the genesis VM
type NativeFunc struct {
	Name      string
	Signature string
	Func      func(call otto.FunctionCall) otto.Value
}

// ParamDef defines basic information about either an argument or return value for NativeFunc
type ParamDef struct {
	Name   string
	GoType string
	JSType string
}
