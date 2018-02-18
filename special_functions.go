// +build !windows

package gscript

import "errors"

func CreateRegKeyAndValue(regHive string, keyPath string, keyObject string, keyValue interface{}) error {
	return errors.New("not implemented for this platform")
}