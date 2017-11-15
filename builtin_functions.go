package gscript

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"time"
	"github.com/davecgh/go-spew/spew"
	"github.com/robertkrimen/otto"
)

func (e *Engine) VMHalt(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMDeleteFile(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMWriteFile(call otto.FunctionCall) otto.Value {
	// Arg1 is a byte array of bytes of the file
	// Arg2 is the string path of the new file to write
	filePath := call.Argument(0)
	fileData := call.Argument(1)
	fileBytes := e.ValueToByteSlice(fileData)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.Logger.Errorf("Function Error: function=%s error=ARY_ARG_NOT_String arg=%s", CalledBy(), spew.Sdump(filePath))
		return otto.FalseValue()
	}
	err = LocalFileCreate(filePathAsString.(string), fileBytes)
	if err != nil {
		e.Logger.Errorf("Error writing the file: function=%s path=%s error=%s", CalledBy(), filePathAsString.(string), err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMCopyFile(call otto.FunctionCall) otto.Value {
	// Arg1 is the string path of the first file we read
	// Arg2 is the string path of the new file to write
	readPath, err := call.ArgumentList[0].ToString()
	if err != nil {
		e.Logger.Errorf("Function Error: function=VMCopyFile() error=ARG_NOT_STRINGABLE arg=%s", spew.Sdump(call.ArgumentList[0]))
		return otto.FalseValue()
	}
	writePath, err := call.ArgumentList[1].ToString()
	if err != nil {
		e.Logger.Errorf("Function Error: function=VMCopyFile() error=ARG_NOT_STRINGABLE arg=%s", spew.Sdump(call.ArgumentList[1]))
		return otto.FalseValue()
	}
	bytes, err := LocalFileRead(readPath)
	if err != nil {
		e.Logger.Errorf("Function Error: function=VMCopyFile() error='There was an error reading the file to copy'")
		return otto.FalseValue()
	}
	//e.Logger.Errorf("Function Error: function=VMCopyFile() error=Debug; read local file at: %s", spew.Sdump(readPath))
	err = LocalFileCreate(writePath, bytes)
	if err != nil {
		e.Logger.Errorf("Function Error: function=VMCopyFile() error='There was an error writing the file to that path'")
		return otto.FalseValue()
	}
	// Testing Debug call function //
	e.Logger.Errorf("Function Error: function=%s error=Debug; wrote local file at: %s", CalledBy(), spew.Sdump(writePath))
	//returnString := fmt.Sprintf("File created at: %s", string(writePath))
	//var ret = otto.Value{}
	//var er error
	//ret, er = otto.ToValue(returnString)
	//if er != nil {
	//	e.Logger.Errorf("Function Error: function=VMCopyFile() error='Error returning value to VM'")
	//	return otto.FalseValue()
	//}
	return otto.TrueValue()
}

func (e *Engine) VMExecuteFile(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMAppendFile(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMReplaceInFile(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMSignal(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMImplode(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMRetrieveFileFromURL(call otto.FunctionCall) otto.Value {
	readURL, err := call.ArgumentList[0].ToString()
	if err != nil {
		e.Logger.Errorf("Function Error: function=%s error=ARG_NOT_STRINGABLE arg=%s", CalledBy(), spew.Sdump(call.ArgumentList[0]))
		return otto.FalseValue()
	}
	bytes, err := HTTPGetFile(readURL)
	if err != nil {
		e.Logger.Errorf("Function Error: function=VMRetrieveFileFromURL() error='There was an error fetching the file: %v'", spew.Sdump(err))
		return otto.FalseValue()
	}
	vmResponse, err := e.VM.ToValue(bytes)
	if err != nil {
		e.Logger.Errorf("Function Error: function=%s error=URL_DOWNLOAD_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(err))
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMDNSQuery(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMHTTPRequest(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMExec(call otto.FunctionCall) otto.Value {
	baseCmd := call.Argument(0)
	if !baseCmd.IsString() {
		e.Logger.Errorf("Function Error: function=%s error=CMD_CALL_NOT_OF_TYPE_STRING arg=%s", CalledBy(), spew.Sdump(baseCmd))
		return otto.FalseValue()
	}
	cmdArgs := call.Argument(1)
	argList := []string{}
	if !cmdArgs.IsNull() {
		argArray, err := cmdArgs.Export()
		if err != nil {
			e.Logger.Errorf("Function Error: function=%s error=CMD_ARGS_NOT_PARSABLE arg=%s", CalledBy(), spew.Sdump(cmdArgs))
			return otto.FalseValue()
		}
		argList = argArray.([]string)
	}
	baseCmdAsString, err := baseCmd.ToString()
	if err != nil {
		e.Logger.Errorf("Function Error: function=%s error=CMD_BASE_NOT_PARSABLE arg=%s", CalledBy(), spew.Sdump(baseCmd))
		return otto.FalseValue()
	}
	cmdOutput := ExecuteCommand(baseCmdAsString, argList...)
	vmResponse, err := e.VM.ToValue(cmdOutput)
	if err != nil {
		e.Logger.Errorf("Function Error: function=%s error=CMD_OUTPUT_OBJECT_CAST_FAILED arg=%s", CalledBy(), spew.Sdump(baseCmd))
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
				e.Logger.Errorf("Function Error: function=MD5() error=ARY_ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
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
				e.Logger.Errorf("Function Error: function=MD5() error=ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
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
		e.Logger.Errorf("Function Error: function=%s error=%s", CalledBy(), err.Error())
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
				e.Logger.Errorf("Function Error: function=SHA1() error=ARY_ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
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
				e.Logger.Errorf("Function Error: function=SHA1() error=ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
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
		e.Logger.Errorf("Function Error: function=%s error=%s", CalledBy(), err.Error())
	}
	return ret
}

func (e *Engine) VMB64Decode(call otto.FunctionCall) otto.Value {
	var NewVal string
	if len(call.ArgumentList) > 0 {
		arg := call.ArgumentList[0]
		val, err := arg.ToString()
		if err != nil {
			e.Logger.Errorf("Function Error: function=B64Decode() error=ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
			return otto.Value{}
		}
		valBytes, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			e.Logger.Errorf("Function Error: function=B64Decode() error=BASE64_DECODE_ERROR error=%s", err.Error())
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
		e.Logger.Errorf("Function Error: function=%s error=%s", CalledBy(), err.Error())
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
				e.Logger.Errorf("Function Error: function=B64Encode() error=ARY_ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
				return otto.Value{}
			}
			for _, i := range expVal.([]interface{}) {
				rawBytes = append(rawBytes, i.(byte))
			}
			EncVal = base64.StdEncoding.EncodeToString([]byte(rawBytes))
		default:
			val, err := call.ArgumentList[0].ToString()
			if err != nil {
				e.Logger.Errorf("Function Error: function=B64Encode() error=ARG_NOT_STRINGABLE arg=%s", spew.Sdump(arg))
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
		e.Logger.Errorf("Function Error: function=%s error=%s", CalledBy(), err.Error())
	}
	return ret
}

func (e *Engine) VMTimestamp(call otto.FunctionCall) otto.Value {
	ts := time.Now().Unix()
	ret, err := otto.ToValue(ts)
	if err != nil {
		e.Logger.Errorf("Function Error: function=%s error=%s", CalledBy(), err.Error())
	}
	return ret
}

func (e *Engine) VMCPUStats(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMMemStats(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMSSHCmd(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMSleep(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 0 {
		arg := call.ArgumentList[0]
		if arg.IsNumber() {
			intArg, err := arg.ToInteger()
			if err != nil {
				e.Logger.Errorf("Function Error: function=%s error=%s", CalledBy(), err.Error())
				return otto.Value{}
			}
			time.Sleep(time.Duration(intArg) * time.Second)
		}
	}
	return otto.Value{}
}

func (e *Engine) VMGetTweet(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMGetDirsInPath(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMEnvVars(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMGetEnv(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMFileCreateTime(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMFileModifyTime(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMLoggedInUsers(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMUsersRunningProcs(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMServeFileOverHTTP(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}
