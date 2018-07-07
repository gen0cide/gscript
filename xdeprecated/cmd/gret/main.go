package main

import (
	"io/ioutil"

	"github.com/gen0cide/gscript/debugger"
	"github.com/robertkrimen/otto"
)

func main() {
	dbg := debugger.New("gret")
	dbg.SetupDebugEngine()
	dbg.Engine.VM.Set("DopeReadFile", vmDopeReadFile)
	dbg.InteractiveSession()
}

type A struct {
	B *B
}

type B struct {
	Name string
}

func (b *B) BobbaFet(s string) (*A, error) {
	a := &A{
		B: &B{
			Name: s,
		},
	}

	return a, nil
}

func (b *B) TestRet() (string, string) {
	return "c", "d"
}

func vmDopeReadFile(call otto.FunctionCall) otto.Value {
	filename, err := call.Argument(0).ToString()
	if err != nil {
		return otto.FalseValue()
	}

	data, funcErr := ioutil.ReadFile(filename)
	jsObj, err := call.Otto.Object(`[]`)
	if err != nil {
		panic(err)
	}

	err = jsObj.Set("0", data)
	if err != nil {
		panic(err)
	}
	err = jsObj.Set("1", funcErr)
	if err != nil {
		panic(err)
	}

	bObj := A{B: &B{Name: "george"}}

	err = jsObj.Set("2", bObj)

	return jsObj.Value()
}
