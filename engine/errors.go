package engine

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

// Raise is a convenience method for throwing a javascript runtime error from go space
func (e *Engine) Raise(name string, format string, args ...interface{}) otto.Value {
	msg := fmt.Sprintf(format, args...)
	e.Logger.Errorf("%s error: %s", name, msg)
	return e.VM.MakeCustomError(name, msg)
}
