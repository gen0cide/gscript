package engine

import (
	"errors"
)

// RetrievePolymorphicData - Retrive data stored within uninitalized space at the end of the gscript binary
//
// Package
//
// polymorph
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  RetrievePolymorphicData()
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.data ([]byte)
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = RetrievePolymorphicData();
//  // obj.data
//  // obj.runtimeError
//
func (e *Engine) RetrievePolymorphicData() ([]byte, error) {
	return []byte{}, errors.New("not implemented for this platform")
}

// WritePolymorphicData - Write data to the uninitalized space at the end of the gscript binary
//
// Package
//
// polymorph
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  WritePolymorphicData(data)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * data ([]byte)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = WritePolymorphicData(data);
//  // obj.runtimeError
//
func (e *Engine) WritePolymorphicData(polymorphicData []byte) error {
	return errors.New("not implemented for this platform")
}
