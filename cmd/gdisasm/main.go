package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	gonet "net"
	"os"
	"sync"
	"time"

	"github.com/gen0cide/gscript/compiler/obfuscator"
)

//FuncDecl

func main() {
	names := []string{}
	_ = names
	fs := token.NewFileSet()
	dir := os.Args[1]
	pkgAST, _ := parser.ParseDir(fs, dir, nil, parser.ParseComments)
	exists := ast.PackageExports(pkgAST["engine"])
	if exists != true {
		panic("there were no exports")
	}
	for _, f := range pkgAST["engine"].Files {
		for _, d := range f.Decls {
			gd, ok := d.(*ast.GenDecl)
			if ok && gd.Tok == token.TYPE {
				for _, s := range gd.Specs {
					ts, ok := s.(*ast.TypeSpec)
					if !ok {
						continue
					}
					typeName := ts.Name.Name
					names = append(names, fmt.Sprintf("engine.%s", typeName))
					names = append(names, fmt.Sprintf("*engine.%s", typeName))
					// structType, ok := ts.Type.(*ast.StructType)
					// if !ok {
					// 	continue
					// }
					// for _, field := range structType.Fields.List {
					// 	for _, fieldName := range field.Names {
					// 		names = append(names, fieldName.Name)
					// 		names = append(names, fmt.Sprintf("*%s", fieldName.Name))
					// 		names = append(names, fmt.Sprintf("%s.%s", typeName, fieldName.Name))
					// 		names = append(names, fmt.Sprintf("*%s.%s", typeName, fieldName.Name))
					// 		names = append(names, fmt.Sprintf("engine.%s", fieldName.Name))
					// 		names = append(names, fmt.Sprintf("*engine.%s", fieldName.Name))
					// 	}
					// }
				}
			}
		}
	}
	for _, s := range names {
		fmt.Println(s)
	}
	obfuscator.GetEngineDecls()
}

//CheckForInUseUDP will send a UDP packet to the local port and see it gets a response or will timeout
func CheckForInUseUDP(port int) (bool, error) {
	timeout, err := time.ParseDuration("50ms")
	if err != nil {
		return false, err
	}

	conn, err := gonet.DialTimeout("udp", fmt.Sprintf("0.0.0.0:%d", port), timeout)
	if err != nil {
		return false, err
	}

	timeNow := time.Now()
	nextTick := timeNow.Add(timeout)

	conn.SetReadDeadline(nextTick)
	conn.SetWriteDeadline(nextTick)

	writeBuf := make([]byte, 1024, 1024)
	for i := 0; i < 1024; i++ {
		writeBuf[i] = 0x00
	}
	retSize, err := conn.Write(writeBuf)
	if err != nil {
		return false, err
	}

	readBuf := make([]byte, 1024, 1024)
	retSize, err = conn.Read(readBuf)
	if err != nil {
		opError, ok := err.(*gonet.OpError)
		if !ok {
			return false, err
		}
		if opError.Timeout() {
			return true, nil
		} else {
			return false, err
		}
	}

	if retSize > 0 {
		return true, nil
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(timeout)
		wg.Done()
	}()
	wg.Wait()

	return false, nil
}
