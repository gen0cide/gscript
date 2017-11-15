package gscript

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"math/rand"
	"time"
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

func LocalCreateFile(path string, bytes []byte) bool {
	if LocalFileExists(path) {
		return false
		//errors.New("The file to create already exists so we won't overwite it")
	}
	err := ioutil.WriteFile(path, bytes, 0700)
	if err != nil {
		return false
	}
	return true
}

// Append adds input to the end of filename.
func LocalAppendFile(filename string, bytes []byte) bool {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false
	}
	file, err := os.OpenFile(filename, os.O_APPEND, fileInfo.Mode())
	if err != nil {
		return false
	}
	if _, err = file.Write(bytes); err != nil {
		return false
	}
	file.Close()
	return true
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

func ExecuteCommand(c string, args ...string) VMExecResponse {
	cmd := exec.Command(c, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	respObj := VMExecResponse{
		Stdout: strings.Split(stdout.String(), "\n"),
		Stderr: strings.Split(stderr.String(), "\n"),
		PID:    cmd.Process.Pid,
	}
	if err != nil {
		respObj.ErrorMsg = err.Error()
		respObj.Success = false
	} else {
		respObj.Success = true
	}
	return respObj
}

func HTTPGetFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	pageData, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return pageData, nil
}

// RandString returns a string the length of strlen
func RandString(strlen int) string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

// RandomInt returns an int inbetween min and max.
func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}
