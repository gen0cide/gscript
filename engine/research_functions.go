package engine

import (
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/robertkrimen/otto"
)

func (e *Engine) VMLocalUserExists(call otto.FunctionCall) otto.Value {
	filePathString := "/etc/passwd"
	search := call.Argument(0)
	searchString, err := search.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
		return otto.FalseValue()
	}
	fileData, err := LocalFileRead(filePathString)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=LocalFileRead_Error arg=%s", err.Error())
		return otto.FalseValue()
	}
	fileStrings := string(fileData)
	if strings.Contains(fileStrings, searchString.(string)) {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMProcExistsWithName(call otto.FunctionCall) otto.Value {
	searchProc := call.Argument(0)
	searchProcString, err := searchProc.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
		return otto.FalseValue()
	}
	ProcPID, err := FindProcessPid(searchProcString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=error_finding_process arg=%s", err.Error())
		return otto.FalseValue()
	}
	ProcExistsResult := ProcExists2(ProcPID)
	if ProcExistsResult {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMCanReadFile(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
		return otto.FalseValue()
	}
	data, err := LocalFileRead(filePathString.(string))
	if data != nil && err == nil {
		e.Logger.WithField("trace", "true").Infof("Function Results: result=%x", spew.Sdump(data))
		return otto.TrueValue()
	} else if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ERR_READING_LOCAL_FILE arg=%s", err.Error())
		return otto.FalseValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMCanWriteFile(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented: %s")
	return otto.FalseValue()
	//Following code breaks building for windows :(
	/*
		filePath := call.Argument(0)
		filePathString, err := filePath.Export()
		if err != nil {
			e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
			return otto.FalseValue()
		}
		result := LocalFileWritable(filePathString.(string))
		if result == true {
			return otto.TrueValue()
		} else {
			return otto.FalseValue()
		}
	*/
}

func (e *Engine) VMCanExecFile(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented: %s")
	return otto.FalseValue()
	//Following code breaks building for windows :(
	/*
		filePath := call.Argument(0)
		filePathString, err := filePath.Export()
		if err != nil {
			e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
			return otto.FalseValue()
		}
		result := LocalFileExecutable(filePathString.(string))
		if result {
			return otto.TrueValue()
		} else {
			return otto.FalseValue()
		}
	*/
}

func (e *Engine) VMFileExists(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
		return otto.FalseValue()
	}
	if LocalFileExists(filePathString.(string)) {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMDirExists(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
		return otto.FalseValue()
	}
	if LocalFileExists(filePathString.(string)) {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMFileContains(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
		return otto.FalseValue()
	}
	search := call.Argument(1)
	searchString, err := search.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
		return otto.FalseValue()
	}
	fileData, err := LocalFileRead(filePathString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=LocalFileRead_Error arg=%s", err.Error())
		return otto.FalseValue()
	}
	fileStrings := string(fileData)
	if strings.Contains(fileStrings, searchString.(string)) {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMIsVM(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented: %s")
	return otto.FalseValue()
}

func (e *Engine) VMIsAWS(call otto.FunctionCall) otto.Value {
	respCode, response, err := HTTPGetFile("http://169.254.169.254/latest/meta-data/")
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=bad_news arg=%s", err.Error())
		return otto.FalseValue()
	} else if respCode == 200 {
		e.Logger.WithField("trace", "true").Infof("Function Results: code=%s result=%+v", respCode, string(response))
		return otto.TrueValue()
	} else {
		e.Logger.WithField("trace", "true").Infof("Function Results: code=%s result=%+v", respCode, string(response))
		return otto.FalseValue()
	}
}

func (e *Engine) VMHasPublicIP(call otto.FunctionCall) otto.Value {
	respCode, response, err := HTTPGetFile("http://icanhazip.com")
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=bad_news arg=%s", err.Error())
		return otto.FalseValue()
	} else if respCode == 200 {
		e.Logger.WithField("trace", "true").Infof("Function Results: code=%s result=%+v", respCode, string(response))
		return otto.TrueValue()
	} else {
		e.Logger.WithField("trace", "true").Infof("Function Results: code=%s result=%+v", respCode, string(response))
		return otto.FalseValue()
	}
}

func (e *Engine) VMCanMakeTCPConn(call otto.FunctionCall) otto.Value {
	ip := call.Argument(0)
	ipString, err := ip.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
		return otto.FalseValue()
	}
	port := call.Argument(1)
	portString, err := port.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
		return otto.FalseValue()
	}
	tcpResponse, err := TCPRead(ipString.(string), portString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
		return otto.FalseValue()
	}
	if tcpResponse != nil {
		e.Logger.WithField("trace", "true").Infof("Function Results: args=%s result=%+v", (ipString.(string) + ":" + portString.(string)), tcpResponse)
		return otto.TrueValue()
	} else {
		e.Logger.WithField("trace", "true").Infof("Function Results: args=%s result=%+v", (ipString.(string) + ":" + portString.(string)), tcpResponse)
		return otto.FalseValue()
	}

}

