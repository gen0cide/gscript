package engine

import (
	"errors"
	"time"

	"github.com/gen0cide/otto"
)

var errTimeout = errors.New("Script Timeout")

// SetEntryPoint sets the function name of the entry point function for the script
func (e *Engine) SetEntryPoint(fnName string) {
	e.EntryPoint = fnName
}

// RunWithTimeout evaluates an expression in the VM that honors the VMs timeout setting
func (e *Engine) RunWithTimeout(command string) (otto.Value, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		if caught := recover(); caught != nil {
			if caught == errTimeout {
				e.Logger.Errorf("VM hit timeout after %v seconds", duration)
				return
			}
		}
	}()

	e.VM.Interrupt = make(chan func(), 1)

	go func() {
		for {
			time.Sleep(time.Duration(e.Timeout) * time.Second)
			if e.Paused {
				e.VM.Interrupt <- func() {
					panic(errTimeout)
				}
				return
			}
		}
	}()

	return e.VM.Run(command)
}
