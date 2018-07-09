package rand

import (
	random "crypto/rand"
	"math/big"
)

//RandomInt Generates a random number between min and max.
func RandomInt(min, max int64) int {
	r, _ := random.Int(random.Reader, big.NewInt(int64(max-min)))
	return int(r.Int64() + min)
}

//GetAlphaNumericString Generates a random alpha numeric string of a specified length
func GetAlphaNumericString(strlen int64) string {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, strlen)
	for i := range result {
		val, err := random.Int(random.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		result[i] = chars[val.Int64()]
	}
	return string(result)
}

//GetAlphaString Generates a random alpha string of a specified length
func GetAlphaString(strlen int64) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, strlen)
	for i := range result {
		val, err := random.Int(random.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		result[i] = chars[val.Int64()]
	}
	return string(result)
}

//GetAlphaNumericSpecialString Generates a random alpha numeric and special char string of a specified length
func GetAlphaNumericSpecialString(strlen int64) string {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()"
	result := make([]byte, strlen)
	for i := range result {
		val, err := random.Int(random.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		result[i] = chars[val.Int64()]
	}
	return string(result)
}

//GetBool returns a random true or false
func GetBool() bool {
	val, _ := random.Int(random.Reader, big.NewInt(int64(2)))
	if int(val.Int64()) == int(0) {
		return false
	}
	return true
}
