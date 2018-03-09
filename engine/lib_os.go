package engine

// func LocalSystemInfo() ([]string, error) {
// 	var InfoDump []string
// 	gi := goInfo.GetInfo()
// 	InfoDump = append(InfoDump, fmt.Sprintf("GoOS: %s", gi.GoOS))
// 	InfoDump = append(InfoDump, fmt.Sprintf("Kernel: %s", gi.Kernel))
// 	InfoDump = append(InfoDump, fmt.Sprintf("Core: %s", gi.Core))
// 	InfoDump = append(InfoDump, fmt.Sprintf("Platform: %s", gi.Platform))
// 	InfoDump = append(InfoDump, fmt.Sprintf("OS: %s", gi.OS))
// 	InfoDump = append(InfoDump, fmt.Sprintf("Hostname: %s", gi.Hostname))
// 	InfoDump = append(InfoDump, fmt.Sprintf("CPUs: %v", gi.CPUs))
// 	if InfoDump != nil {
// 		return InfoDump, nil
// 	}
// 	return nil, errors.New("Failed to retrieve local system information")
// }

// func GetHostname() (string, error) {
// 	gi := goInfo.GetInfo()
// 	hostname := gi.Hostname
// 	if hostname != "" {
// 		return hostname, nil
// 	}
// 	return "", errors.New("Failed to retrieve local hostname")
// }

// func ProcExists2(pidBoi int) bool {
// 	process, err := ps.FindProcess(pidBoi)
// 	if err == nil && process == nil {
// 		return false
// 	} else {
// 		return true
// 	}
// }

// func FindProcessPid(key string) (int, error) {
// 	pid := 0
// 	err := errors.New("Not found")
// 	ps, _ := ps.Processes()
// 	for i, _ := range ps {
// 		if ps[i].Executable() == key {
// 			pid = ps[i].Pid()
// 			err = nil
// 			break
// 		}
// 	}
// 	return pid, err
// }

// func (e *Engine) InstallSystemService(path, name, displayName, description string) error {
// 	c := &services.Config{
// 		Path:        path,
// 		Name:        name,
// 		DisplayName: displayName,
// 		Description: description,
// 	}

// 	s, err := services.NewServiceConfig(c)
// 	if err != nil {
// 		return err
// 	}

// 	err = s.Install()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (e *Engine) StartServiceByName(name string) error {
// 	c := &services.Config{
// 		Name: name,
// 	}

// 	s, err := services.NewServiceConfig(c)
// 	if err != nil {
// 		return err
// 	}

// 	err = s.Start()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (e *Engine) StopServiceByName(name string) error {
// 	c := &services.Config{
// 		Name: name,
// 	}

// 	s, err := services.NewServiceConfig(c)
// 	if err != nil {
// 		return err
// 	}

// 	err = s.Stop()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (e *Engine) RemoveServiceByName(name string) error {
// 	c := &services.Config{
// 		Name: name,
// 	}

// 	s, err := services.NewServiceConfig(c)
// 	if err != nil {
// 		return err
// 	}

// 	err = s.Remove()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (e *Engine) VMSignal(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// }

// func (e *Engine) VMLoggedInUsers(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// }

// func (e *Engine) VMUsersRunningProcs(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// }

// func (e *Engine) VMMemStats(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// }

// func (e *Engine) VMCPUStats(call otto.FunctionCall) otto.Value {
// 	cpuInfo, err := LocalSystemInfo()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("OS Error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	cpuStats := strings.Join(cpuInfo, "\n")
// 	vmResponse, err := otto.ToValue(cpuStats)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return vmResponse
// }

// func (e *Engine) VMEnvVars(call otto.FunctionCall) otto.Value {
// 	rezultz := map[string]string{}
// 	for _, v := range os.Environ() {
// 		pair := strings.Split(v, "=")
// 		rezultz[pair[0]] = pair[1]
// 	}
// 	vmResponse, err := e.VM.ToValue(rezultz)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return vmResponse
// }

// func (e *Engine) VMGetEnv(call otto.FunctionCall) otto.Value {
// 	envVar := call.Argument(0)
// 	envVarAsString, err := envVar.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	finalVal := os.Getenv(envVarAsString.(string))
// 	vmResponse, err := e.VM.ToValue(finalVal)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return vmResponse
// }

// func (e *Engine) VMAddRegKey(call otto.FunctionCall) otto.Value {
// 	regHive := call.Argument(0)
// 	keyPath := call.Argument(1)
// 	keyObject := call.Argument(2)
// 	keyValue := call.Argument(3)
// 	keyValueInterface, err := keyValue.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	regHiveAsString, err := regHive.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	keyPathAsString, err := keyPath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	keyObjectAsString, err := keyObject.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	err = CreateRegKeyAndValue(regHiveAsString.(string), keyPathAsString.(string), keyObjectAsString.(string), keyValueInterface)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Registry error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return otto.TrueValue()
// }

