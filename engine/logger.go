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

// Logger is an interface that allows for generic logging capabilities to be defined
// for the genesis VM. In production, this will be a nil logger, where as while
// developing or using the console, this can be altered without having to include
// the logging dependencies in the base engine.
type Logger interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
}

// NullLogger is a built in type that implements the Logger interface to prevent scripts
// from writing output to the screen during execution (default logging behavior of binary)
type NullLogger struct{}

// Print implements the Logger interface type to prevent debug output
func (n *NullLogger) Print(args ...interface{}) {
	_ = args
	return
}

// Printf implements the Logger interface type to prevent debug output
func (n *NullLogger) Printf(format string, args ...interface{}) {
	_ = format
	_ = args
	return
}

// Println implements the Logger interface type to prevent debug output
func (n *NullLogger) Println(args ...interface{}) {
	_ = args
	return
}

// Debug implements the Logger interface type to prevent debug output
func (n *NullLogger) Debug(args ...interface{}) {
	_ = args
	return
}

// Debugf implements the Logger interface type to prevent debug output
func (n *NullLogger) Debugf(format string, args ...interface{}) {
	_ = format
	_ = args
	return
}

// Debugln implements the Logger interface type to prevent debug output
func (n *NullLogger) Debugln(args ...interface{}) {
	_ = args
	return
}

// Info implements the Logger interface type to prevent debug output
func (n *NullLogger) Info(args ...interface{}) {
	_ = args
	return
}

// Infof implements the Logger interface type to prevent debug output
func (n *NullLogger) Infof(format string, args ...interface{}) {
	_ = format
	_ = args
	return
}

// Infoln implements the Logger interface type to prevent debug output
func (n *NullLogger) Infoln(args ...interface{}) {
	_ = args
	return
}

// Warn implements the Logger interface type to prevent debug output
func (n *NullLogger) Warn(args ...interface{}) {
	_ = args
	return
}

// Warnf implements the Logger interface type to prevent debug output
func (n *NullLogger) Warnf(format string, args ...interface{}) {
	_ = format
	_ = args
	return
}

// Warnln implements the Logger interface type to prevent debug output
func (n *NullLogger) Warnln(args ...interface{}) {
	_ = args
	return
}

// Error implements the Logger interface type to prevent debug output
func (n *NullLogger) Error(args ...interface{}) {
	_ = args
	return
}

// Errorf implements the Logger interface type to prevent debug output
func (n *NullLogger) Errorf(format string, args ...interface{}) {
	_ = format
	_ = args
	return
}

// Errorln implements the Logger interface type to prevent debug output
func (n *NullLogger) Errorln(args ...interface{}) {
	_ = args
	return
}

// Fatal implements the Logger interface type to prevent debug output
func (n *NullLogger) Fatal(args ...interface{}) {
	_ = args
	return
}

// Fatalf implements the Logger interface type to prevent debug output
func (n *NullLogger) Fatalf(format string, args ...interface{}) {
	_ = format
	_ = args
	return
}

// Fatalln implements the Logger interface type to prevent debug output
func (n *NullLogger) Fatalln(args ...interface{}) {
	_ = args
	return
}

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
