package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckForInUseTCP(t *testing.T) {
	inUse, err := CheckForInUseTCP(31337)
	assert.Nil(t, err)
	assert.NotNil(t, inUse)
	//assert.Equal(t, true, inUse, "will be true if running a listener")
	assert.Equal(t, false, inUse, "will be false when the port is not in use")
}

func TestCheckForInUseUDP(t *testing.T) {
	inUse, err := CheckForInUseUDP(2115)
	assert.Nil(t, err)
	assert.NotNil(t, inUse)
	assert.Equal(t, true, inUse, "will be true if something is listening on that udp port")
}
