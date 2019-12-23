package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/profile"

	"bitbucket.org/jeffgreenca/n2t-asm/internal/pkg/assembler"
	"bitbucket.org/jeffgreenca/n2t-asm/internal/pkg/lex"
	"bitbucket.org/jeffgreenca/n2t-asm/internal/pkg/parser"
)

func main() {
	defer profile.Start().Stop()

	r, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer r.Close()

	// tokenize from source file
	var tokens []lex.Token
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		tk, err := lex.Tokenize(scan.Text())
		check(err)
		tokens = append(tokens, tk...)
	}

	// parse tokens
	prog, err := parser.Parse(tokens)
	check(err)

	// assemble HACK machine instructions
	hack, err := assembler.Assemble(prog)
	check(err)

	fmt.Println(strings.Join(hack, "\n"))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
