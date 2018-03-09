package engine

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
	"unicode"
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

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

func (e *Engine) StripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func (e *Engine) DeobfuscateString(Data string) string {
	var ClearText string
	for i := 0; i < len(Data); i++ {
		ClearText += string(int(Data[i]) - 1)
	}
	return ClearText
}

func (e *Engine) ObfuscateString(Data string) string {
	var ObfuscateText string
	for i := 0; i < len(Data); i++ {
		ObfuscateText += string(int(Data[i]) + 1)
	}
	return ObfuscateText
}

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

func (e *Engine) RandomInt(min, max int64) int {
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(r.Int64() + min)
}

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

func (e *Engine) Asset(filename string) ([]byte, error) {
	if dataFunc, ok := e.Imports[filename]; ok {
		byteData := dataFunc()
		return byteData, nil
	}
	e.Logger.WithField("trace", "true").Errorf("Asset File Not Found: %s", filename)
	err := errors.New("Asset not found: " + filename)
	return []byte{}, err
}
