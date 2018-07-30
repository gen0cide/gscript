package exec

import (
	executer "os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteCommand(t *testing.T) {
	s := make([]interface{}, 0)
	pid, stdout, stderr, exitCode, err := ExecuteCommand("whoami", s)
	assert.Nil(t, err)
	assert.Equal(t, "", stderr, "should be no errors")
	assert.Equal(t, 0, exitCode, "should exit cleanly")
	assert.NotNil(t, pid)
	assert.NotEqual(t, "user", stdout, "can test the returned data")
}

func TestExecuteCommandAsync(t *testing.T) {
	s := make([]interface{}, 1)
	s = append(s, "10")
	var retcmd *executer.Cmd
	retcmd, err := ExecuteCommandAsync("sleep", s)
	assert.Nil(t, err)
	assert.NotNil(t, retcmd.Process)
}
