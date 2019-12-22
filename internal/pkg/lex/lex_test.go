package lex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	progAdd = `
	// Computes R0 = 2 + 3  (R0 refers to RAM[0])
		@2
		D=A
		@3
		D=D+A
		@0
		M=D
	(END)
		@END
		0;JMP
`
)

func TestClean(t *testing.T) {
	testCases := map[string]string{
		"D=M+1    ":            "D=M+1",
		" (somelabel)":         "(somelabel)",
		"//some comment":       "",
		"   @100   // comment": "@100",
	}
	for k, v := range testCases {
		actual := clean(k)
		assert.Equal(t, v, actual)
	}
}

func TestTokenizeTypeC(t *testing.T) {
	testCases := map[string][]Token{
		"D=M+1;JNE": {
			{Type: LOCATION, Value: "D"},
			{Type: ASSIGN, Value: "="},
			{Type: LOCATION, Value: "M"},
			{Type: OPERATOR, Value: "+"},
			{Type: NUMBER, Value: "1"},
			{Type: JUMP, Value: "JNE"},
		},
	}

	for k, v := range testCases {
		actual, err := Tokenize(k)
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}
