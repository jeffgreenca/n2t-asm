package parser

import (
	"fmt"
	"strconv"

	"github.com/jeffgreenca/n2t-asm/internal/pkg/lex"
)

type CmdL struct {
	Symbol string
}

type CmdC struct {
	D Dest
	C string
	J string
}

type CmdA struct {
	Address int
	Symbol  string
	Final   bool // true for static address, or resolved Symbol address
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

type Parser struct {
	index   int
	tokens  []lex.Token
	program []Command
	cmdC    CmdC
	cmdA    CmdA
	cmdL    CmdL
}

func New() *Parser {
	return &Parser{}
}

// Parse converts tokens to Commands
func (p *Parser) Parse(t []lex.Token) ([]Command, error) {
	// initialize globals (ick)
	p.program = []Command{}
	p.tokens = t
	p.index = -1

	// read Statements until done
	for !p.end() {
		p.S()
	}

	return p.program, nil
}

func (p *Parser) S() error {
	if p.peek(lex.END) {
		p.accept(lex.END)
		return nil
	}
	if p.peek(lex.LOCATION) || p.peek(lex.OPERATOR) || p.peek(lex.NUMBER) {
		// init
		cmd := Command{Type: C_COMMAND}
		p.cmdC = CmdC{D: Dest{}}
		// parse
		p.C()
		// store
		cmd.C = p.cmdC
		p.program = append(p.program, cmd)
	} else if p.peek(lex.AT) {
		cmd := Command{Type: A_COMMAND}
		p.cmdA = CmdA{}
		p.A()
		cmd.C = p.cmdA
		p.program = append(p.program, cmd)
	} else if p.peek(lex.LABEL) {
		cmd := Command{Type: L_COMMAND}
		p.cmdL = CmdL{}
		p.L()
		cmd.C = p.cmdL
		p.program = append(p.program, cmd)
	} else {
		return fmt.Errorf("unexpected token: %v", p.tokens[p.index+1])
	}
	return p.S()
}

// L parses type L commands, syntax (symbol)
func (p *Parser) L() {
	p.accept(lex.LABEL)
	if p.peek(lex.SYMBOL) {
		p.accept(lex.SYMBOL)
		p.cmdL = CmdL{Symbol: p.tokens[p.index].Value}
	}
	if !p.peek(lex.END) {
		panic("Malforned label syntax")
	}
}

// A parses type A commands, syntax @(symbol|address)
func (p *Parser) A() {
	p.accept(lex.AT)
	if p.peek(lex.ADDRESS) {
		p.accept(lex.ADDRESS)
		i, err := strconv.Atoi(p.tokens[p.index].Value)
		if err != nil {
			panic("Unexpected parse address error")
		}
		p.cmdA = CmdA{Address: i, Final: true}
	} else if p.peek(lex.SYMBOL) {
		p.accept(lex.SYMBOL)
		p.cmdA = CmdA{Symbol: p.tokens[p.index].Value}
	}
	if !p.peek(lex.END) {
		panic("Malforned address syntax (@xxx)")
	}
}

// C parses type C commands, syntax (dest=)comp(;jump)
func (p *Parser) C() {
	if p.peek(lex.OPERATOR) || p.peek(lex.NUMBER) {
		p.COMP()
	} else if p.peek(lex.LOCATION) {
		// maybe this is comp part, or maybe this is dest part
		// lookahead up to 3 for assignment token
		if p.peekN(lex.ASSIGN, 2) || p.peekN(lex.ASSIGN, 3) || p.peekN(lex.ASSIGN, 4) {
			p.DEST()
		} else {
			p.COMP()
		}
	} else if p.peek(lex.JUMP) {
		p.accept(lex.JUMP)
		p.cmdC.J = p.tokens[p.index].Value
		// Done with C() -- ?
	}
}

// COMP is the comp part of C command
func (p *Parser) COMP() {
	if p.peek(lex.LOCATION) || p.peek(lex.OPERATOR) || p.peek(lex.NUMBER) {
		p.acceptAny()
		p.cmdC.C += p.tokens[p.index].Value
		p.COMP()
	} else {
		p.C()
	}
}

// DEST is the dest part of C command
func (p *Parser) DEST() {
	if p.peek(lex.LOCATION) {
		p.accept(lex.LOCATION)
		switch p.tokens[p.index].Value {
		case "A":
			p.cmdC.D.A = true
		case "D":
			p.cmdC.D.D = true
		case "M":
			p.cmdC.D.M = true
		default:
			panic("Unexpected token at DEST()")
		}
	}
	if p.peek(lex.ASSIGN) {
		p.accept(lex.ASSIGN)
		p.COMP()
	} else if p.peek(lex.LOCATION) {
		p.DEST()
	} else {
		panic("Unexpected token at DEST()")
	}
}

// peek returns true if the next token is of type t - no bounds checking
func (p *Parser) peek(t lex.TokenType) bool {
	return p.tokens[p.index+1].Type == t
}

// peekN looks safely ahead N tokens, returning false if out of tokens or not of type t
func (p *Parser) peekN(t lex.TokenType, n int) bool {
	if p.index+n >= len(p.tokens) {
		return false
	}
	return p.tokens[p.index+n].Type == t
}

// accept next token of type t
func (p *Parser) accept(t lex.TokenType) {
	if !p.peek(t) {
		panic("Wrong token type, did you forget to peek")
	}
	p.acceptAny()
}

// end returns true if we have accepted the last token already
func (p *Parser) end() bool {
	return p.index >= len(p.tokens)-1
}

// accept next token of any type
func (p *Parser) acceptAny() {
	p.index++
}
