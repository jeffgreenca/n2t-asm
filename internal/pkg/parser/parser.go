package parser

import (
	"fmt"
	"strconv"

	"github.com/jeffgreenca/n2t-asm/internal/pkg/command"
	"github.com/jeffgreenca/n2t-asm/internal/pkg/token"
)

type state struct {
	index   int
	tokens  []token.Token
	program []command.Any
	cmdC    command.C
	cmdA    command.A
	cmdL    command.L
}

// Parse tokens to commands.
func Parse(tokens []token.Token) ([]command.Any, error) {
	s := &state{
		program: []command.Any{},
		tokens:  tokens,
		index:   -1,
	}
	return s.parse()
}

func (s *state) parse() ([]command.Any, error) {
	for !s.end() {
		err := s.s()
		if err != nil {
			return nil, err
		}
	}
	return s.program, nil
}

func (s *state) s() error {
	if s.peek(token.END) {
		return s.accept(token.END)
	}
	if s.peek(token.LOCATION) || s.peek(token.OPERATOR) || s.peek(token.NUMBER) {
		// init
		s.cmdC = command.C{D: command.Dest{}}
		// parse
		err := s.c()
		if err != nil {
			return fmt.Errorf("parse error for location/operator/number: %v", err)
		}
		// store
		s.program = append(s.program, s.cmdC)
	} else if s.peek(token.AT) {
		s.cmdA = command.A{}
		err := s.a()
		if err != nil {
			return fmt.Errorf("parse error for AT: %v", err)
		}
		s.program = append(s.program, s.cmdA)
	} else if s.peek(token.LABEL) {
		s.cmdL = command.L{}
		err := s.l()
		if err != nil {
			return fmt.Errorf("parse error for label: %v", err)
		}
		s.program = append(s.program, s.cmdL)
	} else {
		return fmt.Errorf("unexpected token, wut: %v", s.peekGet())
	}
	return s.s()
}

// l parses type l commands, syntax (symbol)
func (s *state) l() error {
	err := s.accept(token.LABEL)
	if err != nil {
		return err
	}
	if s.peek(token.SYMBOL) {
		err := s.accept(token.SYMBOL)
		if err != nil {
			return err
		}
		s.cmdL = command.L{Symbol: s.tokens[s.index].Value}
	}
	if !s.peek(token.END) {
		return fmt.Errorf("malforned label syntax, expected END, got: %v", s.peekGet())
	}
	return nil
}

// a parses type a commands, syntax @(symbol|address)
func (s *state) a() error {
	err := s.accept(token.AT)
	if err != nil {
		return err
	}
	if s.peek(token.ADDRESS) {
		err := s.accept(token.ADDRESS)
		if err != nil {
			return err
		}
		i, err := strconv.Atoi(s.tokens[s.index].Value)
		if err != nil {
			return fmt.Errorf("unexpected error parsing address: %v", err)
		}
		s.cmdA = command.A{Address: i, Static: true}
	} else if s.peek(token.SYMBOL) {
		err := s.accept(token.SYMBOL)
		if err != nil {
			return err
		}
		s.cmdA = command.A{Symbol: s.tokens[s.index].Value}
	}
	if !s.peek(token.END) {
		return fmt.Errorf("malforned address syntax (@xxx), expected END got: %v", s.peekGet())
	}
	return nil
}

// c parses type c commands, syntax (dest=)comp(;jump)
func (s *state) c() error {
	if s.peek(token.OPERATOR) || s.peek(token.NUMBER) {
		err := s.comp()
		if err != nil {
			return fmt.Errorf("error calling comp() for operator/number: %v", err)
		}
	} else if s.peek(token.LOCATION) {
		// maybe this is comp part, or maybe this is dest part
		if s.peekRange(token.ASSIGN, 2, 3) {
			err := s.dest()
			if err != nil {
				return fmt.Errorf("error calling dest(): %v", err)
			}
		} else {
			err := s.comp()
			if err != nil {
				return fmt.Errorf("error calling comp() for LOCATION: %v", err)
			}
		}
	} else if s.peek(token.JUMP) {
		err := s.accept(token.JUMP)
		if err != nil {
			return err
		}
		s.cmdC.J = s.tokens[s.index].Value
		// Done with c()
	}
	return nil
}

// comp is the comp part of C command
func (s *state) comp() error {
	if s.peek(token.LOCATION) || s.peek(token.OPERATOR) || s.peek(token.NUMBER) {
		s.acceptAny()
		s.cmdC.C += s.tokens[s.index].Value
		err := s.comp()
		if err != nil {
			return fmt.Errorf("error calling comp(): %v", err)
		}
	} else {
		err := s.c()
		if err != nil {
			return fmt.Errorf("error calling c(): %v", err)
		}
	}
	return nil
}

// dest is the dest part of C command
func (s *state) dest() error {
	if s.peek(token.LOCATION) {
		err := s.accept(token.LOCATION)
		if err != nil {
			return err
		}
		switch s.tokens[s.index].Value {
		case "A":
			s.cmdC.D.A = true
		case "D":
			s.cmdC.D.D = true
		case "M":
			s.cmdC.D.M = true
		default:
			return fmt.Errorf("unexpected value, expected A/D/M got: %v", s.tokens[s.index].Value)
		}
	}
	if s.peek(token.ASSIGN) {
		err := s.accept(token.ASSIGN)
		if err != nil {
			return err
		}
		err = s.comp()
		if err != nil {
			return fmt.Errorf("error calling comp() for ASSIGN: %v", err)
		}
	} else if s.peek(token.LOCATION) {
		err := s.dest()
		if err != nil {
			return fmt.Errorf("dest(): %v", err)
		}
	} else {
		return fmt.Errorf("unexpected token, expected ASSIGN or LOCATION but got: %v", s.peekGet())
	}
	return nil
}

// peek returns true if the next token is of type t - no bounds checking
func (s *state) peek(t token.Type) bool {
	return s.tokens[s.index+1].Type == t
}

// peekGet returns the next token without advancing the counter - no bounds checking
func (s *state) peekGet() token.Token {
	return s.tokens[s.index+1]
}

// peekRange searches ahead, returning true if t is in the next count tokens
func (s *state) peekRange(t token.Type, offset int, count int) bool {
	if len(s.tokens) < s.index+offset+count {
		return false
	}
	for n := 0; n < count; n++ {
		if s.tokens[s.index+offset+n].Type == t {
			return true
		}
	}
	return false
}

// accept next token of type t
func (s *state) accept(t token.Type) error {
	if !s.peek(t) {
		return fmt.Errorf("wrong token type, accept %v but got: %v", t, s.peekGet())
	}
	s.acceptAny()
	return nil
}

// end returns true if we have accepted the last token already
func (s *state) end() bool {
	return s.index >= len(s.tokens)-1
}

// accept next token of any type
func (s *state) acceptAny() {
	s.index++
}
