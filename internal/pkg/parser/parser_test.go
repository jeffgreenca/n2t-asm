package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jeffgreenca/n2t-asm/internal/pkg/command"
	"github.com/jeffgreenca/n2t-asm/internal/pkg/token"
)

func TestCommandTypeL(t *testing.T) {
	type TestCaseL struct {
		tokens   []token.Token
		expected CmdL
	}

	testCases := []TestCaseL{
		{
			tokens: []token.Token{
				{Type: token.LABEL, Value: "("},
				{Type: token.SYMBOL, Value: "foobar"},
				{Type: token.END},
			},
			expected: CmdL{Symbol: "foobar"},
		},
	}

	for _, c := range testCases {
		program, err := Parse(c.tokens)
		assert.NoError(t, err)

		assert.Equal(t, command.L, program[0].Type)
		assert.Equal(t, c.expected, program[0].RealCmd)
	}
}

func TestCommandTypeA(t *testing.T) {
	type TestCaseA struct {
		tokens   []token.Token
		expected CmdA
	}

	testCases := []TestCaseA{
		{
			tokens: []token.Token{
				{Type: token.AT, Value: "@"},
				{Type: token.ADDRESS, Value: "1024"},
				{Type: token.END},
			},
			expected: CmdA{Address: 1024, Final: true, Symbol: ""},
		},
		{
			tokens: []token.Token{
				{Type: token.AT, Value: "@"},
				{Type: token.SYMBOL, Value: "foo"},
				{Type: token.END},
			},
			expected: CmdA{Address: 0, Final: false, Symbol: "foo"},
		},
		{
			tokens: []token.Token{
				{Type: token.AT, Value: "@"},
				{Type: token.SYMBOL, Value: "i"},
				{Type: token.END},
			},
			expected: CmdA{Address: 0, Final: false, Symbol: "i"},
		},
	}

	for _, c := range testCases {
		program, err := Parse(c.tokens)
		assert.NoError(t, err)

		assert.Equal(t, command.A, program[0].Type)
		assert.Equal(t, c.expected, program[0].RealCmd)
	}
}

func TestCommandTypeC(t *testing.T) {
	type TestCaseC struct {
		tokens   []token.Token
		expected CmdC
	}

	testCases := []TestCaseC{
		{
			// dest=comp;jump
			tokens: []token.Token{
				{Type: token.LOCATION, Value: "D"},
				{Type: token.ASSIGN, Value: "="},
				{Type: token.LOCATION, Value: "M"},
				{Type: token.OPERATOR, Value: "+"},
				{Type: token.NUMBER, Value: "1"},
				{Type: token.JUMP, Value: "JNE"},
				{Type: token.END},
			},
			expected: CmdC{
				D: Dest{D: true, A: false, M: false},
				C: "M+1",
				J: "JNE",
			},
		},
		{
			// comp;jump
			tokens: []token.Token{
				{Type: token.NUMBER, Value: "0"},
				{Type: token.JUMP, Value: "JMP"},
				{Type: token.END},
			},
			expected: CmdC{
				D: Dest{D: false, A: false, M: false},
				C: "0",
				J: "JMP",
			},
		},
		{
			// dest=comp (with all 3 destinations, with operator leading comp)
			tokens: []token.Token{
				{Type: token.LOCATION, Value: "D"},
				{Type: token.LOCATION, Value: "M"},
				{Type: token.LOCATION, Value: "A"},
				{Type: token.ASSIGN, Value: "="},
				{Type: token.OPERATOR, Value: "!"},
				{Type: token.LOCATION, Value: "D"},
				{Type: token.END},
			},
			expected: CmdC{
				D: Dest{D: true, A: true, M: true},
				C: "!D",
				J: "",
			},
		},
		{
			// dest=comp (with numeric leading comp)
			tokens: []token.Token{
				{Type: token.LOCATION, Value: "D"},
				{Type: token.LOCATION, Value: "M"},
				{Type: token.LOCATION, Value: "A"},
				{Type: token.ASSIGN, Value: "="},
				{Type: token.OPERATOR, Value: "1"},
				{Type: token.END},
			},
			expected: CmdC{
				D: Dest{D: true, A: true, M: true},
				C: "1",
				J: "",
			},
		},
		{
			// comp only
			tokens: []token.Token{
				{Type: token.OPERATOR, Value: "1"},
				{Type: token.END},
			},
			expected: CmdC{
				D: Dest{D: false, A: false, M: false},
				C: "1",
				J: "",
			},
		},
	}

	for _, c := range testCases {
		program, err := Parse(c.tokens)
		assert.NoError(t, err)

		assert.Equal(t, command.C, program[0].Type)
		assert.Equal(t, c.expected, program[0].RealCmd)
	}
}
