package engine

import (
	"errors"
	"fmt"
	"strings"

	"github.com/robertkrimen/otto"
)

var (
	// the console functions that belong to the javascript console object
	consoleLogLevels = []string{
		"log",
		"debug",
		"info",
		"error",
		"warn",
	}
)

// converts arguments passed into console into helpful string representations
// otto by default prints objects as [Object object] which is incredibly not
// helpful -.-
func normalizeConsoleArgs(c otto.FunctionCall) string {
	o := []string{}
	jsonNs, err := c.Otto.Object("JSON")
	if err != nil {
		return errors.New("runtime error: could not locate the JSON runtime object").Error()
	}
	for _, i := range c.ArgumentList {
		if i.Class() == "Object" {
			i, err = jsonNs.Call("stringify", i)
			if err != nil {
				o = append(o, err.Error())
				continue
			}
		}
		o = append(o, fmt.Sprintf("%v", i))
	}
	return strings.Join(o, " ")

}

// generates native functions for each of the log levels that use the appropriate logger function
// belonging to the provided engine
func logFuncFactory(e *Engine, level string) func(call otto.FunctionCall) otto.Value {
	prefix := fmt.Sprintf("console.%s >>> ", level)
	switch level {
	case "log":
		return func(call otto.FunctionCall) otto.Value {
			e.Logger.Info(prefix, normalizeConsoleArgs(call))
			return otto.Value{}
		}
	case "debug":
		return func(call otto.FunctionCall) otto.Value {
			e.Logger.Debug(prefix, normalizeConsoleArgs(call))
			return otto.Value{}
		}
	case "info":
		return func(call otto.FunctionCall) otto.Value {
			e.Logger.Info(prefix, normalizeConsoleArgs(call))
			return otto.Value{}
		}
	case "error":
		return func(call otto.FunctionCall) otto.Value {
			e.Logger.Error(prefix, normalizeConsoleArgs(call))
			return otto.Value{}
		}
	case "warn":
		return func(call otto.FunctionCall) otto.Value {
			e.Logger.Warn(prefix, normalizeConsoleArgs(call))
			return otto.Value{}
		}
	default:
		return func(call otto.FunctionCall) otto.Value {
			e.Logger.Info(prefix, normalizeConsoleArgs(call))
			return otto.Value{}
		}
	}
}

// HijackConsoleLogging intercepts the javascript runtimes console logging functions (i.e. console.log)
// and dynamically generates new native function implementations of those build ins that use
// the engine object's Logger interface
func HijackConsoleLogging(e *Engine) error {
	c, err := e.VM.Object(`console`)
	if err != nil {
		return err
	}
	for _, l := range consoleLogLevels {
		err = c.Set(l, logFuncFactory(e, l))
		if err != nil {
			return err
		}
	}
	return nil
}
