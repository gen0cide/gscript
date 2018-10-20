package main

import (
	"os"
	"path/filepath"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
	"github.com/ahhh/gopkgs"
)

var (
	packages   = getPackages()
	errSkipDir = errors.New("skipping this directory")
)

func getPackages() map[string]gopkgs.Pkg {
	p, _ := gopkgs.List(gopkgs.Options{NoVendor: true})
	switched := map[string]gopkgs.Pkg{}
	for _, pkg := range p {
		switched[pkg.Dir] = pkg
	}
	return switched
}

func walkWithNative(p string) error {
	err := filepath.Walk(p, func(sp string, fi os.FileInfo, err error) error {
		abspath, _ := filepath.Abs(sp)
		if _, ok := packages[abspath]; ok {
			return filepath.SkipDir
		}
		return nil
	})
	return err
}

func walkWithLib(p string) error {
	err := godirwalk.Walk(p, &godirwalk.Options{
		Callback: func(osp string, de *godirwalk.Dirent) error {
			abspath, _ := filepath.Abs(osp)
			if _, ok := packages[abspath]; ok {
				return errSkipDir
			}
			return nil
		},
		ErrorCallback: func(d string, e error) godirwalk.ErrorAction {
			if errors.Cause(e) == errSkipDir {
				return godirwalk.SkipNode
			}
			return godirwalk.Halt
		},
		Unsorted: true,
	})
	return err
}

func main() {

}
