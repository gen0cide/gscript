package typelib

import "syscall"

func TakePointer(arg uintptr) error {
	return nil
}

func TakeHandle(arg syscall.Signal) error {
	return nil
}
