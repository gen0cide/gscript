package gscript

import (
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/robertkrimen/otto"
)

func (e *Engine) SessionCompleter(d prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func (e *Engine) DebugConsole(call otto.FunctionCall) otto.Value {
	e.InteractiveSession()
	return otto.TrueValue()
}

func (e *Engine) SessionExecutor(in string) {
	newIn := strings.TrimSpace(in)
	if newIn == "exit" || newIn == "quit" {
		os.Exit(0)
	}
	val, err := e.VM.Run(newIn)
	if err != nil {
		e.LogErrorf("Console Error: %s", err.Error())
	}
	retVal, err := val.Export()
	if retVal != nil {
		e.LogInfof(">>> %v", retVal)
	}
}

func (e *Engine) InteractiveSession() {
	if e.Logger == nil {
		e.EnableLogging()
	}
	e.Logger.DisabledInfo = true
	p := prompt.New(
		e.SessionExecutor,
		e.SessionCompleter,
		prompt.OptionPrefix("gscript> "),
		prompt.OptionTitle("Genesis Scripting Engine Console"),
	)
	p.Run()
}
