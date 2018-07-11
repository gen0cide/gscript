package computil

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"

	"github.com/uudashr/gopkgs"
)

var (
	baseImportPath = `github.com/gen0cide/gscript`
	baseRegexpStr  = `github\.com/gen0cide/gscript`
	baseRegexp     = regexp.MustCompile(baseRegexpStr)

	// GenesisLibs is the name of the packages within the genesis standard library
	GenesisLibs = []string{
		"asset",
		"crypto",
		"encoding",
		"exec",
		"file",
		"net",
		"os",
		"rand",
		"requests",
		"time",
	}
)

func regexpForModule(mod ...string) *regexp.Regexp {
	return regexp.MustCompile(filepath.Join(append([]string{baseRegexpStr}, mod...)...))
}

// InstalledGoPackages retrieves a list of all installed go packages in the context of current GOPATH and GOROOT
func InstalledGoPackages() (map[string]gopkgs.Pkg, error) {
	return gopkgs.Packages(gopkgs.Options{NoVendor: true})
}

// ResolveGoPath attempts to resolve the current user's GOPATH
func ResolveGoPath() string {
	gp := os.Getenv("GOPATH")
	if gp != "" {
		return gp
	}
	u, err := user.Current()
	if err != nil {
		// really shouldn't happen
		panic(err)
	}
	return filepath.Join(u.HomeDir, "go")
}

// ResolveGenesisPackageDir attempts to resolve the base directory for the genesis package
func ResolveGenesisPackageDir() (targetDir string, err error) {
	guess := filepath.Join(ResolveGoPath(), baseImportPath)
	if _, ok := os.Stat(guess); ok == nil {
		return guess, nil
	}
	gpkgs, err := InstalledGoPackages()
	if err != nil {
		return targetDir, err
	}
	for name, pkg := range gpkgs {
		if !baseRegexp.MatchString(name) {
			continue
		}
		targetDir = pkg.Dir
	}
	if targetDir == "" {
		return targetDir, errors.New("could not locate the base genesis package")
	}
	return targetDir, nil
}

// ResolveEngineDir attempts to resolve the absolute path of the genesis engine directory
func ResolveEngineDir() (targetDir string, err error) {
	gpkgs, err := InstalledGoPackages()
	if err != nil {
		return targetDir, err
	}
	dirMatch := regexpForModule("engine")
	for name, pkg := range gpkgs {
		if !dirMatch.MatchString(name) {
			continue
		}
		targetDir = pkg.Dir
	}
	if targetDir == "" {
		return targetDir, fmt.Errorf("coult not locate the genesis engine directory")
	}
	return targetDir, nil
}

// ResolveStandardLibraryDir attempts to resolve the absolute path of the specified standard library package
func ResolveStandardLibraryDir(pkg string) (string, error) {
	targetDir := ""
	dirMatch := regexpForModule("stdlib", pkg)
	gpkgs, err := gopkgs.Packages(gopkgs.Options{NoVendor: true})
	if err != nil {
		return targetDir, err
	}
	for name, pkg := range gpkgs {
		if !dirMatch.MatchString(name) {
			continue
		}
		targetDir = pkg.Dir
	}
	if targetDir == "" {
		return targetDir, fmt.Errorf("coult not locate standard library package %s", pkg)
	}
	return targetDir, nil
}
