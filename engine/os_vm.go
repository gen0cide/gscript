package engine

import (
	"os"
	"strings"

	"github.com/robertkrimen/otto"
)

func (e *Engine) VMSignal(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMLoggedInUsers(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMUsersRunningProcs(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMMemStats(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMCPUStats(call otto.FunctionCall) otto.Value {
	cpuInfo, err := LocalSystemInfo()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("OS Error: %s", err.Error())
		return otto.FalseValue()
	}
	cpuStats := strings.Join(cpuInfo, "\n")
	vmResponse, err := otto.ToValue(cpuStats)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMEnvVars(call otto.FunctionCall) otto.Value {
	rezultz := map[string]string{}
	for _, v := range os.Environ() {
		pair := strings.Split(v, "=")
		rezultz[pair[0]] = pair[1]
	}
	vmResponse, err := e.VM.ToValue(rezultz)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMGetEnv(call otto.FunctionCall) otto.Value {
	envVar := call.Argument(0)
	envVarAsString, err := envVar.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	finalVal := os.Getenv(envVarAsString.(string))
	vmResponse, err := e.VM.ToValue(finalVal)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMAddRegKey(call otto.FunctionCall) otto.Value {
	regHive := call.Argument(0)
	keyPath := call.Argument(1)
	keyObject := call.Argument(2)
	keyValue := call.Argument(3)
	keyValueInterface, err := keyValue.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	regHiveAsString, err := regHive.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	keyPathAsString, err := keyPath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	keyObjectAsString, err := keyObject.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	err = CreateRegKeyAndValue(regHiveAsString.(string), keyPathAsString.(string), keyObjectAsString.(string), keyValueInterface)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Registry error: %s", err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMDelRegKey(call otto.FunctionCall) otto.Value {
	regHive := call.Argument(0)
	keyPath := call.Argument(1)
	keyObject := call.Argument(2)
	regHiveAsString, err := regHive.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	keyPathAsString, err := keyPath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	keyObjectAsString, err := keyObject.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	err = DeleteRegKeysValue(regHiveAsString.(string), keyPathAsString.(string), keyObjectAsString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Registry error: %s", err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMQueryRegKey(call otto.FunctionCall) otto.Value {
	regHive := call.Argument(0)
	keyPath := call.Argument(1)
	keyObject := call.Argument(2)
	regHiveAsString, err := regHive.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	keyPathAsString, err := keyPath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	keyObjectAsString, err := keyObject.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	resultStringValue, err := QueryRegKeyString(regHiveAsString.(string), keyPathAsString.(string), keyObjectAsString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Registry error: %s", err.Error())
		return otto.FalseValue()
	}
	vmResponse, err := e.VM.ToValue(resultStringValue)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMLocalUserExists(call otto.FunctionCall) otto.Value {
	filePathString := "/etc/passwd"
	search := call.Argument(0)
	searchString, err := search.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	fileData, err := LocalFileRead(filePathString)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return otto.FalseValue()
	}
	fileStrings := string(fileData)
	if strings.Contains(fileStrings, searchString.(string)) {
		return otto.TrueValue()
	}
	return otto.FalseValue()
}

func (e *Engine) VMProcExistsWithName(call otto.FunctionCall) otto.Value {
	searchProc := call.Argument(0)
	searchProcString, err := searchProc.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	ProcPID, err := FindProcessPid(searchProcString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("OS error: %s", err.Error())
		return otto.FalseValue()
	}
	ProcExistsResult := ProcExists2(ProcPID)
	if ProcExistsResult {
		return otto.TrueValue()
	}
	return otto.FalseValue()
}
