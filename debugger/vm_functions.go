package debugger

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"github.com/robertkrimen/otto"
)

func (d *Debugger) vmDebugConsole(call otto.FunctionCall) otto.Value {
	d.VM.SetLogger(d.Logger)
	d.runDebugger()
	d.VM.SetLogger(d.OldLogger)
	return otto.UndefinedValue()
}

func (d *Debugger) vmVarTable(call otto.FunctionCall) otto.Value {
	vars := d.AvailableVars()
	for ns, va := range vars {
		d.Logger.Infof(">>> %s Package\n\t%s\n", ns, strings.Join(va, "\n\t"))
	}
	return otto.UndefinedValue()
}

func (d *Debugger) vmConstTable(call otto.FunctionCall) otto.Value {
	consts := d.AvailableConsts()
	for ns, cs := range consts {
		d.Logger.Infof(">>> %s Package\n\t%s\n", ns, strings.Join(cs, "\n\t"))
	}
	return otto.UndefinedValue()
}

func (d *Debugger) vmTypeTable(call otto.FunctionCall) otto.Value {
	types := d.AvailableTypes()
	for ns, ts := range types {
		d.Logger.Infof(">>> %s Package\n\t%s\n", ns, strings.Join(ts, "\n\t"))
	}
	return otto.UndefinedValue()
}

func (d *Debugger) vmSymbolTable(call otto.FunctionCall) otto.Value {
	sym := d.AvailableFuncs()
	for ns, funcs := range sym {
		d.Logger.Infof(">>> %s Package\n\t%s\n", ns, strings.Join(funcs, "\n\t"))
	}
	return otto.UndefinedValue()
}

func (d *Debugger) vmPackageDocs(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 1 {
		return d.VM.Raise("arg", "must provide one string argument to Docs()")
	}
	val, err := call.Argument(0).Export()
	if err != nil {
		return d.VM.Raise("jsexport", "coult not convert argument number 0")
	}

	realval, ok := val.(string)
	if !ok {
		return d.VM.Raise("type", "argument was not of type string")
	}

	consts := d.AvailableConsts()
	types := d.AvailableTypes()
	funcs := d.AvailableFuncs()
	vars := d.AvailableVars()

	cl, clok := consts[realval]
	tl, tlok := types[realval]
	fl, flok := funcs[realval]
	vl, vlok := vars[realval]

	if !clok && !tlok && !flok && !vlok {
		return d.VM.Raise("undefined", "package is not defined in this genesis engine")
	}

	title := fmt.Sprintf(">> Package Documentation: %s\n\n", realval)

	contstext := fmt.Sprintf("\n  %s\n\t%s\n", color.HiYellowString("-- CONSTS --"), strings.Join(cl, "\n\t"))
	varstext := fmt.Sprintf("\n  %s\n\t%s\n", color.HiYellowString("-- VARS --"), strings.Join(vl, "\n\t"))
	typestext := fmt.Sprintf("\n  %s\n\t%s\n", color.HiYellowString("-- TYPES --"), strings.Join(tl, "\n\t"))
	funcstext := fmt.Sprintf("\n  %s\n\t%s\n", color.HiYellowString("-- FUNCS --"), strings.Join(fl, "\n\t"))

	finaltext := strings.Join([]string{title, contstext, varstext, typestext, funcstext}, "")

	d.Logger.Infof("%s", finaltext)
	return otto.UndefinedValue()
}

func (d *Debugger) vmTypeChecker(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		return d.VM.Raise("arg", "no argument provided")
	} else if len(call.ArgumentList) == 1 {
		val, err := call.Argument(0).Export()
		if err != nil {
			return d.VM.Raise("jsexport", "could not convert argument number 0")
		}
		retVal, _ := call.Otto.ToValue(spew.Sdump(val))
		return retVal
	} else {
		return d.VM.Raise("arg", "too many arguments provided")
	}
}
