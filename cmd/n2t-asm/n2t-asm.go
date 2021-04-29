package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jeffgreenca/n2t-asm/internal/pkg/assembler"
	"github.com/jeffgreenca/n2t-asm/internal/pkg/lex"
	"github.com/jeffgreenca/n2t-asm/internal/pkg/parser"
)

func usage() {
	fmt.Println("Provide asm via single filename argument or stdin")
}

func main() {
	var r io.Reader

	if len(os.Args) == 1 {
		r = os.Stdin
	} else if len(os.Args) == 2 {
		fr, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer fr.Close()
		r = fr
	} else {
		usage()
		os.Exit(1)
	}

	run(r)
}

func run(r io.Reader) {
	tokens, err := lex.Tokenize(r)
	if err != nil {
		panic(err)
	}

	program, err := parser.Parse(tokens)
	if err != nil {
		panic(err)
	}

	text, err := assembler.Assemble(program)
	if err != nil {
		panic(err)
	}

	for _, s := range text {
		fmt.Println(s)
	}
}
