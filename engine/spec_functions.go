package engine

import "github.com/robertkrimen/otto"

func (e *Engine) VMLogTester(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Debug("testing debug")
	e.Logger.WithField("trace", "true").Info("testing info")
	e.Logger.WithField("trace", "true").Warn("testing warning")
	e.Logger.WithField("trace", "true").Error("testing error")

	return otto.NullValue()
}
