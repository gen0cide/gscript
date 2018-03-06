package compiler

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"strings"
)

type EmbeddedFile struct {
	SourcePath   string
	SourceURL    string
	Filename     string
	NameHash     string
	VariableDef  string
	Uncompressed []byte
	Compressed   []byte
	EmbedData    bytes.Buffer
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
	e.NameHash = RandUpperAlphaString(18)
}

func (e *EmbeddedFile) Embed() {
	e.ResolveFilename()
	e.ResolveVariableName()
	e.ResolveData()
	e.Compress()
	e.GenerateEmbedData()
}

func (e *EmbeddedFile) Data() string {
	return e.EmbedData.String()
}

func (e *EmbeddedFile) GenerateEmbedData() {
	for _, b := range e.Compressed {
		e.EmbedData.WriteString(fmt.Sprintf("\\x%02x", b))
	}
}

func BytesToCompressed(b []byte) []byte {
	buf := new(bytes.Buffer)
	w, _ := gzip.NewWriterLevel(buf, gzip.BestCompression)
	w.Write(b)
	w.Close()
	return buf.Bytes()
}

func CompressedToBytes(b []byte) []byte {
	r, _ := gzip.NewReader(bytes.NewBuffer(b))
	buf, _ := ioutil.ReadAll(r)
	return buf
}

func RandUpperAlphaString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, strlen)
	for i := range result {
		val, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		result[i] = chars[val.Int64()]
	}
	return strings.ToUpper(string(result))
}
