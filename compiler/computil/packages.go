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
	testFileRegexp = regexp.MustCompile(`.*_test\.go$`)

	// GenesisLibs is the name of the packages within the genesis standard library
	GenesisLibs = map[string]bool{
		"crypto":   true,
		"encoding": true,
		"exec":     true,
		"file":     true,
		"net":      true,
		"os":       true,
		"rand":     true,
		"requests": true,
		"time":     true,
	}

	// InstalledGoPackages holds a cache of all currently installed golang libraries
	InstalledGoPackages = GatherInstalledGoPackages()
)

func regexpForModule(mod ...string) *regexp.Regexp {
	return regexp.MustCompile(filepath.Join(append([]string{baseRegexpStr}, mod...)...))
}

// GatherInstalledGoPackages retrieves a list of all installed go packages in the context of current GOPATH and GOROOT
func GatherInstalledGoPackages() map[string]gopkgs.Pkg {
	goPackages, err := gopkgs.Packages(gopkgs.Options{NoVendor: true})
	if err != nil {
		panic(err)
	}
	return goPackages
}

// SourceFileIsTest determines if the given source file is named after the test convention
func SourceFileIsTest(src string) bool {
	return testFileRegexp.MatchString(src)
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
	for name, pkg := range InstalledGoPackages {
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
	dirMatch := regexpForModule("engine")
	for name, pkg := range InstalledGoPackages {
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
func ResolveStandardLibraryDir(pkg string) (*gopkgs.Pkg, error) {
	dirMatch := regexpForModule("stdlib", pkg)
	for name, gpkg := range InstalledGoPackages {
		if !dirMatch.MatchString(name) {
			continue
		}
		return &gpkg, nil
	}
	return nil, fmt.Errorf("could not locate standard library package %s", pkg)
}
