package gscript

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"
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
