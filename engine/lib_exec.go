package engine

import (
	"bytes"
	"os/exec"
	"strings"
)

func (e *Engine) ExecuteCommand(c string, args []string) VMExecResponse {
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

func (e *Engine) ForkExecuteCommand(c string, args []string) (int, error) {
	cmd := exec.Command(c, args...)
	err := cmd.Start()
	if err != nil {
		return 0, err
	}
	pid := cmd.Process.Pid
	return pid, nil
}
