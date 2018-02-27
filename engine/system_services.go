package engine

import (
	services "github.com/gen0cide/service-go"
	"github.com/robertkrimen/otto"
)

func (e *Engine) InstallSystemService(path, name, displayName, description string) error {
	c := &services.Config{
		Path:        path,
		Name:        name,
		DisplayName: displayName,
		Description: description,
	}

	s, err := services.NewServiceConfig(c)
	if err != nil {
		return err
	}

	err = s.Install()
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) StartServiceByName(name string) error {
	c := &services.Config{
		Name: name,
	}

	s, err := services.NewServiceConfig(c)
	if err != nil {
		return err
	}

	err = s.Start()
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) StopServiceByName(name string) error {
	c := &services.Config{
		Name: name,
	}

	s, err := services.NewServiceConfig(c)
	if err != nil {
		return err
	}

	err = s.Stop()
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) RemoveServiceByName(name string) error {
	c := &services.Config{
		Name: name,
	}

	s, err := services.NewServiceConfig(c)
	if err != nil {
		return err
	}

	err = s.Remove()
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) VMInstallSystemService(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 4 {
		e.Logger.Errorf("Not enough arguments provided.")
		return otto.FalseValue()
	}

	path, err := call.Argument(0).ToString()
	if err != nil {
		e.Logger.Errorf("Error converting path to string: %s", err.Error())
		return otto.FalseValue()
	}

	name, err := call.Argument(1).ToString()
	if err != nil {
		e.Logger.Errorf("Error converting name to string: %s", err.Error())
		return otto.FalseValue()
	}

	displayName, err := call.Argument(2).ToString()
	if err != nil {
		e.Logger.Errorf("Error converting displayName to string: %s", err.Error())
		return otto.FalseValue()
	}

	description, err := call.Argument(3).ToString()
	if err != nil {
		e.Logger.Errorf("Error converting description to string: %s", err.Error())
		return otto.FalseValue()
	}

	err = e.InstallSystemService(path, name, displayName, description)
	if err != nil {
		e.Logger.Errorf("Error installing system service: %s", err.Error())
		return otto.FalseValue()
	}

	return otto.TrueValue()
}

func (e *Engine) VMStartServiceByName(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 1 {
		e.Logger.Errorf("Not enough arguments provided.")
		return otto.FalseValue()
	}

	name, err := call.Argument(0).ToString()
	if err != nil {
		e.Logger.Errorf("Error converting name to string: %s", err.Error())
		return otto.FalseValue()
	}

	err = e.StartServiceByName(name)
	if err != nil {
		e.Logger.Errorf("Error starting system service: %s", err.Error())
		return otto.FalseValue()
	}

	return otto.TrueValue()
}

func (e *Engine) VMStopServiceByName(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 1 {
		e.Logger.Errorf("Not enough arguments provided.")
		return otto.FalseValue()
	}

	name, err := call.Argument(0).ToString()
	if err != nil {
		e.Logger.Errorf("Error converting name to string: %s", err.Error())
		return otto.FalseValue()
	}

	err = e.StopServiceByName(name)
	if err != nil {
		e.Logger.Errorf("Error starting system service: %s", err.Error())
		return otto.FalseValue()
	}

	return otto.TrueValue()
}

func (e *Engine) VMRemoveServiceByName(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 1 {
		e.Logger.Errorf("Not enough arguments provided.")
		return otto.FalseValue()
	}

	name, err := call.Argument(0).ToString()
	if err != nil {
		e.Logger.Errorf("Error converting name to string: %s", err.Error())
		return otto.FalseValue()
	}

	err = e.RemoveServiceByName(name)
	if err != nil {
		e.Logger.Errorf("Error starting system service: %s", err.Error())
		return otto.FalseValue()
	}

	return otto.TrueValue()
}
