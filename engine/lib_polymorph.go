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
