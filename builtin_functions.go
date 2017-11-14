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
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
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
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMDNSQuery(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMHTTPRequest(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMCmd(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
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
