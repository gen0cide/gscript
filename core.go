package gscript

import (
	"github.com/robertkrimen/otto"

	// Include Underscore In Otto :)
	_ "github.com/robertkrimen/otto/underscore"
)

var (
	Debugger      = true
	DefaultScript = ` // genesis script

var helloWorld = "helloworld";
var foo = MD5(helloWorld);
VMLogInfo(foo);
Sleep(1);
VMLogInfo("Finished Sleeping");
Halt();

`
)

func (e *Engine) VMBeforeDeploy(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMDeploy(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMAfterDeploy(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMOnError(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}
