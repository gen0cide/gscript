// Package engine is the package that implements the gscript Virtual Machine (VM).
//
// VM Functions implemented:
//
// Library core
//
// Functions in core:
//  Asset(assetName) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.Asset
//  DeobfuscateString(str) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.DeobfuscateString
//  Halt() - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.Halt
//  MD5(data) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.MD5
//  ObfuscateString(str) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.ObfuscateString
//  RandomInt(min, max) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.RandomInt
//  RandomMixedCaseString(strlen) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.RandomMixedCaseString
//  RandomString(strlen) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.RandomString
//  StripSpaces(str) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.StripSpaces
//  Timestamp() - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.Timestamp
//  XorBytes(aByteArray, bByteArray) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.XorBytes
//
// Library exec
//
// Functions in exec:
//  ExecuteCommand(baseCmd, cmdArgs) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.ExecuteCommand
//  ForkExecuteCommand(baseCmd, cmdArgs) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.ForkExecuteCommand
//
// Library file
//
// Functions in file:
//  WriteFile(path, fileData, perms) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.WriteFile
//
package engine

import (
	"github.com/robertkrimen/otto"
)

func (e *Engine) CreateVM() {
	e.VM = otto.New()
	e.injectVars()
	e.VM.Set("Asset", e.vmAsset)
	e.VM.Set("DeobfuscateString", e.vmDeobfuscateString)
	e.VM.Set("ExecuteCommand", e.vmExecuteCommand)
	e.VM.Set("ForkExecuteCommand", e.vmForkExecuteCommand)
	e.VM.Set("Halt", e.vmHalt)
	e.VM.Set("MD5", e.vmMD5)
	e.VM.Set("ObfuscateString", e.vmObfuscateString)
	e.VM.Set("RandomInt", e.vmRandomInt)
	e.VM.Set("RandomMixedCaseString", e.vmRandomMixedCaseString)
	e.VM.Set("RandomString", e.vmRandomString)
	e.VM.Set("StripSpaces", e.vmStripSpaces)
	e.VM.Set("Timestamp", e.vmTimestamp)
	e.VM.Set("WriteFile", e.vmWriteFile)
	e.VM.Set("XorBytes", e.vmXorBytes)
	_, err := e.VM.Run(vmPreload)
	if err != nil {
		e.Logger.WithField("trace", "true").Fatalf("Syntax error in preload: %s", err.Error())
	}
	e.initializeLogger()
}

