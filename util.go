package gscript

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/matishsiao/goInfo"
)

func CalledBy() string {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "Unknown"
	}
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "N/A"
	}
	return fun.Name()
}

func LocalFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func LocalDirCreate(path string) error {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		return err
	}
	return nil
}

func LocalDirRemoveAll(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	err = os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}

func LocalFileDelete(path string) error {
	if LocalFileExists(path) {
		err := os.Remove(path)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("The file dosn't exist to delete")
	}
}

func LocalFileCreate(path string, bytes []byte) error {
	if LocalFileExists(path) {
		return errors.New("The file to create already exists so we won't overwite it")
	}
	err := ioutil.WriteFile(path, bytes, 0700)
	if err != nil {
		return err
	}
	return nil
}

// LocalFileAppendBytes adds bytes to the end of filename's path.
func LocalFileAppendBytes(filename string, bytes []byte) error {
	if LocalFileExists(filename) {
		fileInfo, err := os.Stat(filename)
		if err != nil {
			return err
		}
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, fileInfo.Mode())
		if err != nil {
			return err
		}
		if _, err = file.Write(bytes); err != nil {
			return err
		}
		file.Close()
		// Appened the bytes w/o error
		return nil
	} else {
		err := LocalFileCreate(filename, bytes)
		if err != nil {
			return err
		}
		// Created a new file w/o error
		return nil
	}
}

// LocalFileAppendString adds input as strings to the end of filename's path.
func LocalFileAppendString(input, filename string) error {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_APPEND, fileInfo.Mode())
	if err != nil {
		return err
	}
	if _, err = file.WriteString(input); err != nil {
		return err
	}
	file.Close()
	return nil
}

// Replace will replace all instances of match with replace in file.
func LocalFileReplace(file, match, replacement string) error {
	if LocalFileExists(file) {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return err
		}
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		lines := strings.Split(string(contents), "\n")
		for index, line := range lines {
			if strings.Contains(line, match) {
				lines[index] = strings.Replace(line, match, replacement, 10)
			}
		}

		ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), fileInfo.Mode())
		return nil
	} else {
		return errors.New("The file to read does not exist")
	}
}

// ReplaceMulti will replace all instances of possible matches with replacement in file.
func LocalFileReplaceMulti(file string, matches []string, replacement string) error {
	if LocalFileExists(file) {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return err
		}
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		lines := strings.Split(string(contents), "\n")
		for index, line := range lines {
			for _, match := range matches {
				if strings.Contains(line, match) {
					lines[index] = replacement
				}
			}
		}
		ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), fileInfo.Mode())
		return nil
	} else {
		return errors.New("The file to read does not exist")
	}
}

// LocalReadFile takes a file path and returns the byte array of the file there
func LocalFileRead(path string) ([]byte, error) {
	if LocalFileExists(path) {
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return dat, nil
	}
	return nil, errors.New("The file to read does not exist")
}

func XorFiles(file1 string, file2 string, outPut string) error {
	dat1, err := ioutil.ReadFile(file1)
	if err != nil {
		return err
	}
	dat2, err := ioutil.ReadFile(file2)
	if err != nil {
		return err
	}
	dat3 := XorBytes(dat1[:], dat2[:])
	err = LocalFileCreate(outPut, dat3[:])
	if err != nil {
		return err
	}
	return nil
}

func XorBytes(a []byte, b []byte) []byte {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	var byte_dst [20]byte
	for i := 0; i < n; i++ {
		byte_dst[i] = a[i] ^ b[i]
	}
	return byte_dst[:]
}

func LocalSystemInfo() ([]string, error) {
	var InfoDump []string
	gi := goInfo.GetInfo()
	// later when we define these objects we can just set the values to the object vs the string slice
	InfoDump = append(InfoDump, fmt.Sprintf("GoOS: %s", gi.GoOS))
	InfoDump = append(InfoDump, fmt.Sprintf("Kernel: %s", gi.Kernel))
	InfoDump = append(InfoDump, fmt.Sprintf("Core: %s", gi.Core))
	InfoDump = append(InfoDump, fmt.Sprintf("Platform: %s", gi.Platform))
	InfoDump = append(InfoDump, fmt.Sprintf("OS: %s", gi.OS))
	InfoDump = append(InfoDump, fmt.Sprintf("Hostname: %s", gi.Hostname))
	InfoDump = append(InfoDump, fmt.Sprintf("CPUs: %v", gi.CPUs))
	//gi.VarDump()
	if InfoDump != nil {
		return InfoDump, nil
	} else {
		return nil, errors.New("Failed to retrieve local system information")
	}
}

