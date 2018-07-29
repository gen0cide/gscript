package engine

import "github.com/robertkrimen/otto"

// NativePackage defines a golang library that is being imported into the genesis VM
type NativePackage struct {
	ImportPath  string
	Name        string
	SymbolTable map[string]*NativeFunc
	Types       map[string]*NativeType
	Consts      map[string]*NativeConst
	Vars        map[string]*NativeVar
}

// NativeFunc defines a golang library function that is callable from within the genesis VM
type NativeFunc struct {
	Name      string
	Signature string
	Func      func(call otto.FunctionCall) otto.Value
}

// NativeConst defines a golang const declared within a given library
type NativeConst struct {
	Name  string
	Value interface{}
}

// NativeVar defines a golang top level var declaration within a given library
type NativeVar struct {
	Name      string
	Signature string
	Value     interface{}
}

// ParamDef defines basic information about either an argument or return value for NativeFunc
type ParamDef struct {
	Name   string
	GoType string
	JSType string
}

// NativeType expresses a native Golang type definition that can be used within the engine
type NativeType struct {
	Name    string
	Factory func(call otto.FunctionCall) otto.Value
	Fields  map[string]*NativeField
}

// NativeField expresses a struct field within a Go native type within the engine
type NativeField struct {
	Label     string
	Signature string
}

func (e *Engine) createType(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 {
		return e.Raise("argument", "not enough arguments for call to new()")
	}
	if len(call.ArgumentList) > 1 {
		return e.Raise("argument", "too many arguments for call to new()")
	}
	arg := call.Argument(0)
	if !arg.IsDefined() {
		return e.Raise("type", "invalid type passed to new()")
	}

	var nt *NativeType

	rarg, err := arg.Export()
	if err != nil {
		return e.Raise("jsexport", "could not export argument %d of function %s", 0, "new()")
	}
	switch v := rarg.(type) {
	case *NativeType:
		nt = v
	default:
		return e.Raise("type conversion", "argument type mismatch - expected %s, got %T", "*NativeType", v)
	}

	return nt.Factory(call)
}
