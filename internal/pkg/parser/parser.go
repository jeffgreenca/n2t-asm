package parser

import "bitbucket.org/jeffgreenca/n2t-asm/internal/pkg/lex"

var (
	index   int
	tokens  []lex.Token
	cmdC    = CmdC{}
	program []Command
)

type CmdC struct {
	D Dest
	C string
	J string
}

type CmdA struct {
	address int
	symbol  string
	final   bool // true for static address, or resolved symbol address
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

type Command struct {
	Type CommandType
	C    interface{}
}

// Parse converts tokens to Commands
func Parse(t []lex.Token) ([]Command, error) {
	// initialize globals (ick)
	program = []Command{}
	tokens = t
	index = -1

	// read Statements until done
	for !end() {
		S()
	}

	return program, nil
}

func S() {
	if peek(lex.END) {
		accept(lex.END)
		return
	}
	if peek(lex.LOCATION) || peek(lex.OPERATOR) || peek(lex.NUMBER) {
		cmd := Command{Type: C_COMMAND}
		// initialize C type command global variable
		cmdC = CmdC{D: Dest{}}
		// parse C type command
		C()
		// store in program
		cmd.C = cmdC
		program = append(program, cmd)
	} else {
		panic("unexpected token")
	}
	S()
}

// C parses type C commands, syntax (dest=)comp(;jump)
func C() {
	if peek(lex.OPERATOR) || peek(lex.NUMBER) {
		COMP()
	} else if peek(lex.LOCATION) {
		// maybe this is comp part, or maybe this is dest part
		// lookahead up to 3 for assignment token
		if peekN(lex.ASSIGN, 2) || peekN(lex.ASSIGN, 3) || peekN(lex.ASSIGN, 4) {
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

// COMP is the comp part of C command
func COMP() {
	if peek(lex.LOCATION) || peek(lex.OPERATOR) || peek(lex.NUMBER) {
		acceptAny()
		cmdC.C += tokens[index].Value
		COMP()
	} else {
		C()
	}
}

// DEST is the dest part of C command
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

// peek returns true if the next token is of type t - no bounds checking
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

// accept next token of type t
func accept(t lex.TokenType) {
	if !peek(t) {
		panic("Wrong token type, did you forget to peek")
	}
	acceptAny()
}

// end returns true if we have accepted the last token already
func end() bool {
	return index >= len(tokens)-1
}

// accept next token of any type
func acceptAny() {
	index++
}
