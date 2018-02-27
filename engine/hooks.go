package engine

import (
	"errors"

	"github.com/davecgh/go-spew/spew"
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
		e.Logger.Errorf("Boolean Conversion Error: function=%s result=%s error=%s", CalledBy(), spew.Sdump(result), err.Error())
		return err
	}
	if !boolResult {
		e.Logger.Errorf("BeforeDeploy() returned false.")
		return errors.New("failed BeforeDeploy()")
	}
	e.Logger.Debugf("Completed BeforeDeploy()")
	return nil
}

func (e *Engine) RunDeploy() error {
	result, err := e.VM.Run(`Deploy()`)
	if err != nil {
		return err
	}
	boolResult, err := result.ToBoolean()
	if err != nil {
		e.Logger.Errorf("Boolean Conversion Error: function=%s result=%s error=%s", CalledBy(), spew.Sdump(result), err.Error())
		return err
	}
	if !boolResult {
		e.Logger.Errorf("Deploy() returned false.")
		return errors.New("failed Deploy()")
	}
	e.Logger.Debugf("Completed Deploy()")
	return nil
}

func (e *Engine) RunAfterDeploy() error {
	result, err := e.VM.Run(`AfterDeploy()`)
	if err != nil {
		return err
	}
	boolResult, err := result.ToBoolean()
	if err != nil {
		e.Logger.Errorf("Boolean Conversion Error: function=%s result=%s error=%s", CalledBy(), spew.Sdump(result), err.Error())
		return err
	}
	if !boolResult {
		e.Logger.Errorf("AfterDeploy() returned false.")
		return errors.New("failed AfterDeploy()")
	}
	e.Logger.Debugf("Completed AfterDeploy()")
	return nil
}

func (e *Engine) RunOnError() error {
	result, err := e.VM.Run(`OnError()`)
	if err != nil {
		e.Logger.Errorf("OnError() Error: %s", err.Error())
		return err
	}
	if Debugger {
		e.Logger.Errorf("OnError(): %s", result.String())
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
