package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMD5FromBytes(t *testing.T) {
	hashVal := GetMD5FromBytes([]byte("hello world"))
	assert.NotNil(t, hashVal)
	assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", hashVal, "should be equal")
}

func TestGetMD5FromString(t *testing.T) {
	hashVal := GetMD5FromString("hello world")
	assert.NotNil(t, hashVal)
	assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", hashVal, "should be equal")
}

func TestGetSHA1FromBytes(t *testing.T) {
	hashVal := GetSHA1FromBytes([]byte("hello world"))
	assert.NotNil(t, hashVal)
	assert.Equal(t, "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed", hashVal, "should be equal")

}

func TestGetSHA1FromString(t *testing.T) {
	hashVal := GetSHA1FromString("hello world")
	assert.NotNil(t, hashVal)
	assert.Equal(t, "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed", hashVal, "should be equal")

}

func TestGetSHA256FromBytes(t *testing.T) {
	hashVal := GetSHA256FromBytes([]byte("hello world"))
	assert.NotNil(t, hashVal)
	assert.Equal(t, "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", hashVal, "should be equal")

}

func TestGetSHA256FromString(t *testing.T) {
	hashVal := GetSHA256FromString("hello world")
	assert.NotNil(t, hashVal)
	assert.Equal(t, "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", hashVal, "should be equal")
}

func TestGenerateRSASSHKeyPair(t *testing.T) {
	pubKey, privKey, err := GenerateRSASSHKeyPair(1024)
	assert.Nil(t, err)
	assert.NotNil(t, pubKey)
	assert.NotNil(t, privKey)
	//assert.Equal(t, "", pubKey, "testing")
}
