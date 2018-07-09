package os

import (
	goOs "os"
)

//TerminateVM will halt the execution of the current gscript
func TerminateVM() {
	//Implement later
}

//TerminateSelf will kill the current process
func TerminateSelf() error {
	pid := goOs.Getpid()
	proc, err := goOs.FindProcess(int(pid))
	if err != nil {
		return err
	}
	err = proc.Kill()
	if err != nil {
		return err
	}
	return nil
}
