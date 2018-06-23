// +build !windows

package engine

import "errors"

func (e *Engine) MakeUnDebuggable() (bool, error) {
	return false, errors.New("not implemented for this platform")
}
func (e *Engine) MakeDebuggable() (bool, error) {
	return false, errors.New("not implemented for this platform")
}
func (e *Engine) IsDebuggerPresent() (bool, error) {
	return false, errors.New("not implemented for this platform")
}
func (e *Engine) CheckIfWineGetUnixFileNameExists() (bool, error) {
	return false, errors.New("not implemented for this platform")
}

/*
	func (e *Engine) CheckIfSandboxieExists() (bool, error) {
		return false, errors.New("not implemented for this platform")
	}
*/
