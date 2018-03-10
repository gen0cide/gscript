package engine

import (
	"path/filepath"
	"time"

	"github.com/robertkrimen/otto"
)

func (e *Engine) VMIsVM(call otto.FunctionCall) otto.Value {
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

func (e *Engine) VMHalt(call otto.FunctionCall) otto.Value {
	e.Halted = true
	e.VM.Interrupt <- func() {
		panic(errTimeout)
	}
	return otto.FalseValue()
}

func (e *Engine) VMAsset(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "Asset").Error("Too many arguments in call.")
		return otto.FalseValue()
	}

	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "Asset").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var assetName string

	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "Asset").Error("Could not export field: %s", "assetName")
		return otto.FalseValue()
	}

	switch v := rawArg0.(type) {
	case string:
		assetName = rawArg0.(string)
	default:
		e.Logger.WithField("function", "Asset").Error("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	fileData, err := e.Asset(assetName)

	rawVMRet := VMResponse{}

	rawVMRet["fileData"] = fileData

	rawVMRet["err"] = err

	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "Asset").Error("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}

	return vmRet

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

func (e *Engine) VMPathTester(call otto.FunctionCall) otto.Value {
	arg, err := call.Argument(0).ToString()
	if err != nil {
		e.Logger.Errorf("Path Error: %s", err.Error())
		return otto.NullValue()
	}

	argPath := filepath.Clean(arg)
	e.Logger.Infof("Original: %s", arg)
	e.Logger.Infof("Cleaned:  %s", argPath)
	return otto.NullValue()
}
