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

//ReadFileAsString takes a file path and reads that files contents and returns a string representation of the contents
func ReadFileAsString(readPath string) (string, error) {
	absPath, err := filepath.Abs(readPath)
	if err != nil {
		return "", err
	}
	contents, err := ioutil.ReadFile(absPath)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

//ReadFileAsBytes takes a file path and reads that files contents and returns a byte array of the contents
func ReadFileAsBytes(readPath string) ([]byte, error) {
	absPath, err := filepath.Abs(readPath)
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	return contents, nil
}
