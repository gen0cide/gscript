package debugger

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gen0cide/gscript/engine"
	"github.com/gen0cide/gscript/logger"
	"github.com/robertkrimen/otto"
	readline "gopkg.in/readline.v1"
)

// Debugger is a wrapper type for handling interactive debug consoles in the genesis engine
type Debugger struct {
	VM        *engine.Engine
	Logger    engine.Logger
	OldLogger engine.Logger
}

// New returns a new debugger object wrapping the provided engine
func New(e *engine.Engine) *Debugger {
	dbgLogger := logger.NewStandardLogrusLogger(nil, "debugger", false, true)
	dbg := &Debugger{
		VM:        e,
		Logger:    dbgLogger,
		OldLogger: e.Logger,
	}
	return dbg
}

// InjectDebugConsole injects the DebugConsole command into the runtime
func (d *Debugger) InjectDebugConsole() error {
	return d.VM.VM.Set("DebugConsole", d.vmDebugConsole)
}

func (d *Debugger) vmDebugConsole(call otto.FunctionCall) otto.Value {
	d.VM.SetLogger(d.Logger)
	d.runDebugger()
	d.VM.SetLogger(d.OldLogger)
	return otto.UndefinedValue()
}

func (d *Debugger) runDebugger() error {
	prompt := fmt.Sprintf("%s%s", color.HiRedString("gscript"), color.HiWhiteString("> "))
	c := &readline.Config{
		Prompt: prompt,
	}

	rl, err := readline.NewEx(c)
	if err != nil {
		return err
	}
	logger.PrintLogo()
	title := fmt.Sprintf(
		"%s %s %s %s",
		color.HiWhiteString("***"),
		color.HiRedString("GSCRIPT"),
		color.YellowString("INTERACTIVE SHELL"),
		color.HiWhiteString("***"),
	)
	fmt.Fprintf(color.Output, "%s\n", title)
	rl.Refresh()

	for {
		l, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				if d != nil {
					d = nil
					rl.SetPrompt(prompt)
					rl.Refresh()
					continue
				}
				break
			}
			return err
		}
		if l == "" {
			continue
		}
		if l == "exit" {
			break
		}
		s, err := d.VM.VM.Compile("debugger", l)
		if err != nil {
			d.Logger.Errorf("%v", err)
			rl.SetPrompt(prompt)
			rl.Refresh()
			continue
		}
		v, err := d.VM.VM.Eval(s)
		if err != nil {
			if oerr, ok := err.(*otto.Error); ok {
				d.Logger.Error(oerr.String())
			} else {
				d.Logger.Error(err.Error())
			}
		} else {
			rl.Write([]byte(fmt.Sprintf(">>> %s\n", v.String())))
		}
		rl.Refresh()
	}

	return rl.Close()
}
