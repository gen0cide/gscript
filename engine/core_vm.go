package engine

import (
	"time"

	"github.com/robertkrimen/otto"
)

func (e *Engine) VMIsVM(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMHalt(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMImplode(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMMatches(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMAsset(call otto.FunctionCall) otto.Value {
	filename := call.Argument(0).String()
	if len(filename) == 0 {
		e.Logger.WithField("trace", "true").Errorf("Empty Asset Call")
		return otto.FalseValue()
	}
	if dataFunc, ok := e.Imports[filename]; ok {
		byteData := dataFunc()
		vmVal, err := e.VM.ToValue(byteData)
		if err != nil {
			e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
			return otto.FalseValue()
		}
		return vmVal
	}
	e.Logger.WithField("trace", "true").Errorf("Asset File Not Found: %s", filename)
	return otto.FalseValue()
}

func (e *Engine) VMTimestamp(call otto.FunctionCall) otto.Value {
	ts := time.Now().Unix()
	ret, err := otto.ToValue(ts)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
	}
	return ret
}

func (e *Engine) VMIsAWS(call otto.FunctionCall) otto.Value {
	respCode, _, err := HTTPGetFile("http://169.254.169.254/latest/meta-data/")
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Net error: %s", err.Error())
		return otto.FalseValue()
	} else if respCode == 200 {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMLogTester(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Debug("testing debug")
	e.Logger.WithField("trace", "true").Info("testing info")
	e.Logger.WithField("trace", "true").Warn("testing warning")
	e.Logger.WithField("trace", "true").Error("testing error")

	return otto.NullValue()
}
