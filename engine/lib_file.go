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
	if _, err := os.Stat(path); os.IsNotExist(err) {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return 0, err
	}
	err := ioutil.WriteFile(path, fileData, os.FileMode(uint32(perms)))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Error writing the file: %s", err.Error())
		return 0, err
	}
	return len(fileData), nil
}

// AppendFileBytes - Addes a given byte array to the end of a file
//
// Package
//
// file
//
// Author
//
// - ahhh (https://github.com/ahhh)
//
// Javascript
//
// Here is the Javascript method signature:
//  AppendFileBytes(path, fileData)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * path (string)
//  * fileData ([]byte)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.fileError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = AppendFileBytes(path, fileData);
//  // obj.fileError
//
func (e *Engine) AppendFileBytes(path string, fileData []byte) error {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, fileInfo.Mode())
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return err
	}
	defer file.Close()
	if _, err = file.Write(fileData); err != nil {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return err
	}
	return nil
}

// AppendFileString - Addes a given string to the end of a file
//
// Package
//
// file
//
// Author
//
// - ahhh (https://github.com/ahhh)
//
// Javascript
//
// Here is the Javascript method signature:
//  AppendFileString(path, addString)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * path (string)
//  * addString (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.fileError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = AppendFileString(path, addString);
//  // obj.fileError
//
func (e *Engine) AppendFileString(path, addString string) error {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return err
	}
	file, err := os.OpenFile(path, os.O_APPEND, fileInfo.Mode())
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return err
	}
	if _, err = file.WriteString(addString); err != nil {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return err
	}
	file.Close()
	return nil
}

// DeleteFile - Deletes a file at a given path or returns an error
//
// Package
//
// file
//
// Author
//
// - ahhh (https://github.com/ahhh)
//
// Javascript
//
// Here is the Javascript method signature:
//  DeleteFile(path)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * path (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.fileError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = DeleteFile(path);
//  // obj.fileError
//
func (e *Engine) DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return err
	}
	return nil
}

// ReadFile - Reads a file path and returns a byte array
//
// Package
//
// file
//
// Author
//
// - ahhh (https://github.com/ahhh)
//
// Javascript
//
// Here is the Javascript method signature:
//  ReadFile(path)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * path (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.fileBytes ([]byte)
//  * obj.fileError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = ReadFile(path);
//  // obj.fileBytes
//  // obj.fileError
//
func (e *Engine) ReadFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return []byte{}, err
	}
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Error writing the file: %s", err.Error())
		return nil, err
	}
	return dat, nil
}

// FileExists - Checks if a file exists and returns a bool
//
// Package
//
// file
//
// Author
//
// - ahhh (https://github.com/ahhh)
//
// Javascript
//
// Here is the Javascript method signature:
//  FileExists(path)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * path (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.fileExists (bool)
//  * obj.fileError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = FileExists(path);
//  // obj.fileExists
//  // obj.fileError
//
func (e *Engine) FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return false, err
	}
	return true, nil
}

// CreateDir - Creates a directory at a given path or return an error
//
// Package
//
// file
//
// Author
//
// - ahhh (https://github.com/ahhh)
//
// Javascript
//
// Here is the Javascript method signature:
//  CreateDir(path)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * path (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.fileError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = CreateDir(path);
//  // obj.fileError
//
func (e *Engine) CreateDir(path string) error {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Error writing the file: %s", err.Error())
		return err
	}
	return nil
}
