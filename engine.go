package gscript

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/happierall/l"
	"github.com/robertkrimen/otto"
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
	e.VM.Set("BeforeDeploy", e.VMBeforeDeploy)
	e.VM.Set("Deploy", e.VMDeploy)
	e.VM.Set("AfterDeploy", e.VMAfterDeploy)
	e.VM.Set("OnError", e.VMOnError)
	e.VM.Set("Halt", e.VMHalt)
	e.VM.Set("DeleteFile", e.VMDeleteFile)
	e.VM.Set("CopyFile", e.VMCopyFile)
	e.VM.Set("WriteFile", e.VMWriteFile)
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
	e.VM.Set("FileCreateTime", e.VMFileCreateTime)
	e.VM.Set("FileModifyTime", e.VMFileModifyTime)
	e.VM.Set("LoggedInUsers", e.VMLoggedInUsers)
	e.VM.Set("UsersRunningProcs", e.VMUsersRunningProcs)
	e.VM.Set("ServeFileOverHTTP", e.VMServeFileOverHTTP)
	e.VM.Set("VMLogDebug", e.VMLogDebug)
	e.VM.Set("VMLogInfo", e.VMLogInfo)
	e.VM.Set("VMLogWarn", e.VMLogWarn)
	e.VM.Set("VMLogError", e.VMLogError)
	e.VM.Set("VMLogCrit", e.VMLogCrit)
	_, err := e.VM.Run(VMPreload)
	if err != nil {
		panic(err)
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
			e.LogErrorf("Cannot convert string to byte slice: %s", spew.Sdump(v))
			return valueBytes
		}
		valueBytes = []byte(str.(string))
	} else if v.IsNumber() {
		num, err := v.Export()
		if err != nil {
			e.LogErrorf("Cannot convert string to byte slice: %s", spew.Sdump(v))
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
			e.LogErrorf("Cannot convert array to byte slice: %s", spew.Sdump(v))
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
			e.Logger.Errorf("Failed to cast array to byte slice: function=%s array=%v", CalledBy(), arr)
		}
	} else {
		spew.Dump(v)
		spew.Dump(v.Class())
		e.Logger.Errorf("Unknown class to cast to byte slice: function=%s value=%s class=%s", CalledBy(), spew.Sdump(v), spew.Sdump(v.Class()))
	}

	return valueBytes
}
