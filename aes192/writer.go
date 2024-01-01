package aes192

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"

	"galanov.tech/gocrypt/aes192/util"
)

type aes192Writer struct {
	mode cipher.BlockMode
	iv   []byte

	bufwriter *bufio.Writer
}

func NewAes192Writer(key []byte, iv []byte, plainSize uint32, bufwriter *bufio.Writer) *aes192Writer {
	block, _ := aes.NewCipher(key)
	mode := cipher.NewCBCEncrypter(block, iv)

	filebuffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(filebuffer, plainSize)
	bufwriter.Write(filebuffer)
	bufwriter.Flush()

	bufwriter.Write(iv)
	bufwriter.Flush()

	return &aes192Writer{mode, iv, bufwriter}
}

func (aw *aes192Writer) Write(plainchunk []byte) (int, error) {

	encbuffer := make([]byte, 16*1024)

	if len(plainchunk)%16 != 0 {
		aw.mode.CryptBlocks(encbuffer, util.PaddingPKCS5(plainchunk))
		aw.bufwriter.Write(encbuffer[:len(util.PaddingPKCS5(plainchunk))])
	} else {
		aw.mode.CryptBlocks(encbuffer, plainchunk)
		aw.bufwriter.Write(encbuffer[:len(plainchunk)])
	}

	err := aw.bufwriter.Flush()
	if err != nil {
		fmt.Println("error flush", err)
		return 0, err
	}

	return len(plainchunk), nil
}
