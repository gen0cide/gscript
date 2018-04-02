package engine

import (
	"debug/pe"
	"errors"
	"os"
)

func getPEImgSize(peFile string) (int64, error) {
	// open self
	mySelf, err := pe.Open(peFile)
	if err != nil {
		return int64(0), err
	}
	defer mySelf.Close()

	// get img size
	var imgSize uint32
	switch mySelf.Machine {
	case 0x14c: // x86
		imgSize = mySelf.OptionalHeader.(*pe.OptionalHeader32).SizeOfImage
	case 0x8664: // x86_64
		imgSize = mySelf.OptionalHeader.(*pe.OptionalHeader64).SizeOfImage
	default:
		return int64(0), errors.New("Binary architecture unsupported for this action.")
	}
	return int64(imgSize), nil
}

// RetrievePEPolymorphicData - Retrive data stored within uninitalized space at the end of the gscript binary
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
//  RetrievePEPolymorphicData(peFile)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * peFile (string)
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
//  var obj = RetrievePEPolymorphicData(peFile);
//  // obj.data
//  // obj.runtimeError
//
func (e *Engine) RetrievePEPolymorphicData(peFile string) ([]byte, error) {
	var polymorphicData []byte
	pEImgSize, err := getPEImgSize(peFile)
	if err != nil {
		return polymorphicData, err
	}
	file, err := os.Open(os.Args[0])
	if err != nil {
		return polymorphicData, err
	}
	defer file.Close()
	_, err = file.ReadAt(polymorphicData, pEImgSize)
	return polymorphicData, err
}

// WritePEPolymorphicData - Write data to the uninitalized space at the end of the gscript binary
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
//  WritePEPolymorphicData(peFile, data)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * peFile (string)
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
//  var obj = WritePEPolymorphicData(peFile, data);
//  // obj.runtimeError
//
func (e *Engine) WritePEPolymorphicData(peFile string, polymorphicData []byte) error {
	pEImgSize, err := getPEImgSize(peFile)
	if err != nil {
		return err
	}
	file, err := os.Open(peFile)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteAt(polymorphicData, pEImgSize)
	if err != nil {
		return err
	}
	return nil
}
