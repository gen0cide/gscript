package engine

import (
	"errors"
	"strconv"
	"time"

	"github.com/robertkrimen/otto"
)

var (
	Debugger   = true
	errTimeout = errors.New("GSE VM Timeout")
)

func (e *Engine) RunWithTimeout(command string) (otto.Value, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		if caught := recover(); caught != nil {
			if caught == errTimeout {
				e.Logger.WithField("trace", "true").Errorf("Some code took to long! Stopping after: %v", duration)
				return
			}
			e.Logger.WithField(
				"script",
				e.Name,
			).WithField(
				"line",
				strconv.Itoa(e.VM.Context().Line),
			).Fatalf("Timer experienced fatal error: %s", caught)
		}
		e.Logger.WithField(
			"script",
			e.Name,
		).WithField(
			"line",
			strconv.Itoa(e.VM.Context().Line),
		).Debugf("%s Execution Time: %v", command, duration)
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

func (e *Engine) runBeforeDeploy() error {
	if e.Halted {
		return errors.New("VM Halted")
	}
	result, err := e.RunWithTimeout(`BeforeDeploy()`)
	if err != nil {
		return err
	}
	boolResult, err := result.ToBoolean()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return err
	}
	if !boolResult {
		e.Logger.WithField("trace", "true").Errorf("BeforeDeploy() returned false.")
		return errors.New("failed BeforeDeploy()")
	}
	e.Logger.WithField("trace", "true").Debugf("Completed BeforeDeploy()")
	return nil
}

func (e *Engine) runDeploy() error {
	if e.Halted {
		return errors.New("VM Halted")
	}
	result, err := e.RunWithTimeout(`Deploy()`)
	if err != nil {
		return err
	}
	boolResult, err := result.ToBoolean()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return err
	}
	if !boolResult {
		e.Logger.WithField("trace", "true").Errorf("Deploy() returned false.")
		return errors.New("failed Deploy()")
	}
	e.Logger.WithField("trace", "true").Debugf("Completed Deploy()")
	return nil
}

func (e *Engine) runAfterDeploy() error {
	if e.Halted {
		return errors.New("VM Halted")
	}
	result, err := e.RunWithTimeout(`AfterDeploy()`)
	if err != nil {
		return err
	}
	boolResult, err := result.ToBoolean()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return err
	}
	if !boolResult {
		e.Logger.WithField("trace", "true").Errorf("AfterDeploy() returned false.")
		return errors.New("failed AfterDeploy()")
	}
	e.Logger.WithField("trace", "true").Debugf("Completed AfterDeploy()")
	return nil
}

func (e *Engine) runOnError() error {
	if e.Halted {
		return errors.New("VM Halted")
	}
	result, err := e.RunWithTimeout(`OnError()`)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("OnError() Error: %s", err.Error())
		return err
	}
	if Debugger {
		e.Logger.WithField("trace", "true").Errorf("OnError(): %s", result.String())
		return nil
	}
	return nil
}

func (e *Engine) ExecutePlan() error {
	err := e.runBeforeDeploy()
	if err != nil {
		e.runOnError()
		return err
	}
	err = e.runDeploy()
	if err != nil {
		e.runOnError()
		return err
	}
	err = e.runAfterDeploy()
	if err != nil {
		e.runOnError()
		return err
	}
	return nil
}
