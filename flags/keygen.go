package flags

import (
	"crypto/rand"
	"flag"
	"log"
	"os"

	"galanov.tech/gocrypt/aes192/util"
)

type keygenFlag struct {
	KeygenMode  string
	KeyFileName string
	Key         []byte
}

func NewKeygenFlag() *keygenFlag {
	keygenCmd := flag.NewFlagSet("keygen", flag.ExitOnError)
	keyLength := keygenCmd.Int("length", 256, "key length")
	keyMode := keygenCmd.String("mode", "file", "key store mode")
	keySeed := keygenCmd.String("seed", "", "key seed")
	keyFileName := keygenCmd.String("file-name", "q", "key filename")
	keygenCmd.Parse(os.Args[2:])

	var key []byte
	if *keyLength == 256 {
		key = make([]byte, 32)
		rand.Read(key)
	} else if *keyLength == 192 {
		key = make([]byte, 24)
		rand.Read(key)
	} else {
		log.Fatal("key length not satisfied")
	}

	if len(*keySeed) > 0 {
		key = util.GetSha256Hash(*keySeed, *keyLength/8)
	}

	return &keygenFlag{*keyMode, *keyFileName, key}
}
