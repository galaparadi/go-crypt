package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"galanov.tech/gocrypt/aes192"
	"galanov.tech/gocrypt/aes192/util"
	"galanov.tech/gocrypt/flags"
	coder "galanov.tech/gocrypt/util"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("arg not sattisfied")
		return
	}

	programmode := os.Args[1]

	iv := make([]byte, 16)
	rand.Read(iv)

	switch programmode {
	case "encrypt":
		encryptflag := flags.NewEncryptFlag()
		fmt.Println("encrypting....")

		fsPlain, err := os.Open(*encryptflag.PlainPath)
		if err != nil {
			log.Fatal("error accessing file : ", err)
		}

		fsChiper, err := os.Create(*encryptflag.ChiperPath)
		if err != nil {
			log.Fatal("Error createing file : ", err)
		}
		// fmt.Println("key", encryptflag.Key)
		// fmt.Println("iv", iv)

		filesizebuffer := make([]byte, 4)
		fsinfo, _ := fsPlain.Stat()
		binary.LittleEndian.PutUint32(filesizebuffer, uint32(fsinfo.Size()))

		chiperbufio := bufio.NewWriter(fsChiper)
		plainbufio := bufio.NewReader(fsPlain)

		chiperwriter := aes192.NewAes192Writer(encryptflag.Key, iv, uint32(fsinfo.Size()), chiperbufio)

		plainchunk := make([]byte, 16*1024)
		for {
			n, err := plainbufio.Read(plainchunk)
			if err != nil {
				break
			}
			chiperwriter.Write(plainchunk[:n])
		}

		log.Println("encrypt success")
	case "decrypt":
		decryptflag := flags.NewDecryptFlag()
		fmt.Println("decrypting....")

		fsEnc, err := os.Open(*decryptflag.ChiperPath)
		if err != nil {
			log.Fatal("error akses file", err)
		}

		fsPlain, err := os.Create(*decryptflag.PlainPath)
		if err != nil {
			log.Fatal("error akses file", err)
		}

		fileSizeBuffer := make([]byte, 4)
		fsEnc.Read(fileSizeBuffer)
		filesize := binary.LittleEndian.Uint32(fileSizeBuffer)

		newStart := util.MapChiperDataStartPosition(0)
		offset := util.MapChiperDataEndPosition(filesize - 1)

		ivBuffer := make([]byte, 16)
		fsEnc.Read(ivBuffer)

		partialReader := io.NewSectionReader(fsEnc, int64(newStart)+4+16, int64(offset-uint32(newStart)))
		bReader := bufio.NewReader(partialReader)
		bWriter := bufio.NewWriter(fsPlain)

		// fmt.Println("key", decryptflag.Key)
		// fmt.Println("iv", iv)

		chiperreader := aes192.NewAes192Reader(decryptflag.Key,
			ivBuffer,
			filesize,
			0,
			filesize-1,
			bReader)

		for {
			buffer := make([]byte, 16*1024)
			n, err := chiperreader.Read(buffer)
			if err != nil {
				break
			}
			bWriter.Write(buffer[:n])
			bWriter.Flush()
		}

		log.Println("decrypt success")
		log.Println("the file size", filesize)
	case "check":
		fmt.Println("mengcompare chiper data dengan plain data")

		decryptflag := flags.NewDecryptFlag()
		fmt.Println("decrypting....")

		fsEnc, err := os.Open(*decryptflag.ChiperPath)
		if err != nil {
			log.Fatal("error akses file", err)
		}

		fsPlain, err := os.Open(*decryptflag.PlainPath)
		if err != nil {
			log.Fatal("error akses file", err)
		}

		fileSizeBuffer := make([]byte, 4)
		fsEnc.Read(fileSizeBuffer)
		filesize := binary.LittleEndian.Uint32(fileSizeBuffer)

		plainsize, _ := fsPlain.Stat()
		if filesize != uint32(plainsize.Size()) {
			log.Fatalln("filesize not same")
		}

		newStart := util.MapChiperDataStartPosition(0)
		offset := util.MapChiperDataEndPosition(filesize - 1)

		ivBuffer := make([]byte, 16)
		fsEnc.Read(ivBuffer)

		partialReader := io.NewSectionReader(fsEnc, int64(newStart)+4+16, int64(offset-uint32(newStart)))
		bChiper := bufio.NewReader(partialReader)
		bPlain := bufio.NewReader(fsPlain)

		chiperreader := aes192.NewAes192Reader(decryptflag.Key,
			ivBuffer,
			filesize,
			0,
			filesize-1,
			bChiper)

		for {
			bufferChiper := make([]byte, 16*1024)
			bufferPlain := make([]byte, 16*1024)

			nChiper, err := chiperreader.Read(bufferChiper)
			nPlain, err := bPlain.Read(bufferPlain)

			if nChiper != nPlain {
				log.Fatalln("Test failed. File not same (n value not same)")
			}

			if bytes.Compare(bufferChiper[:nChiper], bufferPlain[:nPlain]) != 0 {
				fmt.Println(bufferChiper[:nChiper])
				fmt.Println(bufferPlain[:nPlain])
				log.Fatalln("Test failed. File not same (bytes value not same)")
			}

			if err != nil {
				break
			}
		}

		log.Println("test success")
	case "keygen":
		keygenflag := flags.NewKeygenFlag()

		if keygenflag.KeygenMode == "file" {
			fsKey, err := os.Create("q") //TODO: check if file already exist
			if err != nil {
				log.Fatal("Error createing file : ", err)
			}
			fsKey.Write(keygenflag.Key)
			fsKey.Close()
			fmt.Println("key telah tersimpan dengan nama file q. Harap simpan key dengan baik")
			fmt.Println(util.GetKeyString(keygenflag.Key))
		} else {
			fmt.Println(util.GetKeyString(keygenflag.Key))
			fmt.Println(util.GetSha256HashHexString(string(keygenflag.Key), 32))
		}
	case "coder":
		coderflag := flags.NewCoderFlag()

		if coderflag.Mode == "encode" {
			fmt.Println("result :", coder.XEncode(coderflag.Input))
		} else if coderflag.Mode == "decode" {
			fmt.Println("result :", coder.XDecode(coderflag.Input))
		} else if coderflag.Mode == "dirdecode" {
			files, _ := os.ReadDir("./")
			for _, v := range files {
				fmt.Print(coder.XDecode(v.Name()))
				fmt.Print(" ----> ")
				fmt.Print(v.Name())
				fmt.Println("")
			}
		}
	case "-h":
		fmt.Println("show help")
	}
}
