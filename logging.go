package gscript

import "github.com/robertkrimen/otto"

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
			e.Logger.Warn(arg.ToString())
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogError(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			e.Logger.Error(arg.ToString())
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogDebug(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			e.Logger.Debug(arg.ToString())
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogCrit(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			e.Logger.Crit(arg.ToString())
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogInfo(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			e.Logger.Log(arg.ToString())
		}
	}
	return otto.Value{}
}
