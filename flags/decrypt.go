package flags

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
)

type decryptFlag struct {
	PlainPath  *string
	ChiperPath *string
	Key        []byte
}

func NewDecryptFlag() *decryptFlag {
	decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
	keyMode := decryptCmd.String("key-mode", "file", "key mode")

	decryptCmd.Parse(os.Args[2:])
	chiperPath := decryptCmd.Args()[0]
	plainPath := decryptCmd.Args()[1]

	var key []byte
	if *keyMode == "file" {
		q, err := os.ReadFile("q")
		if err != nil {
			fmt.Println(err)
		}
		key = q
	} else if *keyMode == "env" {
		//TODO: ambil key dari env variable
		keyString256 := os.Getenv("GO_CRY_256")
		keyString192 := os.Getenv("GO_CRY_192")

		if len(keyString192) > 0 {
			q, err := base64.StdEncoding.DecodeString(keyString192)
			if err != nil {
				log.Fatal("error accessing key 192 from env")
			}
			key = q
		}
		if len(keyString256) > 0 {
			q, err := base64.StdEncoding.DecodeString(keyString256)
			if err != nil {
				log.Fatal("error accessing key 256 from env")
			}
			key = q
		}
	}

	return &decryptFlag{&plainPath, &chiperPath, key}
}
