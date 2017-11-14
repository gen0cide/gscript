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

func (e *Engine) BeforeDeploy(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) Deploy(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) AfterDeploy(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) OnError(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}
