package util

import "bytes"

func PaddingPKCS5(chunk []byte) []byte {
	padByte := ((16 * 1024) - len(chunk)) % 16
	padBytes := bytes.Repeat([]byte{byte(padByte)}, padByte)
	return append(chunk, padBytes...)
}
