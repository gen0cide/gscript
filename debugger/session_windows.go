// +build windows

package debugger

import (
	"github.com/robertkrimen/otto"
)

func (e *Engine) DebugConsole(call otto.FunctionCall) otto.Value {
	e.InteractiveSession()
	return otto.TrueValue()
}

func (e *Engine) SessionExecutor(in string) {
	return
}

func (e *Engine) InteractiveSession() {
	e.LogErrorf("Genesis REPL is not available on Windows.")
}
