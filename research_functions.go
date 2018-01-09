package gscript

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
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	}
	fileData, err := LocalFileRead(filePathString)
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=LocalFileRead_Error arg=%s", CalledBy(), spew.Sdump(err))
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
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	}
	ProcPID, err := FindProcessPid(searchProcString.(string))
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=error_finding_process arg=%s", CalledBy(), spew.Sdump(err))
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
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	}
	data, err := LocalFileRead(filePathString.(string))
	if data != nil && err == nil {
		e.LogInfof("Function Results: function=%s result=%s", CalledBy(), spew.Sdump(data))
		return otto.TrueValue()
	} else if err != nil {
		e.LogErrorf("Function Error: function=%s error=ERR_READING_LOCAL_FILE arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMCanWriteFile(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
	//Following code breaks building for windows :(
	/*
		filePath := call.Argument(0)
		filePathString, err := filePath.Export()
		if err != nil {
			e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
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
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
	//Following code breaks building for windows :(
	/*
		filePath := call.Argument(0)
		filePathString, err := filePath.Export()
		if err != nil {
			e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
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
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
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
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
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
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	}
	search := call.Argument(1)
	searchString, err := search.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	}
	fileData, err := LocalFileRead(filePathString.(string))
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=LocalFileRead_Error arg=%s", CalledBy(), spew.Sdump(err))
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
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMIsAWS(call otto.FunctionCall) otto.Value {
	respCode, response, err := HTTPGetFile("http://169.254.169.254/latest/meta-data/")
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=bad_news arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	} else if respCode == 200 {
		e.LogInfof("Function Results: function=%s code=%s result=%s", CalledBy(), respCode, spew.Sdump(response))
		return otto.TrueValue()
	} else {
		e.LogInfof("Function Results: function=%s code=%s result=%s", CalledBy(), respCode, spew.Sdump(response))
		return otto.FalseValue()
	}
}

func (e *Engine) VMHasPublicIP(call otto.FunctionCall) otto.Value {
	respCode, response, err := HTTPGetFile("http://icanhazip.com")
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=bad_news arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	} else if respCode == 200 {
		e.LogInfof("Function Results: function=%s code=%s result=%s", CalledBy(), respCode, spew.Sdump(response))
		return otto.TrueValue()
	} else {
		e.LogInfof("Function Results: function=%s code=%s result=%s", CalledBy(), respCode, spew.Sdump(response))
		return otto.FalseValue()
	}
}

func (e *Engine) VMCanMakeTCPConn(call otto.FunctionCall) otto.Value {
	ip := call.Argument(0)
	ipString, err := ip.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	}
	port := call.Argument(1)
	portString, err := port.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	}
	tcpResponse, err := TCPRead(ipString.(string), portString.(string))
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	}
	if tcpResponse != nil {
		e.LogInfof("Function Results: function=%s args=%s result=%s", CalledBy(), (ipString.(string) + ":" + portString.(string)), spew.Sdump(tcpResponse))
		return otto.TrueValue()
	} else {
		e.LogInfof("Function Results: function=%s args=%s result=%s", CalledBy(), (ipString.(string) + ":" + portString.(string)), spew.Sdump(tcpResponse))
		return otto.FalseValue()
	}

}

func (e *Engine) VMExpectedDNS(call otto.FunctionCall) otto.Value {
	targetDomain := call.Argument(0)
	queryType := call.Argument(1)
	expectedResult := call.Argument(2)
	targetDomainAsString, err := targetDomain.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(targetDomain))
		return otto.FalseValue()
	}
	queryTypeAsString, err := queryType.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(queryType))
		return otto.FalseValue()
	}
	expectedResultAsString, err := expectedResult.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(expectedResult))
		return otto.FalseValue()
	}
	result, err := DNSQuestion(targetDomainAsString.(string), queryTypeAsString.(string))
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=DNSLookupFailed details=%s args=%s %s", CalledBy(), spew.Sdump(err), spew.Sdump(targetDomainAsString.(string)), queryTypeAsString.(string))
		return otto.FalseValue()
	}
	if expectedResultAsString.(string) == string(result) {
		e.LogInfof("Function: function=%s msg='DNS Results: %s'", CalledBy(), spew.Sdump(string(result)))
		return otto.TrueValue()
	} else {
		e.LogInfof("Function: function=%s msg='DNS Results: %s'", CalledBy(), spew.Sdump(string(result)))
		return otto.FalseValue()
	}
}

func (e *Engine) VMCanMakeHTTPConn(call otto.FunctionCall) otto.Value {
	url1 := call.Argument(0)
	url1String, err := url1.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(url1String))
		return otto.FalseValue()
	}
	respCode, _, err := HTTPGetFile(url1String.(string))
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	} else if respCode != 403 || respCode != 404 || respCode != 500 || respCode != 502 || respCode != 503 || respCode != 504 || respCode != 511 {
		e.LogInfof("Function Results: function=%s args=%s result=%s", CalledBy(), url1String.(string), spew.Sdump(respCode))
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMDetectSSLMITM(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMCmdSuccessful(call otto.FunctionCall) otto.Value {
	cmd := call.Argument(0)
	cmdString, err := cmd.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(cmdString))
		return otto.FalseValue()
	}
	arg := call.Argument(1)
	argString, err := arg.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(argString))
		return otto.FalseValue()
	}
	VMExecResponse := ExecuteCommand(cmdString.(string), argString.(string))
	if VMExecResponse.Success == false {
		e.LogErrorf("Function Error: function=%s error=%s", CalledBy(), spew.Sdump(VMExecResponse.ErrorMsg))
		return otto.FalseValue()
	} else if VMExecResponse.Success == true {
		e.LogInfof("Function Results: function=%s args=%s results=%s", CalledBy(), cmdString.(string), VMExecResponse.Stdout)
		return otto.TrueValue()
	}
	return otto.FalseValue()
}

func (e *Engine) VMCanPing(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMTCPPortInUse(call otto.FunctionCall) otto.Value {
	var minTCPPort int64 = 0
	var maxTCPPort int64 = 65535
	port := call.Argument(0)
	portInt, err := port.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_Int arg=%s", CalledBy(), spew.Sdump(portInt))
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
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_Int arg=%s", CalledBy(), spew.Sdump(portInt))
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
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(cmdString))
		return otto.FalseValue()
	}
	path, err := exec.LookPath(cmdString.(string))
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=PathLookupFailed arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	} else if path != "" {
		e.LogInfof("Function Results: function=%s results=%s", CalledBy(), path)
		return otto.TrueValue()
	} else {
		e.LogErrorf("Function Error: function=%s error=PathLookupFailed arg=%s", CalledBy(), spew.Sdump(cmdString))
		return otto.FalseValue()
	}
}

func (e *Engine) VMCanSudo(call otto.FunctionCall) otto.Value {
	VMExecResponse := ExecuteCommand("sudo", "-v")
	if VMExecResponse.Success == false {
		e.LogErrorf("Function Error: function=%s error=%s", CalledBy(), spew.Sdump(VMExecResponse.ErrorMsg))
		return otto.FalseValue()
	} else if VMExecResponse.Success == true {
		e.LogInfof("Function Results: function=%s results=%s", CalledBy(), VMExecResponse.Stdout)
		return otto.TrueValue()
	}
	return otto.FalseValue()
}

func (e *Engine) VMMatches(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMCanSSHLogin(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}
