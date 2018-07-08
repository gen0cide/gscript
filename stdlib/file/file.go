package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// WriteFileFromBytes writes data from a byte array to a dest filepath with the dest parent dirs permissions.
func WriteFileFromBytes(destPath string, fileData []byte) error {
	absDir, err := filepath.Abs(filepath.Dir(destPath))
	if err != nil {
		return err
	}
	dirInfo, err := os.Stat(absDir)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(destPath, fileData, os.FileMode(dirInfo.Mode()))
	if err != nil {
		return err
	}
	return nil
}

// WriteFileFromString writes data from a string to a dest filepath with the dest parent dirs permissions.
func WriteFileFromString(destPath string, fileData string) error {
	absDir, err := filepath.Abs(filepath.Dir(destPath))
	if err != nil {
		return err
	}
	dirInfo, err := os.Stat(absDir)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(destPath, []byte(fileData), os.FileMode(dirInfo.Mode()))
	if err != nil {
		return err
	}
	return nil
}
