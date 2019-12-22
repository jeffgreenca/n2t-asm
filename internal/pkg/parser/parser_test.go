package parser

import (
	"testing"

	"bitbucket.org/jeffgreenca/n2t-asm/internal/pkg/lex"
	"github.com/stretchr/testify/assert"
)

var (
	tokenizedCommandC1 = []lex.Token{
		{Type: lex.LOCATION, Value: "D"},
		{Type: lex.ASSIGN, Value: "="},
		{Type: lex.LOCATION, Value: "M"},
		{Type: lex.OPERATOR, Value: "+"},
		{Type: lex.NUMBER, Value: "1"},
		{Type: lex.JUMP, Value: "JNE"},
		{Type: lex.END},
	}

	tokenizedCommandC2 = []lex.Token{
		{Type: lex.NUMBER, Value: "0"},
		{Type: lex.JUMP, Value: "JMP"},
		{Type: lex.END},
	}
)

func TestC1(t *testing.T) {
	_, err := Parse(tokenizedCommandC1)
	assert.NoError(t, err)

	expected := CmdC{
		D: Dest{D: true, A: false, M: false},
		C: "M+1",
		J: "JNE",
	}

	assert.Equal(t, C_COMMAND, cmdType)
	assert.Equal(t, expected, cmdC)
}

func TestC2(t *testing.T) {
	_, err := Parse(tokenizedCommandC2)
	assert.NoError(t, err)

	expected := CmdC{
		D: Dest{D: false, A: false, M: false},
		C: "0",
		J: "JMP",
	}

	assert.Equal(t, C_COMMAND, cmdType)
	assert.Equal(t, expected, cmdC)
}
