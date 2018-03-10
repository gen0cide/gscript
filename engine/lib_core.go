package engine

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"time"
	"unicode"
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// XorBytes - XOR two byte arrays together.
//
// Package
//
// core
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  XorBytes(aByteArray, bByteArray)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * aByteArray ([]byte)
//  * bByteArray ([]byte)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.value ([]byte)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = XorBytes(aByteArray, bByteArray);
//  // obj.value
//
func (e *Engine) XorBytes(a []byte, b []byte) []byte {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	var byteDst [20]byte
	for i := 0; i < n; i++ {
		byteDst[i] = a[i] ^ b[i]
	}
	return byteDst[:]
}

// StripSpaces - Strip any unicode characters out of a string.
//
// Package
//
// core
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  StripSpaces(str)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * str (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.value (string)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = StripSpaces(str);
//  // obj.value
//
func (e *Engine) StripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

// DeobfuscateString - Basic string deobfuscator function.
//
// Package
//
// core
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  DeobfuscateString(str)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * str (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.value (string)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = DeobfuscateString(str);
//  // obj.value
//
func (e *Engine) DeobfuscateString(Data string) string {
	var ClearText string
	for i := 0; i < len(Data); i++ {
		ClearText += string(int(Data[i]) - 1)
	}
	return ClearText
}

// ObfuscateString - Basic string obfuscator function.
//
// Package
//
// core
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  ObfuscateString(str)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * str (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.value (string)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = ObfuscateString(str);
//  // obj.value
//
func (e *Engine) ObfuscateString(Data string) string {
	var ObfuscateText string
	for i := 0; i < len(Data); i++ {
		ObfuscateText += string(int(Data[i]) + 1)
	}
	return ObfuscateText
}

// RandomString - Generates a random alpha numeric string of a specified length.
//
// Package
//
// core
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  RandomString(strlen)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * strlen (int64)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.value (string)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = RandomString(strlen);
//  // obj.value
//
func (e *Engine) RandomString(strlen int64) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		val, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		result[i] = chars[val.Int64()]
	}
	return string(result)
}

// RandomInt - Generates a random number between min and max.
//
// Package
//
// core
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  RandomInt(min, max)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * min (int64)
//  * max (int64)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.value (int)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = RandomInt(min, max);
//  // obj.value
//
func (e *Engine) RandomInt(min, max int64) int {
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(r.Int64() + min)
}

// RandomMixedCaseString - Generates a random mixed case alpha numeric string of a specified length.
//
// Package
//
// core
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  RandomMixedCaseString(strlen)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * strlen (int64)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.value (string)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = RandomMixedCaseString(strlen);
//  // obj.value
//
func (e *Engine) RandomMixedCaseString(n int64) string {
	b := make([]rune, n)
	for i := range b {
		val, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			panic(err)
		}
		b[i] = letterRunes[val.Int64()]
	}
	return string(b)
}

// Asset - Retrieves a packed asset from the VM embedded file store.
//
// Package
//
// core
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  Asset(assetName)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * assetName (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.fileData ([]byte)
//  * obj.err (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = Asset(assetName);
//  // obj.fileData
//  // obj.err
//
func (e *Engine) Asset(filename string) ([]byte, error) {
	if dataFunc, ok := e.Imports[filename]; ok {
		byteData := dataFunc()
		return byteData, nil
	}
	e.Logger.WithField("trace", "true").Errorf("Asset File Not Found: %s", filename)
	err := errors.New("Asset not found: " + filename)
	return []byte{}, err
}

// Timestamp - Get the system's current timestamp in epoch format.
//
// Package
//
// core
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  Timestamp()
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.value (int64)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = Timestamp();
//  // obj.value
//
func (e *Engine) Timestamp() int64 {
	return time.Now().Unix()
}

// Halt - Stop the current VM from continuing execution.
//
// Package
//
// core
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  Halt()
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.value (bool)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = Halt();
//  // obj.value
//
func (e *Engine) Halt() bool {
	e.Halted = true
	e.VM.Interrupt <- func() {
		panic(errTimeout)
	}
	return true
}

// MD5 - Perform an MD5() hash on data.
//
// Package
//
// core
//
// Author
//
// - gen0cide (https://github.com/gen0cide)
//
// Javascript
//
// Here is the Javascript method signature:
//  MD5(data)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * data ([]byte)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.value (string)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = MD5(data);
//  // obj.value
//
func (e *Engine) MD5(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}
