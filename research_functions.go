package gscript

import (
	"github.com/robertkrimen/otto"
	"github.com/davecgh/go-spew/spew"

)

func (e *Engine) VMLocalUserExists(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMProcExistsWithName(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMCanReadFile(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMCanWriteFile(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMCanExecFile(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMFileExists(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMDirExists(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMFileContains(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMIsVM(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMIsAWS(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMHasPublicIP(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMCanMakeTCPConn(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMExpectedDNS(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMCanMakeHTTPConn(call otto.FunctionCall) otto.Value {
	url1 := call.Argument(0)
	url1String, err := url1.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(url1String))
		return otto.FalseValue()
	}
	respCode, _, err := HTTPGetFile(url1String.(string))
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	} else if (respCode != 403 || respCode != 404 || respCode != 500 || respCode != 502 || respCode != 503 || respCode != 504 || respCode != 511) {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMDetectSSLMITM(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMCmdSuccessful(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMCanPing(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMTCPPortInUse(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMUDPPortInUse(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMExistsInPath(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMCanSudo(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMMatches(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMCanSSHLogin(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}
