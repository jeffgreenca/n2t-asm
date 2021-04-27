package token

type Token struct {
	Value string
	Type  Type
}

type Type int

const (
	UNKNOWN Type = iota
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

// Commonly used fixed token types
var (
	End = Token{Type: END}
)
