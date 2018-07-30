package debugger

import (
	"fmt"

	"github.com/gen0cide/gscript/engine"
	"github.com/gen0cide/gscript/logger"
	"github.com/gen0cide/gscript/logger/standard"
)

// Debugger is a wrapper type for handling interactive debug consoles in the genesis engine
type Debugger struct {
	VM        *engine.Engine
	Logger    logger.Logger
	OldLogger logger.Logger
}

// New returns a new debugger object wrapping the provided engine
func New(e *engine.Engine) *Debugger {
	dbgLogger := standard.NewStandardLogger(nil, "debugger", false, true)
	dbg := &Debugger{
		VM:        e,
		Logger:    dbgLogger,
		OldLogger: e.Logger,
	}
	return dbg
}

// InjectDebugConsole injects the DebugConsole command into the runtime
func (d *Debugger) InjectDebugConsole() error {
	d.VM.VM.Set("_DEBUGGER", d)
	err := d.VM.VM.Set("DebugConsole", d.vmDebugConsole)
	if err != nil {
		return err
	}
	err = d.VM.VM.Set("SymbolTable", d.vmSymbolTable)
	if err != nil {
		return err
	}
	err = d.VM.VM.Set("TypeTable", d.vmTypeTable)
	if err != nil {
		return err
	}
	err = d.VM.VM.Set("ConstTable", d.vmConstTable)
	if err != nil {
		return err
	}
	err = d.VM.VM.Set("VarTable", d.vmVarTable)
	if err != nil {
		return err
	}
	err = d.VM.VM.Set("Docs", d.vmPackageDocs)
	if err != nil {
		return err
	}

	return d.VM.VM.Set("TypeOf", d.vmTypeChecker)
}

// AvailableFuncs returns the current debugger's symbol table
func (d *Debugger) AvailableFuncs() map[string][]string {
	ret := map[string][]string{}
	for name, p := range d.VM.Packages {
		flist := []string{}
		idx := 0
		for _, f := range p.SymbolTable {
			flist = append(flist, fmt.Sprintf("%d) %s", idx, f.Signature))
			idx++
		}
		if len(flist) > 0 {
			ret[name] = flist
		}
	}
	return ret
}

// AvailableTypes generates a type table for the debugger
func (d *Debugger) AvailableTypes() map[string][]string {
	ret := map[string][]string{}
	for name, p := range d.VM.Packages {
		tlist := []string{}
		idx := 0
		for tn := range p.Types {
			tlist = append(tlist, fmt.Sprintf("%d) %s.%s", idx, name, tn))
			idx++
		}
		if len(tlist) > 0 {
			ret[name] = tlist
		}
	}
	return ret
}

// AvailableConsts generates a const symbol table for the debugger
func (d *Debugger) AvailableConsts() map[string][]string {
	ret := map[string][]string{}
	for name, p := range d.VM.Packages {
		clist := []string{}
		idx := 0
		for c, cv := range p.Consts {
			clist = append(clist, fmt.Sprintf("%d) %s.%s = %v", idx, name, c, cv.Value))
			idx++
		}
		if len(clist) > 0 {
			ret[name] = clist
		}
	}
	return ret
}

// AvailableVars generates a var symbol table for the debugger
func (d *Debugger) AvailableVars() map[string][]string {
	ret := map[string][]string{}
	for name, p := range d.VM.Packages {
		vlist := []string{}
		idx := 0
		for vname, va := range p.Vars {
			vlist = append(vlist, fmt.Sprintf("%d) %s.%s (%s)", idx, name, vname, va.Signature))
			idx++
		}
		if len(vlist) > 0 {
			ret[name] = vlist
		}
	}
	return ret
}
