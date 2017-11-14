package gscript

import "runtime"

func CalledBy() string {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "Unknown"
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "N/A"
	}

	return fun.Name()
}
