package gscript

import (
	"fmt"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

var g_file_1 = fmt.Sprintf("/tmp/%s", RandString(6))
var g_file_2 = fmt.Sprintf("/tmp/%s", RandString(6))
var g_file_3 = fmt.Sprintf("/tmp/%s", RandString(6))

func TestVMMD5(t *testing.T) {
	testScript := `
    var hash_val = "helloworld";
    var return_value = MD5(hash_val);
  `
	// "helloworld" = fc5e038d38a57032085441e7fe7010b0

	e := New("MD5")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "fc5e038d38a57032085441e7fe7010b0", retValAsString)
}

func TestVMCopyFile(t *testing.T) {
	file_2 := g_file_1
	testScript := fmt.Sprintf(`
    var file_1 = "/etc/passwd";
    var file_2 = "%s";
    var return_value = CopyFile(file_1, file_2);
  `, file_2)

	e := New("CopyFile")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString)
}

func TestVMAppendFile(t *testing.T) {
	bytes := "60,104,116,109,108,62,10,32,32,60,98,111,100,121,62,10,32,32,32,32"
	testScript := fmt.Sprintf(`
    var file_1 = "%s";
		var file_2 = "%s";
    var bytes = [%s];
    var return_value1 = AppendFile(file_1, bytes);
		var return_value2 = AppendFile(file_2, bytes);
  `, g_file_1, g_file_2, bytes)

	e := New("AppendFile")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	e.LogInfof("Function: function=%s msg='Appended local file at: %s'", CalledBy(), spew.Sdump(g_file_1))
	e.LogInfof("Function: function=%s msg='Appended local file at: %s'", CalledBy(), spew.Sdump(g_file_2))
	retVal, err := e.VM.Get("return_value1")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString)
	retVal2, err := e.VM.Get("return_value2")
	assert.Nil(t, err)
	retValAsString2, err := retVal2.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString2)
}

func TestVMReplaceInFile(t *testing.T) {
	string01 := "root"
	string02 := "lol"
	testScript := fmt.Sprintf(`
    var file_1 = "%s";
    var string01 = "%s";
		var string02 = "%s";
    var return_value1 = ReplaceInFile(file_1, string01, string02);
  `, g_file_1, string01, string02)

	e := New("ReplaceInFile")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value1")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString)
}

func TestVMRetrieveFileFromURL(t *testing.T) {
	url := "https://alexlevinson.com/"
	file_3 := g_file_3
	testScript2 := fmt.Sprintf(`
	  var url = "%s";
		var file_3 = "%s";
	  var response2 = RetrieveFileFromURL(url);
	  var return_value2 = response2;
		var response3 = WriteFile(file_3, response2);
  `, url, file_3)
	e := New("RetrieveFileFromURL")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript2)
	e.LogInfof("Function: function=%s msg='wrote local file at: %s'", CalledBy(), spew.Sdump(file_3))
	retVal, err := e.VM.Get("return_value2")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "60,104,116,109,108,62,10,32,32,60,98,111,100,121,62,10,32,32,32,32,60,99,101,110,116,101,114,62,10,32,32,32,32,32,32,60,105,109,103,32,115,114,99,61,34,114,111,111,116,46,106,112,103,34,32,47,62,10,32,32,32,32,60,47,99,101,110,116,101,114,62,10,32,32,60,47,98,111,100,121,62,10,60,47,104,116,109,108,62,10", retValAsString)
}

func TestVMDeleteFile(t *testing.T) {
	testScript := fmt.Sprintf(`
		var file_1 = "%s";
    var return_value1 = DeleteFile(file_1);
    var file_2 = "%s";
    var return_value2 = DeleteFile(file_2);
		var file_3 = "%s";
    var return_value3 = DeleteFile(file_3);
  `, g_file_1, g_file_2, g_file_3)

	e := New("DeleteFile")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value1")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString)
	retVal2, err := e.VM.Get("return_value2")
	assert.Nil(t, err)
	retValAsString2, err := retVal2.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString2)
	retVal3, err := e.VM.Get("return_value3")
	assert.Nil(t, err)
	retValAsString3, err := retVal3.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString3)
}

