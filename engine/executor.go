package engine

import (
	"errors"
	"time"

	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/ast"
)

var errTimeout = errors.New("Script Timeout")

// SetEntryPoint sets the function name of the entry point function for the script
func (e *Engine) SetEntryPoint(fnName string) {
	e.EntryPoint = fnName
}

// LoadScriptWithTimeout evaluates an expression in the VM that honors the VMs timeout setting
func (e *Engine) LoadScriptWithTimeout(script *ast.Program) (otto.Value, error) {
	start := time.Now()
	defer e.recoveryHandler(start)
	e.VM.Interrupt = make(chan func(), 1)
	go e.timeoutMonitor()
	return e.VM.Eval(script)
}

// CallFunctionWithTimeout calls a given top level function in the VM that honors the VMs timeout setting
func (e *Engine) CallFunctionWithTimeout(fn string) (otto.Value, error) {
	start := time.Now()
	defer e.recoveryHandler(start)
	e.VM.Interrupt = make(chan func(), 1)
	go e.timeoutMonitor()
	return e.VM.Call(fn, nil, nil)
}

func (e *Engine) recoveryHandler(begin time.Time) {
	duration := time.Since(begin)
	if caught := recover(); caught != nil {
		if caught == errTimeout {
			e.Logger.Errorf("VM hit timeout after %v seconds", duration)
			return
		}
		e.Logger.Errorf("VM encountered unexpected error: %v", caught)
		return
	}
	return
}

func (e *Engine) timeoutMonitor() {
	for {
		time.Sleep(time.Duration(e.Timeout) * time.Second)
		if !e.Paused {
			e.VM.Interrupt <- func() {
				panic(errTimeout)
			}
			return
		} else {
			continue
		}
	}
}
