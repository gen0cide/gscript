package computil

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	lowercaseAlphaNumericChars = "abcdefghijklmnopqrstuvwxyz0123456789"
	lowerAlphaChars            = "abcdefghijklmnopqrstuvwxyz"
	mixedAlphaNumericChars     = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	mixedAlphaChars            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	upperAlphaChars            = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// RandAlphaNumericString creates a random lowercase alpha-numeric string of a given length
func RandAlphaNumericString(strlen int) string {
	result := make([]byte, strlen)
	for i := range result {
		var val *big.Int
		var err error
		if i == 0 {
			val, err = rand.Int(rand.Reader, big.NewInt(int64(len(lowerAlphaChars))))
		} else {
			val, err = rand.Int(rand.Reader, big.NewInt(int64(len(lowercaseAlphaNumericChars))))
		}
		if err != nil {
			panic(err)
		}
		result[i] = lowercaseAlphaNumericChars[val.Int64()]
	}
	return string(result)
}

// RandUpperAlphaNumericString creates a random uppercase alpha-numeric string of a given length
func RandUpperAlphaNumericString(strlen int) string {
	return strings.ToUpper(RandAlphaNumericString(strlen))
}

// RandomInt returns a random integer between a min and max value
func RandomInt(min, max int) int {
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(r.Int64()) + min
}

// RandMixedAlphaNumericString creates a random mixed-case alpha-numeric string of a given length
func RandMixedAlphaNumericString(n int) string {
	b := make([]byte, n)
	for i := range b {
		var val *big.Int
		var err error
		if i == 0 {
			val, err = rand.Int(rand.Reader, big.NewInt(int64(len(mixedAlphaChars))))
		} else {
			val, err = rand.Int(rand.Reader, big.NewInt(int64(len(mixedAlphaNumericChars))))
		}
		if err != nil {
			panic(err)
		}
		b[i] = mixedAlphaNumericChars[val.Int64()]
	}
	return string(b)
}

// RandUpperAlphaString creates a random uppercase alpha-only string of a given length
func RandUpperAlphaString(strlen int) string {
	return strings.ToUpper(RandLowerAlphaString(strlen))
}

// RandLowerAlphaString creates a random lowercase alpha-only string of a given length
func RandLowerAlphaString(strlen int) string {
	result := make([]byte, strlen)
	for i := range result {
		val, err := rand.Int(rand.Reader, big.NewInt(int64(len(lowerAlphaChars))))
		if err != nil {
			panic(err)
		}
		result[i] = lowerAlphaChars[val.Int64()]
	}
	return string(result)
}
