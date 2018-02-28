package engine

import (
	"io/ioutil"

	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
)

type Engine struct {
	VM              *otto.Otto
	Logger          *logrus.Logger
	Imports         map[string]func() []byte
	Name            string
	DebuggerEnabled bool
	Timeout         int
	Halted          bool
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
		Halted:          false,
		Timeout:         30,
	}
}

func (e *Engine) SetLogger(logger *logrus.Logger) {
	e.Logger = logger
}

func (e *Engine) SetTimeout(timeout int) {
	e.Timeout = timeout
}

func (e *Engine) AddImport(name string, data func() []byte) {
	e.Imports[name] = data
}

func (e *Engine) SetName(name string) {
	e.Name = name
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
	e.VM.Set("LogDebug", e.VMLogDebug)
	e.VM.Set("LogInfo", e.VMLogInfo)
	e.VM.Set("LogWarn", e.VMLogWarn)
	e.VM.Set("LogError", e.VMLogError)
	e.VM.Set("LogFatal", e.VMLogCrit)
	e.VM.Set("ForkExec", e.VMForkExec)
	e.VM.Set("ShellcodeExec", e.VMShellcodeExec)
	e.VM.Set("AddRegKey", e.VMAddRegKey)
	e.VM.Set("QueryRegKey", e.VMQueryRegKey)
	e.VM.Set("DelRegKey", e.VMDelRegKey)
	e.VM.Set("GetHost", e.VMGetHostname)
	e.VM.Set("LogTester", e.VMLogTester)
	e.VM.Set("InstallSystemService", e.VMInstallSystemService)
	e.VM.Set("StartServiceByName", e.VMStartServiceByName)
	e.VM.Set("StopServiceByName", e.VMStopServiceByName)
	e.VM.Set("RemoveServiceByName", e.VMRemoveServiceByName)
	_, err := e.VM.Run(VMPreload)
	if err != nil {
		e.Logger.WithField("trace", "true").Fatalf("Syntax error in preload: %s", err.Error())
	}
}
