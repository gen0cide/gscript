package gscript

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/robertkrimen/otto"
)

func (e *Engine) LogWarn(i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Warn(i...)
	}
}

func (e *Engine) LogError(i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Error(i...)
	}
}

func (e *Engine) LogDebug(i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Debug(i...)
	}
}

func (e *Engine) LogCrit(i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Crit(i...)
	}
}

func (e *Engine) LogInfo(i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Log(i...)
	}
}

func (e *Engine) LogWarnf(fmtString string, i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Warnf(fmtString, i...)
	}
}

func (e *Engine) LogErrorf(fmtString string, i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Errorf(fmtString, i...)
	}
}

func (e *Engine) LogDebugf(fmtString string, i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Debugf(fmtString, i...)
	}
}

func (e *Engine) LogCritf(fmtString string, i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Critf(fmtString, i...)
	}
}

func (e *Engine) LogInfof(fmtString string, i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Logf(fmtString, i...)
	}
}

func (e *Engine) VMLogWarn(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.LogErrorf("Logging Error - argument couldn't be converted to string: %s", spew.Sdump(arg))
				continue
			}
			e.Logger.Warnf("%s", newArg)
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogError(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.LogErrorf("Logging Error - argument couldn't be converted to string: %s", spew.Sdump(arg))
				continue
			}
			e.Logger.Errorf("%s", newArg)
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogDebug(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.LogErrorf("Logging Error - argument couldn't be converted to string: %s", spew.Sdump(arg))
				continue
			}
			e.Logger.Debugf("%s", newArg)
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogCrit(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.LogErrorf("Logging Error - argument couldn't be converted to string: %s", spew.Sdump(arg))
				continue
			}
			e.Logger.Critf("%s", newArg)
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogInfo(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.LogErrorf("Logging Error - argument couldn't be converted to string: %s", spew.Sdump(arg))
				continue
			}
			e.Logger.Logf("%s", newArg)
		}
	}
	return otto.Value{}
}
