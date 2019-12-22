package parser

import "bitbucket.org/jeffgreenca/n2t-asm/internal/pkg/lex"

var (
	index   int
	tokens  []lex.Token
	cmdC    = CmdC{}
	cmdType CommandType
)

type CmdC struct {
	D Dest
	C string
	J string
}

type Dest struct {
	A bool
	D bool
	M bool
}

type CommandType int

const (
	UNUSED CommandType = iota
	L_COMMAND
	C_COMMAND
	A_COMMAND
)

// Parse takes a list of tokens and returns one or more HACK statements ?
func Parse(t []lex.Token) ([]string, error) {
	tokens = t
	index = -1
	S()

	return nil, nil
}

func S() {
	cmdC = CmdC{D: Dest{}}
	if peek(lex.LOCATION) || peek(lex.OPERATOR) || peek(lex.NUMBER) {
		cmdType = C_COMMAND
		C()
	}
}

// C parses type C commands, syntax (dest=)comp(;jump)
func C() {
	if peek(lex.OPERATOR) || peek(lex.NUMBER) {
		COMP()
	} else if peek(lex.LOCATION) {
		// maybe this is comp part, or maybe this is dest part
		// lookahead up to 3 for assignment token
		if peekN(lex.ASSIGN, 1) || peekN(lex.ASSIGN, 2) || peekN(lex.ASSIGN, 3) {
			DEST()
		} else {
			COMP()
		}
	} else if peek(lex.JUMP) {
		accept(lex.JUMP)
		cmdC.J = tokens[index].Value
		// Done with C() -- ?
	}
}

func COMP() {
	if peek(lex.LOCATION) || peek(lex.OPERATOR) || peek(lex.NUMBER) {
		acceptAny()
		cmdC.C += tokens[index].Value
		COMP()
	} else {
		C()
	}
}

func DEST() {
	if peek(lex.LOCATION) {
		accept(lex.LOCATION)
		switch tokens[index].Value {
		case "A":
			cmdC.D.A = true
		case "D":
			cmdC.D.D = true
		case "M":
			cmdC.D.M = true
		default:
			panic("Unexpected token at DEST()")
		}
	}
	if peek(lex.ASSIGN) {
		accept(lex.ASSIGN)
		COMP()
	} else if peek(lex.LOCATION) {
		DEST()
	} else {
		panic("Unexpected token at DEST()")
	}
}

// peek returns true if the next token is of type t
func peek(t lex.TokenType) bool {
	return tokens[index+1].Type == t
}

// peekN looks safely ahead N tokens, returning false if out of tokens or not of type t
func peekN(t lex.TokenType, n int) bool {
	if index+n >= len(tokens) {
		return false
	}
	return tokens[index+n].Type == t
}

func accept(t lex.TokenType) {
	if !peek(t) {
		panic("Wrong token type, did you forget to peek")
	}
	acceptAny()
}

func end() bool {
	return index >= len(tokens)-1
}

func acceptAny() {
	index++
}
