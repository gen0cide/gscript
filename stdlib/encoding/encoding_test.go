package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBase64(t *testing.T) {
	newString, err := DecodeBase64("aGVsbG8gd29ybGQK")
	assert.Nil(t, err)
	assert.Equal(t, newString, "hello world\n", "Should be equal")
}

func TestEncodeBase64(t *testing.T) {
	newString := EncodeBase64("hello world")
	assert.Equal(t, newString, "aGVsbG8gd29ybGQ=", "should be equal")
}

func TestEncodeStringAsBytes(t *testing.T) {
	newBytes := EncodeStringAsBytes("hello world")
	assert.Equal(t, newBytes, ([]byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64}), "should be equal")
}

func TestEncodeBytesAsString(t *testing.T) {
	newString := EncodeBytesAsString([]byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64})
	assert.Equal(t, newString, "hello world", "should be equal")
}
