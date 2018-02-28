package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCPUStats(t *testing.T) {
	//resultz := CPUStats()
	testScript := `
          var results = CPUStats();
      `
	e := New("CPUStats")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("results")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString)
}

func TestVMEnvVars(t *testing.T) {
	testScript := `
          var results = EnvVars();
      `
	e := New("EnvVarsTest")
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

func TestVMLocalUserExists(t *testing.T) {
	testScript := `
    var user = "root";
    var return_value = LocalUserExists(user);
  `
	e := New("LocaUserTest")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMProcExistsWithName(t *testing.T) {
	testScript := `
    var name = "notbash";
    var return_value = ProcExistsWithName(name);
  `
	e := New("TestProcExists")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "false", retValAsString)
}
