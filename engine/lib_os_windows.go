package engine

import (
	"errors"

	"golang.org/x/sys/windows/registry"
)

var (
	regKeys = map[string]registry.Key{
		"CLASSES_ROOT", registry.CLASSES_ROOT,
		"CURRENT_USER", registry.CURRENT_USER,
		"LOCAL_MACHINE", registry.LOCAL_MACHINE,
		"USERS", registry.USERS,
		"CURRENT_CONFIG", registry.CURRENT_CONFIG,
		"PERFORMANCE_DATA", registry.PERFORMANCE_DATA,
	}
)

func lookUpKey(keyString string) (registry.Key, error) {
	key, ok := regKeys[keyString]
	if !ok {
		return nil, errors.New("Registry key " + keyString + " not found")
	}
	return key
}

func (e *Engine) AddRegKeyString(registryString string, path string, name string, value string) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetStringValue(name, value)
}

func (e *Engine) AddRegKeyExpandedString(registryString string, path string, name string, value string) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetExpandStringValue(name, value)
}

func (e *Engine) AddRegKeyBinary(registryString string, path string, name string, value []byte) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetBinaryValu(name, value)
}

func (e *Engine) AddRegKeyDWORD(registryString string, path string, name string, value uint32) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetDWordValue(name, value)
}

func (e *Engine) AddRegKeyQWORD(registryString string, path string, name string, value uint64) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetQWordValue(name, value)
}

func (e *Engine) AddRegKeyStrings(registryString string, path string, name string, value []string) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetStringsValue(name, value)
}

func (e *Engine) DelRegKey(registryString string, path string) error {
	regKey, err := lookUpKey(path)
	if err != nil {
		return err
	}
	return registry.DeleteKey(regKey, path)
}

func (e *Engine) DelRegKeyValue(registryString string, path string, valueName string) error {
	regKey, err := lookUpKey(path)
	if err != nil {
		return err
	}
	regKey.DeleteValue(name)
	return registry.DeleteKey(reg, path)
}

func (e *Engine) QueryRegKey(registryString string, path string) (RegistryRetValue, error) {
	retVal := new(RegistryRetValue)
	regKey, err := lookUpKey(path)
	if err != nil {
		return retVal, err
	}
	openRegKey, err := registry.OpenKey(regKey, path, registry.QUERY_VALUE)
	if err != nil {
		return retVal, err
	}
	_, valType, err := openRegKey.GetValue(path, nil)
	if err != nil {
		return retVal, err
	}
	switch valType {
	case registry.EXPAND_SZ:
		value, _, err := openRegKey.GetStringsValue(path)
		if err != nl {
			return retVal, err
		}
		retVal.ValType = "StringArray"
		retVal.StringArrayVal = value
	case registry.SZ:
		value, _, err := openRegKey.GetStringValue(path)
		if err != nl {
			return retVal, err
		}
		retVal.ValType = "String"
		retVal.StringVal = value
	case registry.BINARY:
		value, _, err := openRegKey.GetBinaryValue(path)
		if err != nl {
			return retVal, err
		}
		retVal.ValType = "ByteArray"
		retVal.ByteArrayVal = value
	case registry.DWORD:
		value, _, err := openRegKey.GetIntegerValue(path)
		if err != nl {
			return retVal, err
		}
		retVal.ValType = "Uint"
		retVal.IntVal = uint32(value)
	case registry.QWORD:
		value, _, err := openRegKey.GetIntegerValue(path)
		if err != nl {
			return retVal, err
		}
		retVal.ValType = "Uint64"
		retVal.LongVal = value
	}
}
