package engine

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/gen0cide/gscript/logger"
	"github.com/gen0cide/gscript/logger/null"
	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/file"
	"github.com/robertkrimen/otto/parser"
)

// Engine defines the virtual machine type for the genesis scripting engine
type Engine struct {
	// javascript V8 runtime
	VM *otto.Otto

	// logger interface for any output
	Logger logger.Logger

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

	// backwards compatibility to tell the planner whether to execute the before hook
	BeforeHook bool

	// backwards compatibility option to tell the planner whether to execute the after hook
	AfterHook bool

	// defines the entry point function for execution of the script
	EntryPoint string

	// used to map the various javascript files
	fileSet *file.FileSet
}

// New returns a new genesis virtual machine with the given parameters (does not run, just returns the container object)
func New(name, id string, timeout int, entrypoint string) *Engine {
	e := &Engine{
		Name:       name,
		ID:         id,
		Timeout:    timeout,
		EntryPoint: entrypoint,
		AfterHook:  false,
		BeforeHook: false,
		fileSet:    &file.FileSet{},
		Imports:    map[string]func() []byte{},
		Packages:   map[string]*NativePackage{},
	}
	e.InitVM()
	e.SetLogger(&null.Logger{})
	e.setGlobalRef()
	return e
}

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
	for n, t := range pkg.Types {
		err = ns.Set(n, t)
		if err != nil {
			return err
		}
	}
	for n, c := range pkg.Consts {
		err = ns.Set(n, c.Value)
		if err != nil {
			return err
		}
	}
	for n, va := range pkg.Vars {
		err = ns.Set(n, va.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// ImportStandardLibrary injects all provided native packages into the standard libraries namespace within the engine
func (e *Engine) ImportStandardLibrary(pkgs []*NativePackage) error {
	for _, p := range pkgs {
		pkgName := p.Name
		nsObj, err := e.VM.Object(fmt.Sprintf("G.%s = {}", pkgName))
		if err != nil {
			return err
		}
		for n, f := range p.SymbolTable {
			err = nsObj.Set(n, f.Func)
			if err != nil {
				return err
			}
		}
		e.Packages[fmt.Sprintf("G.%s", pkgName)] = p
		// err = ns.Set(pkgName, pkgNs)
		// if err != nil {
		// 	return err
		// }
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
	program, err := parser.ParseFile(e.fileSet, filename, source, 2)
	if err != nil {
		e.Logger.Error(err)
		return err
	}
	_, err = e.LoadScriptWithTimeout(program)
	if err != nil {
		e.Logger.Error(err)
	}
	return err
}

// Exec takes a single string of javascript and evaluates it within the VMs current context.
// It will return both a javascript value object as well as an error if one was encountered
func (e *Engine) Exec(fn string) (otto.Value, error) {
	return e.CallFunctionWithTimeout(fn)
}

// SetLogger overrides the logging interface for this virtual machine.
func (e *Engine) SetLogger(l logger.Logger) error {
	e.Logger = l
	return HijackConsoleLogging(e)
}

func (e *Engine) setGlobalRef() error {
	_, err := e.DeclareNamespace("G")
	if err != nil {
		return err
	}
	err = e.SetConst("_ENGINE", e)
	if err != nil {
		return err
	}
	err = e.SetConst("OS", runtime.GOOS)
	if err != nil {
		return err
	}
	err = e.SetConst("ARCH", runtime.GOARCH)
	if err != nil {
		return err
	}
	err = e.VM.Set("Create", e.createType)
	if err != nil {
		return err
	}
	return nil
}

// EnableAssets injects the core asset handling functions into the engine's runtime
// TODO (gen0cide): Fix asset retrieval to call from vm functions, not the raw translations
func (e *Engine) EnableAssets() error {
	err := e.VM.Set("GetAssetAsString", e.retrieveAssetAsString)
	if err != nil {
		return err
	}
	err = e.VM.Set("GetAssetAsBytes", e.retrieveAssetAsBytes)
	return err
}
