package gscript

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/djherbis/times"
	"github.com/robertkrimen/otto"
)

func (e *Engine) VMHalt(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMDeleteFile(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(filePath))
		return otto.FalseValue()
	}
	err = LocalFileDelete(filePathAsString.(string))
	if err != nil {
		e.LogErrorf("Error deleting the file: function=%s path=%s error=%s", CalledBy(), filePathAsString.(string), err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMWriteFile(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	fileData := call.Argument(1)
	fileBytes := e.ValueToByteSlice(fileData)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(filePath))
		return otto.FalseValue()
	}
	err = LocalFileCreate(filePathAsString.(string), fileBytes)
	if err != nil {
		e.LogErrorf("Error writing the file: function=%s path=%s error=%s", CalledBy(), filePathAsString.(string), err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMReadFile(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(filePath))
		return otto.FalseValue()
	}
	bytes, err := LocalFileRead(filePathAsString.(string))
	if err != nil {
		e.LogErrorf("Error reading the file: function=%s path=%s error=%s", CalledBy(), filePathAsString.(string), err.Error())
		return otto.FalseValue()
	}

	vmResponse, err := e.VM.ToValue(string(bytes))

  return vmResponse;
}

func (e *Engine) VMCopyFile(call otto.FunctionCall) otto.Value {
	readPath, err := call.ArgumentList[0].ToString()
	if err != nil {
		e.LogErrorf("Function Error: function=VMCopyFile() error=ARG_NOT_STRINGABLE arg=%s", spew.Sdump(call.ArgumentList[0]))
		return otto.FalseValue()
	}
	writePath, err := call.ArgumentList[1].ToString()
	if err != nil {
		e.LogErrorf("Function Error: function=VMCopyFile() error=ARG_NOT_STRINGABLE arg=%s", spew.Sdump(call.ArgumentList[1]))
		return otto.FalseValue()
	}
	bytes, err := LocalFileRead(readPath)
	if err != nil {
		e.LogErrorf("Function Error: function=VMCopyFile() error='There was an error reading the file to copy'")
		return otto.FalseValue()
	}
	err = LocalFileCreate(writePath, bytes)
	if err != nil {
		e.LogErrorf("Function Error: function=VMCopyFile() error='There was an error writing the file to that path'")
		return otto.FalseValue()
	}
	e.LogInfof("Function: function=%s msg='wrote local file at: %s'", CalledBy(), spew.Sdump(writePath))
	return otto.TrueValue()
}

func (e *Engine) VMExecuteFile(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.ToString()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=CMD_BASE_NOT_PARSABLE arg=%s", CalledBy(), spew.Sdump(filePath))
		return otto.FalseValue()
	}
	cmdArgs := call.Argument(1)
	argList := []string{}
	if !cmdArgs.IsNull() {
		argArray, err := cmdArgs.Export()
		if err != nil {
			e.LogErrorf("Function Error: function=%s error=CMD_ARGS_NOT_PARSABLE arg=%s", CalledBy(), spew.Sdump(cmdArgs))
			return otto.FalseValue()
		}
		argList = argArray.([]string)
	}

	cmdOutput := ExecuteCommand(filePathAsString, argList...)
	vmResponse, err := e.VM.ToValue(cmdOutput)
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=CMD_OUTPUT_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMAsset(call otto.FunctionCall) otto.Value {
	filename := call.Argument(0).String()
	if len(filename) == 0 {
		e.LogErrorf("Empty Asset Call")
		return otto.FalseValue()
	}
	if dataFunc, ok := e.Imports[filename]; ok {
		byteData := dataFunc()
		vmVal, err := e.VM.ToValue(byteData)
		if err != nil {
			e.LogErrorf("Function Error: function=%s error=VM_BYTE_CONVERSION_FAILED details=%s", CalledBy(), spew.Sdump(err))
			return otto.FalseValue()
		}
		return vmVal
	}
	e.LogErrorf("Asset File Not Found: asset=%s", filename)
	e.LogErrorf("All Assets: %s", e.Imports)
	return otto.FalseValue()
}

func (e *Engine) VMAppendFile(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	fileData := call.Argument(1)
	fileBytes := e.ValueToByteSlice(fileData)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(filePath))
		return otto.FalseValue()
	}
	err = LocalFileAppendBytes(filePathAsString.(string), fileBytes)
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=LocalFileAppendBytesFailed details=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMReplaceInFile(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(filePath))
		return otto.FalseValue()
	}
	sFind := call.Argument(1)
	sFindAsString, err := sFind.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(sFind))
		return otto.FalseValue()
	}
	sReplace := call.Argument(2)
	sReplaceAsString, err := sReplace.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(sReplace))
		return otto.FalseValue()
	}
	err = LocalFileReplace(filePathAsString.(string), sFindAsString.(string), sReplaceAsString.(string))
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=Failed to run LocalFileReplace info=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMSignal(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMImplode(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMRetrieveFileFromURL(call otto.FunctionCall) otto.Value {
	readURL, err := call.ArgumentList[0].ToString()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARG_NOT_STRINGABLE arg=%s", CalledBy(), spew.Sdump(call.ArgumentList[0]))
		return otto.FalseValue()
	}
	_, bytes, err := HTTPGetFile(readURL)
	if err != nil {
		e.LogErrorf("Function Error: function=VMRetrieveFileFromURL() error='There was an error fetching the file: %v'", spew.Sdump(err))
		return otto.FalseValue()
	}
	vmResponse, err := e.VM.ToValue(bytes)
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=URL_DOWNLOAD_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMDNSQuery(call otto.FunctionCall) otto.Value {
	targetDomain := call.Argument(0)
	queryType := call.Argument(1)
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
	result, err := DNSQuestion(targetDomainAsString.(string), queryTypeAsString.(string))
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=DNSLookupFailed details=%s args=%s %s", CalledBy(), spew.Sdump(err), spew.Sdump(targetDomainAsString.(string)), queryTypeAsString.(string))
		return otto.FalseValue()
	}
	e.LogInfof("Function: function=%s msg='DNS Results: %s'", CalledBy(), spew.Sdump(string(result)))
	return otto.TrueValue()
}

func (e *Engine) VMHTTPRequest(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMExec(call otto.FunctionCall) otto.Value {
	baseCmd := call.Argument(0)
	if !baseCmd.IsString() {
		e.LogErrorf("Function Error: function=%s error=CMD_CALL_NOT_OF_TYPE_STRING arg=%s", CalledBy(), spew.Sdump(baseCmd))
		return otto.FalseValue()
	}
	cmdArgs := call.Argument(1)
	argList := []string{}
	if !cmdArgs.IsNull() {
		argArray, err := cmdArgs.Export()
		if err != nil {
			e.LogErrorf("Function Error: function=%s error=CMD_ARGS_NOT_PARSABLE arg=%s", CalledBy(), spew.Sdump(cmdArgs))
			return otto.FalseValue()
		}
		argList = argArray.([]string)
	}
	baseCmdAsString, err := baseCmd.ToString()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=CMD_BASE_NOT_PARSABLE arg=%s", CalledBy(), spew.Sdump(baseCmd))
		return otto.FalseValue()
	}
	cmdOutput := ExecuteCommand(baseCmdAsString, argList...)
	vmResponse, err := e.VM.ToValue(cmdOutput)
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=CMD_OUTPUT_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(baseCmd))
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMMD5(call otto.FunctionCall) otto.Value {
	var HashVal string
	if len(call.ArgumentList) > 0 {
		arg := call.ArgumentList[0]
		switch arg.Class() {
		case "Array":
			rawBytes := []byte{}
			expVal, err := arg.Export()
			if err != nil {
				e.LogErrorf("Function Error: function=MD5() error=ARY_ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
				return otto.Value{}
			}
			for _, i := range expVal.([]interface{}) {
				rawBytes = append(rawBytes, i.(byte))
			}
			hasher := md5.New()
			hasher.Write(rawBytes)
			HashVal = hex.EncodeToString(hasher.Sum(nil))
		default:
			val, err := call.ArgumentList[0].ToString()
			if err != nil {
				e.LogErrorf("Function Error: function=MD5() error=ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
				return otto.Value{}
			}
			hasher := md5.New()
			hasher.Write([]byte(val))
			HashVal = hex.EncodeToString(hasher.Sum(nil))
		}
	}
	var ret = otto.Value{}
	var err error
	if len(HashVal) > 0 {
		ret, err = otto.ToValue(HashVal)
	} else {
		ret, err = otto.ToValue(nil)
	}
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=%s", CalledBy(), err.Error())
	}
	return ret
}

func (e *Engine) VMSHA1(call otto.FunctionCall) otto.Value {
	var HashVal string
	if len(call.ArgumentList) > 0 {
		arg := call.ArgumentList[0]
		switch arg.Class() {
		case "Array":
			rawBytes := []byte{}
			expVal, err := arg.Export()
			if err != nil {
				e.LogErrorf("Function Error: function=SHA1() error=ARY_ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
				return otto.Value{}
			}
			for _, i := range expVal.([]interface{}) {
				rawBytes = append(rawBytes, i.(byte))
			}
			hasher := sha1.New()
			hasher.Write(rawBytes)
			HashVal = hex.EncodeToString(hasher.Sum(nil))
		default:
			val, err := call.ArgumentList[0].ToString()
			if err != nil {
				e.LogErrorf("Function Error: function=SHA1() error=ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
				return otto.Value{}
			}
			hasher := sha1.New()
			hasher.Write([]byte(val))
			HashVal = hex.EncodeToString(hasher.Sum(nil))
		}
	}
	var ret = otto.Value{}
	var err error
	if len(HashVal) > 0 {
		ret, err = otto.ToValue(HashVal)
	} else {
		ret, err = otto.ToValue(nil)
	}
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=%s", CalledBy(), err.Error())
	}
	return ret
}

func (e *Engine) VMB64Decode(call otto.FunctionCall) otto.Value {
	var NewVal string
	if len(call.ArgumentList) > 0 {
		arg := call.ArgumentList[0]
		val, err := arg.ToString()
		if err != nil {
			e.LogErrorf("Function Error: function=B64Decode() error=ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
			return otto.Value{}
		}
		valBytes, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			e.LogErrorf("Function Error: function=B64Decode() error=BASE64_DECODE_ERROR error=%s", err.Error())
			return otto.Value{}
		}
		NewVal = string(valBytes)
	}
	var ret = otto.Value{}
	var err error
	if len(NewVal) > 0 {
		ret, err = otto.ToValue(NewVal)
	} else {
		ret, err = otto.ToValue(nil)
	}
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=%s", CalledBy(), err.Error())
	}
	return ret
}

func (e *Engine) VMB64Encode(call otto.FunctionCall) otto.Value {
	var EncVal string
	if len(call.ArgumentList) > 0 {
		arg := call.ArgumentList[0]
		switch arg.Class() {
		case "Array":
			rawBytes := []byte{}
			expVal, err := arg.Export()
			if err != nil {
				e.LogErrorf("Function Error: function=B64Encode() error=ARY_ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
				return otto.Value{}
			}
			for _, i := range expVal.([]interface{}) {
				rawBytes = append(rawBytes, i.(byte))
			}
			EncVal = base64.StdEncoding.EncodeToString([]byte(rawBytes))
		default:
			val, err := call.ArgumentList[0].ToString()
			if err != nil {
				e.LogErrorf("Function Error: function=B64Encode() error=ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
				return otto.Value{}
			}
			EncVal = base64.StdEncoding.EncodeToString([]byte(val))
		}
	}
	var ret = otto.Value{}
	var err error
	if len(EncVal) > 0 {
		ret, err = otto.ToValue(EncVal)
	} else {
		ret, err = otto.ToValue(nil)
	}
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=%s", CalledBy(), err.Error())
	}
	return ret
}

func (e *Engine) VMTimestamp(call otto.FunctionCall) otto.Value {
	ts := time.Now().Unix()
	ret, err := otto.ToValue(ts)
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=%s", CalledBy(), err.Error())
	}
	return ret
}

func (e *Engine) VMCPUStats(call otto.FunctionCall) otto.Value {
	CPUInfo, err := LocalSystemInfo()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=%s", CalledBy(), err.Error())
	} else {
		e.LogInfof("Function: function=%s msg='CPU Stats: %s'", CalledBy(), spew.Sdump(CPUInfo))
	}
	return otto.TrueValue()
}

func (e *Engine) VMMemStats(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMSSHCmd(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMSleep(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 0 {
		arg := call.ArgumentList[0]
		if arg.IsNumber() {
			intArg, err := arg.ToInteger()
			if err != nil {
				e.LogErrorf("Function Error: function=%s error=%s", CalledBy(), err.Error())
				return otto.Value{}
			}
			time.Sleep(time.Duration(intArg) * time.Second)
		}
	}
	return otto.Value{}
}

func (e *Engine) VMGetTweet(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMGetDirsInPath(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMEnvVars(call otto.FunctionCall) otto.Value {
	rezultz := map[string]string{}
	for _, v := range os.Environ() {
		pair := strings.Split(v, "=")
		rezultz[pair[0]] = pair[1]
		//e.LogInfof("Function: function=%s msg='wrote local file at: %s'", CalledBy(), spew.Sdump(pair))
	}
	vmResponse, err := e.VM.ToValue(rezultz)
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=CMD_OUTPUT_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(rezultz))
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMGetEnv(call otto.FunctionCall) otto.Value {
	envVar := call.Argument(0)
	envVarAsString, err := envVar.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(envVar))
		return otto.FalseValue()
	}
	finalVal := os.Getenv(envVarAsString.(string))
	vmResponse, err := e.VM.ToValue(finalVal)
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=CMD_OUTPUT_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(finalVal))
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMFileChangeTime(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(filePath))
		return otto.FalseValue()
	}
	t, err := times.Stat(filePathAsString.(string))
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=CMD_OUTPUT_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(filePathAsString))
		return otto.FalseValue()
	}
	if t.HasChangeTime() {
		vmResponse, err := e.VM.ToValue(t.ChangeTime().String())
		if err != nil {
			e.LogErrorf("Function Error: function=%s error=CMD_OUTPUT_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(t))
			return otto.FalseValue()
		}
		return vmResponse
	} else {
		e.LogErrorf("Function Error: function=%s error=FILE_HAS_NO_CTIME arg=%s", CalledBy(), spew.Sdump(t))
		return otto.FalseValue()
	}
}

