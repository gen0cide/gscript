// +build windows

package gscript

import (
	"errors"
	"fmt"
	reg "golang.org/x/sys/windows/registry"
)

// CreateRegKeyAndValue creates a new regestry key in a dynamic hive, dynamic path, dynamic object, and infers the type as either string or uint32 and creats the correct key type accordingly
func CreateRegKeyAndValue(regHive string, keyPath string, keyObject string, keyValue interface{}) error {
	var hive reg.Key
	if regHive == "LM" {
		hive = reg.LOCAL_MACHINE
	} else if regHive == "CU" {
		hive = reg.CURRENT_USER
	} else {
		return errors.New("invalid reg hive provided")
	}
	// Create our key or see if it exists
	k, _, err := reg.CreateKey(hive, keyPath, reg.ALL_ACCESS)
	if err != nil {
		return errors.New("create key error")
	}
	// regardless if its new or created it we set it to our value
	defer k.Close()
	// switch on keyValue type to create different key type values
	switch v := keyValue.(type) {
	case string:
		keyValueSZ := keyValue.(string)
		k.SetStringValue(keyObject, keyValueSZ)
	case uint32:
		keyValueUI32 := keyValue.(uint32)
		k.SetDWordValue(keyObject, keyValueUI32)
	default:
		vstring := fmt.Sprintf("Info: %v", v)
		return errors.New("add value error, Value "+vstring)
	}
	return nil
}

// DeleteRegKeysValue deletes a dynamic keyobject in a dynamic hive and path
func DeleteRegKeysValue(regHive string, keypath string, keyobject string) error {
	var hive reg.Key
	if regHive == "LM" {
		hive = reg.LOCAL_MACHINE
	} else if regHive == "CU" {
		hive = reg.CURRENT_USER
	} else {
		return errors.New("invalid reg hive provided")
	}
	k, err := reg.OpenKey(hive, keypath, reg.ALL_ACCESS)
	if err != nil {
		//log.Fatal(err)
		return errors.New("couldn't open key")
	}
	defer k.Close()
	k.DeleteValue(keyobject)
	return nil
}

// QueryRegKeyString reads and returns a registry key for windows
func QueryRegKeyString(regHive string, keypath string, keyobject string) (string, error) {
	var hive reg.Key
	if regHive == "LM" {
		hive = reg.LOCAL_MACHINE
	} else if regHive == "CU" {
		hive = reg.CURRENT_USER
	} else {
		return "error", errors.New("invalid reg hive provided")
	}
	k, err := reg.OpenKey(hive, keypath, reg.QUERY_VALUE)
	if err != nil {
		return "error", errors.New("invalid reg path")
	}
	defer k.Close()
	s, _, err := k.GetStringValue(keyobject)
	if err != nil {
		return "error", errors.New("invalid key to query for")
	}
	return s, nil
}