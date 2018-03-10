package engine

import (
	"io/ioutil"
	"os"
)

// WriteFile - Writes data from a byte array to a file with the given permissions.
//
// Package
//
// file
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  WriteFile(path, fileData, perms)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * path (string)
//  * fileData ([]byte)
//  * perms (int64)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.bytesWritten (int)
//  * obj.fileError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = WriteFile(path, fileData, perms);
//  // obj.bytesWritten
//  // obj.fileError
//
func (e *Engine) WriteFile(path string, fileData []byte, perms int64) (int, error) {
	err := ioutil.WriteFile(path, fileData, os.FileMode(uint32(perms)))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Error writing the file: %s", err.Error())
		return 0, err
	}
	return len(fileData), nil
}