func (e *Engine) VMExpectedDNS(call otto.FunctionCall) otto.Value {
	targetDomain := call.Argument(0)
	queryType := call.Argument(1)
	expectedResult := call.Argument(2)
	targetDomainAsString, err := targetDomain.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARG_NOT_String arg=%+v", targetDomain)
		return otto.FalseValue()
	}
	queryTypeAsString, err := queryType.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARG_NOT_String arg=%+v", queryType)
		return otto.FalseValue()
	}
	expectedResultAsString, err := expectedResult.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARG_NOT_String arg=%+v", expectedResult)
		return otto.FalseValue()
	}
	result, err := DNSQuestion(targetDomainAsString.(string), queryTypeAsString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=DNSLookupFailed details=%s args=%s query_type=%s", err.Error(), targetDomainAsString, queryTypeAsString.(string))
		return otto.FalseValue()
	}
	if expectedResultAsString.(string) == string(result) {
		e.Logger.WithField("trace", "true").Infof("Function: msg='DNS Results: %s'", string(result))
		return otto.TrueValue()
	} else {
		e.Logger.WithField("trace", "true").Infof("Function: msg='DNS Results: %s'", string(result))
		return otto.FalseValue()
	}
}

func (e *Engine) VMCanMakeHTTPConn(call otto.FunctionCall) otto.Value {
	url1 := call.Argument(0)
	url1String, err := url1.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%+v", url1String)
		return otto.FalseValue()
	}
	respCode, _, err := HTTPGetFile(url1String.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARG_NOT_String arg=%s", err.Error())
		return otto.FalseValue()
	} else if respCode != 403 || respCode != 404 || respCode != 500 || respCode != 502 || respCode != 503 || respCode != 504 || respCode != 511 {
		e.Logger.WithField("trace", "true").Infof("Function Results: args=%+v result=%+v", url1String, respCode)
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMDetectSSLMITM(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented: %s")
	return otto.FalseValue()
}

func (e *Engine) VMCmdSuccessful(call otto.FunctionCall) otto.Value {
	cmd := call.Argument(0)
	cmdString, err := cmd.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%+v", cmdString)
		return otto.FalseValue()
	}
	arg := call.Argument(1)
	argString, err := arg.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%+v", argString)
		return otto.FalseValue()
	}
	VMExecResponse := ExecuteCommand(cmdString.(string), argString.(string))
	if VMExecResponse.Success == false {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=%s", VMExecResponse.ErrorMsg)
		return otto.FalseValue()
	} else if VMExecResponse.Success == true {
		e.Logger.WithField("trace", "true").Infof("Function Results: args=%s results=%s", cmdString, VMExecResponse.Stdout)
		return otto.TrueValue()
	}
	return otto.FalseValue()
}

func (e *Engine) VMCanPing(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented: %s")
	return otto.FalseValue()
}

func (e *Engine) VMTCPPortInUse(call otto.FunctionCall) otto.Value {
	var minTCPPort int64 = 0
	var maxTCPPort int64 = 65535
	port := call.Argument(0)
	portInt, err := port.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_Int arg=%+v", portInt)
		return otto.FalseValue()
	}
	if portInt.(int64) < minTCPPort || portInt.(int64) > maxTCPPort {
		return otto.FalseValue()
	}
	conn, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", portInt.(int64)))
	if err != nil {
		return otto.TrueValue()
	}
	conn.Close()
	return otto.FalseValue()
}

func (e *Engine) VMUDPPortInUse(call otto.FunctionCall) otto.Value {
	var minUDPPort int64 = 0
	var maxUDPPort int64 = 65535
	port := call.Argument(0)
	portInt, err := port.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_Int arg=%+v", portInt)
		return otto.FalseValue()
	}
	if portInt.(int64) < minUDPPort || portInt.(int64) > maxUDPPort {
		return otto.FalseValue()
	}
	conn, err := net.Listen("udp", fmt.Sprintf("127.0.0.1:%d", portInt.(int64)))
	if err != nil {
		return otto.TrueValue()
	}
	conn.Close()
	return otto.FalseValue()
}

func (e *Engine) VMExistsInPath(call otto.FunctionCall) otto.Value {
	cmd := call.Argument(0)
	cmdString, err := cmd.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%+v", cmdString)
		return otto.FalseValue()
	}
	path, err := exec.LookPath(cmdString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=PathLookupFailed arg=%s", err.Error())
		return otto.FalseValue()
	} else if path != "" {
		e.Logger.WithField("trace", "true").Infof("Function Results: results=%s", path)
		return otto.TrueValue()
	} else {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=PathLookupFailed arg=%+v", cmdString)
		return otto.FalseValue()
	}
}

func (e *Engine) VMCanSudo(call otto.FunctionCall) otto.Value {
	VMExecResponse := ExecuteCommand("sudo", "-v")
	if VMExecResponse.Success == false {
		e.Logger.WithField("trace", "true").Errorf("Function Error: error=%s", VMExecResponse.ErrorMsg)
		return otto.FalseValue()
	} else if VMExecResponse.Success == true {
		e.Logger.WithField("trace", "true").Infof("Function Results: results=%s", VMExecResponse.Stdout)
		return otto.TrueValue()
	}
	return otto.FalseValue()
}

func (e *Engine) VMMatches(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented: %s", "VMMatches")
	return otto.FalseValue()
}

func (e *Engine) VMCanSSHLogin(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented: %s", "VMCanSSHLogin")
	return otto.FalseValue()
}
