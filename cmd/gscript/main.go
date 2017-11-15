package main

import (
	"github.com/gen0cide/gscript"
)

func main() {
	a := gscript.New()
	a.EnableLogging()
	a.CreateVM()
	a.VM.Run(gscript.DefaultScript)

}
