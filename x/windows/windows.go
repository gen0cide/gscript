// +build windows

package windows

import (
	"errors"

	"golang.org/x/sys/windows/registry"
)

var (
	regKeys = map[string]registry.Key{
		"CLASSES_ROOT":     registry.CLASSES_ROOT,
		"CURRENT_USER":     registry.CURRENT_USER,
		"LOCAL_MACHINE":    registry.LOCAL_MACHINE,
		"USERS":            registry.USERS,
		"CURRENT_CONFIG":   registry.CURRENT_CONFIG,
		"PERFORMANCE_DATA": registry.PERFORMANCE_DATA,
	}
)

type RegistryRetValue struct {
	ValType        string   `json:"return_type"`
	StringVal      string   `json:"string_val"`
	StringArrayVal []string `json:"string_array_val"`
	ByteArrayVal   []byte   `json:"byte_array_val"`
	IntVal         uint32   `json:"int_val"`
	LongVal        uint64   `json:"long_val"`
}

func lookUpKey(keyString string) (registry.Key, error) {
	key, ok := regKeys[keyString]
	if !ok {
		// lol, picking a key at random because fuck golang return types
		return registry.CLASSES_ROOT, errors.New("Registry key " + keyString + " not found")
	}
	return key, nil
}

func AddRegKeyString(registryString string, path string, name string, value string) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, _, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetStringValue(name, value)
}

func AddRegKeyExpandedString(registryString string, path string, name string, value string) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, _, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetExpandStringValue(name, value)
}

func AddRegKeyBinary(registryString string, path string, name string, value []byte) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, _, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetBinaryValue(name, value)
}

func AddRegKeyDWORD(registryString string, path string, name string, value uint32) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, _, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetDWordValue(name, value)
}

func AddRegKeyQWORD(registryString string, path string, name string, value uint64) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, _, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetQWordValue(name, value)
}

func AddRegKeyStrings(registryString string, path string, name string, value []string) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, _, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetStringsValue(name, value)
}

func DelRegKey(registryString string, path string) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	return registry.DeleteKey(regKey, path)
}

func DelRegKeyValue(registryString string, path string, valueName string) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	regKey.DeleteValue(path)
	return registry.DeleteKey(regKey, path)
}

func QueryRegKey(registryString string, path string, key string) (RegistryRetValue, error) {
	retVal := RegistryRetValue{}
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return retVal, err
	}
	openRegKey, err := registry.OpenKey(regKey, path, registry.QUERY_VALUE)
	if err != nil {
		return retVal, err
	}
	_, valType, err := openRegKey.GetValue(key, nil)
	if err != nil {
		return retVal, err
	}
	switch valType {
	case registry.EXPAND_SZ:
		value, _, err := openRegKey.GetStringsValue(key)
		if err != nil {
			return retVal, err
		}
		retVal.ValType = "StringArray"
		retVal.StringArrayVal = value
	case registry.SZ:
		value, _, err := openRegKey.GetStringValue(key)
		if err != nil {
			return retVal, err
		}
		retVal.ValType = "String"
		retVal.StringVal = value
	case registry.BINARY:
		value, _, err := openRegKey.GetBinaryValue(key)
		if err != nil {
			return retVal, err
		}
		retVal.ValType = "ByteArray"
		retVal.ByteArrayVal = value
	case registry.DWORD:
		value, _, err := openRegKey.GetIntegerValue(key)
		if err != nil {
			return retVal, err
		}
		retVal.ValType = "Uint"
		retVal.IntVal = uint32(value)
	case registry.QWORD:
		value, _, err := openRegKey.GetIntegerValue(key)
		if err != nil {
			return retVal, err
		}
		retVal.ValType = "Uint64"
		retVal.LongVal = value
	}
	return retVal, nil
}
