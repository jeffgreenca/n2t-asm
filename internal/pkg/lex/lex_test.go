package lex

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jeffgreenca/n2t-asm/internal/pkg/token"
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

func TestTokenizeTypeL(t *testing.T) {
	testCases := map[string][]token.Token{
		"(foobar)": {
			{Type: token.LABEL, Value: "("},
			{Type: token.SYMBOL, Value: "foobar"},
			{Type: token.END, Value: ""},
		},
	}

	for k, v := range testCases {
		actual, err := tokenize(k)
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}

func TestTokenizeTypeC(t *testing.T) {
	testCases := map[string][]token.Token{
		"D=M+1;JNE": {
			{Type: token.LOCATION, Value: "D"},
			{Type: token.ASSIGN, Value: "="},
			{Type: token.LOCATION, Value: "M"},
			{Type: token.OPERATOR, Value: "+"},
			{Type: token.NUMBER, Value: "1"},
			{Type: token.JUMP, Value: "JNE"},
			{Type: token.END, Value: ""},
		},
		"D;JGT": {
			{Type: token.LOCATION, Value: "D"},
			{Type: token.JUMP, Value: "JGT"},
			{Type: token.END, Value: ""},
		},
	}

	for k, v := range testCases {
		actual, err := tokenize(k)
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}

func TestTokenizeTypeA(t *testing.T) {
	testCases := map[string][]token.Token{
		"@100": {
			{Type: token.AT, Value: "@"},
			{Type: token.ADDRESS, Value: "100"},
			{Type: token.END, Value: ""},
		},
		"@i": {
			{Type: token.AT, Value: "@"},
			{Type: token.SYMBOL, Value: "i"},
			{Type: token.END, Value: ""},
		},
		"@foo": {
			{Type: token.AT, Value: "@"},
			{Type: token.SYMBOL, Value: "foo"},
			{Type: token.END, Value: ""},
		},
	}

	for k, v := range testCases {
		actual, err := tokenize(k)
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}
