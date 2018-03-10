// +build !windows

package engine

import "errors"

func InjectIntoProc(shellcode string, proccessID int64) error {
	return errors.New("not implemented for this platform")
}
