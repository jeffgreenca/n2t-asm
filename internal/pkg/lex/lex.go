package lex

import (
	"bufio"
	"fmt"
	"os"
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
	globalTokens = make([]Token, 0, 20)
	lexCTokens   = make([]Token, 0, 20)

	end = Token{Type: END}
)

// TokenizeFile reads file and returns tokenized form or error
func TokenizeFile(f *os.File) ([]Token, error) {
	var result []Token

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		tokens, err := Tokenize(line)
		if err != nil {
			return []Token{}, fmt.Errorf("failed to tokenize line '%s': %v", line, err)
		}
		result = append(result, tokens...)
	}
	return result, nil
}

// Tokenize one line of nand2tetris assembly statement
func Tokenize(s string) ([]Token, error) {
	s = clean(s)
	if s == "" {
		return nil, nil
	}

	// convert string to a sequence of tokens
	// zero out globalTokens slice, re-using memory space
	globalTokens = globalTokens[:0]
	var tokens []Token
	var err error
	switch {
	case s[0] == '@':
		tokens, err = lexA(s)
	case s[0] == '(':
		tokens, err = lexL(s)
	case isC(s):
		tokens, err = lexC(s)
	default:
		return []Token{}, fmt.Errorf("unrecognized symbol: %s", s)
	}
	if err != nil {
		return []Token{}, fmt.Errorf("error lexing '%s': %v", s, err)
	}
	globalTokens = append(globalTokens, tokens...)
	return globalTokens, nil
}

func lexL(s string) ([]Token, error) {
	label := strings.Trim(s, "()")
	if len(s)-2 != len(label) {
		return []Token{}, fmt.Errorf("malformed label: %v", s)
	}
	tokens := []Token{{Value: "(", Type: LABEL}, Token{Value: label, Type: SYMBOL}, end}
	return tokens, nil
}

func lexA(s string) ([]Token, error) {
	if len(s) < 2 {
		return []Token{}, fmt.Errorf("malformed '@' command, too short: %s", s)
	}
	v := s[1:]
	tokens := []Token{{Value: "@", Type: AT}, Token{Value: v, Type: typeFromVal(v)}, end}
	return tokens, nil
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
			return []Token{}, fmt.Errorf("unexpected rune '%v' in: %s", ch, s)
		}
	}
	if len(jump) > 0 {
		switch jump {
		case
			"JGT",
			"JEQ",
			"JGE",
			"JLT",
			"JNE",
			"JLE",
			"JMP":
			lexCTokens[len(lexCTokens)-2] = Token{Value: jump, Type: JUMP}
		}
	}

	lexCTokens[len(lexCTokens)-1] = end
	return lexCTokens, nil
}

func clean(s string) string {
	i := strings.Index(s, "//")
	if i > -1 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}

func isC(s string) bool {
	for _, ch := range s {
		if ch == '=' {
			return true
		}
		if ch == ';' {
			return true
		}
	}
	return false
}

func typeFromVal(s string) TokenType {
	if isNum(s) {
		return ADDRESS
	}
	return SYMBOL
}

func isNum(s string) bool {
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
