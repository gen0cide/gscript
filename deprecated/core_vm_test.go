package engine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVMTimestamp(t *testing.T) {
	currTime := time.Now().Unix()

	testScript := `
        var test_time = Timestamp();
  `
	e := New("Timestamp")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("test_time")
	assert.Nil(t, err)
	assert.True(t, retVal.IsNumber())
	retValAsNumber, err := retVal.ToInteger()
	assert.Nil(t, err)
	assert.True(t, (retValAsNumber >= currTime))
}

func TestVMIsAWS(t *testing.T) {
	testScript := `
    var return_value = IsAWS();
  `
	e := New("AWSTest")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "false", retValAsString)
}
