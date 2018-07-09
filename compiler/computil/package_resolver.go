package computil

import (
	"errors"
	"regexp"

	"github.com/uudashr/gopkgs"
)

var (
	packageRegexp = regexp.MustCompile(`github\.com/gen0cide/gscript/engine`)
)

func ResolveEngineLibrary() (string, error) {
	targetDir := ""
	gpkgs, err := gopkgs.Packages(gopkgs.Options{NoVendor: true})
	if err != nil {
		return targetDir, err
	}
	for name, pkg := range gpkgs {
		if !packageRegexp.MatchString(name) {
			continue
		}
		targetDir = pkg.Dir
	}
	if targetDir == "" {
		return targetDir, errors.New("could not locate gscript engine library")
	}
	return targetDir, nil
}
