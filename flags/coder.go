package flags

import (
	"flag"
	"os"
)

type coderFlag struct {
	Mode  string
	Input string
}

func NewCoderFlag() *coderFlag {
	coderCmd := flag.NewFlagSet("coder", flag.ExitOnError)
	Mode := coderCmd.String("mode", "encode", "encode/decode")

	coderCmd.Parse(os.Args[2:])

	var input string
	if len(coderCmd.Args()) <= 0 {
		input = ""
	} else {
		input = coderCmd.Args()[0]
	}

	return &coderFlag{*Mode, input}
}
