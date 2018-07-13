package exec

import (
	"bytes"
	executer "os/exec"
	"os/signal"
	"syscall"
)

// ExecuteCommand executes the given command and waits for it to complete, returning pid, stdout, stderr, exitCode, or any errors
func ExecuteCommand(c string, args []string) (int, string, string, int, error) {
	cmd := executer.Command(c, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	pid := cmd.Process.Pid
	sysStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
	if sysStatus.ExitStatus() != 1 {
		return pid, stdout.String(), stderr.String(), sysStatus.ExitStatus(), err
	}
	return pid, stdout.String(), stderr.String(), sysStatus.ExitStatus(), nil
}

// ExecuteCommandAsync runs the command and does not wait for it to return.
func ExecuteCommandAsync(c string, args []string) (*executer.Cmd, error) {
	signal.Ignore(syscall.SIGHUP)
	cmd := executer.Command(c, args...)
	err := cmd.Start()
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}
