package assembler

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jeffgreenca/n2t-asm/internal/pkg/command"
)

func TestA(t *testing.T) {
	prog := command.Program{command.A{Address: 7, Static: true}}
	o, err := Assemble(prog)
	assert.NoError(t, err)
	assert.Equal(t, []string{"0000000000000111"}, o)
}
func TestC(t *testing.T) {
	prog := command.Program{command.C{
		D: command.Dest{M: true, D: true},
		C: "M+1",
		J: "",
	}}
	o, err := Assemble(prog)
	assert.NoError(t, err)
	assert.Equal(t, []string{"1111110111011000"}, o)
}
func TestCWithJump(t *testing.T) {
	prog := command.Program{command.C{
		D: command.Dest{M: true, D: true},
		C: "M+1",
		J: "JMP",
	}}
	o, err := Assemble(prog)
	assert.NoError(t, err)
	assert.Equal(t, []string{"1111110111011111"}, o)
}
