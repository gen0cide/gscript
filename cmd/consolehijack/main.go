package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/file"
	"github.com/robertkrimen/otto/parser"
)

var (
	script = `
	var foo = "bar";

	function poop() {
		return null;
	}

	function Deploy() {
		poop();
		return "happy";
	}

	function AfterDeploy() {
		return true;
	}
	`

	callables = map[string]string{
		"BeforeDeploy": "no",
		"Deploy":       "yes",
		"AfterDeploy":  "no",
	}
)

type walker struct {
	source string
	shift  file.Idx
}

func (w *walker) Exit(n ast.Node) {
	return
}

func (w *walker) Enter(n ast.Node) ast.Visitor {
	spew.Dump(n)
	fmt.Println("==============================")
	return w
}

func main() {
	// logger := logger.NewStandardLogrusLogger(nil, "testhijack", false, false)
	// e := engine.New("testhijack", "RANDOMVMID", 30, "Deploy")
	// e.SetLogger(logger)
	// err := e.LoadScript("test.gs", []byte(script))
	// if err != nil {
	// 	logger.Error("Not continuing since script load failed.")
	// 	return
	// }
	// val, err := e.Exec("Deploy")
	// if err != nil {
	// 	panic(err)
	// }
	// spew.Dump(val)

	callableFuncs := map[string]bool{}
	prog, err := parser.ParseFile(nil, "test", script, 2)
	if err != nil {
		panic(err)
	}
	for _, s := range prog.Body {
		funcStmt, ok := s.(*ast.FunctionStatement)
		if !ok {
			continue
		}
		fnLabel := funcStmt.Function.Name.Name
		if callables[fnLabel] != "" {
			fmt.Printf("Found Entrypoint: %s()\n", fnLabel)
			callableFuncs[fnLabel] = true
		}
	}

	if len(callableFuncs) == 3 {
		fmt.Println("valid legacy script")
		return
	}

	if len(callableFuncs) == 1 && callableFuncs["Deploy"] == true {
		fmt.Println("valid v2 script")
		return
	}

	fmt.Println("not a valid legacy script!")
	return

}
