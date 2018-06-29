package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"

	"github.com/gen0cide/gscript/cmd/gcomp/testlib"
	"github.com/gen0cide/gscript/engine"
	"github.com/robertkrimen/otto"
)

const (
	igbwffbybohxbaejsf = `<nil>`

	ymodgqdddystjokxaa = `<nil>`
)

func main() {
	return
}

type twpnuljxkqie struct {
	E *engine.Engine
}

func Newtwpnuljxkqie() *twpnuljxkqie {
	te := engine.New("test.gs", "twpnuljxkqie", 30, "Execute")
	o := &twpnuljxkqie{
		E: te,
	}
	return o
}

func (o *twpnuljxkqie) FKDDGHTJYERK() error {

	o.E.Imports["__ENTRYPOINT"] = func() []byte {
		return _twpnuljxkqieD(igbwffbybohxbaejsf)
	}

	o.E.Imports["__PRELOAD"] = func() []byte {
		return _twpnuljxkqieD(ymodgqdddystjokxaa)
	}

	return nil
}

func (o *twpnuljxkqie) LPVVQPVWMMKE() error {
	return o.E.LoadScript("preload.js", o.E.Imports["__PRELOAD"]())
}

func (o *twpnuljxkqie) XRMTUHSNFWCP() error {
	return o.E.LoadScript("test.gs", o.E.Imports["__ENTRYPOINT"]())
}

func (o *twpnuljxkqie) JGZUARRBZLJG() error {
	var err error
	_ = err

	_nptestlib := &engine.NativePackage{
		ImportPath:  "github.com/gen0cide/gscript/cmd/gcomp/testlib",
		Name:        "testlib",
		SymbolTable: map[string]*engine.NativeFunc{},
	}

	_nf0 := &engine.NativeFunc{
		Name: "Test1",
		Func: o.twpnuljxkqieuympgcddhaanxevm,
	}
	_nptestlib.SymbolTable["Test1"] = _nf0

	err = o.E.ImportNativePackage("testlib", _nptestlib)
	if err != nil {
		return err
	}

	return nil
}

func _twpnuljxkqieD(s string) []byte {
	db := new(bytes.Buffer)
	src := bytes.NewReader([]byte(s))
	decoder := base64.NewDecoder(base64.StdEncoding, src)
	gzr, err := gzip.NewReader(decoder)
	if err != nil {
		return []byte{}
	}
	_, err = io.Copy(db, gzr)
	if err != nil {
		return []byte{}
	}
	gzr.Close()
	return db.Bytes()
}

func (o *twpnuljxkqie) twpnuljxkqieuympgcddhaanxevm(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		o.E.Logger.Errorf("too many arguments passed to function %s at %s", "Test1", call.CallerLocation())
		return call.Otto.MakeCustomError("function error", "too many arguments passed into function")
	}
	if len(call.ArgumentList) < 1 {
		o.E.Logger.Errorf("too few arguments passed to function %s at %s", "Test1", call.CallerLocation())
		return call.Otto.MakeCustomError("function error", "too few arguments passed into function")
	}
	var a0 string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		o.E.Logger.Errorf("could not export argument %d of function %s at %s", 0, "Test1", call.CallerLocation())
		o.E.Logger.Error(err)
		return call.Otto.MakeCustomError("function error", fmt.Sprintf("could not translate argument %d into go value", 0))
	}
	switch v := rawArg0.(type) {
	case string:
		a0 = rawArg0.(string)
	default:
		errMsg := fmt.Sprintf("Argument type mismatch: expected %s, got %T", "string", v)
		o.E.Logger.Errorf("argument type conversion error: %s", errMsg)
		return call.Otto.MakeCustomError("function error", errMsg)
	}

	var r0 *url.URL

	var r1 error

	r0, r1 = testlib.Test1(a0)

	jsObj, err := call.Otto.Object(`[]`)
	if err != nil {
		errMsg := fmt.Sprintf("could not make an array object for multiple assignment return")
		o.E.Logger.Errorf(errMsg)
		return call.Otto.MakeCustomError("runtime error", errMsg)
	}

	err = jsObj.Set("0", r0)
	if err != nil {
		errMsg := fmt.Sprintf("could not add element %d of type %T to js return (err: %s)", 0, r0, err.Error())
		o.E.Logger.Errorf(errMsg)
		return call.Otto.MakeCustomError("runtime error", errMsg)
	}

	err = jsObj.Set("1", r1)
	if err != nil {
		errMsg := fmt.Sprintf("could not add element %d of type %T to js return (err: %s)", 1, r1, err.Error())
		o.E.Logger.Errorf(errMsg)
		return call.Otto.MakeCustomError("runtime error", errMsg)
	}

	return jsObj.Value()

}
