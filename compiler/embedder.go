package compiler

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gen0cide/gscript/compiler/computil"
)

const (
	// Preload is the static preload that is injected into the runtime before execution
	Preload = `
		function StringToByteArray(s) {
			var data = [];
			for (var i = 0; i < s.length; i++) {
				data.push(s.charCodeAt(i));
			}
			return data;
		}

		function ByteArrayToString(a) {
			return String.fromCharCode.apply(String, a);
		}

		function Dump(obj) {
			return "\n" + JSON.stringify(obj, null, 2);
		}

		function BeforeDeploy() {
			return true;
		}

		function Deploy() {
			return false;
		}

		function AfterDeploy() {
			return true;
		}

		function OnError() {
			return false;
		}

		function Sleep(seconds) {
			var start = new Date().getTime();
			for (var i = 0; i < 1e7; i++) {
				if ((new Date().getTime() - start) > (seconds * 1000)) {
					break;
				}
			}
		}

		function DebugConsole() {
			return true;
		}
	`
)

// EmbeddedFile is an object that manages the lifecycle of resolving and translating
// embedded assets referenced in the Genesis VM into callable values that are
// embedded by the compiler
type EmbeddedFile struct {
	// local file path to the target file
	SourcePath string

	// url to the target file (optional)
	CachedPath string

	// the composite ($ID_$OrigName) filename of the referenced file
	Filename string

	// the original basename of the referenced file
	OrigName string

	// unique identifier that is used by the compiler to reference file contents
	ID string

	// unique identifier that is used by the compiler to swizzle the file decoding and decrypting into a function pointer
	FuncID string

	// uncompressed embedded file data
	Uncompressed []byte

	// compressed embedded file data
	Compressed []byte

	// used to AES encrypt the embedded assets
	EncryptionKey []byte

	// temporary buffer used to generate the intermediate representation of the compressed data
	EmbedData *bytes.Buffer
}

// NewEmbeddedFile takes a path on the local file system and returns an EmbeddedFile object reference
func NewEmbeddedFile(source string) (*EmbeddedFile, error) {
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return nil, err
	}
	absPath, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}
	id := computil.RandUpperAlphaString(18)
	ef := &EmbeddedFile{
		SourcePath: absPath,
		OrigName:   filepath.Base(source),
		Filename:   fmt.Sprintf("%s_%s", id, filepath.Base(source)),
		ID:         id,
	}
	return ef, nil
}

// CacheFile attempts to copy the files original location (e.SourcePath) to the
// destination cacheDir provided as an argument to this function call
func (e *EmbeddedFile) CacheFile(cacheDir string) error {
	dstAbsPath := filepath.Join(cacheDir, e.Filename)
	fileData, err := ioutil.ReadFile(e.SourcePath)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dstAbsPath, fileData, 0644)
	if err != nil {
		return err
	}
	e.CachedPath = dstAbsPath
	return nil
}

// // Compress is used to gzip the embedded file's uncompressed data
// func (e *EmbeddedFile) Compress() {
// 	if len(e.Uncompressed) > 0 {
// 		e.Compressed = BytesToCompressed(e.Uncompressed)
// 	}
// }

// // ResolveData is used to load the file's contents into the compiler
// func (e *EmbeddedFile) ResolveData() {
// 	d, _ := ioutil.ReadFile(e.SourcePath)
// 	e.Uncompressed = d
// }

// // ResolveFilename gathers the base name of the file for pointer reference
// // in the VM bundle's import map
// func (e *EmbeddedFile) ResolveFilename() {
// 	e.Filename = filepath.Base(e.SourcePath)
// }

// // ResolveVariableName generates a unique identifier for this embed used by the compiler
// func (e *EmbeddedFile) ResolveVariableName() {
// 	e.NameHash = RandLowerAlphaString(18)
// }

// Embed performs all of the embed functions required to resolve and generate a compressed EmbeddedFile
func (e *EmbeddedFile) Embed() {
	// e.ResolveFilename()
	// e.ResolveVariableName()
	// e.ResolveData()
	// e.Compress()
	// e.GenerateEmbedData()
}

// Data retrieves the current EmbedData's buffer as a string
func (e *EmbeddedFile) Data() string {
	return e.EmbedData.String()

}

// GenerateEmbedData enumerates the compressed embed and creates a byte slice representation of it
func (e *EmbeddedFile) GenerateEmbedData() error {
	pipeBuf := new(bytes.Buffer)
	ioReader, err := os.Open(e.CachedPath)
	if err != nil {
		return err
	}
	w, err := gzip.NewWriterLevel(pipeBuf, gzip.BestCompression)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, ioReader)
	if err != nil {
		return err
	}
	ioReader.Close()
	w.Close()
	block, err := aes.NewCipher(e.EncryptionKey)
	if err != nil {
		return err
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])
	e.EmbedData = new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, e.EmbedData)
	encWriter := &cipher.StreamWriter{S: stream, W: encoder}
	if _, err := io.Copy(encWriter, pipeBuf); err != nil {
		return err
	}
	encoder.Close()
	return nil
}

// func (e *EmbeddedFile) GenerateEmbedData() error {
// 	ioReader, ioWriter := io.Pipe()
// 	ioReader, err := os.Open(e.CachedPath)
// 	if err != nil {
// 		return err
// 	}
// 	w := gzip.NewWriterLevel(e.EmbedData, gzip.BestCompression)
// 	_, err := io.Copy(w, ioReader)
// 	if err != nil {
// 		return err
// 	}
// 	for _, b := range e.Compressed {
// 		e.EmbedData.WriteString(fmt.Sprintf("0x%02x, ", b))
// 	}
// 	e.Compressed = []byte{}
// 	e.Uncompressed = []
// }

// // BytesToCompressed compresses a byte array into a gzip compressed byte array
// func BytesToCompressed(b []byte) []byte {
// 	buf := new(bytes.Buffer)
// 	w, _ := gzip.NewWriterLevel(buf, gzip.BestCompression)
// 	w.Write(b)
// 	w.Close()
// 	return buf.Bytes()
// }

// // CompressedToBytes uncompresses a byte array using gzip and returns it's original data
// func CompressedToBytes(b []byte) []byte {
// 	r, _ := gzip.NewReader(bytes.NewBuffer(b))
// 	buf, _ := ioutil.ReadAll(r)
// 	return buf
// }

// ExampleDecodeEmbed is a reference implementation of how embedded assets should be unpacked
// inside of a genesis engine
func ExampleDecodeEmbed(b64encoded string, key string) []byte {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	db1 := new(bytes.Buffer)
	db2 := new(bytes.Buffer)
	src := bytes.NewReader([]byte(b64encoded))
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])
	decoder := base64.NewDecoder(base64.StdEncoding, src)
	encReader := &cipher.StreamReader{S: stream, R: decoder}
	if _, err := io.Copy(db1, encReader); err != nil {
		panic(err)
	}
	gzr, err := gzip.NewReader(db1)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(db2, gzr)
	if err != nil {
		panic(err)
	}
	gzr.Close()
	return db2.Bytes()
}
