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
