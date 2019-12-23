package parser

import (
	"testing"

	"bitbucket.org/jeffgreenca/n2t-asm/internal/pkg/lex"
	"github.com/stretchr/testify/assert"
)

func TestCommandTypeA(t *testing.T) {
	type TestCaseA struct {
		tokens   []lex.Token
		expected CmdA
	}

	testCases := []TestCaseA{
		{
			tokens: []lex.Token{
				{Type: lex.AT, Value: "@"},
				{Type: lex.ADDRESS, Value: "1024"},
				{Type: lex.END},
			},
			expected: CmdA{Address: 1024, Final: true, Symbol: ""},
		},
		{
			tokens: []lex.Token{
				{Type: lex.AT, Value: "@"},
				{Type: lex.SYMBOL, Value: "foo"},
				{Type: lex.END},
			},
			expected: CmdA{Address: 0, Final: false, Symbol: "foo"},
		},
		{
			tokens: []lex.Token{
				{Type: lex.AT, Value: "@"},
				{Type: lex.SYMBOL, Value: "i"},
				{Type: lex.END},
			},
			expected: CmdA{Address: 0, Final: false, Symbol: "i"},
		},
	}

	for _, c := range testCases {
		program, err := Parse(c.tokens)
		assert.NoError(t, err)

		assert.Equal(t, A_COMMAND, program[0].Type)
		assert.Equal(t, c.expected, program[0].C)
	}
}

func TestCommandTypeC(t *testing.T) {
	type TestCaseC struct {
		tokens   []lex.Token
		expected CmdC
	}

	testCases := []TestCaseC{
		{
			// dest=comp;jump
			tokens: []lex.Token{
				{Type: lex.LOCATION, Value: "D"},
				{Type: lex.ASSIGN, Value: "="},
				{Type: lex.LOCATION, Value: "M"},
				{Type: lex.OPERATOR, Value: "+"},
				{Type: lex.NUMBER, Value: "1"},
				{Type: lex.JUMP, Value: "JNE"},
				{Type: lex.END},
			},
			expected: CmdC{
				D: Dest{D: true, A: false, M: false},
				C: "M+1",
				J: "JNE",
			},
		},
		{
			// comp;jump
			tokens: []lex.Token{
				{Type: lex.NUMBER, Value: "0"},
				{Type: lex.JUMP, Value: "JMP"},
				{Type: lex.END},
			},
			expected: CmdC{
				D: Dest{D: false, A: false, M: false},
				C: "0",
				J: "JMP",
			},
		},
		{
			// dest=comp (with all 3 destinations, with operator leading comp)
			tokens: []lex.Token{
				{Type: lex.LOCATION, Value: "D"},
				{Type: lex.LOCATION, Value: "M"},
				{Type: lex.LOCATION, Value: "A"},
				{Type: lex.ASSIGN, Value: "="},
				{Type: lex.OPERATOR, Value: "!"},
				{Type: lex.LOCATION, Value: "D"},
				{Type: lex.END},
			},
			expected: CmdC{
				D: Dest{D: true, A: true, M: true},
				C: "!D",
				J: "",
			},
		},
		{
			// dest=comp (with numeric leading comp)
			tokens: []lex.Token{
				{Type: lex.LOCATION, Value: "D"},
				{Type: lex.LOCATION, Value: "M"},
				{Type: lex.LOCATION, Value: "A"},
				{Type: lex.ASSIGN, Value: "="},
				{Type: lex.OPERATOR, Value: "1"},
				{Type: lex.END},
			},
			expected: CmdC{
				D: Dest{D: true, A: true, M: true},
				C: "1",
				J: "",
			},
		},
		{
			// comp only
			tokens: []lex.Token{
				{Type: lex.OPERATOR, Value: "1"},
				{Type: lex.END},
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

		assert.Equal(t, C_COMMAND, program[0].Type)
		assert.Equal(t, c.expected, program[0].C)
	}
}
