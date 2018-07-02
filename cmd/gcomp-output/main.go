package main

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"

	"github.com/gen0cide/gscript/engine"
	"github.com/robertkrimen/otto"

	// Importing dependency package testlib for native dependency
	"github.com/gen0cide/gscript/cmd/gcomp/testlib"
)

const (
	// IEDGHJMU9OYSQ3DVH8 holds the data for embedded file test.gs
	IEDGHJMU9OYSQ3DVH8 = `Riv02fG5PuZmdcrVhOZgoIS/9vH8YSRDNgyv/P69yD0MTRr4dBMERHpie99Ij7oli/eY7T8uEp6xof8YuTcd+rqgxoL35vadrmxTZpkWRFeyhrWxZzKaTLikMGZ3IwveknW4Tp5TV89ScyP0QNss9A==`
	// S2816GXGFCHJ6NHAH7 holds the data for embedded file preload.gs
	S2816GXGFCHJ6NHAH7 = `Riv02fG5PuZmdQRvBOKMKVj2Sn2/joXcKfEhbIS5HIciOaDd7VTCQSXgrH6inwOEDPuoHg/txyyJ/1NnycyFJSVzdoe1WNqTpgQLpj0+eoB2CPLxbw/h3ibmazj14t9IN/00kBuXCf/yBXXvcI5r0s3iD/STPKCBqZvdDSkCl1/X+5kLzQmvxfzF5GPNpdpJgwHxwwBsB+ZWCKchXVvJbEQLhSHZKXoNBMx7HzpSNV8Uu50KqagamnlxMO8QSufN83hlell5LxuTH0coAlB2vR7TdRkrxXywTgx5ndN6xRMjkw5YmUFnJ5TNk3mitgqORCzwty6DoEIjf8BhohNKTdQsxFJHUmxWVLxYs551mbs557edR4xhqMtuwzy0hkGcTJEKCxJ3A4TparTEWB4XqIc4RQ==`
)

// PJXTFFGXHCLQTP wraps the genesis VM for test.gs
type PJXTFFGXHCLQTP struct {
	E *engine.Engine
	K []byte
}

// NewPJXTFFGXHCLQTP creates the genesis runtim for script test.gs
func NewPJXTFFGXHCLQTP() *PJXTFFGXHCLQTP {
	te := engine.New("test.gs", "PJXTFFGXHCLQTP", 30, "Execute")
	o := &PJXTFFGXHCLQTP{
		E: te,
	}
	return o
}

// OSDWJMYSZNGJ imports assets into the genesis runtime for script test.gs
func (o *PJXTFFGXHCLQTP) OSDWJMYSZNGJ() error {
	// __ENTRYPOINT = test.gs
	o.E.Imports["__ENTRYPOINT"] = func() []byte {
		// Unpacker wrapper with const declared above
		return PJXTFFGXHCLQTPD(IEDGHJMU9OYSQ3DVH8)
	}
	// __PRELOAD = preload.gs
	o.E.Imports["__PRELOAD"] = func() []byte {
		// Unpacker wrapper with const declared above
		return PJXTFFGXHCLQTPD(S2816GXGFCHJ6NHAH7)
	}
	return nil
}

// EWOZDXVKGBCL imports the runtime preload library for script test.gs
func (o *PJXTFFGXHCLQTP) EWOZDXVKGBCL() error {
	return o.E.LoadScript("preload.js", o.E.Imports["__PRELOAD"]())
}

// WTZYBFXADKKL loads the script test.gs into the genesis VM
func (o *PJXTFFGXHCLQTP) WTZYBFXADKKL() error {
	return o.E.LoadScript("test.gs", o.E.Imports["__ENTRYPOINT"]())
}

// RNJBZXSUVYXD injects the dynamically linked native functions into the genesis VM for script test.gs
func (o *PJXTFFGXHCLQTP) RNJBZXSUVYXD() error {
	var err error
	err = nil

	// -- BEGIN NATIVE PACKAGE IMPORTS
	// Importing github.com/gen0cide/gscript/cmd/gcomp/testlib under namespace testlib
	_nptestlib := &engine.NativePackage{
		ImportPath:  "github.com/gen0cide/gscript/cmd/gcomp/testlib",
		Name:        "testlib",
		SymbolTable: map[string]*engine.NativeFunc{},
	}
	// Adding function pointer for native function testlib.Test1 to symbol table
	_nf0 := &engine.NativeFunc{
		Name: "Test1",
		Func: o.PJXTFFGXHCLQTPntjqpgtlctrieoeq,
	}
	_nptestlib.SymbolTable["Test1"] = _nf0
	// Injecting native package into the genesis VM
	err = o.E.ImportNativePackage("testlib", _nptestlib)
	if err != nil {
		return err
	}
	// -- END NATIVE PACKAGE IMPORTS
	return err
}

