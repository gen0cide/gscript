package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
)

var (
	script = `// a new script
	//go_import:github.com/tdewolff/minify as minify
	//go_import:github.com/tdewolff/minify/js as js
	var foo = "bar";

	function poop() {
		return null;
	}

	function Deploy() {
		poop();
		minfier = minify.New()
		minifier.AddFunc("text/javascript", js.Minify)
		return true;
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

// type walker struct {
// 	source string
// 	shift  file.Idx
// }

// func (w *walker) Exit(n ast.Node) {
// 	return
// }

// func (w *walker) Enter(n ast.Node) ast.Visitor {
// 	spew.Dump(n)
// 	fmt.Println("==============================")
// 	return w
// }

type Embed struct {
	CachedPath string
	EmbedData  *bytes.Buffer
}

func (e *Embed) GenerateEmbedData() error {
	ioReader, err := os.Open(e.CachedPath)
	if err != nil {
		return err
	}
	buf1 := new(bytes.Buffer)
	w, err := gzip.NewWriterLevel(buf1, gzip.BestCompression)
	if err != nil {
		return err
	}
	bw, err := io.Copy(w, ioReader)
	fmt.Printf("bw1 = %d\n", bw)
	ioReader.Close()
	w.Close()
	encoder := base64.NewEncoder(base64.StdEncoding, e.EmbedData)
	buf1.WriteTo(encoder)
	encoder.Close()
	return nil

	// ioReader, err := os.Open(e.CachedPath)
	// if err != nil {
	// 	return err
	// }
	// defer ioReader.Close()
	// w, err := gzip.NewWriterLevel(encoder, gzip.BestCompression)
	// if err != nil {
	// 	return err
	// }
	// defer w.Close()
	// bytesWritten, err := io.Copy(w, ioReader)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("Bytes written during encode: %d\n", bytesWritten)
	// defer encoder.Close()
	// return nil
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
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

	// callableFuncs := map[string]bool{}
	// prog, err := parser.ParseFile(nil, "test", script, 2)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, s := range prog.Body {
	// 	funcStmt, ok := s.(*ast.FunctionStatement)
	// 	if !ok {
	// 		continue
	// 	}
	// 	fnLabel := funcStmt.Function.Name.Name
	// 	if callables[fnLabel] != "" {
	// 		fmt.Printf("Found Entrypoint: %s()\n", fnLabel)
	// 		callableFuncs[fnLabel] = true
	// 	}
	// }

	// if len(callableFuncs) == 3 {
	// 	fmt.Println("valid legacy script")
	// 	return
	// }

	// if len(callableFuncs) == 1 && callableFuncs["Deploy"] == true {
	// 	fmt.Println("valid v2 script")
	// 	return
	// }

	// fmt.Println("not a valid legacy script!")
	// return

	e := Embed{
		CachedPath: os.Args[1],
		EmbedData:  new(bytes.Buffer),
	}
	err := e.GenerateEmbedData()
	if err != nil {
		panic(err)
	}

	fmt.Printf("BUF SIZE: %d\n", e.EmbedData.Len())

	eData := e.EmbedData.String()

	fmt.Printf("STRING: %s\n", eData)
	fmt.Printf("LEN: %d\n", len(eData))

	db := new(bytes.Buffer)
	sr := bytes.NewReader([]byte(eData))
	decoder := base64.NewDecoder(base64.StdEncoding, sr)
	gz, err := gzip.NewReader(decoder)
	if err != nil {
		panic(err)
	}
	bw2, err := io.Copy(db, gz)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bytes written during decode: %d\n", bw2)
	gz.Close()
	spew.Dump(db)

	// data, err := base64.StdEncoding.DecodeString(eData)
	// if err != nil {
	// 	spew.Dump(e.EmbedData.Bytes())
	// 	panic(err)
	// }
	// gz, err := gzip.NewReader(bytes.NewReader(data))
	// if err != nil {
	// 	panic(err)
	// }
	// s, err := ioutil.ReadAll(gz)

	// if err != nil {
	// 	panic(err)
	// }

	fmt.Printf("DECODED!\n\nOUTPUT:\n%s\n", db.String())

	// m := minify.New()
	// m.AddFunc("text/javascript", js.Minify)
	// miniVersion := new(bytes.Buffer)
	// r := bytes.NewReader([]byte(script))
	// if err := m.Minify("text/javascript", miniVersion, r); err != nil {
	// 	fmt.Printf("minification error\n")
	// 	panic(err)
	// }
	// miniFinal := miniVersion.Bytes()
	// fmt.Printf("Original Size: %d bytes\n", len([]byte(script)))
	// fmt.Printf("Minified Size: %d bytes\n", len(miniFinal))
	// fmt.Printf("NEW SOURCE:\n%s\n\n", string(miniFinal))

}
