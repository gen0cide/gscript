package engine

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
)

type Engine struct {
	VM              *otto.Otto
	Logger          *logrus.Logger
	Imports         map[string]func() []byte
	Name            string
	DebuggerEnabled bool
}

func New(name string) *Engine {
	logger := logrus.New()
	logger.Formatter = new(logrus.TextFormatter)
	logger.Out = ioutil.Discard
	return &Engine{
		Name:            name,
		Imports:         map[string]func() []byte{},
		Logger:          logger,
		DebuggerEnabled: false,
	}
}

func (e *Engine) SetLogger(logger *logrus.Logger) {
	e.Logger = logger
}

func (e *Engine) AddImport(name string, data func() []byte) {
	e.Imports[name] = data
}

func (e *Engine) SetName(name string) {
	e.Name = name
}

func (e *Engine) CurrentUser() map[string]string {
	userInfo := map[string]string{}
	u, err := user.Current()
	if err != nil {
		e.Logger.Errorf("User Loading Error: %s", err.Error())
		return userInfo
	}
	userInfo["uid"] = u.Uid
	userInfo["gid"] = u.Gid
	userInfo["username"] = u.Username
	userInfo["name"] = u.Name
	userInfo["home_dir"] = u.HomeDir
	groups, err := u.GroupIds()
	if err != nil {
		e.Logger.Errorf("Group Loading Error: %s", err.Error())
		return userInfo
	}
	userInfo["groups"] = strings.Join(groups, ":")
	return userInfo
}

func (e *Engine) InjectVars() {
	userInfo, err := e.VM.ToValue(e.CurrentUser())
	if err != nil {
		e.Logger.Errorf("Could not inject user info into VM: %s", err.Error())
	} else {
		e.VM.Set("USER_INFO", userInfo)
	}
	osVal, err := e.VM.ToValue(runtime.GOOS)
	if err != nil {
		e.Logger.Errorf("Could not inject os info into VM: %s", err.Error())
	} else {
		e.VM.Set("OS", osVal)
	}
	hn, err := os.Hostname()
	if err != nil {
		e.Logger.Errorf("Could not obtain hostname info: %s", err.Error())
	} else {
		hostnameVal, err := e.VM.ToValue(hn)
		if err != nil {
			e.Logger.Errorf("Could not inject hostname info into VM: %s", err.Error())
		} else {
			e.VM.Set("HOSTNAME", hostnameVal)
		}
	}
	archVal, err := e.VM.ToValue(runtime.GOARCH)
	if err != nil {
		e.Logger.Errorf("Could not inject arch info into VM: %s", err.Error())
	} else {
		e.VM.Set("ARCH", archVal)
	}
	ipVals, err := e.VM.ToValue(GetLocalIPs())
	if err != nil {
		e.Logger.Errorf("Could not inject ip info into VM: %s", err.Error())
	} else {
		e.VM.Set("IP_ADDRS", ipVals)
	}
}