// PJXTFFGXHCLQTPODK returns the decryption key for the genesis VM assets
func PJXTFFGXHCLQTPODK() []byte {
	return []byte("8A4i0SzL8ZhsohQvFzrLKtLs55eVwWcK")
}

// PJXTFFGXHCLQTPD is the decoding function for embedded assets for script test.gs
func PJXTFFGXHCLQTPD(s string) []byte {
	b, err := aes.NewCipher(PJXTFFGXHCLQTPODK())
	if err != nil {
		return []byte{}
	}
	db1 := new(bytes.Buffer)
	db2 := new(bytes.Buffer)
	src := bytes.NewReader([]byte(s))
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])
	decoder := base64.NewDecoder(base64.StdEncoding, src)
	encReader := &cipher.StreamReader{S: stream, R: decoder}
	if _, err := io.Copy(db1, encReader); err != nil {
		return []byte{}
	}
	gzr, err := gzip.NewReader(db1)
	if err != nil {
		return []byte{}
	}
	_, err = io.Copy(db2, gzr)
	if err != nil {
		return []byte{}
	}
	gzr.Close()
	return db2.Bytes()
}

// Declaring Dynamically Linked Function Wrappers

// PJXTFFGXHCLQTPntjqpgtlctrieoeq is the linker function for testlib.Test1
func (o *PJXTFFGXHCLQTP) PJXTFFGXHCLQTPntjqpgtlctrieoeq(call otto.FunctionCall) otto.Value {
	// Argument Sanity Checks
	if len(call.ArgumentList) > 1 {
		o.E.Logger.Errorf("too many arguments passed to function %s at %s", "Test1", call.CallerLocation())
		return call.Otto.MakeCustomError("function error", "too many arguments passed into function")
	}
	if len(call.ArgumentList) < 1 {
		o.E.Logger.Errorf("too few arguments passed to function %s at %s", "Test1", call.CallerLocation())
		return call.Otto.MakeCustomError("function error", "too few arguments passed into function")
	}

	// Native Function Argument #0
	var a0 string
	ra0, err := call.Argument(0).Export()
	if err != nil {
		o.E.Logger.Errorf("could not export argument %d of function %s at %s", 0, "Test1", call.CallerLocation())
		o.E.Logger.Error(err)
		return call.Otto.MakeCustomError("function error", fmt.Sprintf("could not translate argument %d into go value", 0))
	}
	switch v := ra0.(type) {
	case string:
		a0 = ra0.(string)
	default:
		errMsg := fmt.Sprintf("Argument type mismatch: expected %s, got %T", "string", v)
		o.E.Logger.Errorf("argument type conversion error: %s", errMsg)
		return call.Otto.MakeCustomError("function error", errMsg)
	}

	// Native Function Return #0
	var r0 *url.URL
	// Native Function Return #1
	var r1 error

	// Call the native function
	r0, r1 = testlib.Test1(a0)

	// This function has multiple returns - injecting into a JS array for single return context compatibility
	jsObj, err := call.Otto.Object(`[]`)
	if err != nil {
		errMsg := fmt.Sprintf("could not make an array object for multiple assignment return")
		o.E.Logger.Errorf(errMsg)
		return call.Otto.MakeCustomError("runtime error", errMsg)
	}

	// Return Value #0
	err = jsObj.Set("0", r0)
	if err != nil {
		errMsg := fmt.Sprintf("could not add element %d of type %T to js return (err: %s)", 0, r0, err.Error())
		o.E.Logger.Errorf(errMsg)
		return call.Otto.MakeCustomError("runtime error", errMsg)
	}
	// Return Value #1
	err = jsObj.Set("1", r1)
	if err != nil {
		errMsg := fmt.Sprintf("could not add element %d of type %T to js return (err: %s)", 1, r1, err.Error())
		o.E.Logger.Errorf(errMsg)
		return call.Otto.MakeCustomError("runtime error", errMsg)
	}

	// Return the generated object
	return jsObj.Value()
}
