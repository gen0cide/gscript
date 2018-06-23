// +build !windows

package engine

import "errors"

// InjectIntoProc - Inject shellcode into a provided PID
//
// Package
//
// injection
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  InjectIntoProc(shellcode, proccessID)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * shellcode (string)
//  * proccessID (int64)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = InjectIntoProc(shellcode, proccessID);
//  // obj.runtimeError
//
func (e *Engine) InjectIntoProc(shellcode string, proccessID int64) error {
	return errors.New("not implemented for this platform")
}

// InjectIntoSelf - Inject shellcode into the parrent proccess
//
// Package
//
// injection
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  InjectIntoSelf(shellcode)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * shellcode (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = InjectIntoSelf(shellcode);
//  // obj.runtimeError
//
func (e *Engine) InjectIntoSelf(shellcode string) error {
	return errors.New("not implemented for this platform")
}
