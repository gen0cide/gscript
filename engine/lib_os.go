package engine

import (
	"errors"
	"os"
	"strings"

	services "github.com/gen0cide/service-go"
	"github.com/mitchellh/go-ps"
)

type RegistryRetValue struct {
	ValType        string   `json:"return_type"`
	StringVal      string   `json:"string_val"`
	StringArrayVal []string `json:"string_array_val"`
	ByteArrayVal   []byte   `json:"byte_array_val"`
	IntVal         uint32   `json:"int_val"`
	LongVal        uint64   `json:"long_val"`
}

func (e *Engine) FindProcByName(procName string) (int, error) {
	// * FIXME: this function currently matches against the name of the executible which is NOT technically the proccess name
	procs, err := ps.Processes()
	if err != nil {
		return -1, err
	}
	for _, proc := range procs {
		if procName == proc.Executable() {
			return proc.Pid(), nil
		}
	}
	return -1, errors.New("processes name not found")
}

func (e *Engine) InstallSystemService(path, name, displayName, description string) error {
	c := &services.Config{
		Path:        path,
		Name:        name,
		DisplayName: displayName,
		Description: description,
	}

	s, err := services.NewServiceConfig(c)
	if err != nil {
		return err
	}

	err = s.Install()
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) StartServiceByName(name string) error {
	c := &services.Config{
		Name: name,
	}

	s, err := services.NewServiceConfig(c)
	if err != nil {
		return err
	}

	err = s.Start()
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) StopServiceByName(name string) error {
	c := &services.Config{
		Name: name,
	}

	s, err := services.NewServiceConfig(c)
	if err != nil {
		return err
	}

	err = s.Stop()
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) RemoveServiceByName(name string) error {
	c := &services.Config{
		Name: name,
	}

	s, err := services.NewServiceConfig(c)
	if err != nil {
		return err
	}

	err = s.Remove()
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) Signal(proc int, sig int) error {
	foundProc, err := os.FindProcess(proc)
	if err != nil {
		return err
	}
	return foundProc.Signal(sig)
}

//func (e *Engine) LoggedInUsers() ([]string, error) {}

func (e *Engine) RunningProcs() []int {
	var procs []int
	for _, proc := range ps.Processes() {
		procs = append(procs, proc.Pid())
	}
	return procs
}

func (e *Engine) GetProcName(pid int) (string, error) {
	// * FIXME: this function currently returns the name of the executible which is NOT technically the proccess name
	proc, err := ps.FindProcess(pid)
	if err != nil {
		return nil, err
	}
	return proc.Executable().nil
}

//func (e *Engine) UsersRunningProcs() ([]string, error) {}

func (e *Engine) EnvVars() map[string]string {
	vars := make(map[string]string)
	for _, eVar := range os.Environ() {
		eVarSegments := strings.Split(eVar, "=")
		if len(eVarSegments) > 1 {
			vars[eVarSegments[0]] = eVarSegments[1]
		}
	}
	return vars
}

func (e *Engine) GetEnvVar(eVar string) string {
	return os.Getenv(eVar)
}

// func (e *Engine) LocalUserExists(user string) bool

// stubs for windows only funcs
func (e *Engine) AddRegKeyString(registryString string, path string, name string, value string) error {
	return errors.New("this function is unimplemented on non windows platforms")
}
func (e *Engine) AddRegKeyExpandedString(registryString string, path string, name string, value string) error {
	return errors.New("this function is unimplemented on non windows platforms")
}
func (e *Engine) AddRegKeyBinary(registryString string, path string, name string, value []byte) error {
	return errors.New("this function is unimplemented on non windows platforms")
}
func (e *Engine) AddRegKeyDWORD(registryString string, path string, name string, value uint32) error {
	return errors.New("this function is unimplemented on non windows platforms")
}
func (e *Engine) AddRegKeyQWORD(registryString string, path string, name string, value uint64) error {
	return errors.New("this function is unimplemented on non windows platforms")
}
func (e *Engine) AddRegKeyStrings(registryString string, path string, name string, value []string) error {
	return errors.New("this function is unimplemented on non windows platforms")
}
func (e *Engine) DelRegKey(registryString string, path string) error {
	return errors.New("this function is unimplemented on non windows platforms")
}
func (e *Engine) DelRegKeyValue(registryString string, path string, valueName string) error {
	return errors.New("this function is unimplemented on non windows platforms")
}
func (e *Engine) QueryRegKey(key string) (RegistryRetValue, error) {
	return errors.New("this function is unimplemented on non windows platforms")
}
