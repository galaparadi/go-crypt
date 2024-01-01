package util

import "crypto/aes"

func MapChiperDataStartPosition(start uint32) uint32 {
	return (start / aes.BlockSize) * aes.BlockSize
}

func MapChiperDataEndPosition(end uint32) uint32 {
	return (((end / aes.BlockSize) + 1) * aes.BlockSize)
}

func MapIvStartPosition(start uint32) uint32 {
	if start < 16 {
		return 4
	} else {
		return (start/16)*16 + 4
	}
}
