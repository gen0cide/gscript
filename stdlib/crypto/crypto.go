package crypto

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
)

//GetMD5FromBytes takes a byte array and returns the md5 string of it
func GetMD5FromBytes(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

//GetMD5FromString takes a string and returns the md5 string of it
func GetMD5FromString(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

//GetSHA1FromBytes takes a byte array and returns the md5 string of it
func GetSHA1FromBytes(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

//GetSHA1FromString takes a byte array and returns the md5 string of it
func GetSHA1FromString(data string) string {
	hasher := sha1.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

//GetSHA256FromBytes takes a byte array and returns the md5 string of it
func GetSHA256FromBytes(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

//GetSHA256FromString takes a byte array and returns the md5 string of it
func GetSHA256FromString(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

//GenerateRSASSHKeyPair Generates a new rsa key pair the size of the arg, and returns the pub and priv key as strings
func GenerateRSASSHKeyPair(size int) (string, string, error) {
	reader := rand.Reader
	key, err := rsa.GenerateKey(reader, size)
	if err != nil {
		return "", "", err
	}
	privKeyDer := x509.MarshalPKCS1PrivateKey(key)
	privKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privKeyDer,
	}
	privKeyPem := string(pem.EncodeToMemory(&privKeyBlock))
	pubKey := key.PublicKey
	pubKeyDer, err := x509.MarshalPKIXPublicKey(&pubKey)
	if err != nil {
		return "", "", err
	}
	pubKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubKeyDer,
	}
	pubKeyPem := string(pem.EncodeToMemory(&pubKeyBlock))
	return pubKeyPem, privKeyPem, nil
}
