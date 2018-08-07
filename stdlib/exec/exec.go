package exec

import (
	"bytes"
	"fmt"
	executer "os/exec"
	"os/signal"
	"syscall"
)

// ExecuteCommand executes the given command and waits for it to complete, returning pid, stdout, stderr, exitCode, or any errors
func ExecuteCommand(c string, args []interface{}) (int, string, string, int, error) {
	cmd := executer.Command(c, getArgs(args)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	var pid, exitCode int
	if cmd.Process == nil {
		pid = -1
	} else {
		pid = cmd.Process.Pid
	}
	if cmd.ProcessState == nil {
		exitCode = 1
	} else {
		ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)
		if !ok {
			exitCode = 0
		} else {
			exitCode = ws.ExitStatus()
		}
	}
	return pid, stdout.String(), stderr.String(), exitCode, err
}

// ExecuteCommandAsync runs the command and does not wait for it to return.
func ExecuteCommandAsync(c string, args []interface{}) (*executer.Cmd, error) {
	signal.Ignore(syscall.SIGHUP)
	cmd := executer.Command(c, getArgs(args)...)
	err := cmd.Start()
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func getArgs(a []interface{}) []string {
	ret := []string{}
	for _, s := range a {
		strVal, ok := s.(string)
		if !ok {
			ret = append(ret, fmt.Sprintf("%v", s))
		} else {
			if len(strVal) > 0 {
				ret = append(ret, strVal)
			}
		}
	}
	return ret
}
