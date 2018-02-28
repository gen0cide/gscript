package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVMExec(t *testing.T) {
	testCmd := ExecuteCommand("ls", "-lah")

	testScript := `
          var test_exec = Exec("ls", ["-lah"]);
      `
	e := New("Exec")
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

func TestVMForkExec(t *testing.T) {

	testScript := `
          var test_exec = ForkExec("nc", ["-l", "8080"]);
      `
	e := New("ForkExec")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("test_exec")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString)
}

func TestVMCanSudo(t *testing.T) {
	testScript := fmt.Sprintf(`
    var return_value2 = CanSudo();
  `)

	e := New("SudoTest")
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
	e.CreateVM()

	e.VM.Run(testScript)
	retVal2, err := e.VM.Get("return_value2")
	assert.Nil(t, err)
	retValAsString2, err := retVal2.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString2)
}

func TestVMCmdSuccessful(t *testing.T) {
	testScript := `
    var name = "ls";
    var arg = "-al"
    var return_value = CmdSuccessful(name, arg);
  `
	e := New("TestCmdSuccessful")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}
