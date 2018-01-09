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
	"github.com/mitchellh/go-ps"
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

//These two functions break compiling on windows
/*
func LocalFileWritable(path string) bool {
	return unix.Access(path, unix.W_OK) == nil
}

func LocalFileExecutable(path string) bool {
	return unix.Access(path, unix.X_OK) == nil
}
*/

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
	}
	return errors.New("The file dosn't exist to delete")
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
		return nil
	}
	err := LocalFileCreate(filename, bytes)
	if err != nil {
		return err
	}
	return nil
}

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
	var byteDst [20]byte
	for i := 0; i < n; i++ {
		byteDst[i] = a[i] ^ b[i]
	}
	return byteDst[:]
}

func LocalSystemInfo() ([]string, error) {
	var InfoDump []string
	gi := goInfo.GetInfo()
	InfoDump = append(InfoDump, fmt.Sprintf("GoOS: %s", gi.GoOS))
	InfoDump = append(InfoDump, fmt.Sprintf("Kernel: %s", gi.Kernel))
	InfoDump = append(InfoDump, fmt.Sprintf("Core: %s", gi.Core))
	InfoDump = append(InfoDump, fmt.Sprintf("Platform: %s", gi.Platform))
	InfoDump = append(InfoDump, fmt.Sprintf("OS: %s", gi.OS))
	InfoDump = append(InfoDump, fmt.Sprintf("Hostname: %s", gi.Hostname))
	InfoDump = append(InfoDump, fmt.Sprintf("CPUs: %v", gi.CPUs))
	if InfoDump != nil {
		return InfoDump, nil
	}
	return nil, errors.New("Failed to retrieve local system information")
}

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

func ForkExecuteCommand(c string, args ...string) (int, error) {
	cmd := exec.Command(c, args...)
	err := cmd.Start()
	if err != nil {
		return 0, err
	}
	pid := cmd.Process.Pid
	return pid, nil
}

func DNSQuestion(target, request string) (string, error) {
	if request == "A" {
		var stringAnswerArray []string
		answerPTR, err := net.LookupIP(target)
		if err != nil {
			return "failed", err
		}
		for _, answrPTR := range answerPTR {
			stringAnswerArray = append(stringAnswerArray, answrPTR.String())
		}
		stringAnswer := strings.Join(stringAnswerArray, "/n")
		return stringAnswer, nil
	} else if request == "TXT" {
		answerTXT, err := net.LookupTXT(target)
		if err != nil {
			return "failed", err
		}
		stringAnswer := strings.Join(answerTXT, "/n")
		return stringAnswer, nil
	} else if request == "PTR" {
		answerA, err := net.LookupAddr(target)
		if err != nil {
			return "failed", err
		}
		stringAnswer := strings.Join(answerA, "/n")
		return stringAnswer, nil
	} else if request == "MX" {
		var stringAnswerArray []string
		answerMX, err := net.LookupMX(target)
		if err != nil {
			return "failed", err
		}
		for _, answrMX := range answerMX {
			stringAnswerArray = append(stringAnswerArray, answrMX.Host)
		}
		stringAnswer := strings.Join(stringAnswerArray, "/n")
		return stringAnswer, nil
	} else if request == "NS" {
		var stringAnswerArray []string
		answerNS, err := net.LookupNS(target)
		if err != nil {
			return "failed", err
		}
		for _, answrNS := range answerNS {
			stringAnswerArray = append(stringAnswerArray, answrNS.Host)
		}
		stringAnswer := strings.Join(stringAnswerArray, "/n")
		return stringAnswer, nil
	} else if request == "CNAME" {
		answerCNAME, err := net.LookupCNAME(target)
		if err != nil {
			return "failed", err
		}
		return string(answerCNAME), nil
	} else {
		answerA, err := net.LookupHost(target)
		if err != nil {
			return "failed", err
		}
		stringAnswer := strings.Join(answerA, "/n")
		return stringAnswer, nil
	}
}

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

func TCPRead(ip, port string) ([]byte, error) {
	host := ip + ":" + port
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	return buffer, nil
}

func TCPWrite(writeData []byte, ip, port string) ([]byte, error) {
	host := ip + ":" + port
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	conn.Write(writeData)
	buffer2 := make([]byte, 1024)
	conn.Read(buffer2)
	return buffer2, nil
}

func UDPWrite(writeData []byte, ip, port string) error {
	host := ip + ":" + port
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.Write(writeData)
	return nil
}

//Function breaks compiling on windows
/*
func ProcExists1(pidBoi int) bool {
	killErr := syscall.Kill(pidBoi, syscall.Signal(0))
	procExistsBoi := (killErr == nil || killErr == syscall.EPERM)
	return procExistsBoi
}
*/

func ProcExists2(pidBoi int) bool {
	process, err := ps.FindProcess(pidBoi)
	if err == nil && process == nil {
		return false
	} else {
		return true
	}
}

func FindProcessPid(key string) (int, error) {
	pid := 0
	err := errors.New("Not found")
	ps, _ := ps.Processes()
	for i, _ := range ps {
		if ps[i].Executable() == key {
			pid = ps[i].Pid()
			err = nil
			break
		}
	}
	return pid, err
}

func StripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func DeobfuscateString(Data string) string {
	var ClearText string
	for i := 0; i < len(Data); i++ {
		ClearText += string(int(Data[i]) - 1)
	}
	return ClearText
}

func ObfuscateString(Data string) string {
	var ObfuscateText string
	for i := 0; i < len(Data); i++ {
		ObfuscateText += string(int(Data[i]) + 1)
	}
	return ObfuscateText
}

func RandString(strlen int) string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

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
