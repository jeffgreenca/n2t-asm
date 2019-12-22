package lex

import (
	"errors"
	"fmt"
	"strings"
)

type Token struct {
	Value string
	Type  TokenType
}

type TokenType int

const (
	UNKNOWN TokenType = iota
	LOCATION
	ASSIGN
	OPERATOR
	NUMBER
	JUMP
	END
	AT
	SYMBOL
	ADDRESS
)

// Tokenize takes as input one line of nand2tetris assembly statement, and returns tokenized formt.
func Tokenize(s string) ([]Token, error) {
	s = clean(s)
	if s == "" {
		return []Token{}, nil
	}

	// convert string to a sequence of tokens
	var tokens = []Token{}
	switch {
	case s[0] == '@':
		// TODO A command
	case strings.Contains(s, "="):
		t, err := lexC(s)
		tokens = append(tokens, t...)
		if err != nil {
			return tokens, err
		}
	case s[0] == '(':
		// TODO L command
	default:
		// TODO - unknown
	}
	return tokens, nil
}

func lexC(s string) ([]Token, error) {
	var tokens = []Token{}
	parts := strings.SplitN(s, ";", 2)
	for _, ch := range parts[0] {
		switch ch {
		case '=':
			tokens = append(tokens, Token{Value: "=", Type: ASSIGN})
		case '0':
			tokens = append(tokens, Token{Value: "0", Type: NUMBER})
		case '1':
			tokens = append(tokens, Token{Value: "1", Type: NUMBER})
		case '+':
			tokens = append(tokens, Token{Value: "+", Type: OPERATOR})
		case '-':
			tokens = append(tokens, Token{Value: "-", Type: OPERATOR})
		case '!':
			tokens = append(tokens, Token{Value: "!", Type: OPERATOR})
		case '&':
			tokens = append(tokens, Token{Value: "&", Type: OPERATOR})
		case '|':
			tokens = append(tokens, Token{Value: "|", Type: OPERATOR})
		case 'D':
			tokens = append(tokens, Token{Value: "D", Type: LOCATION})
		case 'M':
			tokens = append(tokens, Token{Value: "M", Type: LOCATION})
		case 'A':
			tokens = append(tokens, Token{Value: "A", Type: LOCATION})
		default:
			return tokens, errors.New(fmt.Sprintf("Unexpected token in: %s", s))
		}
	}
	if len(parts) == 2 {
		switch parts[1] {
		case "JMP", "JLT", "JNE":
			tokens = append(tokens, Token{Value: parts[1], Type: JUMP})
		}
	}

	tokens = append(tokens, Token{Type: END})
	return tokens, nil
}

// clean deletes comments and whitespace
func clean(s string) string {
	s = strings.Split(s, "//")[0]
	s = strings.TrimSpace(s)
	return s
}
