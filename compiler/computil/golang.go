package computil

import (
	"path/filepath"
	"regexp"
	"strings"
)

const (
	goosList   = `android darwin dragonfly freebsd js linux nacl netbsd openbsd plan9 solaris windows zos`
	goarchList = `386 amd64 amd64p32 arm armbe arm64 arm64be ppc64 ppc64le mips mipsle mips64 mips64le mips64p32 mips64p32le ppc riscv riscv64 s390 s390x sparc sparc64 wasm`
)

var (
	// GOOS maintains a cache of valid go OS's
	GOOS = GOOSList()

	// GOARCH maintains a cache of valid go architectures
	GOARCH = GOARCHList()

	archRegexp = regexp.MustCompile(strings.Join(GOARCH, "|"))
	osRegexp   = regexp.MustCompile(strings.Join(GOOS, "|"))
)

// GOOSList returns a slice of all possible go architectures
func GOOSList() []string {
	return strings.Split(goosList, " ")
}

// GOARCHList returns a slice of all possible go architectures
func GOARCHList() []string {
	return strings.Split(goarchList, " ")
}

// IsBuildSpecificFile tests a file to see if it's got platform specific naming to it's convention
func IsBuildSpecificFile(fn string) bool {
	base := filepath.Base(fn)
	if archRegexp.MatchString(base) {
		return true
	}
	if osRegexp.MatchString(base) {
		return true
	}
	return false
}
