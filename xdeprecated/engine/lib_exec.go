package engine

import (
	"bytes"
	"os/exec"
	"strings"
)

// ExecuteCommand - Executes system commands.
//
// Package
//
// exec
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  ExecuteCommand(baseCmd, cmdArgs)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * baseCmd (string)
//  * cmdArgs ([]string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.retObject (VMExecResponse)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = ExecuteCommand(baseCmd, cmdArgs);
//  // obj.retObject
//
func (e *Engine) ExecuteCommand(c string, args []string) VMExecResponse {
	cmd := exec.Command(c, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	respObj := VMExecResponse{
		Stdout: strings.Split(stdout.String(), "\n"),
		Stderr: strings.Split(stderr.String(), "\n"),
		PID:    cmd.Process.Pid,
	}
	if err != nil {
		respObj.ErrorMsg = err.Error()
		respObj.Success = false
	} else {
		respObj.Success = true
	}
	return respObj
}

// ForkExecuteCommand - Executes system commands via a forked call.
//
// Package
//
// exec
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  ForkExecuteCommand(baseCmd, cmdArgs)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * baseCmd (string)
//  * cmdArgs ([]string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.pid (int)
//  * obj.execError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = ForkExecuteCommand(baseCmd, cmdArgs);
//  // obj.pid
//  // obj.execError
//
func (e *Engine) ForkExecuteCommand(c string, args []string) (int, error) {
	cmd := exec.Command(c, args...)
	err := cmd.Start()
	if err != nil {
		return 0, err
	}
	pid := cmd.Process.Pid
	return pid, nil
}

// TODO: Migrate to generated functions
// func (e *Engine) VMCanSudo(call otto.FunctionCall) otto.Value {
//  resp := ExecuteCommand("sudo", "-v")
//  if resp.Success == false {
//    e.Logger.WithField("trace", "true").Errorf("OS error: %s", resp.ErrorMsg)
//    return otto.FalseValue()
//  } else if resp.Success == true {
//    return otto.TrueValue()
//  }
//  return otto.FalseValue()
// }

// func (e *Engine) VMExistsInPath(call otto.FunctionCall) otto.Value {
//  cmd := call.Argument(0)
//  cmdString, err := cmd.Export()
//  if err != nil {
//    e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
//    return otto.FalseValue()
//  }
//  path, err := exec.LookPath(cmdString.(string))
//  if err != nil {
//    e.Logger.WithField("trace", "true").Errorf("OS error: %s", err.Error())
//    return otto.FalseValue()
//  } else if path != "" {
//    return otto.TrueValue()
//  } else {
//    return otto.FalseValue()
//  }
// }

// func (e *Engine) VMCmdSuccessful(call otto.FunctionCall) otto.Value {
//  cmd := call.Argument(0)
//  cmdString, err := cmd.Export()
//  if err != nil {
//    e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
//    return otto.FalseValue()
//  }
//  arg := call.Argument(1)
//  argString, err := arg.Export()
//  if err != nil {
//    e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
//    return otto.FalseValue()
//  }
//  resp := ExecuteCommand(cmdString.(string), argString.(string))
//  if resp.Success == false {
//    e.Logger.WithField("trace", "true").Errorf("OS error: %s", resp.ErrorMsg)
//    return otto.FalseValue()
//  } else if resp.Success == true {
//    return otto.TrueValue()
//  }
//  return otto.FalseValue()
// }

// func (e *Engine) VMShellcodeExec(call otto.FunctionCall) otto.Value {
//  shellcode, err := call.Argument(0).ToString()
//  if err != nil {
//    e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
//    return otto.FalseValue()
//  }

//  pid, err := call.Argument(1).ToInteger()
//  if err != nil {
//    e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
//    return otto.FalseValue()
//  }

//  err = InjectIntoProc(shellcode, pid)
//  if err != nil {
//    e.Logger.WithField("trace", "true").Errorf("Injection error: %s", err.Error())
//    return otto.FalseValue()
//  }
//  return otto.TrueValue()
// }
