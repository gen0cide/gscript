package gscript

import (
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
