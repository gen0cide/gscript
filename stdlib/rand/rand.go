package rand

import (
	random "crypto/rand"
	"math/big"
)

//RandomInt Generates a random number between min and max.
func RandomInt(min, max int) int {
	r, _ := random.Int(random.Reader, big.NewInt(int64(max-min)))
	return int(r.Int64() + int64(min))
}

//GetAlphaNumericString Generates a random alpha numeric string of a specified length
func GetAlphaNumericString(strlen int) string {
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

//GetAlphaString generates a random alpha string of a specified length
func GetAlphaString(strlen int) string {
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

//GetAlphaNumericSpecialString generates a random alpha numeric and special char string of a specified length
func GetAlphaNumericSpecialString(strlen int) string {
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
