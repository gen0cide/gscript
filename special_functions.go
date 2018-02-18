// +build !windows

package gscript

import "errors"

func CreateRegKeyAndValue(regHive string, keyPath string, keyObject string, keyValue interface{}) error {
	return errors.New("not implemented for this platform")
}

func DeleteRegKeysValue(regHive string, keypath string, keyobject string) error {
	return errors.New("not implemented for this platform")
}

func QueryRegKeyString(regHive string, keypath string, keyobject string) (string, error) {
	return "error", errors.New("not implemented for this platform")
}