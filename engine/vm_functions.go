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
//  AppendFileBytes(path, fileData) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.AppendFileBytes
//  AppendFileString(path, addString) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.AppendFileString
//  CreateDir(path) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.CreateDir
//  DeleteFile(path) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.DeleteFile
//  FileExists(path) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.FileExists
//  ReadFile(path) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.ReadFile
//  WriteFile(path, fileData, perms) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.WriteFile
//
// Library os
//
// Functions in os:
//  Chmod(path, perms) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.Chmod
//  EnvVars() - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.EnvVars
//  FindProcByName(procName) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.FindProcByName
//  GetEnvVar(vars) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.GetEnvVar
//  GetProcName(pid) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.GetProcName
//  InstallSystemService(path, name, displayName, description) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.InstallSystemService
//  ModTime(path) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.ModTime
//  ModifyTimestamp(path, accessTime, modifyTime) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.ModifyTimestamp
//  RemoveServiceByName(name) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.RemoveServiceByName
//  RunningProcs() - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.RunningProcs
//  SelfPath() - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.SelfPath
//  Signal(signal, pid) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.Signal
//  StartServiceByName(name) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.StartServiceByName
//  StopServiceByName(name) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.StopServiceByName
//
// Library registry
//
// Functions in registry:
//  AddRegKeyBinary(registryString, path, name, value) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.AddRegKeyBinary
//  AddRegKeyDWORD(registryString, path, name, value) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.AddRegKeyDWORD
//  AddRegKeyExpandedString(registryString, path, name, value) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.AddRegKeyExpandedString
//  AddRegKeyQWORD(registryString, path, name, value) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.AddRegKeyQWORD
//  AddRegKeyString(registryString, path, name, value) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.AddRegKeyString
//  AddRegKeyStrings(registryString, path, name, value) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.AddRegKeyStrings
//  DelRegKey(registryString, path) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.DelRegKey
//  DelRegKeyValue(registryString, path, value) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.DelRegKeyValue
//  QueryRegKey(registryString, path) - https://godoc.org/github.com/gen0cide/gscript/engine/#Engine.QueryRegKey
//
package engine

import (
	"path/filepath"

	"github.com/robertkrimen/otto"
)

