package engine

import (
	"errors"
	"fmt"

	"github.com/gen0cide/otto"
)

// Engine defines the virtual machine type for the genesis scripting engine
type Engine struct {
	// javascript V8 runtime
	VM *otto.Otto

	// logger interface for any output
	Logger Logger

	// maps the asset names to the functions that return their bytes
	Imports map[string]func() []byte

	// maps the namespaces to native packages created by the compiler
	Packages map[string]*NativePackage

	// plaintext name of the VM - usually the script file basename
	Name string

	// unique identifier for this VM (unique per build)
	ID string

	// flag to denote whether the debugger is enabled
	DebuggerEnabled bool

	// timeout in seconds
	Timeout int

	// is the genesis VM halted
	Halted bool

	// is the genesis VM paused
	Paused bool

	// defines the entry point function for execution of the script
	EntryPoint string
}

// New returns a new genesis virtual machine with the given parameters (does not run, just returns the container object)
func New(name, id string, timeout int, entrypoint string) *Engine {
	e := &Engine{
		Name:       name,
		ID:         id,
		Timeout:    timeout,
		EntryPoint: entrypoint,
	}

	e.VM = otto.New()
	e.SetLogger(&NullLogger{})
}

// // NativeFunc is a type alias for native function wrappers
// type NativeFunc func(call otto.FunctionCall) otto.Value

// // SymbolTable maps a native function wrapper to it's caller reference for use by the VM
// type SymbolTable map[string]NativeFunc

// DeclareNamespace adds an empty namespace to the virtual machine.
// Caution! will overwrite any values at existing namespace.
func (e *Engine) DeclareNamespace(namespace string) (*otto.Object, error) {
	ns, err := e.VM.Get(namespace)
	if ns.IsUndefined() != true {
		return nil, fmt.Errorf("namespace %s is already defined in virtual machine %s", namespace, e.Name)
	} else if ns.IsUndefined() == true && err != nil {
		return nil, err
	}
	nsObj, err := e.VM.Object(fmt.Sprintf("%s = {}", namespace))
	if err != nil {
		return nil, err
	}
	return nsObj, nil
}

// ImportNativePackage adds a golang native package to the virtual machine's runtime at a specified namespace
func (e *Engine) ImportNativePackage(namespace string, pkg *NativePackage) error {
	ns, err := e.DeclareNamespace(namespace)
	if err != nil {
		return err
	}
	e.Packages[namespace] = pkg
	for n, f := range pkg.SymbolTable {
		err = ns.Set(n, f.Func)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetConst defines a const value within the virtual machine
func (e *Engine) SetConst(name string, value interface{}) error {
	val, err := e.VM.ToValue(value)
	if err != nil {
		return err
	}
	return e.VM.Set(name, val)
}

// SetTimeout sets the timeout in seconds for the virtual machine
func (e *Engine) SetTimeout(t int) {
	e.Timeout = t
}

// AddImport maps an asset to a filename in the VMs virtual file system
func (e *Engine) AddImport(filename string, loader func() []byte) {
	e.Imports[filename] = loader
}

// SetName sets the VM's human readable name
func (e *Engine) SetName(n string) {
	e.Name = n
}

// SetID sets the VM's unique ID
func (e *Engine) SetID(id string) {
	e.ID = id
}

// InitVM initializes the Engine's javascript virtual machine
func (e *Engine) InitVM() error {
	e.VM = otto.New()
	if e.VM != nil {
		return errors.New("could not initialize virtual machine")
	}
	return nil
}

// LoadScript takes a script (source) with a corrasponding filename for debugging purposes and
// checks it for syntax errors before evaluating the script within the virtual machine's
// current scope
func (e *Engine) LoadScript(filename string, source []byte) error {
	script, err := e.VM.Compile(filename, source)
	if err != nil {
		return err
	}
	_, err = e.VM.Eval(script)
	return err
}

// SetLogger overrides the logging interface for this virtual machine.
func (e *Engine) SetLogger(l Logger) {
	e.Logger = l
}
