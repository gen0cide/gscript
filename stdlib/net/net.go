package net

import (
	"fmt"
	gonet "net"
	"strconv"
	"strings"
	"sync"
	"time"
)

//CheckForInUseTCP is a function that checks all local IPv4 interfaces for to see if something is listening on the specified TCP port will timeout after 3 seconds
func CheckForInUseTCP(port int) (bool, error) {
	timeout, err := time.ParseDuration("50ms")
	portString := strconv.Itoa(port)
	addr := "0.0.0.0:" + portString
	fmt.Println(addr)
	conn, err := gonet.DialTimeout("tcp", addr, timeout)
	if err != nil {
		if strings.ContainsAny(err.Error(), "connection refused") {
			return false, nil
		}
		return false, err
	}
	if conn != nil {
		conn.Close()
		return true, nil
	}
	return false, nil
}

//CheckForInUseUDP will send a UDP packet to the local port and see it gets a response or will timeout
func CheckForInUseUDP(port int) (bool, error) {
	timeout, err := time.ParseDuration("50ms")
	if err != nil {
		return false, err
	}
	conn, err := gonet.DialTimeout("udp", fmt.Sprintf("0.0.0.0:%d", port), timeout)
	if err != nil {
		return false, err
	}
	timeNow := time.Now()
	nextTick := timeNow.Add(timeout)
	conn.SetReadDeadline(nextTick)
	conn.SetWriteDeadline(nextTick)
	writeBuf := make([]byte, 1024, 1024)
	for i := 0; i < 1024; i++ {
		writeBuf[i] = 0x00
	}
	retSize, err := conn.Write(writeBuf)
	if err != nil {
		return false, err
	}
	readBuf := make([]byte, 1024, 1024)
	retSize, err = conn.Read(readBuf)
	if err != nil {
		opError, ok := err.(*gonet.OpError)
		if !ok {
			return false, err
		}
		if opError.Timeout() {
			return true, nil
		}
		return false, err
	}
	if retSize > 0 {
		return true, nil
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(timeout)
		wg.Done()
	}()
	wg.Wait()
	return false, nil
}
