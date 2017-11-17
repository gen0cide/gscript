package gscript

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

type EmbeddedFile struct {
	SourcePath   string
	SourceURL    string
	Filename     string
	NameHash     string
	VariableDef  string
	Uncompressed []byte
	Compressed   []byte
}

func (e *EmbeddedFile) Compress() {
	if len(e.Uncompressed) > 0 {
		e.Compressed = BytesToCompressed(e.Uncompressed)
	}
}

func (e *EmbeddedFile) ResolveData() {
	d, _ := ioutil.ReadFile(e.SourcePath)
	e.Uncompressed = d
}

func (e *EmbeddedFile) ResolveFilename() {
	e.Filename = filepath.Base(e.SourcePath)
}

func (e *EmbeddedFile) ResolveVariableName() {
	e.NameHash = RandLowerAlphaString(18)
}

func (e *EmbeddedFile) GenerateVariableDef() {
	if len(e.Compressed) == 0 {
		return
	}
	function := ""
	fHeader := "var %s = func() []byte {\n\ts, _ := ioutil.ReadAll(flate.NewReader(bytes.NewBuffer([]byte{"
	function += fmt.Sprintf(fHeader, e.NameHash)
	for idx, b := range e.Compressed {
		if idx%15 == 0 {
			function += "\n\t\t"
		}
		function += fmt.Sprintf("0x%02x, ", b)
	}
	function += "\n\t"
	fCloser := "})))\n\treturn s\n}\n"
	function += fCloser
	e.VariableDef = function
}

func (e *EmbeddedFile) Embed() {
	e.ResolveFilename()
	e.ResolveVariableName()
	e.ResolveData()
	e.Compress()
	e.GenerateVariableDef()
}

func BytesToCompressed(b []byte) []byte {
	buf := new(bytes.Buffer)
	w, _ := flate.NewWriter(buf, flate.BestCompression)
	w.Write(b)
	w.Close()
	return buf.Bytes()
}

func CompressedToBytes(b []byte) []byte {
	buf, _ := ioutil.ReadAll(flate.NewReader(bytes.NewBuffer(b)))
	return buf
}

func RandLowerAlphaString(strlen int) string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	const chars = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return strings.ToUpper(string(result))
}
