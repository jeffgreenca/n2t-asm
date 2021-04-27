package main

import (
	"fmt"
	"os"

	"github.com/jeffgreenca/n2t-asm/internal/pkg/assembler"
	"github.com/jeffgreenca/n2t-asm/internal/pkg/lex"
	"github.com/jeffgreenca/n2t-asm/internal/pkg/parser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: n2t-asm <file>")
		return
	}

	r, err := os.Open(os.Args[1])
	defer r.Close()
	if err != nil {
		panic(err)
	}

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
