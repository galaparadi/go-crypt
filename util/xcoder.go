package util

import (
	"encoding/hex"
)

func XEncode(src string) string {
	bytesSrc := []byte(src)
	for i, v := range bytesSrc {
		bytesSrc[i] = ^v
	}

	return hex.EncodeToString(bytesSrc)
}

func XDecode(src string) string {
	dst, _ := hex.DecodeString(src)

	for i, v := range dst {
		dst[i] = ^v
	}

	return string(dst[:])
}
