package compiler

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// EmbeddedFile is an object that manages the lifecycle of resolving and translating
// embedded assets referenced in the Genesis VM into callable values that are
// embedded by the compiler
type EmbeddedFile struct {
	// local file path to the target file
	SourcePath string

	// url to the target file (optional)
	SourceURL string

	// basename of the referenced file
	Filename string

	// unique identifier that is used by the compiler to reference file contents
	NameHash string

	// uncompressed embedded file data
	Uncompressed []byte

	// compressed embedded file data
	Compressed []byte

	// temporary buffer used to generate the intermediate representation of the compressed data
	EmbedData bytes.Buffer
}

// Compress is used to gzip the embedded file's uncompressed data
func (e *EmbeddedFile) Compress() {
	if len(e.Uncompressed) > 0 {
		e.Compressed = BytesToCompressed(e.Uncompressed)
	}
}

// ResolveData is used to load the file's contents into the compiler
func (e *EmbeddedFile) ResolveData() {
	d, _ := ioutil.ReadFile(e.SourcePath)
	e.Uncompressed = d
}

// ResolveFilename gathers the base name of the file for pointer reference
// in the VM bundle's import map
func (e *EmbeddedFile) ResolveFilename() {
	e.Filename = filepath.Base(e.SourcePath)
}

// ResolveVariableName generates a unique identifier for this embed used by the compiler
func (e *EmbeddedFile) ResolveVariableName() {
	e.NameHash = RandUpperAlphaString(18)
}

// Embed performs all of the embed functions required to resolve and generate a compressed EmbeddedFile
func (e *EmbeddedFile) Embed() {
	e.ResolveFilename()
	e.ResolveVariableName()
	e.ResolveData()
	e.Compress()
	e.GenerateEmbedData()
}

// Data retrieves the current EmbedData's buffer as a string
func (e *EmbeddedFile) Data() string {
	return e.EmbedData.String()
}

// GenerateEmbedData enumerates the compressed embed and creates a byte slice representation of it
func (e *EmbeddedFile) GenerateEmbedData() {
	for _, b := range e.Compressed {
		e.EmbedData.WriteString(fmt.Sprintf("0x%02x, ", b))
	}
}

// BytesToCompressed compresses a byte array into a gzip compressed byte array
func BytesToCompressed(b []byte) []byte {
	buf := new(bytes.Buffer)
	w, _ := gzip.NewWriterLevel(buf, gzip.BestCompression)
	w.Write(b)
	w.Close()
	return buf.Bytes()
}

// CompressedToBytes uncompresses a byte array using gzip and returns it's original data
func CompressedToBytes(b []byte) []byte {
	r, _ := gzip.NewReader(bytes.NewBuffer(b))
	buf, _ := ioutil.ReadAll(r)
	return buf
}
