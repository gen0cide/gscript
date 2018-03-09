package engine

import (
	"strconv"

	"github.com/robertkrimen/otto"
)

func (e *Engine) VMLogWarn(call otto.FunctionCall) otto.Value {
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

func (e *Engine) VMLogError(call otto.FunctionCall) otto.Value {
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

func (e *Engine) VMLogDebug(call otto.FunctionCall) otto.Value {
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

func (e *Engine) VMLogCrit(call otto.FunctionCall) otto.Value {
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

func (e *Engine) VMLogInfo(call otto.FunctionCall) otto.Value {
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
