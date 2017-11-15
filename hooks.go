package gscript

import (
	"errors"

	"github.com/davecgh/go-spew/spew"

	// Include Underscore In Otto :)
	_ "github.com/robertkrimen/otto/underscore"
)

var (
	Debugger = true
)

func (e *Engine) RunBeforeDeploy() error {
	result, err := e.VM.Run(`BeforeDeploy()`)
	if err != nil {
		return err
	}
	boolResult, err := result.ToBoolean()
	if err != nil {
		e.LogErrorf("Boolean Conversion Error: function=%s result=%s error=%s", CalledBy(), spew.Sdump(result), err.Error())
		return err
	}
	if !boolResult {
		e.LogErrorf("BeforeDeploy() returned false.")
		return errors.New("failed BeforeDeploy()")
	}
	e.LogDebugf("Completed BeforeDeploy()")
	return nil
}

func (e *Engine) RunDeploy() error {
	result, err := e.VM.Run(`Deploy()`)
	if err != nil {
		return err
	}
	boolResult, err := result.ToBoolean()
	if err != nil {
		e.LogErrorf("Boolean Conversion Error: function=%s result=%s error=%s", CalledBy(), spew.Sdump(result), err.Error())
		return err
	}
	if !boolResult {
		e.LogErrorf("Deploy() returned false.")
		return errors.New("failed Deploy()")
	}
	e.LogDebugf("Completed Deploy()")
	return nil
}

func (e *Engine) RunAfterDeploy() error {
	result, err := e.VM.Run(`AfterDeploy()`)
	if err != nil {
		return err
	}
	boolResult, err := result.ToBoolean()
	if err != nil {
		e.LogErrorf("Boolean Conversion Error: function=%s result=%s error=%s", CalledBy(), spew.Sdump(result), err.Error())
		return err
	}
	if !boolResult {
		e.LogErrorf("AfterDeploy() returned false.")
		return errors.New("failed AfterDeploy()")
	}
	e.LogDebugf("Completed AfterDeploy()")
	return nil
}

func (e *Engine) RunOnError() error {
	result, err := e.VM.Run(`OnError()`)
	if err != nil {
		e.LogErrorf("OnError() Error: %s", err.Error())
		return err
	}
	if Debugger {
		e.LogErrorf("OnError(): %s", result.String())
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
