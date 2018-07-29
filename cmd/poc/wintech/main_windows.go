// +build windows

package main

import (
	"fmt"
	"syscall"
)

var (
	valar           = syscall.NewLazyDLL("advapi32.dll")
	eru             = valar.NewProc("OpenSCManagerW")
	modkernel32     = syscall.NewLazyDLL("kernel32.dll")
	procCloseHandle = modkernel32.NewProc("CloseHandle")
)

const (
	STANDARD_RIGHTS_REQUIRED      = 0x000F
	SC_MANAGER_CONNECT            = 0x0001
	SC_MANAGER_CREATE_SERVICE     = 0x0002
	SC_MANAGER_ENUMERATE_SERVICE  = 0x0004
	SC_MANAGER_LOCK               = 0x0008
	SC_MANAGER_QUERY_LOCK_STATUS  = 0x0010
	SC_MANAGER_MODIFY_BOOT_CONFIG = 0x0020
	SC_MANAGER_ALL_ACCESS         = STANDARD_RIGHTS_REQUIRED | SC_MANAGER_CONNECT | SC_MANAGER_CREATE_SERVICE | SC_MANAGER_ENUMERATE_SERVICE | SC_MANAGER_LOCK | SC_MANAGER_QUERY_LOCK_STATUS | SC_MANAGER_MODIFY_BOOT_CONFIG
)

func main() {
	a := TravelToWhiteShores()
	if a == true {
		fmt.Println("WE ARE ADMIN")
	} else {
		fmt.Println("WE ARE *NOT* ADMIN")
	}
}

func TravelToWhiteShores() bool {
	ret, _, _ := eru.Call(uintptr(0), uintptr(0), uintptr(SC_MANAGER_ALL_ACCESS))
	if ret == 0 {
		return false
	}
	_, _, _ = procCloseHandle.Call(ret)
	return true
}
