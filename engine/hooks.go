package engine

import (
	"errors"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
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
		).Infof("%s Execution Time: %v", command, duration)
	}()

	e.VM.Interrupt = make(chan func(), 1) // The buffer prevents blocking

	go func() {
		time.Sleep(time.Duration(e.Timeout) * time.Second) // Stop after two seconds
		e.VM.Interrupt <- func() {
			panic(errTimeout)
		}
	}()

	return e.VM.Run(command)
}

func (e *Engine) RunBeforeDeploy() error {
	result, err := e.RunWithTimeout(`BeforeDeploy()`)
	if err != nil {
		return err
	}
	boolResult, err := result.ToBoolean()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Boolean Conversion Error: function=%s result=%s error=%s", CalledBy(), spew.Sdump(result), err.Error())
		return err
	}
	if !boolResult {
		e.Logger.WithField("trace", "true").Errorf("BeforeDeploy() returned false.")
		return errors.New("failed BeforeDeploy()")
	}
	e.Logger.WithField("trace", "true").Debugf("Completed BeforeDeploy()")
	return nil
}

func (e *Engine) RunDeploy() error {
	result, err := e.RunWithTimeout(`Deploy()`)
	if err != nil {
		return err
	}
	boolResult, err := result.ToBoolean()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Boolean Conversion Error: function=%s result=%s error=%s", CalledBy(), spew.Sdump(result), err.Error())
		return err
	}
	if !boolResult {
		e.Logger.WithField("trace", "true").Errorf("Deploy() returned false.")
		return errors.New("failed Deploy()")
	}
	e.Logger.WithField("trace", "true").Debugf("Completed Deploy()")
	return nil
}

func (e *Engine) RunAfterDeploy() error {
	result, err := e.RunWithTimeout(`AfterDeploy()`)
	if err != nil {
		return err
	}
	boolResult, err := result.ToBoolean()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Boolean Conversion Error: function=%s result=%s error=%s", CalledBy(), spew.Sdump(result), err.Error())
		return err
	}
	if !boolResult {
		e.Logger.WithField("trace", "true").Errorf("AfterDeploy() returned false.")
		return errors.New("failed AfterDeploy()")
	}
	e.Logger.WithField("trace", "true").Debugf("Completed AfterDeploy()")
	return nil
}

func (e *Engine) RunOnError() error {
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
	err := e.RunBeforeDeploy()
	if err != nil {
		e.RunOnError()
		return err
	}
	err = e.RunDeploy()
	if err != nil {
		e.RunOnError()
		return err
	}
	err = e.RunAfterDeploy()
	if err != nil {
		e.RunOnError()
		return err
	}
	return nil
}