func TestVMDNSQuery(t *testing.T) {
	testScript := fmt.Sprintf(`
		var url = "google.com";
		var ip = "8.8.8.8";
		var type1 = "A";
    var return_value1 = DNSQuery(url, type1);
		var type2 = "CNAME";
    var return_value2 = DNSQuery(url, type2);
		var type3 = "TXT";
    var return_value3 = DNSQuery(url, type3);
		var type4 = "MX";
    var return_value4 = DNSQuery(url, type4);
		var type5 = "NS";
    var return_value5 = DNSQuery(url, type5);
		var type6 = "PTR";
    var return_value6 = DNSQuery(ip, type6);
  `)

	e := New("DNSQuery")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value1")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString)
	retVal2, err := e.VM.Get("return_value2")
	assert.Nil(t, err)
	retValAsString2, err := retVal2.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString2)
	retVal3, err := e.VM.Get("return_value3")
	assert.Nil(t, err)
	retValAsString3, err := retVal3.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString3)
	retVal4, err := e.VM.Get("return_value4")
	assert.Nil(t, err)
	retValAsString4, err := retVal4.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString4)
	retVal5, err := e.VM.Get("return_value5")
	assert.Nil(t, err)
	retValAsString5, err := retVal5.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString5)
	retVal6, err := e.VM.Get("return_value6")
	assert.Nil(t, err)
	retValAsString6, err := retVal6.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString6)
}

func TestVMTimestamp(t *testing.T) {
	currTime := time.Now().Unix()

	testScript := `
    var test_time = Timestamp();
  `
	e := New("Timestamp")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("test_time")
	assert.Nil(t, err)
	assert.True(t, retVal.IsNumber())
	retValAsNumber, err := retVal.ToInteger()
	assert.Nil(t, err)
	assert.True(t, (retValAsNumber >= currTime))
}

func TestExec(t *testing.T) {
	testCmd := ExecuteCommand("ls", "-lah")

	testScript := `
      var test_exec = Exec("ls", ["-lah"]);
    `
	e := New("Exec")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("test_exec")
	assert.Nil(t, err)
	assert.True(t, retVal.IsObject())
	retValAsInterface, err := retVal.Export()
	assert.Nil(t, err)
	realRetVal := retValAsInterface.(VMExecResponse)

	assert.Equal(t, testCmd.Stdout, realRetVal.Stdout)
}

func TestCPUStats(t *testing.T) {
	//resultz := CPUStats()
	testScript := `
      var results = CPUStats();
    `
	e := New("CPUStats")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("results")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString)
}

func TestVMExecuteFile(t *testing.T) {
	testScript := `
			var file_path = "uname";
			var args = ["-o"]
      var results = ExecuteFile(file_path, args);
    `
	e := New("ExecuteFileTest")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("results")
	assert.Nil(t, err)
	assert.True(t, retVal.IsObject())
	retValAsInterface, err := retVal.Export()
	assert.Nil(t, err)
	realRetVal := retValAsInterface.(VMExecResponse)
	assert.Nil(t, err)
	assert.Equal(t, "GNU/Linux", realRetVal.Stdout[0])
}

func TestVMEnvVars(t *testing.T) {
	testScript := `
      var results = EnvVars();
    `
	e := New("EnvVarsTest")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("results")
	assert.Nil(t, err)
	assert.True(t, retVal.IsObject())
	retValAsInterface, err := retVal.Export()
	assert.Nil(t, err)
	realRetVal := retValAsInterface.(map[string]string)
	assert.Nil(t, err)
	assert.Equal(t, "root", realRetVal["LOGNAME"])
}

func TestVMEGetEnv(t *testing.T) {
	testScript := `
			var envvar1 = "USERNAME";
			var envvar2 = "DECKARDCAIN"
			var results1 = GetEnv(envvar1);
			var results2 = GetEnv(envvar2);
    `
	e := New("GetEnvVarTest")
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal1, err := e.VM.Get("results1")
	assert.Nil(t, err)
	retValAsString1, err := retVal1.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "root", retValAsString1)
	retVal2, err := e.VM.Get("results2")
	assert.Nil(t, err)
	retValAsString2, err := retVal2.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "", retValAsString2)
}
