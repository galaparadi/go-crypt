package util

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func GetSha256Hash(s string, l int) []byte {
	hash := sha256.New()
	hash.Write([]byte(s))

	return hash.Sum(nil)[:l]
}

func GetSha256HashBase64(s string, l int) string {
	hash := sha256.New()
	hash.Write([]byte(s))

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)[:l])
}

func GetSha256HashHexString(s string, l int) string {
	hash := sha256.New()
	hash.Write([]byte(s))

	return hex.EncodeToString(hash.Sum(nil)[:l])
}

func GetKeyString(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
