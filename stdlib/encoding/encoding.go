package encoding

import (
	"encoding/base64"
)

// DecodeBase64 decodes a base64 string and returns a string
func DecodeBase64(data string) (string, error) {
	valBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(valBytes), nil
}

//EncodeBase64 takes a string and turns it into a base64 string
func EncodeBase64(data string) string {
	return string(base64.StdEncoding.EncodeToString([]byte(data)))
}

//EncodeStringAsBytes takes a string and returns the bytes that make it
func EncodeStringAsBytes(data string) []byte {
	return []byte(data)
}

//EncodeBytesAsString takes a byte array and returns a string representation
func EncodeBytesAsString(data []byte) string {
	return string(data)
}
