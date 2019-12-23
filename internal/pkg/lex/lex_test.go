package lex

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			{Type: END, Value: ""},
		},
	}

	for k, v := range testCases {
		actual, err := Tokenize(k)
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}

func TestTokenizeTypeA(t *testing.T) {
	testCases := map[string][]Token{
		"@100": {
			{Type: AT, Value: "@"},
			{Type: ADDRESS, Value: "100"},
			{Type: END, Value: ""},
		},
		"@i": {
			{Type: AT, Value: "@"},
			{Type: SYMBOL, Value: "i"},
			{Type: END, Value: ""},
		},
		"@foo": {
			{Type: AT, Value: "@"},
			{Type: SYMBOL, Value: "foo"},
			{Type: END, Value: ""},
		},
	}

	for k, v := range testCases {
		actual, err := Tokenize(k)
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}