func (e *Engine) vmAsset(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "Asset").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "Asset").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var assetName string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "Asset").WithField("trace", "true").Errorf("Could not export field: %s", "assetName")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		assetName = rawArg0.(string)
	default:
		e.Logger.WithField("function", "Asset").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	fileData, err := e.Asset(assetName)
	rawVMRet := VMResponse{}
	rawVMRet["fileData"] = fileData
	rawVMRet["err"] = err
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "Asset").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmDeobfuscateString(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "DeobfuscateString").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "DeobfuscateString").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var str string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "DeobfuscateString").WithField("trace", "true").Errorf("Could not export field: %s", "str")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		str = rawArg0.(string)
	default:
		e.Logger.WithField("function", "DeobfuscateString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	value := e.DeobfuscateString(str)
	rawVMRet := VMResponse{}
	rawVMRet["value"] = value
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "DeobfuscateString").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmExecuteCommand(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 2 {
		e.Logger.WithField("function", "ExecuteCommand").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 2 {
		e.Logger.WithField("function", "ExecuteCommand").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var baseCmd string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "ExecuteCommand").WithField("trace", "true").Errorf("Could not export field: %s", "baseCmd")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		baseCmd = rawArg0.(string)
	default:
		e.Logger.WithField("function", "ExecuteCommand").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var cmdArgs []string
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "ExecuteCommand").WithField("trace", "true").Errorf("Could not export field: %s", "cmdArgs")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case []string:
		cmdArgs = rawArg1.([]string)
	default:
		e.Logger.WithField("function", "ExecuteCommand").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "[]string", v)
		return otto.FalseValue()
	}
	retObject := e.ExecuteCommand(baseCmd, cmdArgs)
	rawVMRet := VMResponse{}
	rawVMRet["retObject"] = retObject
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "ExecuteCommand").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmForkExecuteCommand(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 2 {
		e.Logger.WithField("function", "ForkExecuteCommand").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 2 {
		e.Logger.WithField("function", "ForkExecuteCommand").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var baseCmd string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "ForkExecuteCommand").WithField("trace", "true").Errorf("Could not export field: %s", "baseCmd")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		baseCmd = rawArg0.(string)
	default:
		e.Logger.WithField("function", "ForkExecuteCommand").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var cmdArgs []string
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "ForkExecuteCommand").WithField("trace", "true").Errorf("Could not export field: %s", "cmdArgs")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case []string:
		cmdArgs = rawArg1.([]string)
	default:
		e.Logger.WithField("function", "ForkExecuteCommand").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "[]string", v)
		return otto.FalseValue()
	}
	pid, execError := e.ForkExecuteCommand(baseCmd, cmdArgs)
	rawVMRet := VMResponse{}
	rawVMRet["pid"] = pid
	rawVMRet["execError"] = execError
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "ForkExecuteCommand").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmHalt(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 0 {
		e.Logger.WithField("function", "Halt").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 0 {
		e.Logger.WithField("function", "Halt").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}
	value := e.Halt()
	rawVMRet := VMResponse{}
	rawVMRet["value"] = value
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "Halt").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmMD5(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "MD5").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "MD5").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	data := e.ValueToByteSlice(call.Argument(0))
	value := e.MD5(data)
	rawVMRet := VMResponse{}
	rawVMRet["value"] = value
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "MD5").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmObfuscateString(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "ObfuscateString").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "ObfuscateString").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var str string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "ObfuscateString").WithField("trace", "true").Errorf("Could not export field: %s", "str")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		str = rawArg0.(string)
	default:
		e.Logger.WithField("function", "ObfuscateString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	value := e.ObfuscateString(str)
	rawVMRet := VMResponse{}
	rawVMRet["value"] = value
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "ObfuscateString").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmRandomInt(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 2 {
		e.Logger.WithField("function", "RandomInt").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 2 {
		e.Logger.WithField("function", "RandomInt").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var min int64
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "RandomInt").WithField("trace", "true").Errorf("Could not export field: %s", "min")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case int64:
		min = rawArg0.(int64)
	default:
		e.Logger.WithField("function", "RandomInt").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "int64", v)
		return otto.FalseValue()
	}

	var max int64
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "RandomInt").WithField("trace", "true").Errorf("Could not export field: %s", "max")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case int64:
		max = rawArg1.(int64)
	default:
		e.Logger.WithField("function", "RandomInt").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "int64", v)
		return otto.FalseValue()
	}
	value := e.RandomInt(min, max)
	rawVMRet := VMResponse{}
	rawVMRet["value"] = value
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "RandomInt").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmRandomMixedCaseString(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "RandomMixedCaseString").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "RandomMixedCaseString").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var strlen int64
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "RandomMixedCaseString").WithField("trace", "true").Errorf("Could not export field: %s", "strlen")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case int64:
		strlen = rawArg0.(int64)
	default:
		e.Logger.WithField("function", "RandomMixedCaseString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "int64", v)
		return otto.FalseValue()
	}
	value := e.RandomMixedCaseString(strlen)
	rawVMRet := VMResponse{}
	rawVMRet["value"] = value
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "RandomMixedCaseString").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmRandomString(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "RandomString").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "RandomString").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var strlen int64
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "RandomString").WithField("trace", "true").Errorf("Could not export field: %s", "strlen")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case int64:
		strlen = rawArg0.(int64)
	default:
		e.Logger.WithField("function", "RandomString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "int64", v)
		return otto.FalseValue()
	}
	value := e.RandomString(strlen)
	rawVMRet := VMResponse{}
	rawVMRet["value"] = value
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "RandomString").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmStripSpaces(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "StripSpaces").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "StripSpaces").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var str string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "StripSpaces").WithField("trace", "true").Errorf("Could not export field: %s", "str")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		str = rawArg0.(string)
	default:
		e.Logger.WithField("function", "StripSpaces").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	value := e.StripSpaces(str)
	rawVMRet := VMResponse{}
	rawVMRet["value"] = value
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "StripSpaces").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmTimestamp(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 0 {
		e.Logger.WithField("function", "Timestamp").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 0 {
		e.Logger.WithField("function", "Timestamp").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}
	value := e.Timestamp()
	rawVMRet := VMResponse{}
	rawVMRet["value"] = value
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "Timestamp").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmWriteFile(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 3 {
		e.Logger.WithField("function", "WriteFile").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 3 {
		e.Logger.WithField("function", "WriteFile").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var path string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "WriteFile").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		path = rawArg0.(string)
	default:
		e.Logger.WithField("function", "WriteFile").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	fileData := e.ValueToByteSlice(call.Argument(1))

	var perms int64
	rawArg2, err := call.Argument(2).Export()
	if err != nil {
		e.Logger.WithField("function", "WriteFile").WithField("trace", "true").Errorf("Could not export field: %s", "perms")
		return otto.FalseValue()
	}
	switch v := rawArg2.(type) {
	case int64:
		perms = rawArg2.(int64)
	default:
		e.Logger.WithField("function", "WriteFile").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "int64", v)
		return otto.FalseValue()
	}
	bytesWritten, fileError := e.WriteFile(path, fileData, perms)
	rawVMRet := VMResponse{}
	rawVMRet["bytesWritten"] = bytesWritten
	rawVMRet["fileError"] = fileError
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "WriteFile").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmXorBytes(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 2 {
		e.Logger.WithField("function", "XorBytes").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 2 {
		e.Logger.WithField("function", "XorBytes").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	aByteArray := e.ValueToByteSlice(call.Argument(0))

	bByteArray := e.ValueToByteSlice(call.Argument(1))
	value := e.XorBytes(aByteArray, bByteArray)
	rawVMRet := VMResponse{}
	rawVMRet["value"] = value
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "XorBytes").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}
