package engine

import (
	"os"
	"syscall"
)

func (e *Engine) MakeUnDebuggable() (bool, error) {
	debuggerPresent, err := e.IsDebuggerPresent()
	if err != nil {
		return false, err
	}
	if !debuggerPresent {
		funcAddr, err := syscall.MustLoadDLL("Kernel32.dll").FindProc("DebugActiveProcess")
		if err != nil {
			return false, err
		}
		retVal, _, err := funcAddr.Call(uintptr(os.Getpid()), 0, 0, 0, 0)
		if err != nil {
			return false, err
		}
		if retVal > 0 {
			return true, nil
		}
	}
	return true, nil
}

func (e *Engine) MakeDebuggable() (bool, error) {
	debuggerPresent, err := e.IsDebuggerPresent()
	if err != nil {
		return false, err
	}
	if !debuggerPresent {
		funcAddr, err := syscall.MustLoadDLL("Kernel32.dll").FindProc("DebugActiveProcessStop")
		if err != nil {
			return false, err
		}
		retVal, _, err := funcAddr.Call(uintptr(os.Getpid()), 0, 0, 0, 0)
		if err != nil {
			return false, err
		}
		if retVal > 0 {
			return true, nil
		}
	}
	return true, nil
}

func (e *Engine) IsDebuggerPresent() (bool, error) {
	funcAddr, err := syscall.MustLoadDLL("Kernel32.dll").FindProc("IsDebuggerPresent")
	if err != nil {
		return false, err
	}
	retVal, _, err := funcAddr.Call(0, 0, 0, 0, 0)
	if err != nil {
		return false, err
	}
	if retVal > 0 {
		return true, nil
	}
	return false, nil
}

func (e *Engine) CheckIfWineGetUnixFileNameExists() (bool, error) {
	_, err := syscall.MustLoadDLL("Kernel32.dll").FindProc("wine_get_unix_file_name")
	if err != nil {
		return false, err
	}
	return true, nil
}

/*
	// need to fix this so that rather than checking if the dll is loadble, it checks to see if the dll IS loaded
	func (e *Engine) CheckIfSandboxieExists() (bool, error) {
		_, err := syscall.LoadDLL("sbiedll.dll")
		if err != nil {
			return false, err
		}
		return true, nil
	}
*/
