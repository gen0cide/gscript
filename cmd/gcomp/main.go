package main

import (
	"io/ioutil"

	"github.com/gen0cide/gscript/compiler"
	"github.com/gen0cide/gscript/logger"
)

func main() {
	tmplFile := "../../compiler/templates/vm_file.go.tmpl"
	tmplData, err := ioutil.ReadFile(tmplFile)
	if err != nil {
		panic(err)
	}
	logga := logger.NewStandardLogrusLogger(nil, "gcomp", false, false)
	c := compiler.NewWithDefault()
	c.SetLogger(logga)
	c.AddScript("./test.gs")
	err = c.Do()
	if err != nil {
		panic(err)
	}
	vms := c.GetVMs()
	v1 := vms[0]
	v1.GenerateFunctionKeys()
	err = v1.RenderVMBundle(string(tmplData))
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("../gcomp-output/main.go", v1.GenesisFile.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
	//spew.Dump(v1.Linker.Funcs["Test1"].GoReturns)
}
