package engine

import (
	"errors"
	"fmt"

	"github.com/matishsiao/goInfo"
	ps "github.com/mitchellh/go-ps"
)

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

func GetHostname() (string, error) {
	gi := goInfo.GetInfo()
	hostname := gi.Hostname
	if hostname != "" {
		return hostname, nil
	}
	return "", errors.New("Failed to retrieve local hostname")
}

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
