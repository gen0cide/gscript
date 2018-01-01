package gscript

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVMCanMakeHTTPConn(t *testing.T) {
	testScript := `
    var url1 = "https://www.google.com";
    var return_value = CanMakeHTTPConn(url1);
  `

	e := New("HTTPTest")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMCanMakeTCPConn(t *testing.T) {
	testScript := `
    var ip = "towel.blinkenlights.nl";
		var port = "23";
    var return_value = CanMakeTCPConn(ip,port);
  `
	e := New("TCPTest")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMHasPublicIP(t *testing.T) {
	testScript := `
    var return_value = HasPublicIP();
  `
	e := New("PublicIPTest")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

/*
func TestVMIsAWS(t *testing.T) {
	testScript := `
    var return_value = IsAWS();
  `
	e := New("AWSTest")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "false", retValAsString)
}
*/

func TestVMFileContains(t *testing.T) {
	testScript := `
    var wordz = "root";
		var file = "/etc/passwd";
		var return_value = FileContains(file, wordz);
  `
	e := New("FileContainsTest")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMLocalUserExists(t *testing.T) {
	testScript := `
    var user = "root";
		var return_value = LocalUserExists(user);
  `
	e := New("LocaUserTest")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMTCPPortInUse(t *testing.T) {
	testScript := `
    var port = 8080;
		var return_value = TCPPortInUse(port);
  `
	e := New("TCPPortInUseTest")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "false", retValAsString)
}

func TestVMProcExistsWithName(t *testing.T) {
	testScript := `
    var name = "notbash";
		var return_value = ProcExistsWithName(name);
  `
	e := New("TestProcExists")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "false", retValAsString)
}

func TestVMDirExists(t *testing.T) {
	testScript := `
    var name = "/etc/";
		var return_value = DirExists(name);
  `
	e := New("TestDirExists")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMFileExists(t *testing.T) {
	testScript := `
    var name = "/etc/passwd";
		var return_value = FileExists(name);
  `
	e := New("TestFileExists")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMCanReadFile(t *testing.T) {
	testScript := `
    var name = "/etc/passwd";
		var return_value = CanReadFile(name);
  `
	e := New("TestCanReadFile")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMCanWriteFile(t *testing.T) {
	testScript := `
    var name = "/etc/passwd";
		var return_value = CanWriteFile(name);
  `
	e := New("TestCanWriteFile")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMCanExecFile(t *testing.T) {
	testScript := `
    var name = "/bin/bash";
		var return_value = CanExecFile(name);
  `
	e := New("TestCanExecFile")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMCmdSuccessful(t *testing.T) {
	testScript := `
    var name = "ls";
		var arg = "-al"
		var return_value = CmdSuccessful(name, arg);
  `
	e := New("TestCmdSuccessful")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMExpectedDNS(t *testing.T) {
	testScript := fmt.Sprintf(`
		var url = "google.com";
		var type2 = "CNAME";
    var return_value2 = ExpectedDNS(url, type2, "google.com.");
  `)

	e := New("DNSQueryTest")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal2, err := e.VM.Get("return_value2")
	assert.Nil(t, err)
	retValAsString2, err := retVal2.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString2)
}

func TestVMCanSudo(t *testing.T) {
	testScript := fmt.Sprintf(`
    var return_value2 = CanSudo();
  `)

	e := New("SudoTest")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal2, err := e.VM.Get("return_value2")
	assert.Nil(t, err)
	retValAsString2, err := retVal2.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString2)
}

func TestVMExistsInPath(t *testing.T) {
	testScript := fmt.Sprintf(`
		var cmd = "whoami"
    var return_value2 = ExistsInPath(cmd);
  `)

	e := New("CmdInPathTest")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal2, err := e.VM.Get("return_value2")
	assert.Nil(t, err)
	retValAsString2, err := retVal2.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString2)
}
