package gscript

import (
	"testing"
	"time"
  "fmt"
	"github.com/stretchr/testify/assert"
)

func TestVMMD5(t *testing.T) {
	testScript := `
    var hash_val = "helloworld";
    var return_value = MD5(hash_val);
  `
	// "helloworld" = fc5e038d38a57032085441e7fe7010b0

	e := New()
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
	file_2 := "/tmp/copyFile"
	testScript := fmt.Sprintf(`
    var file_1 = "/etc/passwd";
    var file_2 = "%s";
    var return_value = CopyFile(file_1, file_2);
  `, file_2)
	e := New()
	e.EnableLogging()
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

  finalDestination := fmt.Sprintf("File created at: %s", file_2)
	assert.Equal(t, finalDestination, retValAsString)
}

func TestVMTimestamp(t *testing.T) {
	currTime := time.Now().Unix()

	testScript := `
    var test_time = Timestamp();
  `
	e := New()
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
