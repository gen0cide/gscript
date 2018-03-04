package compiler

import (
	"io/ioutil"
	"path/filepath"

	"github.com/robertkrimen/otto/parser"
)

func ValidateAST(source []byte) error {
	_, err := parser.ParseFile(nil, "", source, 0)
	return err
}

// ParseMacros normalizes the import files into localized assets
func (c *Compiler) ParseMacros(vm *VMBundle) []string {
	imports := []string{}
	script, err := ioutil.ReadFile(vm.ScriptFile)
	if err != nil {
		c.Logger.WithField("file", filepath.Base(vm.ScriptFile)).Fatalf("Error reading genesis script: %s", err.Error())
	}

	macroList := ParseMacros(string(script), c.Logger.WithField("file", filepath.Base(vm.ScriptFile)))
	if macroList == nil {
		c.Logger.WithField("file", filepath.Base(vm.ScriptFile)).Fatalf("Could not parse macros for script!")
	}

	vm.Timeout = macroList.Timeout
	vm.Priority = macroList.Priority

	for _, i := range macroList.LocalFiles {
		imports = append(imports, i)
	}

	for _, i := range macroList.RemoteFiles {
		imports = append(imports, i)
	}

	return imports
}
