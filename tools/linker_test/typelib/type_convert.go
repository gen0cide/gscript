package typelib

import "syscall"

// TakePointer is a test function for translating uintptr types
func TakePointer(arg uintptr) error {
	return nil
}

// TakeHandle is a test function for converting aliased types
func TakeHandle(arg syscall.Signal) error {
	return nil
}
