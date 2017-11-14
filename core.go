package gscript

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"runtime"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/happierall/l"
	"github.com/robertkrimen/otto"

	// Include Underscore In Otto :)
	_ "github.com/robertkrimen/otto/underscore"
)

var (
	Debugger      = true
	DefaultScript = ` // genesis script

var helloWorld = "helloworld";
var foo = MD5(helloWorld);
Halt();
VMLogInfo(foo);

`
)

type Engine struct {
	VM     *otto.Otto
	Logger *l.Logger
}

func New() *Engine {
	return &Engine{}
}

func (e *Engine) EnableLogging() {
	e.Logger = l.New()
	e.Logger.Prefix = "[GENESIS] "
	e.Logger.DisabledInfo = false
}

func (e *Engine) CreateVM() {
	e.VM = otto.New()
	e.VM.Set("BeforeDeploy", e.BeforeDeploy)
	e.VM.Set("Deploy", e.Deploy)
	e.VM.Set("AfterDeploy", e.AfterDeploy)
	e.VM.Set("OnError", e.OnError)
	e.VM.Set("Halt", e.Halt)
	e.VM.Set("DeleteFile", e.DeleteFile)
	e.VM.Set("WriteFile", e.WriteFile)
	e.VM.Set("ExecuteFile", e.ExecuteFile)
	e.VM.Set("AppendFile", e.AppendFile)
	e.VM.Set("ReplaceInFile", e.ReplaceInFile)
	e.VM.Set("Signal", e.Signal)
	e.VM.Set("Implode", e.Implode)
	e.VM.Set("LocalUserExists", e.LocalUserExists)
	e.VM.Set("ProcExistsWithName", e.ProcExistsWithName)
	e.VM.Set("CanReadFile", e.CanReadFile)
	e.VM.Set("CanWriteFile", e.CanWriteFile)
	e.VM.Set("CanExecFile", e.CanExecFile)
	e.VM.Set("FileExists", e.FileExists)
	e.VM.Set("DirExists", e.DirExists)
	e.VM.Set("FileContains", e.FileContains)
	e.VM.Set("IsVM", e.IsVM)
	e.VM.Set("IsAWS", e.IsAWS)
	e.VM.Set("HasPublicIP", e.HasPublicIP)
	e.VM.Set("CanMakeTCPConn", e.CanMakeTCPConn)
	e.VM.Set("ExpectedDNS", e.ExpectedDNS)
	e.VM.Set("CanMakeHTTPConn", e.CanMakeHTTPConn)
	e.VM.Set("DetectSSLMITM", e.DetectSSLMITM)
	e.VM.Set("CmdSuccessful", e.CmdSuccessful)
	e.VM.Set("CanPing", e.CanPing)
	e.VM.Set("TCPPortInUse", e.TCPPortInUse)
	e.VM.Set("UDPPortInUse", e.UDPPortInUse)
	e.VM.Set("ExistsInPath", e.ExistsInPath)
	e.VM.Set("CanSudo", e.CanSudo)
	e.VM.Set("Matches", e.Matches)
	e.VM.Set("CanSSHLogin", e.CanSSHLogin)
	e.VM.Set("RetrieveFileFromURL", e.RetrieveFileFromURL)
	e.VM.Set("DNSQuery", e.DNSQuery)
	e.VM.Set("HTTPRequest", e.HTTPRequest)
	e.VM.Set("Cmd", e.Cmd)
	e.VM.Set("MD5", e.MD5)
	e.VM.Set("SHA1", e.SHA1)
	e.VM.Set("B64Decode", e.B64Decode)
	e.VM.Set("B64Encode", e.B64Encode)
	e.VM.Set("Timestamp", e.Timestamp)
	e.VM.Set("CPUStats", e.CPUStats)
	e.VM.Set("MemStats", e.MemStats)
	e.VM.Set("SSHCmd", e.SSHCmd)
	e.VM.Set("Sleep", e.Sleep)
	e.VM.Set("GetTweet", e.GetTweet)
	e.VM.Set("GetDirsInPath", e.GetDirsInPath)
	e.VM.Set("EnvVars", e.EnvVars)
	e.VM.Set("GetEnv", e.GetEnv)
	e.VM.Set("FileCreateTime", e.FileCreateTime)
	e.VM.Set("FileModifyTime", e.FileModifyTime)
	e.VM.Set("LoggedInUsers", e.LoggedInUsers)
	e.VM.Set("UsersRunningProcs", e.UsersRunningProcs)
	e.VM.Set("ServeFileOverHTTP", e.ServeFileOverHTTP)
	e.VM.Set("VMLogDebug", e.VMLogDebug)
	e.VM.Set("VMLogInfo", e.VMLogInfo)
	e.VM.Set("VMLogWarn", e.VMLogWarn)
	e.VM.Set("VMLogError", e.VMLogError)
	e.VM.Set("VMLogCrit", e.VMLogCrit)
}

