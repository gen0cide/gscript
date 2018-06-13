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
	sandboxUsers := []string{
		"sandbox",
		"virus",
		"malware",
		".bin",
		".elf",
		".exe",
		"cuckoo",
		"vagrant",
	}

	userObj, err := user.Current()
	if err != nil {
		return false, err
	}

	for _, value := range sandboxUsers {
		if strings.Contains(
			strings.ToLower(userObj.Name),
			strings.ToLower(value)) {
			return true, nil
		}
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
