package aes192

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
)

type aes192Reader struct {
	mode     cipher.BlockMode
	iv       []byte
	filesize uint32
	start    uint32
	end      uint32

	readenchunk uint32
	firstchunk  bool
	encreader   *bufio.Reader
}

func NewAes192Reader(key []byte, iv []byte, filesize uint32, start uint32, end uint32, encreader *bufio.Reader) *aes192Reader {
	readenchunk := uint32(0)
	firstchunk := true

	block, _ := aes.NewCipher(key)
	mode := cipher.NewCBCDecrypter(block, iv)

	if start > end {
		panic("start is larger than end")
	}

	if end >= filesize {
		panic("end larger than filesize")
	}

	return &aes192Reader{mode, iv, filesize, start, end, readenchunk, firstchunk, encreader}
}

func (ar *aes192Reader) Read(result []byte) (int, error) {

	chunkbuffer := make([]byte, 16*1024)
	n, err := ar.encreader.Read(chunkbuffer)
	if err != nil {
		return 0, err
	}

	ar.readenchunk += uint32(n)
	plainbuffer := make([]byte, 16*1024)
	ar.mode.CryptBlocks(plainbuffer, chunkbuffer)

	if ar.firstchunk && (ar.readenchunk > (ar.end - ar.start)) {
		copy(result, plainbuffer[ar.start%16:uint32(n)-(16-(ar.end%16))+1])
		ar.firstchunk = false
		return (n - (16 - (int(ar.end) % 16)) + 1) - (int(ar.start) % 16), nil
	}

	if ar.firstchunk {
		copy(result, plainbuffer[ar.start%16:n])
		ar.firstchunk = false
		return n - (int(ar.start) % 16), nil
	}

	if ar.readenchunk > (ar.end - ar.start) {
		copy(result, plainbuffer[:uint32(n)-(16-(ar.end%16))+1])
		return n - (16 - (int(ar.end) % 16)) + 1, nil //
	}

	copy(result, plainbuffer[:n])
	return n, nil
}
