package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func client() {
	fw, err := os.Create("/tmp/logger")
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetLevel(logrus.DebugLevel)
	logger.Out = fw

	secDur, err := time.ParseDuration("1s")
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(secDur)
		currTime := time.Now()
		logger.Infof("Current Time: %v", currTime)
	}
}

func spawner() {
	signal.Ignore(syscall.SIGHUP)
	currentBin, err := os.Executable()
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(currentBin, "client")
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	fmt.Println("spawned... sleeping 10s to see what happens...")
	secDur, err := time.ParseDuration("30s")
	time.Sleep(secDur)
	if err != nil {
		panic(err)
	}
	fmt.Println("finished spawner sleeping, moving on")
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("must supply a run type")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "spawner":
		spawner()
	case "client":
		client()
	default:
		fmt.Println("not a valid run type (spawner or client)")
		os.Exit(1)
	}

	// timeout, err := time.ParseDuration("50ms")
	// if err != nil {
	// 	panic(err)
	// }

	// conn, err := net.DialTimeout("udp", "0.0.0.0:53", timeout)
	// if err != nil {
	// 	panic(err)
	// }

	// timeNow := time.Now()
	// nextTick := timeNow.Add(timeout)

	// conn.SetReadDeadline(nextTick)
	// conn.SetWriteDeadline(nextTick)

	// writeBuf := make([]byte, 1024, 1024)
	// for i := 0; i < 1024; i++ {
	// 	writeBuf[i] = 0x00
	// }
	// retSize, err := conn.Write(writeBuf)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Sent Size: %d\n", retSize)

	// readBuf := make([]byte, 1024, 1024)

	// retSize, err = conn.Read(readBuf)
	// if err != nil {
	// 	opError, ok := err.(*net.OpError)
	// 	if !ok {
	// 		fmt.Printf("error is type of %T\n", err)
	// 		panic(err)
	// 	}
	// 	if opError.Timeout() {
	// 		fmt.Printf("conn read timed out, but conn was not closed, open just hanging")
	// 		return
	// 	} else {
	// 		panic(err)
	// 	}
	// }

	// fmt.Printf("Return Size: %d\n", retSize)

	// var wg sync.WaitGroup

	// wg.Add(1)
	// go func() {
	// 	time.Sleep(timeout)
	// 	wg.Done()
	// }()
	// wg.Wait()

	// fmt.Printf("Why did we get here?\n")

}
