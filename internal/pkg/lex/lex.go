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
	LABEL
)

var (
	globalTokens []Token
)

func init() {
	// globalTokens is used each Tokenize call to allow reusing memory space.
	// this optimization saved about 25% time assembling a 28k line asm program
	globalTokens = make([]Token, 0, 20)
}

// Tokenize takes as input one line of nand2tetris assembly statement, and returns tokenized formt.
func Tokenize(s string) ([]Token, error) {
	s = clean(s)
	if s == "" {
		return nil, nil
	}

	// convert string to a sequence of tokens
	// zero out globalTokens slice, re-using memory space
	globalTokens = globalTokens[:0]
	switch {
	case s[0] == '@':
		t, err := lexA(s)
		globalTokens = append(globalTokens, t...)
		if err != nil {
			return globalTokens, err
		}
	case strings.Contains(s, "="), strings.Contains(s, ";"):
		t, err := lexC(s)
		globalTokens = append(globalTokens, t...)
		if err != nil {
			return globalTokens, err
		}
	case s[0] == '(':
		t, err := lexL(s)
		globalTokens = append(globalTokens, t...)
		if err != nil {
			return globalTokens, err
		}
	default:
		panic(fmt.Sprintf("Unrecognized symbols: %v", s))
	}
	return globalTokens, nil
}

func lexL(s string) ([]Token, error) {
	var tokens = []Token{{Value: "(", Type: LABEL}}
	if s[len(s)-1] != ')' {
		return []Token{}, errors.New("Malformed label")
	}
	tokens = append(tokens, Token{Value: s[1 : len(s)-1], Type: SYMBOL})

	tokens = append(tokens, Token{Type: END})
	return tokens, nil
}

func lexA(s string) ([]Token, error) {
	var tokens = []Token{{Value: "@", Type: AT}}
	if len(s) < 2 {
		return tokens, errors.New("Malformed @ command - too short")
	}

	val := s[1:]
	if isNumeric(val) {
		tokens = append(tokens, Token{Value: val, Type: ADDRESS}, Token{Type: END})
	} else {
		tokens = append(tokens, Token{Value: val, Type: SYMBOL}, Token{Type: END})
	}

	return tokens, nil
}

// isNumeric returns true if s contains only digits 0-9, and is faster than using a regex.
func isNumeric(s string) bool {
	for _, ch := range s {
		if !('0' <= ch && ch <= '9') {
			return false
		}
	}
	return true
}

var (
	lexCTokens []Token
)

func init() {
	lexCTokens = make([]Token, 0, 20)
}

func lexC(s string) ([]Token, error) {
	split := strings.IndexRune(s, ';')

	// optimization - we know slice size needed based on if a jump field exists,
	// and that each comp rune becomes 1 token.
	var comp string
	var jump string
	var size int
	if split != -1 {
		comp = s[:split]
		jump = s[split+1:]
		size = len(comp) + 2
	} else {
		comp = s
		jump = ""
		size = len(comp) + 1
	}

	lexCTokens = lexCTokens[:size]
	for i, ch := range comp {
		switch ch {
		case '=':
			lexCTokens[i] = Token{Value: "=", Type: ASSIGN}
		case '0':
			lexCTokens[i] = Token{Value: "0", Type: NUMBER}
		case '1':
			lexCTokens[i] = Token{Value: "1", Type: NUMBER}
		case '+':
			lexCTokens[i] = Token{Value: "+", Type: OPERATOR}
		case '-':
			lexCTokens[i] = Token{Value: "-", Type: OPERATOR}
		case '!':
			lexCTokens[i] = Token{Value: "!", Type: OPERATOR}
		case '&':
			lexCTokens[i] = Token{Value: "&", Type: OPERATOR}
		case '|':
			lexCTokens[i] = Token{Value: "|", Type: OPERATOR}
		case 'D':
			lexCTokens[i] = Token{Value: "D", Type: LOCATION}
		case 'M':
			lexCTokens[i] = Token{Value: "M", Type: LOCATION}
		case 'A':
			lexCTokens[i] = Token{Value: "A", Type: LOCATION}
		default:
			return lexCTokens, fmt.Errorf("Unexpected token in: %s", s)
		}
	}
	if jump != "" {
		switch jump {
		case "JGT", "JEQ", "JGE", "JLT", "JNE", "JLE", "JMP":
			lexCTokens[len(lexCTokens)-2] = Token{Value: jump, Type: JUMP}
		}
	}

	lexCTokens[len(lexCTokens)-1] = Token{Type: END}
	return lexCTokens, nil
}

// clean deletes comments and whitespace
func clean(s string) string {
	i := strings.Index(s, "//")
	if i != -1 {
		s = s[:i]
	}
	s = strings.TrimSpace(s)
	return s
}