func (e *Engine) LogWarn(i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Warn(i...)
	}
}

func (e *Engine) LogError(i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Error(i...)
	}
}

func (e *Engine) LogDebug(i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Debug(i...)
	}
}

func (e *Engine) LogCrit(i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Crit(i...)
	}
}

func (e *Engine) LogInfo(i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Log(i...)
	}
}

func (e *Engine) LogWarnf(fmtString string, i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Warnf(fmtString, i...)
	}
}

func (e *Engine) LogErrorf(fmtString string, i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Errorf(fmtString, i...)
	}
}

func (e *Engine) LogDebugf(fmtString string, i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Debugf(fmtString, i...)
	}
}

func (e *Engine) LogCritf(fmtString string, i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Critf(fmtString, i...)
	}
}

func (e *Engine) LogInfof(fmtString string, i ...interface{}) {
	if e.Logger != nil {
		e.Logger.Logf(fmtString, i...)
	}
}

func (e *Engine) BeforeDeploy(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) Deploy(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) AfterDeploy(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) OnError(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) Halt(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) DeleteFile(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) WriteFile(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) ExecuteFile(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) AppendFile(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) ReplaceInFile(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) Signal(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) Implode(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) LocalUserExists(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) ProcExistsWithName(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) CanReadFile(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) CanWriteFile(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) CanExecFile(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) FileExists(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) DirExists(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) FileContains(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) IsVM(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) IsAWS(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) HasPublicIP(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) CanMakeTCPConn(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) ExpectedDNS(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) CanMakeHTTPConn(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) DetectSSLMITM(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) CmdSuccessful(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) CanPing(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) TCPPortInUse(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) UDPPortInUse(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) ExistsInPath(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) CanSudo(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) Matches(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) CanSSHLogin(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) RetrieveFileFromURL(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) DNSQuery(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) HTTPRequest(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) Cmd(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) MD5(call otto.FunctionCall) otto.Value {
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
		e.Logger.Errorf("Function Error: function=MD5() error=%s", err.Error())
	}
	return ret
}

func (e *Engine) SHA1(call otto.FunctionCall) otto.Value {
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
		e.Logger.Errorf("Function Error: function=SHA1() error=%s", err.Error())
	}
	return ret
}

func (e *Engine) B64Decode(call otto.FunctionCall) otto.Value {
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
		e.Logger.Errorf("Function Error: function=B64Encode() error=%s", err.Error())
	}
	return ret
}

func (e *Engine) B64Encode(call otto.FunctionCall) otto.Value {
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
		e.Logger.Errorf("Function Error: function=B64Encode() error=%s", err.Error())
	}
	return ret
}

func (e *Engine) Timestamp(call otto.FunctionCall) otto.Value {
	ts := time.Now().Unix()
	ret, err := otto.ToValue(ts)
	if err != nil {
		e.Logger.Errorf("Function Error: function=SHA1() error=%s", err.Error())
	}
	return ret
}

func (e *Engine) CPUStats(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) MemStats(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) SSHCmd(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) Sleep(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) > 0 {
		arg := call.ArgumentList[0]
		if arg.IsNumber() {
			intArg, err := arg.ToInteger()
			if err != nil {
				e.Logger.Errorf("Function Error: function=Sleep() error=%s", err.Error())
				return otto.Value{}
			}
			time.Sleep(time.Duration(intArg) * time.Second)
		}
	}
	return otto.Value{}
}

func (e *Engine) GetTweet(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) GetDirsInPath(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) EnvVars(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) GetEnv(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) FileCreateTime(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) FileModifyTime(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) LoggedInUsers(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) UsersRunningProcs(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) ServeFileOverHTTP(call otto.FunctionCall) otto.Value {
	e.LogCritf("Function Not Implemented: %s", CalledBy())
	return otto.FalseValue()
}

func (e *Engine) VMLogWarn(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			e.Logger.Warn(arg.ToString())
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogError(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			e.Logger.Error(arg.ToString())
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogDebug(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			e.Logger.Debug(arg.ToString())
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogCrit(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			e.Logger.Crit(arg.ToString())
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogInfo(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			e.Logger.Log(arg.ToString())
		}
	}
	return otto.Value{}
}

func CalledBy() string {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "Unknown"
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "N/A"
	}

	return fun.Name()
}