// ExecuteCommand function
func ExecuteCommand(c string, args ...string) VMExecResponse {
	cmd := exec.Command(c, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	respObj := VMExecResponse{
		Stdout: strings.Split(stdout.String(), "\n"),
		Stderr: strings.Split(stderr.String(), "\n"),
		PID:    cmd.Process.Pid,
	}
	if err != nil {
		respObj.ErrorMsg = err.Error()
		respObj.Success = false
	} else {
		respObj.Success = true
	}
	return respObj
}

func DNSQuestion(target, request string) (string, error) {
	if request == "A" {
		var stringAnswerArray []string
		answerPTR, err := net.LookupIP(target)
		if err != nil {
			return "failed", err
		}
		// Formating output and debug strings here
		for _, answrPTR := range answerPTR {
			stringAnswerArray = append(stringAnswerArray, answrPTR.String())
		}
		stringAnswer := strings.Join(stringAnswerArray, "/n")
		//fmt.Println(stringAnswer)
		return stringAnswer, nil
	} else if request == "TXT" {
		answerTXT, err := net.LookupTXT(target)
		if err != nil {
			return "failed", err
		}
		// Formating output and debug strings here
		stringAnswer := strings.Join(answerTXT, "/n")
		//fmt.Println(stringAnswer)
		return stringAnswer, nil
	} else if request == "PTR" {
		answerA, err := net.LookupAddr(target)
		if err != nil {
			return "failed", err
		}
		// Formating output and debug strings here
		stringAnswer := strings.Join(answerA, "/n")
		//fmt.Println(stringAnswer)
		return stringAnswer, nil
	} else if request == "MX" {
		var stringAnswerArray []string
		answerMX, err := net.LookupMX(target)
		if err != nil {
			return "failed", err
		}
		// Formating output and debug strings here
		for _, answrMX := range answerMX {
			stringAnswerArray = append(stringAnswerArray, answrMX.Host)
		}
		stringAnswer := strings.Join(stringAnswerArray, "/n")
		//fmt.Println(stringAnswer)
		return stringAnswer, nil
	} else if request == "NS" {
		var stringAnswerArray []string
		answerNS, err := net.LookupNS(target)
		if err != nil {
			return "failed", err
		}
		// Formating output and debug strings here
		for _, answrNS := range answerNS {
			stringAnswerArray = append(stringAnswerArray, answrNS.Host)
		}
		stringAnswer := strings.Join(stringAnswerArray, "/n")
		//fmt.Println(stringAnswer)
		return stringAnswer, nil
	} else if request == "CNAME" {
		answerCNAME, err := net.LookupCNAME(target)
		if err != nil {
			return "failed", err
		}
		// Formating output and debug strings here
		//fmt.Println(string(answerCNAME))
		return string(answerCNAME), nil
	} else {
		answerA, err := net.LookupHost(target)
		if err != nil {
			return "failed", err
		}
		// Formating output and debug strings here
		stringAnswer := strings.Join(answerA, "/n")
		//fmt.Println(stringAnswer)
		return stringAnswer, nil
	}
}

// HTTPGetFile takes a url and returns a status code, a byte slice of the file there, and an error
func HTTPGetFile(url string) (int, []byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, nil, err
	}
	respCode := resp.StatusCode
	pageData, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return respCode, pageData, nil
}

// TCPRead returns a byte slice from the initial connection and an error
func TCPRead(ip, port string) ([]byte, error) {
	host := ip+":"+port
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	return buffer, nil
}

// TCPWrite returns a byte slice which is the response to our write and an error
func TCPWrite(writeData []byte, ip, port string) ([]byte, error) {
	host := ip+":"+port
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	// Write our response
	conn.Write(writeData)
	return buffer, nil
	// Read the reply
	buffer2 := make([]byte, 1024)
	conn.Read(buffer2)
	return buffer2, nil
}

// StripSpaces will remove the spaces from a single string and return the new string
func StripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

// DeobfuscateString from https://github.com/SaturnsVoid/GoBot2/blob/7d6609cd49f006f5aee76a4ffd97eb25d12a1a9b/components/Cryptography.go#L44
func DeobfuscateString(Data string) string {
	var ClearText string
	for i := 0; i < len(Data); i++ {
		ClearText += string(int(Data[i]) - 1)
	}
	return ClearText
}

// ObfuscateString from https://github.com/SaturnsVoid/GoBot2/blob/7d6609cd49f006f5aee76a4ffd97eb25d12a1a9b/components/Cryptography.go#L52
func ObfuscateString(Data string) string {
	var ObfuscateText string
	for i := 0; i < len(Data); i++ {
		ObfuscateText += string(int(Data[i]) + 1)
	}
	return ObfuscateText
}

// RandString returns a string the length of strlen
func RandString(strlen int) string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

// RandomInt returns an int inbetween min and max.
func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func LocalCopyFile(src, dst string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}

func GetLocalIPs() []string {
	addresses := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return addresses
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addresses = append(addresses, ipnet.IP.String())
			}
		}
	}
	return addresses
}