// func (e *Engine) VMDelRegKey(call otto.FunctionCall) otto.Value {
// 	regHive := call.Argument(0)
// 	keyPath := call.Argument(1)
// 	keyObject := call.Argument(2)
// 	regHiveAsString, err := regHive.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	keyPathAsString, err := keyPath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	keyObjectAsString, err := keyObject.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	err = DeleteRegKeysValue(regHiveAsString.(string), keyPathAsString.(string), keyObjectAsString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Registry error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return otto.TrueValue()
// }

// func (e *Engine) VMQueryRegKey(call otto.FunctionCall) otto.Value {
// 	regHive := call.Argument(0)
// 	keyPath := call.Argument(1)
// 	keyObject := call.Argument(2)
// 	regHiveAsString, err := regHive.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	keyPathAsString, err := keyPath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	keyObjectAsString, err := keyObject.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	resultStringValue, err := QueryRegKeyString(regHiveAsString.(string), keyPathAsString.(string), keyObjectAsString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Registry error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	vmResponse, err := e.VM.ToValue(resultStringValue)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return vmResponse
// }

// func (e *Engine) VMLocalUserExists(call otto.FunctionCall) otto.Value {
// 	filePathString := "/etc/passwd"
// 	search := call.Argument(0)
// 	searchString, err := search.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	fileData, err := LocalFileRead(filePathString)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	fileStrings := string(fileData)
// 	if strings.Contains(fileStrings, searchString.(string)) {
// 		return otto.TrueValue()
// 	}
// 	return otto.FalseValue()
// }

// func (e *Engine) VMProcExistsWithName(call otto.FunctionCall) otto.Value {
// 	searchProc := call.Argument(0)
// 	searchProcString, err := searchProc.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	ProcPID, err := FindProcessPid(searchProcString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("OS error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	ProcExistsResult := ProcExists2(ProcPID)
// 	if ProcExistsResult {
// 		return otto.TrueValue()
// 	}
// 	return otto.FalseValue()
// }

// func (e *Engine) VMInstallSystemService(call otto.FunctionCall) otto.Value {
// 	if len(call.ArgumentList) != 4 {
// 		e.Logger.Errorf("Not enough arguments provided.")
// 		return otto.FalseValue()
// 	}

// 	path, err := call.Argument(0).ToString()
// 	if err != nil {
// 		e.Logger.Errorf("Error converting path to string: %s", err.Error())
// 		return otto.FalseValue()
// 	}

// 	name, err := call.Argument(1).ToString()
// 	if err != nil {
// 		e.Logger.Errorf("Error converting name to string: %s", err.Error())
// 		return otto.FalseValue()
// 	}

// 	displayName, err := call.Argument(2).ToString()
// 	if err != nil {
// 		e.Logger.Errorf("Error converting displayName to string: %s", err.Error())
// 		return otto.FalseValue()
// 	}

// 	description, err := call.Argument(3).ToString()
// 	if err != nil {
// 		e.Logger.Errorf("Error converting description to string: %s", err.Error())
// 		return otto.FalseValue()
// 	}

// 	err = e.InstallSystemService(path, name, displayName, description)
// 	if err != nil {
// 		e.Logger.Errorf("Error installing system service: %s", err.Error())
// 		return otto.FalseValue()
// 	}

// 	return otto.TrueValue()
// }

// func (e *Engine) VMStartServiceByName(call otto.FunctionCall) otto.Value {
// 	if len(call.ArgumentList) != 1 {
// 		e.Logger.Errorf("Not enough arguments provided.")
// 		return otto.FalseValue()
// 	}

// 	name, err := call.Argument(0).ToString()
// 	if err != nil {
// 		e.Logger.Errorf("Error converting name to string: %s", err.Error())
// 		return otto.FalseValue()
// 	}

// 	err = e.StartServiceByName(name)
// 	if err != nil {
// 		e.Logger.Errorf("Error starting system service: %s", err.Error())
// 		return otto.FalseValue()
// 	}

// 	return otto.TrueValue()
// }

// func (e *Engine) VMStopServiceByName(call otto.FunctionCall) otto.Value {
// 	if len(call.ArgumentList) != 1 {
// 		e.Logger.Errorf("Not enough arguments provided.")
// 		return otto.FalseValue()
// 	}

// 	name, err := call.Argument(0).ToString()
// 	if err != nil {
// 		e.Logger.Errorf("Error converting name to string: %s", err.Error())
// 		return otto.FalseValue()
// 	}

// 	err = e.StopServiceByName(name)
// 	if err != nil {
// 		e.Logger.Errorf("Error starting system service: %s", err.Error())
// 		return otto.FalseValue()
// 	}

// 	return otto.TrueValue()
// }

// func (e *Engine) VMRemoveServiceByName(call otto.FunctionCall) otto.Value {
// 	if len(call.ArgumentList) != 1 {
// 		e.Logger.Errorf("Not enough arguments provided.")
// 		return otto.FalseValue()
// 	}

// 	name, err := call.Argument(0).ToString()
// 	if err != nil {
// 		e.Logger.Errorf("Error converting name to string: %s", err.Error())
// 		return otto.FalseValue()
// 	}

// 	err = e.RemoveServiceByName(name)
// 	if err != nil {
// 		e.Logger.Errorf("Error starting system service: %s", err.Error())
// 		return otto.FalseValue()
// 	}

// 	return otto.TrueValue()
// }
