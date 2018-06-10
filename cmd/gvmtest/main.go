package main

import (
	"net/url"

	"github.com/gen0cide/gscript/debugger"
	"github.com/robertkrimen/otto"
)

func main() {
	dbg := debugger.New("")
	dbg.SetupDebugEngine()
	dbg.Engine.VM.Set("ParseURL", vmParseURL)
	dbg.Engine.VM.Set("FixURL", vmFixURL)
	dbg.InteractiveSession()
}

func FixURL(u *url.URL) {
	u.Host = "getrekt.com"
}

func vmParseURL(call otto.FunctionCall) otto.Value {
	u, _ := call.Argument(0).ToString()
	newURL, _ := url.Parse(u)
	retVal, _ := call.Otto.ToValue(newURL)
	return retVal
}

func vmFixURL(call otto.FunctionCall) otto.Value {
	u, _ := call.Argument(0).Export()
	FixURL(u.(*url.URL))
	return otto.TrueValue()
}
