// +build !windows

package engine

import "errors"

// AddRegKeyString - Add a string registry key
//
// Package
//
// registry
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  AddRegKeyString(registryString, path, name, value)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * registryString (string)
//  * path (string)
//  * name (string)
//  * value (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = AddRegKeyString(registryString, path, name, value);
//  // obj.runtimeError
//
func (e *Engine) AddRegKeyString(registryString string, path string, name string, value string) error {
	return errors.New("this function is unimplemented on non windows platforms")
}

// AddRegKeyExpandedString - Add an expanded string registry key
//
// Package
//
// registry
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  AddRegKeyExpandedString(registryString, path, name, value)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * registryString (string)
//  * path (string)
//  * name (string)
//  * value (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = AddRegKeyExpandedString(registryString, path, name, value);
//  // obj.runtimeError
//
func (e *Engine) AddRegKeyExpandedString(registryString string, path string, name string, value string) error {
	return errors.New("this function is unimplemented on non windows platforms")
}

// AddRegKeyBinary - Add a binary registry key
//
// Package
//
// registry
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  AddRegKeyBinary(registryString, path, name, value)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * registryString (string)
//  * path (string)
//  * name (string)
//  * value ([]byte)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = AddRegKeyBinary(registryString, path, name, value);
//  // obj.runtimeError
//
func (e *Engine) AddRegKeyBinary(registryString string, path string, name string, value []byte) error {
	return errors.New("this function is unimplemented on non windows platforms")
}

// AddRegKeyDWORD - Add a DWORD registry key
//
// Package
//
// registry
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  AddRegKeyDWORD(registryString, path, name, value)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * registryString (string)
//  * path (string)
//  * name (string)
//  * value (uint32)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = AddRegKeyDWORD(registryString, path, name, value);
//  // obj.runtimeError
//
func (e *Engine) AddRegKeyDWORD(registryString string, path string, name string, value uint32) error {
	return errors.New("this function is unimplemented on non windows platforms")
}

// AddRegKeyQWORD - Add a QWORD registry key
//
// Package
//
// registry
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  AddRegKeyQWORD(registryString, path, name, value)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * registryString (string)
//  * path (string)
//  * name (string)
//  * value (uint64)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = AddRegKeyQWORD(registryString, path, name, value);
//  // obj.runtimeError
//
func (e *Engine) AddRegKeyQWORD(registryString string, path string, name string, value uint64) error {
	return errors.New("this function is unimplemented on non windows platforms")
}

// AddRegKeyStrings - Add a registry key of type string(s)
//
// Package
//
// registry
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  AddRegKeyStrings(registryString, path, name, value)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * registryString (string)
//  * path (string)
//  * name (string)
//  * value ([]string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = AddRegKeyStrings(registryString, path, name, value);
//  // obj.runtimeError
//
func (e *Engine) AddRegKeyStrings(registryString string, path string, name string, value []string) error {
	return errors.New("this function is unimplemented on non windows platforms")
}

// DelRegKey - Delete a registry key
//
// Package
//
// registry
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  DelRegKey(registryString, path)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * registryString (string)
//  * path (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = DelRegKey(registryString, path);
//  // obj.runtimeError
//
func (e *Engine) DelRegKey(registryString string, path string) error {
	return errors.New("this function is unimplemented on non windows platforms")
}

// DelRegKeyValue - Delete a registry key value
//
// Package
//
// registry
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  DelRegKeyValue(registryString, path, value)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * registryString (string)
//  * path (string)
//  * value (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = DelRegKeyValue(registryString, path, value);
//  // obj.runtimeError
//
func (e *Engine) DelRegKeyValue(registryString string, path string, valueName string) error {
	return errors.New("this function is unimplemented on non windows platforms")
}

// QueryRegKey - Retrive a registry key
//
// Package
//
// registry
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  QueryRegKey(registryString, path)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * registryString (string)
//  * path (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.keyObj (RegistryRetValue)
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = QueryRegKey(registryString, path);
//  // obj.keyObj
//  // obj.runtimeError
//
func (e *Engine) QueryRegKey(registryString string, path string) (RegistryRetValue, error) {
	return RegistryRetValue{}, errors.New("this function is unimplemented on non windows platforms")
}
