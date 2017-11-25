package gscript

import (
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
