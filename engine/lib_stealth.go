package engine

import (
	"os/user"
	"runtime"
	"strings"
)

// CheckSandboxUsernames - Checks to see if the current username contains a series of typical sandbox strings.
//
// Package
//
// stealth
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  CheckSandboxUsernames()
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.areWeInASandbox (bool)
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = CheckSandboxUsernames();
//  // obj.areWeInASandbox
//  // obj.runtimeError
//
func (e *Engine) CheckSandboxUsernames() (bool, error) {
	userObj, err := user.Current()
	if err != nil {
		return false, err
	}
	if strings.Contains(userObj.Name, "SANDBOX") {
		return true, nil
	}
	if strings.Contains(userObj.Name, "VIRUS") {
		return true, nil
	}
	if strings.Contains(userObj.Name, "MALWARE") {
		return true, nil
	}
	if strings.Contains(userObj.Name, "malware") {
		return true, nil
	}
	if strings.Contains(userObj.Name, "virus") {
		return true, nil
	}
	if strings.Contains(userObj.Name, "sandbox") {
		return true, nil
	}
	if strings.Contains(userObj.Name, ".bin") {
		return true, nil
	}
	if strings.Contains(userObj.Name, ".elf") {
		return true, nil
	}
	if strings.Contains(userObj.Name, ".exe") {
		return true, nil
	}
	return false, nil
}

// CheckIfCPUCountIsHigherThanOne - Checks to see if the system we are on has at least 2 CPUs (or cores) (sandbox check).
//
// Package
//
// stealth
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  CheckIfCPUCountIsHigherThanOne()
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.areWeInASandbox (bool)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = CheckIfCPUCountIsHigherThanOne();
//  // obj.areWeInASandbox
//
func (e *Engine) CheckIfCPUCountIsHigherThanOne() bool {
	if runtime.NumCPU() > 1 {
		return true
	}
	return false
}

// CheckIfRAMAmountIsBelow1GB - Checks to see if the system we are on has at least 1GB of ram (sandbox check).
//
// Package
//
// stealth
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  CheckIfRAMAmountIsBelow1GB()
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.areWeInASandbox (bool)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = CheckIfRAMAmountIsBelow1GB();
//  // obj.areWeInASandbox
//
func (e *Engine) CheckIfRAMAmountIsBelow1GB() bool {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	if memStats.Sys > uint64(1000000000) {
		return true
	}
	return false
}