func (e *Engine) VMFileBirthTime(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(filePath))
		return otto.FalseValue()
	}
	t, err := times.Stat(filePathAsString.(string))
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=CMD_OUTPUT_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(filePathAsString))
		return otto.FalseValue()
	}
	if t.HasBirthTime() {
		vmResponse, err := e.VM.ToValue(t.BirthTime().String())
		if err != nil {
			e.LogErrorf("Function Error: function=%s error=CMD_OUTPUT_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(t))
			return otto.FalseValue()
		}
		return vmResponse
	} else {
		e.LogErrorf("Function Error: function=%s error=FILE_HAS_NO_CTIME arg=%s", CalledBy(), spew.Sdump(t))
		return otto.FalseValue()
	}
}

func (e *Engine) VMFileModifyTime(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(filePath))
		return otto.FalseValue()
	}
	t, err := times.Stat(filePathAsString.(string))
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=CMD_OUTPUT_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(filePathAsString))
		return otto.FalseValue()
	}
	vmResponse, err := e.VM.ToValue(t.ModTime().String())
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=CMD_OUTPUT_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(t))
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMFileAccessTime(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(filePath))
		return otto.FalseValue()
	}
	t, err := times.Stat(filePathAsString.(string))
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=CMD_OUTPUT_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(filePathAsString))
		return otto.FalseValue()
	}
	vmResponse, err := e.VM.ToValue(t.AccessTime().String())
	if err != nil {
		e.LogErrorf("Function Error: function=%s error=CMD_OUTPUT_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(t))
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMLoggedInUsers(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMUsersRunningProcs(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMServeFileOverHTTP(call otto.FunctionCall) otto.Value {
	e.LogErrorf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}
