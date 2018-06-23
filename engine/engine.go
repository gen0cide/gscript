package engine

import (
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

// NativeFunc is a type alias for native function wrappers
type NativeFunc func(call otto.FunctionCall) otto.Value

// DeclareNamespace adds an empty namespace to the virtual machine.
// Caution! will overwrite any values at existing namespace.
func (e *Engine) DeclareNamespace(namespace string) error {
	_, err := e.VM.Object(fmt.Sprintf("%s = {}", namespace))
	return err
}

// AddNativeFuncToNamespace adds func (fn) to namespace (namespace) with name (name) within the virtual machine
func (e *Engine) AddNativeFuncToNamespace(namespace, name string, fn NativeFunc) error {
	return nil
}

// SetConst defines a const value within the virtual machine
func (e *Engine) SetConst(name string, value interface{}) error {
	return nil
}