func (e *Engine) CreateVM() {
	e.VM = otto.New()
	e.injectVars()
	e.VM.Set("AddRegKeyBinary", e.vmAddRegKeyBinary)
	e.VM.Set("AddRegKeyDWORD", e.vmAddRegKeyDWORD)
	e.VM.Set("AddRegKeyExpandedString", e.vmAddRegKeyExpandedString)
	e.VM.Set("AddRegKeyQWORD", e.vmAddRegKeyQWORD)
	e.VM.Set("AddRegKeyString", e.vmAddRegKeyString)
	e.VM.Set("AddRegKeyStrings", e.vmAddRegKeyStrings)
	e.VM.Set("AppendFileBytes", e.vmAppendFileBytes)
	e.VM.Set("AppendFileString", e.vmAppendFileString)
	e.VM.Set("Asset", e.vmAsset)
	e.VM.Set("Chmod", e.vmChmod)
	e.VM.Set("CreateDir", e.vmCreateDir)
	e.VM.Set("DelRegKey", e.vmDelRegKey)
	e.VM.Set("DelRegKeyValue", e.vmDelRegKeyValue)
	e.VM.Set("DeleteFile", e.vmDeleteFile)
	e.VM.Set("DeobfuscateString", e.vmDeobfuscateString)
	e.VM.Set("EnvVars", e.vmEnvVars)
	e.VM.Set("ExecuteCommand", e.vmExecuteCommand)
	e.VM.Set("FileExists", e.vmFileExists)
	e.VM.Set("FindProcByName", e.vmFindProcByName)
	e.VM.Set("ForkExecuteCommand", e.vmForkExecuteCommand)
	e.VM.Set("GetEnvVar", e.vmGetEnvVar)
	e.VM.Set("GetProcName", e.vmGetProcName)
	e.VM.Set("Halt", e.vmHalt)
	e.VM.Set("InstallSystemService", e.vmInstallSystemService)
	e.VM.Set("MD5", e.vmMD5)
	e.VM.Set("ModTime", e.vmModTime)
	e.VM.Set("ModifyTimestamp", e.vmModifyTimestamp)
	e.VM.Set("ObfuscateString", e.vmObfuscateString)
	e.VM.Set("QueryRegKey", e.vmQueryRegKey)
	e.VM.Set("RandomInt", e.vmRandomInt)
	e.VM.Set("RandomMixedCaseString", e.vmRandomMixedCaseString)
	e.VM.Set("RandomString", e.vmRandomString)
	e.VM.Set("ReadFile", e.vmReadFile)
	e.VM.Set("RemoveServiceByName", e.vmRemoveServiceByName)
	e.VM.Set("RunningProcs", e.vmRunningProcs)
	e.VM.Set("SelfPath", e.vmSelfPath)
	e.VM.Set("Signal", e.vmSignal)
	e.VM.Set("StartServiceByName", e.vmStartServiceByName)
	e.VM.Set("StopServiceByName", e.vmStopServiceByName)
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

func (e *Engine) vmAddRegKeyBinary(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 4 {
		e.Logger.WithField("function", "AddRegKeyBinary").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 4 {
		e.Logger.WithField("function", "AddRegKeyBinary").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var registryString string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyBinary").WithField("trace", "true").Errorf("Could not export field: %s", "registryString")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		registryString = rawArg0.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyBinary").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var path string
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyBinary").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case string:
		path = filepath.Clean(rawArg1.(string))
	default:
		e.Logger.WithField("function", "AddRegKeyBinary").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var name string
	rawArg2, err := call.Argument(2).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyBinary").WithField("trace", "true").Errorf("Could not export field: %s", "name")
		return otto.FalseValue()
	}
	switch v := rawArg2.(type) {
	case string:
		name = rawArg2.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyBinary").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	value := e.ValueToByteSlice(call.Argument(3))
	runtimeError := e.AddRegKeyBinary(registryString, path, name, value)
	rawVMRet := VMResponse{}

	if runtimeError != nil {
		e.Logger.WithField("function", "AddRegKeyBinary").WithField("trace", "true").Errorf("<function error> %s", runtimeError.Error())
		rawVMRet["runtimeError"] = runtimeError.Error()
	} else {
		rawVMRet["runtimeError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "AddRegKeyBinary").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmAddRegKeyDWORD(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 4 {
		e.Logger.WithField("function", "AddRegKeyDWORD").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 4 {
		e.Logger.WithField("function", "AddRegKeyDWORD").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var registryString string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyDWORD").WithField("trace", "true").Errorf("Could not export field: %s", "registryString")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		registryString = rawArg0.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyDWORD").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var path string
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyDWORD").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case string:
		path = filepath.Clean(rawArg1.(string))
	default:
		e.Logger.WithField("function", "AddRegKeyDWORD").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var name string
	rawArg2, err := call.Argument(2).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyDWORD").WithField("trace", "true").Errorf("Could not export field: %s", "name")
		return otto.FalseValue()
	}
	switch v := rawArg2.(type) {
	case string:
		name = rawArg2.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyDWORD").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var value uint32
	rawArg3, err := call.Argument(3).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyDWORD").WithField("trace", "true").Errorf("Could not export field: %s", "value")
		return otto.FalseValue()
	}
	switch v := rawArg3.(type) {
	case uint32:
		value = rawArg3.(uint32)
	default:
		e.Logger.WithField("function", "AddRegKeyDWORD").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "uint32", v)
		return otto.FalseValue()
	}
	runtimeError := e.AddRegKeyDWORD(registryString, path, name, value)
	rawVMRet := VMResponse{}

	if runtimeError != nil {
		e.Logger.WithField("function", "AddRegKeyDWORD").WithField("trace", "true").Errorf("<function error> %s", runtimeError.Error())
		rawVMRet["runtimeError"] = runtimeError.Error()
	} else {
		rawVMRet["runtimeError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "AddRegKeyDWORD").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmAddRegKeyExpandedString(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 4 {
		e.Logger.WithField("function", "AddRegKeyExpandedString").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 4 {
		e.Logger.WithField("function", "AddRegKeyExpandedString").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var registryString string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyExpandedString").WithField("trace", "true").Errorf("Could not export field: %s", "registryString")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		registryString = rawArg0.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyExpandedString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var path string
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyExpandedString").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case string:
		path = filepath.Clean(rawArg1.(string))
	default:
		e.Logger.WithField("function", "AddRegKeyExpandedString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var name string
	rawArg2, err := call.Argument(2).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyExpandedString").WithField("trace", "true").Errorf("Could not export field: %s", "name")
		return otto.FalseValue()
	}
	switch v := rawArg2.(type) {
	case string:
		name = rawArg2.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyExpandedString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var value string
	rawArg3, err := call.Argument(3).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyExpandedString").WithField("trace", "true").Errorf("Could not export field: %s", "value")
		return otto.FalseValue()
	}
	switch v := rawArg3.(type) {
	case string:
		value = rawArg3.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyExpandedString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	runtimeError := e.AddRegKeyExpandedString(registryString, path, name, value)
	rawVMRet := VMResponse{}

	if runtimeError != nil {
		e.Logger.WithField("function", "AddRegKeyExpandedString").WithField("trace", "true").Errorf("<function error> %s", runtimeError.Error())
		rawVMRet["runtimeError"] = runtimeError.Error()
	} else {
		rawVMRet["runtimeError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "AddRegKeyExpandedString").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmAddRegKeyQWORD(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 4 {
		e.Logger.WithField("function", "AddRegKeyQWORD").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 4 {
		e.Logger.WithField("function", "AddRegKeyQWORD").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var registryString string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyQWORD").WithField("trace", "true").Errorf("Could not export field: %s", "registryString")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		registryString = rawArg0.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyQWORD").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var path string
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyQWORD").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case string:
		path = filepath.Clean(rawArg1.(string))
	default:
		e.Logger.WithField("function", "AddRegKeyQWORD").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var name string
	rawArg2, err := call.Argument(2).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyQWORD").WithField("trace", "true").Errorf("Could not export field: %s", "name")
		return otto.FalseValue()
	}
	switch v := rawArg2.(type) {
	case string:
		name = rawArg2.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyQWORD").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var value uint64
	rawArg3, err := call.Argument(3).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyQWORD").WithField("trace", "true").Errorf("Could not export field: %s", "value")
		return otto.FalseValue()
	}
	switch v := rawArg3.(type) {
	case uint64:
		value = rawArg3.(uint64)
	default:
		e.Logger.WithField("function", "AddRegKeyQWORD").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "uint64", v)
		return otto.FalseValue()
	}
	runtimeError := e.AddRegKeyQWORD(registryString, path, name, value)
	rawVMRet := VMResponse{}

	if runtimeError != nil {
		e.Logger.WithField("function", "AddRegKeyQWORD").WithField("trace", "true").Errorf("<function error> %s", runtimeError.Error())
		rawVMRet["runtimeError"] = runtimeError.Error()
	} else {
		rawVMRet["runtimeError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "AddRegKeyQWORD").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmAddRegKeyString(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 4 {
		e.Logger.WithField("function", "AddRegKeyString").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 4 {
		e.Logger.WithField("function", "AddRegKeyString").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var registryString string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyString").WithField("trace", "true").Errorf("Could not export field: %s", "registryString")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		registryString = rawArg0.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var path string
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyString").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case string:
		path = filepath.Clean(rawArg1.(string))
	default:
		e.Logger.WithField("function", "AddRegKeyString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var name string
	rawArg2, err := call.Argument(2).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyString").WithField("trace", "true").Errorf("Could not export field: %s", "name")
		return otto.FalseValue()
	}
	switch v := rawArg2.(type) {
	case string:
		name = rawArg2.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var value string
	rawArg3, err := call.Argument(3).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyString").WithField("trace", "true").Errorf("Could not export field: %s", "value")
		return otto.FalseValue()
	}
	switch v := rawArg3.(type) {
	case string:
		value = rawArg3.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	runtimeError := e.AddRegKeyString(registryString, path, name, value)
	rawVMRet := VMResponse{}

	if runtimeError != nil {
		e.Logger.WithField("function", "AddRegKeyString").WithField("trace", "true").Errorf("<function error> %s", runtimeError.Error())
		rawVMRet["runtimeError"] = runtimeError.Error()
	} else {
		rawVMRet["runtimeError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "AddRegKeyString").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmAddRegKeyStrings(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 4 {
		e.Logger.WithField("function", "AddRegKeyStrings").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 4 {
		e.Logger.WithField("function", "AddRegKeyStrings").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var registryString string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyStrings").WithField("trace", "true").Errorf("Could not export field: %s", "registryString")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		registryString = rawArg0.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyStrings").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var path string
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyStrings").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case string:
		path = filepath.Clean(rawArg1.(string))
	default:
		e.Logger.WithField("function", "AddRegKeyStrings").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var name string
	rawArg2, err := call.Argument(2).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyStrings").WithField("trace", "true").Errorf("Could not export field: %s", "name")
		return otto.FalseValue()
	}
	switch v := rawArg2.(type) {
	case string:
		name = rawArg2.(string)
	default:
		e.Logger.WithField("function", "AddRegKeyStrings").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var value []string
	rawArg3, err := call.Argument(3).Export()
	if err != nil {
		e.Logger.WithField("function", "AddRegKeyStrings").WithField("trace", "true").Errorf("Could not export field: %s", "value")
		return otto.FalseValue()
	}
	switch v := rawArg3.(type) {
	case []string:
		value = rawArg3.([]string)
	default:
		e.Logger.WithField("function", "AddRegKeyStrings").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "[]string", v)
		return otto.FalseValue()
	}
	runtimeError := e.AddRegKeyStrings(registryString, path, name, value)
	rawVMRet := VMResponse{}

	if runtimeError != nil {
		e.Logger.WithField("function", "AddRegKeyStrings").WithField("trace", "true").Errorf("<function error> %s", runtimeError.Error())
		rawVMRet["runtimeError"] = runtimeError.Error()
	} else {
		rawVMRet["runtimeError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "AddRegKeyStrings").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmAppendFileBytes(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 2 {
		e.Logger.WithField("function", "AppendFileBytes").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 2 {
		e.Logger.WithField("function", "AppendFileBytes").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var path string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "AppendFileBytes").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		path = filepath.Clean(rawArg0.(string))
	default:
		e.Logger.WithField("function", "AppendFileBytes").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	fileData := e.ValueToByteSlice(call.Argument(1))
	fileError := e.AppendFileBytes(path, fileData)
	rawVMRet := VMResponse{}

	if fileError != nil {
		e.Logger.WithField("function", "AppendFileBytes").WithField("trace", "true").Errorf("<function error> %s", fileError.Error())
		rawVMRet["fileError"] = fileError.Error()
	} else {
		rawVMRet["fileError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "AppendFileBytes").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmAppendFileString(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 2 {
		e.Logger.WithField("function", "AppendFileString").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 2 {
		e.Logger.WithField("function", "AppendFileString").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var path string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "AppendFileString").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		path = filepath.Clean(rawArg0.(string))
	default:
		e.Logger.WithField("function", "AppendFileString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var addString string
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "AppendFileString").WithField("trace", "true").Errorf("Could not export field: %s", "addString")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case string:
		addString = rawArg1.(string)
	default:
		e.Logger.WithField("function", "AppendFileString").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	fileError := e.AppendFileString(path, addString)
	rawVMRet := VMResponse{}

	if fileError != nil {
		e.Logger.WithField("function", "AppendFileString").WithField("trace", "true").Errorf("<function error> %s", fileError.Error())
		rawVMRet["fileError"] = fileError.Error()
	} else {
		rawVMRet["fileError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "AppendFileString").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
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

	if err != nil {
		e.Logger.WithField("function", "Asset").WithField("trace", "true").Errorf("<function error> %s", err.Error())
		rawVMRet["err"] = err.Error()
	} else {
		rawVMRet["err"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "Asset").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmChmod(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 2 {
		e.Logger.WithField("function", "Chmod").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 2 {
		e.Logger.WithField("function", "Chmod").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var path string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "Chmod").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		path = filepath.Clean(rawArg0.(string))
	default:
		e.Logger.WithField("function", "Chmod").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var perms int64
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "Chmod").WithField("trace", "true").Errorf("Could not export field: %s", "perms")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case int64:
		perms = rawArg1.(int64)
	default:
		e.Logger.WithField("function", "Chmod").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "int64", v)
		return otto.FalseValue()
	}
	osError := e.Chmod(path, perms)
	rawVMRet := VMResponse{}

	if osError != nil {
		e.Logger.WithField("function", "Chmod").WithField("trace", "true").Errorf("<function error> %s", osError.Error())
		rawVMRet["osError"] = osError.Error()
	} else {
		rawVMRet["osError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "Chmod").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmCreateDir(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "CreateDir").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "CreateDir").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var path string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "CreateDir").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		path = filepath.Clean(rawArg0.(string))
	default:
		e.Logger.WithField("function", "CreateDir").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	fileError := e.CreateDir(path)
	rawVMRet := VMResponse{}

	if fileError != nil {
		e.Logger.WithField("function", "CreateDir").WithField("trace", "true").Errorf("<function error> %s", fileError.Error())
		rawVMRet["fileError"] = fileError.Error()
	} else {
		rawVMRet["fileError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "CreateDir").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmDelRegKey(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 2 {
		e.Logger.WithField("function", "DelRegKey").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 2 {
		e.Logger.WithField("function", "DelRegKey").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var registryString string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "DelRegKey").WithField("trace", "true").Errorf("Could not export field: %s", "registryString")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		registryString = rawArg0.(string)
	default:
		e.Logger.WithField("function", "DelRegKey").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var path string
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "DelRegKey").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case string:
		path = filepath.Clean(rawArg1.(string))
	default:
		e.Logger.WithField("function", "DelRegKey").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	runtimeError := e.DelRegKey(registryString, path)
	rawVMRet := VMResponse{}

	if runtimeError != nil {
		e.Logger.WithField("function", "DelRegKey").WithField("trace", "true").Errorf("<function error> %s", runtimeError.Error())
		rawVMRet["runtimeError"] = runtimeError.Error()
	} else {
		rawVMRet["runtimeError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "DelRegKey").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmDelRegKeyValue(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 3 {
		e.Logger.WithField("function", "DelRegKeyValue").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 3 {
		e.Logger.WithField("function", "DelRegKeyValue").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var registryString string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "DelRegKeyValue").WithField("trace", "true").Errorf("Could not export field: %s", "registryString")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		registryString = rawArg0.(string)
	default:
		e.Logger.WithField("function", "DelRegKeyValue").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var path string
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "DelRegKeyValue").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case string:
		path = filepath.Clean(rawArg1.(string))
	default:
		e.Logger.WithField("function", "DelRegKeyValue").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var value string
	rawArg2, err := call.Argument(2).Export()
	if err != nil {
		e.Logger.WithField("function", "DelRegKeyValue").WithField("trace", "true").Errorf("Could not export field: %s", "value")
		return otto.FalseValue()
	}
	switch v := rawArg2.(type) {
	case string:
		value = rawArg2.(string)
	default:
		e.Logger.WithField("function", "DelRegKeyValue").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	runtimeError := e.DelRegKeyValue(registryString, path, value)
	rawVMRet := VMResponse{}

	if runtimeError != nil {
		e.Logger.WithField("function", "DelRegKeyValue").WithField("trace", "true").Errorf("<function error> %s", runtimeError.Error())
		rawVMRet["runtimeError"] = runtimeError.Error()
	} else {
		rawVMRet["runtimeError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "DelRegKeyValue").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmDeleteFile(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "DeleteFile").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "DeleteFile").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var path string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "DeleteFile").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		path = filepath.Clean(rawArg0.(string))
	default:
		e.Logger.WithField("function", "DeleteFile").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	fileError := e.DeleteFile(path)
	rawVMRet := VMResponse{}

	if fileError != nil {
		e.Logger.WithField("function", "DeleteFile").WithField("trace", "true").Errorf("<function error> %s", fileError.Error())
		rawVMRet["fileError"] = fileError.Error()
	} else {
		rawVMRet["fileError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "DeleteFile").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
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

func (e *Engine) vmEnvVars(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 0 {
		e.Logger.WithField("function", "EnvVars").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 0 {
		e.Logger.WithField("function", "EnvVars").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}
	vars := e.EnvVars()
	rawVMRet := VMResponse{}

	rawVMRet["vars"] = vars
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "EnvVars").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
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

func (e *Engine) vmFileExists(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "FileExists").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "FileExists").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var path string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "FileExists").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		path = filepath.Clean(rawArg0.(string))
	default:
		e.Logger.WithField("function", "FileExists").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	fileExists, fileError := e.FileExists(path)
	rawVMRet := VMResponse{}

	rawVMRet["fileExists"] = fileExists

	if fileError != nil {
		e.Logger.WithField("function", "FileExists").WithField("trace", "true").Errorf("<function error> %s", fileError.Error())
		rawVMRet["fileError"] = fileError.Error()
	} else {
		rawVMRet["fileError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "FileExists").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmFindProcByName(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "FindProcByName").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "FindProcByName").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var procName string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "FindProcByName").WithField("trace", "true").Errorf("Could not export field: %s", "procName")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		procName = rawArg0.(string)
	default:
		e.Logger.WithField("function", "FindProcByName").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	pid, procError := e.FindProcByName(procName)
	rawVMRet := VMResponse{}

	rawVMRet["pid"] = pid

	if procError != nil {
		e.Logger.WithField("function", "FindProcByName").WithField("trace", "true").Errorf("<function error> %s", procError.Error())
		rawVMRet["procError"] = procError.Error()
	} else {
		rawVMRet["procError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "FindProcByName").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
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

	if execError != nil {
		e.Logger.WithField("function", "ForkExecuteCommand").WithField("trace", "true").Errorf("<function error> %s", execError.Error())
		rawVMRet["execError"] = execError.Error()
	} else {
		rawVMRet["execError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "ForkExecuteCommand").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmGetEnvVar(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "GetEnvVar").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "GetEnvVar").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var vars string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "GetEnvVar").WithField("trace", "true").Errorf("Could not export field: %s", "vars")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		vars = rawArg0.(string)
	default:
		e.Logger.WithField("function", "GetEnvVar").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	value := e.GetEnvVar(vars)
	rawVMRet := VMResponse{}

	rawVMRet["value"] = value
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "GetEnvVar").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmGetProcName(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "GetProcName").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "GetProcName").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var pid int
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "GetProcName").WithField("trace", "true").Errorf("Could not export field: %s", "pid")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case int:
		pid = rawArg0.(int)
	default:
		e.Logger.WithField("function", "GetProcName").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "int", v)
		return otto.FalseValue()
	}
	procName, runtimeError := e.GetProcName(pid)
	rawVMRet := VMResponse{}

	rawVMRet["procName"] = procName

	if runtimeError != nil {
		e.Logger.WithField("function", "GetProcName").WithField("trace", "true").Errorf("<function error> %s", runtimeError.Error())
		rawVMRet["runtimeError"] = runtimeError.Error()
	} else {
		rawVMRet["runtimeError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "GetProcName").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
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

func (e *Engine) vmInstallSystemService(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 4 {
		e.Logger.WithField("function", "InstallSystemService").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 4 {
		e.Logger.WithField("function", "InstallSystemService").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var path string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "InstallSystemService").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		path = filepath.Clean(rawArg0.(string))
	default:
		e.Logger.WithField("function", "InstallSystemService").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var name string
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "InstallSystemService").WithField("trace", "true").Errorf("Could not export field: %s", "name")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case string:
		name = rawArg1.(string)
	default:
		e.Logger.WithField("function", "InstallSystemService").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var displayName string
	rawArg2, err := call.Argument(2).Export()
	if err != nil {
		e.Logger.WithField("function", "InstallSystemService").WithField("trace", "true").Errorf("Could not export field: %s", "displayName")
		return otto.FalseValue()
	}
	switch v := rawArg2.(type) {
	case string:
		displayName = rawArg2.(string)
	default:
		e.Logger.WithField("function", "InstallSystemService").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var description string
	rawArg3, err := call.Argument(3).Export()
	if err != nil {
		e.Logger.WithField("function", "InstallSystemService").WithField("trace", "true").Errorf("Could not export field: %s", "description")
		return otto.FalseValue()
	}
	switch v := rawArg3.(type) {
	case string:
		description = rawArg3.(string)
	default:
		e.Logger.WithField("function", "InstallSystemService").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	installError := e.InstallSystemService(path, name, displayName, description)
	rawVMRet := VMResponse{}

	if installError != nil {
		e.Logger.WithField("function", "InstallSystemService").WithField("trace", "true").Errorf("<function error> %s", installError.Error())
		rawVMRet["installError"] = installError.Error()
	} else {
		rawVMRet["installError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "InstallSystemService").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
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

func (e *Engine) vmModTime(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "ModTime").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "ModTime").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var path string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "ModTime").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		path = filepath.Clean(rawArg0.(string))
	default:
		e.Logger.WithField("function", "ModTime").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	modTime, fileError := e.ModTime(path)
	rawVMRet := VMResponse{}

	rawVMRet["modTime"] = modTime

	if fileError != nil {
		e.Logger.WithField("function", "ModTime").WithField("trace", "true").Errorf("<function error> %s", fileError.Error())
		rawVMRet["fileError"] = fileError.Error()
	} else {
		rawVMRet["fileError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "ModTime").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmModifyTimestamp(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 3 {
		e.Logger.WithField("function", "ModifyTimestamp").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 3 {
		e.Logger.WithField("function", "ModifyTimestamp").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var path string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "ModifyTimestamp").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		path = filepath.Clean(rawArg0.(string))
	default:
		e.Logger.WithField("function", "ModifyTimestamp").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var accessTime int64
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "ModifyTimestamp").WithField("trace", "true").Errorf("Could not export field: %s", "accessTime")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case int64:
		accessTime = rawArg1.(int64)
	default:
		e.Logger.WithField("function", "ModifyTimestamp").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "int64", v)
		return otto.FalseValue()
	}

	var modifyTime int64
	rawArg2, err := call.Argument(2).Export()
	if err != nil {
		e.Logger.WithField("function", "ModifyTimestamp").WithField("trace", "true").Errorf("Could not export field: %s", "modifyTime")
		return otto.FalseValue()
	}
	switch v := rawArg2.(type) {
	case int64:
		modifyTime = rawArg2.(int64)
	default:
		e.Logger.WithField("function", "ModifyTimestamp").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "int64", v)
		return otto.FalseValue()
	}
	fileError := e.ModifyTimestamp(path, accessTime, modifyTime)
	rawVMRet := VMResponse{}

	if fileError != nil {
		e.Logger.WithField("function", "ModifyTimestamp").WithField("trace", "true").Errorf("<function error> %s", fileError.Error())
		rawVMRet["fileError"] = fileError.Error()
	} else {
		rawVMRet["fileError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "ModifyTimestamp").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
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

func (e *Engine) vmQueryRegKey(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 2 {
		e.Logger.WithField("function", "QueryRegKey").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 2 {
		e.Logger.WithField("function", "QueryRegKey").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var registryString string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "QueryRegKey").WithField("trace", "true").Errorf("Could not export field: %s", "registryString")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		registryString = rawArg0.(string)
	default:
		e.Logger.WithField("function", "QueryRegKey").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}

	var path string
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "QueryRegKey").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case string:
		path = filepath.Clean(rawArg1.(string))
	default:
		e.Logger.WithField("function", "QueryRegKey").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	keyObj, runtimeError := e.QueryRegKey(registryString, path)
	rawVMRet := VMResponse{}

	rawVMRet["keyObj"] = keyObj

	if runtimeError != nil {
		e.Logger.WithField("function", "QueryRegKey").WithField("trace", "true").Errorf("<function error> %s", runtimeError.Error())
		rawVMRet["runtimeError"] = runtimeError.Error()
	} else {
		rawVMRet["runtimeError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "QueryRegKey").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
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

func (e *Engine) vmReadFile(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "ReadFile").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "ReadFile").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var path string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "ReadFile").WithField("trace", "true").Errorf("Could not export field: %s", "path")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		path = filepath.Clean(rawArg0.(string))
	default:
		e.Logger.WithField("function", "ReadFile").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	fileBytes, fileError := e.ReadFile(path)
	rawVMRet := VMResponse{}

	rawVMRet["fileBytes"] = fileBytes

	if fileError != nil {
		e.Logger.WithField("function", "ReadFile").WithField("trace", "true").Errorf("<function error> %s", fileError.Error())
		rawVMRet["fileError"] = fileError.Error()
	} else {
		rawVMRet["fileError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "ReadFile").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmRemoveServiceByName(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "RemoveServiceByName").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "RemoveServiceByName").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var name string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "RemoveServiceByName").WithField("trace", "true").Errorf("Could not export field: %s", "name")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		name = rawArg0.(string)
	default:
		e.Logger.WithField("function", "RemoveServiceByName").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	removalError := e.RemoveServiceByName(name)
	rawVMRet := VMResponse{}

	if removalError != nil {
		e.Logger.WithField("function", "RemoveServiceByName").WithField("trace", "true").Errorf("<function error> %s", removalError.Error())
		rawVMRet["removalError"] = removalError.Error()
	} else {
		rawVMRet["removalError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "RemoveServiceByName").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmRunningProcs(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 0 {
		e.Logger.WithField("function", "RunningProcs").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 0 {
		e.Logger.WithField("function", "RunningProcs").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}
	pids, runtimeError := e.RunningProcs()
	rawVMRet := VMResponse{}

	rawVMRet["pids"] = pids

	if runtimeError != nil {
		e.Logger.WithField("function", "RunningProcs").WithField("trace", "true").Errorf("<function error> %s", runtimeError.Error())
		rawVMRet["runtimeError"] = runtimeError.Error()
	} else {
		rawVMRet["runtimeError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "RunningProcs").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmSelfPath(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 0 {
		e.Logger.WithField("function", "SelfPath").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 0 {
		e.Logger.WithField("function", "SelfPath").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}
	path, osError := e.SelfPath()
	rawVMRet := VMResponse{}

	rawVMRet["path"] = path

	if osError != nil {
		e.Logger.WithField("function", "SelfPath").WithField("trace", "true").Errorf("<function error> %s", osError.Error())
		rawVMRet["osError"] = osError.Error()
	} else {
		rawVMRet["osError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "SelfPath").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmSignal(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 2 {
		e.Logger.WithField("function", "Signal").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 2 {
		e.Logger.WithField("function", "Signal").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var signal int
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "Signal").WithField("trace", "true").Errorf("Could not export field: %s", "signal")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case int:
		signal = rawArg0.(int)
	default:
		e.Logger.WithField("function", "Signal").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "int", v)
		return otto.FalseValue()
	}

	var pid int
	rawArg1, err := call.Argument(1).Export()
	if err != nil {
		e.Logger.WithField("function", "Signal").WithField("trace", "true").Errorf("Could not export field: %s", "pid")
		return otto.FalseValue()
	}
	switch v := rawArg1.(type) {
	case int:
		pid = rawArg1.(int)
	default:
		e.Logger.WithField("function", "Signal").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "int", v)
		return otto.FalseValue()
	}
	runtimeError := e.Signal(signal, pid)
	rawVMRet := VMResponse{}

	if runtimeError != nil {
		e.Logger.WithField("function", "Signal").WithField("trace", "true").Errorf("<function error> %s", runtimeError.Error())
		rawVMRet["runtimeError"] = runtimeError.Error()
	} else {
		rawVMRet["runtimeError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "Signal").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmStartServiceByName(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "StartServiceByName").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "StartServiceByName").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var name string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "StartServiceByName").WithField("trace", "true").Errorf("Could not export field: %s", "name")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		name = rawArg0.(string)
	default:
		e.Logger.WithField("function", "StartServiceByName").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	startError := e.StartServiceByName(name)
	rawVMRet := VMResponse{}

	if startError != nil {
		e.Logger.WithField("function", "StartServiceByName").WithField("trace", "true").Errorf("<function error> %s", startError.Error())
		rawVMRet["startError"] = startError.Error()
	} else {
		rawVMRet["startError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "StartServiceByName").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
		return otto.FalseValue()
	}
	return vmRet
}

func (e *Engine) vmStopServiceByName(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 1 {
		e.Logger.WithField("function", "StopServiceByName").WithField("trace", "true").Error("Too many arguments in call.")
		return otto.FalseValue()
	}
	if len(call.ArgumentList) < 1 {
		e.Logger.WithField("function", "StopServiceByName").WithField("trace", "true").Error("Too few arguments in call.")
		return otto.FalseValue()
	}

	var name string
	rawArg0, err := call.Argument(0).Export()
	if err != nil {
		e.Logger.WithField("function", "StopServiceByName").WithField("trace", "true").Errorf("Could not export field: %s", "name")
		return otto.FalseValue()
	}
	switch v := rawArg0.(type) {
	case string:
		name = rawArg0.(string)
	default:
		e.Logger.WithField("function", "StopServiceByName").WithField("trace", "true").Errorf("Argument type mismatch: expected %s, got %T", "string", v)
		return otto.FalseValue()
	}
	installError := e.StopServiceByName(name)
	rawVMRet := VMResponse{}

	if installError != nil {
		e.Logger.WithField("function", "StopServiceByName").WithField("trace", "true").Errorf("<function error> %s", installError.Error())
		rawVMRet["installError"] = installError.Error()
	} else {
		rawVMRet["installError"] = nil
	}
	vmRet, vmRetError := e.VM.ToValue(rawVMRet)
	if vmRetError != nil {
		e.Logger.WithField("function", "StopServiceByName").WithField("trace", "true").Errorf("Return conversion failed: %s", vmRetError.Error())
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
		path = filepath.Clean(rawArg0.(string))
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

	if fileError != nil {
		e.Logger.WithField("function", "WriteFile").WithField("trace", "true").Errorf("<function error> %s", fileError.Error())
		rawVMRet["fileError"] = fileError.Error()
	} else {
		rawVMRet["fileError"] = nil
	}
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
