package net

import (
	"fmt"
	gonet "net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckForInUseTCP(t *testing.T) {
	t.Helper()
	ln, err := gonet.Listen("tcp", "127.0.0.1:31337")
	assert.Nil(t, err)
	assert.NotNil(t, ln)
	var server gonet.Conn
	_ = server
	go func() {
		defer ln.Close()
		server, err = ln.Accept()
	}()
	inUse, err := CheckForInUseTCP(31337)
	assert.Nil(t, err)
	assert.NotNil(t, inUse)
	assert.Equal(t, true, inUse, "will be true if running a listener")
	notInUse, err := CheckForInUseTCP(31338)
	assert.Nil(t, err)
	assert.NotNil(t, notInUse)
	assert.Equal(t, false, notInUse, "will be false when the port is not in use")
}

func TestCheckForInUseUDP(t *testing.T) {
	t.Helper()
	addr, err := gonet.ResolveUDPAddr("udp", "127.0.0.1:29201")
	assert.Nil(t, err)
	assert.NotNil(t, addr)
	serverConn, err := gonet.ListenUDP("udp", addr)
	assert.Nil(t, err)
	assert.NotNil(t, serverConn)
	go func() {
		defer serverConn.Close()

		buf := make([]byte, 1024)
		for {
			n, addr, err := serverConn.ReadFromUDP(buf)
			fmt.Println("Received ", string(buf[0:n]), " from ", addr)

			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
	}()
	inUse, err := CheckForInUseUDP(29201)
	assert.Nil(t, err)
	assert.NotNil(t, inUse)
	assert.Equal(t, true, inUse, "will be true if something is listening on that udp port")
	notInUse, err := CheckForInUseUDP(67)
	assert.NotNil(t, err)
	assert.NotNil(t, notInUse)
	assert.Equal(t, false, notInUse, "will be true if something is not in use")
}
