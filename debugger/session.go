package debugger

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"github.com/gen0cide/gscript/compiler"
	"github.com/gen0cide/gscript/engine"
	"github.com/gen0cide/gscript/logging"
	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
)

type Debugger struct {
	Engine *engine.Engine
	Logger *logrus.Logger
	Prompt *prompt.Prompt
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
	val, err := d.Engine.VM.Eval(newIn)
	if err != nil {
		d.Engine.Logger.Errorf("Console Error: %s", err.Error())
	}
	retVal, _ := val.Export()
	if retVal != nil {
		fmt.Printf(">>> %v\n", retVal)
	}
}

func (d *Debugger) LoadScript(script, filename string) error {
	ml := compiler.ParseMacros(string(script), d.Logger.WithField("file", filename))
	if ml == nil {
		d.Logger.WithField("file", filename).Fatalf("Compiler could not parse macros!")
	}
	d.Engine.SetTimeout(ml.Timeout)
	if ml.Timeout != compiler.DEFAULT_TIMEOUT {
		d.Logger.Infof("Timeout has been changed. Timeout is now %d seconds.", ml.Timeout)
	}
	for _, i := range ml.LocalFiles {
		name := filepath.Base(i)
		data, err := ioutil.ReadFile(i)
		if err != nil {
			d.Logger.WithField("file", filename).Fatalf("Compiler could not load file: %s", err.Error())
		}
		d.Engine.AddImport(name, func() []byte {
			return data
		})
		d.Logger.WithField("file", filename).Debugf("Importing Local File: %s", name)
	}
	for _, i := range ml.RemoteFiles {
		name := filepath.Base(i)
		data, err := ioutil.ReadFile(i)
		if err != nil {
			d.Logger.WithField("file", filename).Fatalf("Compiler could not load file: %s", err.Error())
		}
		d.Engine.AddImport(name, func() []byte {
			return data
		})
		d.Logger.WithField("file", filename).Debugf("Importing Remote File: %s", name)
	}
	return d.Engine.LoadScript([]byte(script))
}

func (d *Debugger) InteractiveSession() {
	p := prompt.New(
		d.SessionExecutor,
		d.SessionCompleter,
		prompt.OptionPrefix("gscript> "),
		prompt.OptionPrefixTextColor(prompt.Red),
		prompt.OptionTitle("Genesis Scripting Engine Console"),
	)
	d.Prompt = p
	entryText := []string{
		fmt.Sprintf(
			"%s %s %s %s",
			color.HiWhiteString("***"),
			color.HiRedString("GSCRIPT"),
			color.YellowString("INTERACTIVE SHELL"),
			color.HiWhiteString("***"),
		),
		fmt.Sprintf("%s %s", color.YellowString("NOTE:"), "To exit the debugger, use CONTROL+D"),
	}
	for _, l := range entryText {
		fmt.Fprintf(color.Output, "%s\n", l)
	}
	p.Run()
}
