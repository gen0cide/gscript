package engine

import (
	"errors"
	"os"
	"strings"
	"syscall"

	services "github.com/gen0cide/service-go"
	"github.com/mitchellh/go-ps"
)

type RegistryRetValue struct {
	ValType        string   `json:"return_type"`
	StringVal      string   `json:"string_val"`
	StringArrayVal []string `json:"string_array_val"`
	ByteArrayVal   []byte   `json:"byte_array_val"`
	IntVal         uint32   `json:"int_val"`
	LongVal        uint64   `json:"long_val"`
}

// FindProcByName - Returns the Pid of a given proccess, if the proccess can not be found, an error is returned
//
// Package
//
// os
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  FindProcByName(procName)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * procName (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.pid (int)
//  * obj.procError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = FindProcByName(procName);
//  // obj.pid
//  // obj.procError
//
func (e *Engine) FindProcByName(procName string) (int, error) {
	// * FIXME: this function currently matches against the name of the executible which is NOT technically the proccess name
	procs, err := ps.Processes()
	if err != nil {
		return -1, err
	}
	for _, proc := range procs {
		if procName == proc.Executable() {
			return proc.Pid(), nil
		}
	}
	return -1, errors.New("processes name not found")
}

// InstallSystemService - Installs a target binary as a system service
//
// Package
//
// os
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  InstallSystemService(path, name, displayName, description)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * path (string)
//  * name (string)
//  * displayName (string)
//  * description (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.installError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = InstallSystemService(path, name, displayName, description);
//  // obj.installError
//
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

// StartServiceByName - Starts a system service
//
// Package
//
// os
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  StartServiceByName(name)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * name (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.installError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = StartServiceByName(name);
//  // obj.installError
//
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

// StopServiceByName - Stops a system service
//
// Package
//
// os
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  StopServiceByName(name)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * name (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.installError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = StopServiceByName(name);
//  // obj.installError
//
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

// RemoveServiceByName - Uninstalls a system service
//
// Package
//
// os
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  RemoveServiceByName(name)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * name (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.removealError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = RemoveServiceByName(name);
//  // obj.removealError
//
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

// Signal - Sends a signal to a target proccess
//
// Package
//
// os
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  Signal(signal, pid)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * signal (int)
//  * pid (int)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = Signal(signal, pid);
//  // obj.runtimeError
//
func (e *Engine) Signal(proc int, sig int) error {
	foundProc, err := os.FindProcess(proc)
	if err != nil {
		return err
	}
	return foundProc.Signal(syscall.Signal(sig))
}

// RunningProcs - Returns an array of int's representing active PIDs currently running
//
// Package
//
// os
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  RunningProcs()
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.pids ([]int)
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = RunningProcs();
//  // obj.pids
//  // obj.runtimeError
//
func (e *Engine) RunningProcs() ([]int, error) {
	var pids []int
	procs, err := ps.Processes()
	if err != nil {
		return pids, err
	}
	for _, proc := range procs {
		pids = append(pids, proc.Pid())
	}
	return pids, nil
}

// GetProcName - Returns the name of a target proccess
//
// Package
//
// os
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  GetProcName(pid)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * pid (int)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.procName (string)
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = GetProcName(pid);
//  // obj.procName
//  // obj.runtimeError
//
func (e *Engine) GetProcName(pid int) (string, error) {
	// * FIXME: this function currently returns the name of the executible which is NOT technically the proccess name
	proc, err := ps.FindProcess(pid)
	if err != nil {
		return "", err
	}
	return proc.Executable(), nil
}

// EnvVars - Returns a map of enviornment variable names to their corrisponding values.
//
// Package
//
// os
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  EnvVars()
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.vars (map[string]string)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = EnvVars();
//  // obj.vars
//
func (e *Engine) EnvVars() map[string]string {
	vars := make(map[string]string)
	for _, eVar := range os.Environ() {
		eVarSegments := strings.Split(eVar, "=")
		if len(eVarSegments) > 1 {
			vars[eVarSegments[0]] = eVarSegments[1]
		}
	}
	return vars
}

// GetEnvVar - Returns the value of a given enviornment variable
//
// Package
//
// os
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  GetEnvVar(vars)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * vars (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.value (string)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = GetEnvVar(vars);
//  // obj.value
//
func (e *Engine) GetEnvVar(eVar string) string {
	return os.Getenv(eVar)
}
