package engine

import (
	"strconv"

	"github.com/robertkrimen/otto"
)

func (e *Engine) initializeLogger() {
	e.VM.Set("LogDebug", e.vmLogDebug)
	e.VM.Set("LogInfo", e.vmLogInfo)
	e.VM.Set("LogWarn", e.vmLogWarn)
	e.VM.Set("LogError", e.vmLogError)
	e.VM.Set("LogFatal", e.vmLogFatal)
}

func (e *Engine) vmLogWarn(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.Logger.WithField(
					"script",
					e.Name,
				).WithField(
					"line",
					strconv.Itoa(e.VM.Context().Line),
				).WithField(
					"caller",
					e.VM.Context().Callee,
				).Errorf("Parameter parsing error: %s", err.Error())
				continue
			}
			e.Logger.WithField(
				"script",
				e.Name,
			).WithField(
				"line",
				strconv.Itoa(e.VM.Context().Line),
			).WithField(
				"caller",
				e.VM.Context().Callee,
			).Warnf("%s", newArg)
		}
	}
	return otto.Value{}
}

func (e *Engine) vmLogError(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.Logger.WithField(
					"script",
					e.Name,
				).WithField(
					"line",
					strconv.Itoa(e.VM.Context().Line),
				).WithField(
					"caller",
					e.VM.Context().Callee,
				).Errorf("Parameter parsing error: %s", err.Error())
				continue
			}
			e.Logger.WithField(
				"script",
				e.Name,
			).WithField(
				"line",
				strconv.Itoa(e.VM.Context().Line),
			).WithField(
				"caller",
				e.VM.Context().Callee,
			).Errorf("%s", newArg)
		}
	}
	return otto.Value{}
}

func (e *Engine) vmLogDebug(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.Logger.WithField(
					"script",
					e.Name,
				).WithField(
					"line",
					strconv.Itoa(e.VM.Context().Line),
				).WithField(
					"caller",
					e.VM.Context().Callee,
				).Errorf("Parameter parsing error: %s", err.Error())
				continue
			}
			e.Logger.WithField(
				"script",
				e.Name,
			).WithField(
				"line",
				strconv.Itoa(e.VM.Context().Line),
			).WithField(
				"caller",
				e.VM.Context().Callee,
			).Debugf("%s", newArg)
		}
	}
	return otto.Value{}
}

func (e *Engine) vmLogFatal(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.Logger.WithField(
					"script",
					e.Name,
				).WithField(
					"line",
					strconv.Itoa(e.VM.Context().Line),
				).WithField(
					"caller",
					e.VM.Context().Callee,
				).Errorf("Parameter parsing error: %s", err.Error())
				continue
			}
			e.Logger.WithField(
				"script",
				e.Name,
			).WithField(
				"line",
				strconv.Itoa(e.VM.Context().Line),
			).WithField(
				"caller",
				e.VM.Context().Callee,
			).Fatalf("%s", newArg)
		}
	}
	return otto.Value{}
}

func (e *Engine) vmLogInfo(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.Logger.WithField(
					"script",
					e.Name,
				).WithField(
					"line",
					strconv.Itoa(e.VM.Context().Line),
				).WithField(
					"caller",
					e.VM.Context().Callee,
				).Errorf("Parameter parsing error: %s", err.Error())
				continue
			}
			e.Logger.WithField(
				"script",
				e.Name,
			).WithField(
				"line",
				strconv.Itoa(e.VM.Context().Line),
			).WithField(
				"caller",
				e.VM.Context().Callee,
			).Infof("%s", newArg)
		}
	}
	return otto.Value{}
}