func (e *Engine) CreateVM() {
	e.VM = otto.New()
	e.InjectVars()
	e.VM.Set("Halt", e.VMHalt)
	e.VM.Set("Asset", e.VMAsset)
	e.VM.Set("DeleteFile", e.VMDeleteFile)
	e.VM.Set("CopyFile", e.VMCopyFile)
	e.VM.Set("WriteFile", e.VMWriteFile)
	e.VM.Set("ReadFile", e.VMReadFile)
	e.VM.Set("ExecuteFile", e.VMExecuteFile)
	e.VM.Set("AppendFile", e.VMAppendFile)
	e.VM.Set("ReplaceInFile", e.VMReplaceInFile)
	e.VM.Set("Signal", e.VMSignal)
	e.VM.Set("Implode", e.VMImplode)
	e.VM.Set("LocalUserExists", e.VMLocalUserExists)
	e.VM.Set("ProcExistsWithName", e.VMProcExistsWithName)
	e.VM.Set("CanReadFile", e.VMCanReadFile)
	e.VM.Set("CanWriteFile", e.VMCanWriteFile)
	e.VM.Set("CanExecFile", e.VMCanExecFile)
	e.VM.Set("FileExists", e.VMFileExists)
	e.VM.Set("DirExists", e.VMDirExists)
	e.VM.Set("FileContains", e.VMFileContains)
	e.VM.Set("IsVM", e.VMIsVM)
	e.VM.Set("IsAWS", e.VMIsAWS)
	e.VM.Set("HasPublicIP", e.VMHasPublicIP)
	e.VM.Set("CanMakeTCPConn", e.VMCanMakeTCPConn)
	e.VM.Set("ExpectedDNS", e.VMExpectedDNS)
	e.VM.Set("CanMakeHTTPConn", e.VMCanMakeHTTPConn)
	e.VM.Set("DetectSSLMITM", e.VMDetectSSLMITM)
	e.VM.Set("CmdSuccessful", e.VMCmdSuccessful)
	e.VM.Set("CanPing", e.VMCanPing)
	e.VM.Set("TCPPortInUse", e.VMTCPPortInUse)
	e.VM.Set("UDPPortInUse", e.VMUDPPortInUse)
	e.VM.Set("ExistsInPath", e.VMExistsInPath)
	e.VM.Set("CanSudo", e.VMCanSudo)
	e.VM.Set("Matches", e.VMMatches)
	e.VM.Set("CanSSHLogin", e.VMCanSSHLogin)
	e.VM.Set("RetrieveFileFromURL", e.VMRetrieveFileFromURL)
	e.VM.Set("DNSQuery", e.VMDNSQuery)
	e.VM.Set("HTTPRequest", e.VMHTTPRequest)
	e.VM.Set("Exec", e.VMExec)
	e.VM.Set("MD5", e.VMMD5)
	e.VM.Set("SHA1", e.VMSHA1)
	e.VM.Set("B64Decode", e.VMB64Decode)
	e.VM.Set("B64Encode", e.VMB64Encode)
	e.VM.Set("Timestamp", e.VMTimestamp)
	e.VM.Set("CPUStats", e.VMCPUStats)
	e.VM.Set("MemStats", e.VMMemStats)
	e.VM.Set("SSHCmd", e.VMSSHCmd)
	e.VM.Set("Sleep", e.VMSleep)
	e.VM.Set("GetTweet", e.VMGetTweet)
	e.VM.Set("GetDirsInPath", e.VMGetDirsInPath)
	e.VM.Set("EnvVars", e.VMEnvVars)
	e.VM.Set("GetEnv", e.VMGetEnv)
	e.VM.Set("FileChangeTime", e.VMFileChangeTime)
	e.VM.Set("FileModifyTime", e.VMFileModifyTime)
	e.VM.Set("FileAccessTime", e.VMFileAccessTime)
	e.VM.Set("FileBirthTime", e.VMFileBirthTime)
	e.VM.Set("LoggedInUsers", e.VMLoggedInUsers)
	e.VM.Set("UsersRunningProcs", e.VMUsersRunningProcs)
	e.VM.Set("ServeFileOverHTTP", e.VMServeFileOverHTTP)
	e.VM.Set("VMLogDebug", e.VMLogDebug)
	e.VM.Set("VMLogInfo", e.VMLogInfo)
	e.VM.Set("VMLogWarn", e.VMLogWarn)
	e.VM.Set("VMLogError", e.VMLogError)
	e.VM.Set("VMLogCrit", e.VMLogCrit)
	e.VM.Set("ForkExec", e.VMForkExec)
	e.VM.Set("ShellcodeExec", e.VMShellcodeExec)
	e.VM.Set("AddRegKey", e.VMAddRegKey)
	e.VM.Set("QueryRegKey", e.VMQueryRegKey)
	e.VM.Set("DelRegKey", e.VMDelRegKey)
	e.VM.Set("GetHost", e.VMGetHostname)
	_, err := e.VM.Run(VMPreload)
	if err != nil {
		e.Logger.Fatalf("Syntax error in preload: %s", err.Error())
	}
}

func (e *Engine) ValueToByteSlice(v otto.Value) []byte {
	valueBytes := []byte{}
	if v.IsNull() || v.IsUndefined() {
		return valueBytes
	}
	if v.IsString() {
		str, err := v.Export()
		if err != nil {
			e.Logger.Errorf("Cannot convert string to byte slice: %s", v)
			return valueBytes
		}
		valueBytes = []byte(str.(string))
	} else if v.IsNumber() {
		num, err := v.Export()
		if err != nil {
			e.Logger.Errorf("Cannot convert string to byte slice: %s", v)
			return valueBytes
		}
		buf := new(bytes.Buffer)
		err = binary.Write(buf, binary.LittleEndian, num)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
		valueBytes = buf.Bytes()
	} else if v.Class() == "Array" || v.Class() == "GoArray" {
		arr, err := v.Export()
		if err != nil {
			e.Logger.Errorf("Cannot convert array to byte slice: %x", v)
			return valueBytes
		}
		switch t := arr.(type) {
		case []uint:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []uint8:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []uint16:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []uint32:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []uint64:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []int:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []int16:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []int32:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []int64:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []float32:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []float64:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []string:
			for _, i := range t {
				for _, c := range i {
					valueBytes = append(valueBytes, byte(c))
				}
			}
		default:
			_ = t
			e.Logger.Errorf("Failed to cast array to byte slice array=%v", arr)
		}
	} else {
		e.Logger.Errorf("Unknown class to cast to byte slice")
	}

	return valueBytes
}

func (e *Engine) VMLogWarn(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.Logger.Errorf("Logging Error - argument couldn't be converted to string: %s", spew.Sdump(arg))
				continue
			}
			e.Logger.Warnf("%s", newArg)
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogError(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.Logger.Errorf("Logging Error - argument couldn't be converted to string: %s", spew.Sdump(arg))
				continue
			}
			e.Logger.Errorf("%s", newArg)
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogDebug(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.Logger.Errorf("Logging Error - argument couldn't be converted to string: %s", spew.Sdump(arg))
				continue
			}
			e.Logger.Debugf("%s", newArg)
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogCrit(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.Logger.Errorf("Logging Error - argument couldn't be converted to string: %s", spew.Sdump(arg))
				continue
			}
			e.Logger.Fatalf("%s", newArg)
		}
	}
	return otto.Value{}
}

func (e *Engine) VMLogInfo(call otto.FunctionCall) otto.Value {
	if e.Logger != nil {
		for _, arg := range call.ArgumentList {
			newArg, err := arg.ToString()
			if err != nil {
				e.Logger.Errorf("Logging Error - argument couldn't be converted to string: %s", spew.Sdump(arg))
				continue
			}
			e.Logger.Infof("%s", newArg)
		}
	}
	return otto.Value{}
}
