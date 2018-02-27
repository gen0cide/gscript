package debugger

import (
	"fmt"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/gen0cide/gscript/engine"
	"github.com/gen0cide/gscript/logging"
	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
)

type Debugger struct {
	Engine *engine.Engine
	Logger *logrus.Logger
}

func New(name string) *Debugger {
	logger := logrus.New()
	logger.Formatter = &logging.GSEFormatter{}
	logger.Out = logging.LogWriter{Name: name}
	logger.Level = logrus.DebugLevel
	gse := engine.New(name)
	return &Debugger{
		Engine: gse,
		Logger: logger,
	}
}

func (d *Debugger) SetupDebugEngine() {
	d.Engine.SetLogger(d.Logger)
	d.Engine.CreateVM()
	d.Engine.VM.Set("DebugConsole", d.DebugConsole)
}

func (d *Debugger) SessionCompleter(p prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func (d *Debugger) DebugConsole(call otto.FunctionCall) otto.Value {
	d.InteractiveSession()
	return otto.TrueValue()
}

func (d *Debugger) SessionExecutor(in string) {
	newIn := strings.TrimSpace(in)
	if newIn == "exit" || newIn == "quit" {
		os.Exit(0)
	}
	val, err := d.Engine.VM.Run(newIn)
	if err != nil {
		d.Engine.Logger.Errorf("Console Error: %s", err.Error())
	}
	retVal, _ := val.Export()
	if retVal != nil {
		fmt.Printf(">>> %v\n", retVal)
	}
}

func (d *Debugger) InteractiveSession() {
	p := prompt.New(
		d.SessionExecutor,
		d.SessionCompleter,
		prompt.OptionPrefix("gscript> "),
		prompt.OptionTitle("Genesis Scripting Engine Console"),
	)
	p.Run()
}
