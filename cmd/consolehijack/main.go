package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"

	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/repl"
)

var (
	l      = logrus.New()
	levels = []string{
		"log",
		"debug",
		"info",
		"error",
		"warn",
	}
)

type Eng struct {
	Logger *logrus.Logger
	VM     *otto.Otto
}

func normalizeConsoleArgs(c otto.FunctionCall) string {
	o := []string{}
	jsonNs, err := c.Otto.Object("JSON")
	if err != nil {
		return errors.New("runtime error: could not locate the JSON runtime object").Error()
	}
	for _, i := range c.ArgumentList {
		if i.Class() == "Object" {
			i, err = jsonNs.Call("stringify", i, nil, 2)
			if err != nil {
				o = append(o, err.Error())
				continue
			}
		}
		o = append(o, fmt.Sprintf("%v", i))
	}
	return strings.Join(o, " ")

}

func main() {
	vm := &Eng{
		Logger: logrus.New(),
		VM:     otto.New(),
	}
	hijackConsole(vm)
	err := repl.Run(vm.VM)
	if err != nil {
		panic(err)
	}
}

func hijackConsole(e *Eng) {
	c, err := e.VM.Object(`console`)
	if err != nil {
		panic(err)
	}
	for _, l := range levels {
		err = c.Set(l, loggerFactory(e, l))
		if err != nil {
			panic(err)
		}
	}
}

func loggerFactory(e *Eng, level string) func(call otto.FunctionCall) otto.Value {
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
