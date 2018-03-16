package engine

import (
	"os/user"
	"runtime"
	"strings"
)

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

func (e *Engine) CheckIfCPUCountIsHigherThanOne() bool {
	if runtime.NumCPU() > 1 {
		return true
	}
	return false
}
func (e *Engine) CheckIfRAMAmountIsBelow1GB() bool {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	if memStats.Sys > uint64(1000000000) {
		return true
	}
	return false
}
