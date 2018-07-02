package main

import (
	"fmt"

	"github.com/gen0cide/gscript/compiler"
	"github.com/gen0cide/gscript/logger"
)

func main() {
	logga := logger.NewStandardLogrusLogger(nil, "gcomp", false, false)
	c := compiler.NewWithDefault()
	c.SetLogger(logga)
	c.AddScript("./test.gs")
	c.AddScript("./test.gs")
	err := c.Do()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Binary: %s\n", c.OutputFile)
}
