package assembler

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jeffgreenca/n2t-asm/internal/pkg/command"
	"github.com/jeffgreenca/n2t-asm/internal/pkg/parser"
)

func TestHelloWorld(t *testing.T) {
	// t.Fatal("not implemented")
}

func TestA(t *testing.T) {
	prog := []parser.Command{{Type: command.TypeA, RealCmd: parser.CmdA{Address: 7, Final: true}}}
	o, err := Assemble(prog)
	assert.NoError(t, err)
	assert.Equal(t, []string{"0000000000000111"}, o)
}
func TestC(t *testing.T) {
	prog := []parser.Command{{Type: command.TypeC,
		RealCmd: parser.CmdC{
			D: parser.Dest{M: true, D: true},
			C: "M+1",
			J: "",
		},
	}}
	o, err := Assemble(prog)
	assert.NoError(t, err)
	assert.Equal(t, []string{"1111110111011000"}, o)
}
func TestCWithJump(t *testing.T) {
	prog := []parser.Command{{Type: command.TypeC,
		RealCmd: parser.CmdC{
			D: parser.Dest{M: true, D: true},
			C: "M+1",
			J: "JMP",
		},
	}}
	o, err := Assemble(prog)
	assert.NoError(t, err)
	assert.Equal(t, []string{"1111110111011111"}, o)
}
