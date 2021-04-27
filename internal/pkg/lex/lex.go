package lex

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/jeffgreenca/n2t-asm/internal/pkg/token"
)

// TODO encapsulate
var (
	globalTokens = make([]token.Token, 0, 20)
	lexCTokens   = make([]token.Token, 0, 20)
)

// Tokenize
func Tokenize(r io.Reader) ([]token.Token, error) {
	var result []token.Token

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		tokens, err := tokenize(line)
		if err != nil {
			return []token.Token{}, fmt.Errorf("failed to tokenize line '%s': %v", line, err)
		}
		result = append(result, tokens...)
	}
	return result, nil
}

// tokenize one line of nand2tetris assembly statement
func tokenize(s string) ([]token.Token, error) {
	s = clean(s)
	if s == "" {
		return nil, nil
	}

	// convert string to a sequence of tokens
	// zero out globalTokens slice, re-using memory space
	globalTokens = globalTokens[:0]
	var tokens []token.Token
	var err error
	switch {
	case s[0] == '@':
		tokens, err = lexA(s)
	case s[0] == '(':
		tokens, err = lexL(s)
	case isC(s):
		tokens, err = lexC(s)
	default:
		return []token.Token{}, fmt.Errorf("unrecognized symbol: %s", s)
	}
	if err != nil {
		return []token.Token{}, fmt.Errorf("error lexing '%s': %v", s, err)
	}
	globalTokens = append(globalTokens, tokens...)
	return globalTokens, nil
}

func lexL(s string) ([]token.Token, error) {
	label := strings.Trim(s, "()")
	if len(s)-2 != len(label) {
		return []token.Token{}, fmt.Errorf("malformed label: %v", s)
	}
	tokens := []token.Token{{Value: "(", Type: token.LABEL}, token.Token{Value: label, Type: token.SYMBOL}, token.End}
	return tokens, nil
}

func lexA(s string) ([]token.Token, error) {
	if len(s) < 2 {
		return []token.Token{}, fmt.Errorf("malformed '@' command, too short: %s", s)
	}
	v := s[1:]
	tokens := []token.Token{{Value: "@", Type: token.AT}, token.Token{Value: v, Type: typeFromVal(v)}, token.End}
	return tokens, nil
}
func lexC(s string) ([]token.Token, error) {
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
			lexCTokens[i] = token.Token{Value: "=", Type: token.ASSIGN}
		case '0':
			lexCTokens[i] = token.Token{Value: "0", Type: token.NUMBER}
		case '1':
			lexCTokens[i] = token.Token{Value: "1", Type: token.NUMBER}
		case '+':
			lexCTokens[i] = token.Token{Value: "+", Type: token.OPERATOR}
		case '-':
			lexCTokens[i] = token.Token{Value: "-", Type: token.OPERATOR}
		case '!':
			lexCTokens[i] = token.Token{Value: "!", Type: token.OPERATOR}
		case '&':
			lexCTokens[i] = token.Token{Value: "&", Type: token.OPERATOR}
		case '|':
			lexCTokens[i] = token.Token{Value: "|", Type: token.OPERATOR}
		case 'D':
			lexCTokens[i] = token.Token{Value: "D", Type: token.LOCATION}
		case 'M':
			lexCTokens[i] = token.Token{Value: "M", Type: token.LOCATION}
		case 'A':
			lexCTokens[i] = token.Token{Value: "A", Type: token.LOCATION}
		default:
			return []token.Token{}, fmt.Errorf("unexpected rune '%v' in: %s", ch, s)
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
			lexCTokens[len(lexCTokens)-2] = token.Token{Value: jump, Type: token.JUMP}
		}
	}

	lexCTokens[len(lexCTokens)-1] = token.End
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

func typeFromVal(s string) token.Type {
	if isNum(s) {
		return token.ADDRESS
	}
	return token.SYMBOL
}

func isNum(s string) bool {
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
