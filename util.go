package gscript

import (
	"runtime"
	"io/ioutil"
	"os"
	"errors"
)

func CalledBy() string {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "Unknown"
	}
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "N/A"
	}
	return fun.Name()
}

func LocalFileExists(path string) bool {
  _, err := os.Stat(path)
  if err == nil {
    return true
  }
  return false
}

func LocalCreateFile(bytes []byte, path string) error {
	if LocalFileExists(path) {
    return errors.New("The file to create already exists so we won't overwite it")
  }
  err := ioutil.WriteFile(path, bytes, 0700)
  if err != nil {
      return err
  }
  return nil
}

func LocalReadFile(path string) ([]byte, error) {
	if LocalFileExists(path) {
    dat, err := ioutil.ReadFile(path)
    if err != nil {
      return nil, err
    }
  return dat, nil
 }
 return nil, errors.New("The file to read does not exist")
}
