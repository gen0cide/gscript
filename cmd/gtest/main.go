package main

import (
	"fmt"

	"github.com/gen0cide/gscript/compiler/computil"
)

func main() {
	p, err := computil.ResolveGenesisPackageDir()
	if err != nil {
		panic(err)
	}
	fmt.Println(p)
}
