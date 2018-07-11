package engine

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

func (e *Engine) retrieveAssetAsString(name string) (string, error) {
	if e.Imports[name] == nil {
		return "", fmt.Errorf("could not locate asset %s", name)
	}
	ret := e.Imports[name]()
	return string(ret), nil
}

func (e *Engine) retrieveAssetAsBytes(name string) ([]byte, error) {
	ret := []byte{}
	if e.Imports[name] == nil {
		return ret, fmt.Errorf("could not locate asset %s", name)
	}
	ret = e.Imports[name]()
	return ret, nil
}

func (e *Engine) vmRetrieveAssetAsString(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 1 {
		return e.Raise("argument", "argument count mismatch for %s", "GetAssetAsString")
	}

	var a0 string
	ra0, err := call.Argument(0).Export()
	if err != nil {
		return e.Raise("jsexport", "could not export argument %d of function %s", 0, "GetAssetAsString")
	}
	switch v := ra0.(type) {
	case string:
		a0 = v
	default:
		return e.Raise("type conversion", "argument type mismatch - expected %s, got %T", "string", v)
	}

	var r0 string
	var r1 error

	r0, r1 = e.retrieveAssetAsString(a0)

	if r1 != nil {
		return e.Raise("asset", "retrieving asset from asset table failed - %v", r1)
	}

	retVal, err := call.Otto.ToValue(r0)
	if err != nil {
		return e.Raise("return", "conversion failed for return 0 (type=%T) - %v", r0, err)
	}

	return retVal
}

func (e *Engine) vmRetrieveAssetAsBytes(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 1 {
		return e.Raise("argument", "argument count mismatch for %s", "GetAssetAsBytes")
	}

	var a0 string
	ra0, err := call.Argument(0).Export()
	if err != nil {
		return e.Raise("jsexport", "could not export argument %d of function %s", 0, "GetAssetAsBytes")
	}
	switch v := ra0.(type) {
	case string:
		a0 = v
	default:
		return e.Raise("type conversion", "argument type mismatch - expected %s, got %T", "string", v)
	}

	var r0 []byte
	var r1 error

	r0, r1 = e.retrieveAssetAsBytes(a0)

	if r1 != nil {
		return e.Raise("asset", "retrieving asset from asset table failed - %v", r1)
	}

	retVal, err := call.Otto.ToValue(r0)
	if err != nil {
		return e.Raise("return", "conversion failed for return 0 (type=%T) - %v", r0, err)
	}

	return retVal
}
