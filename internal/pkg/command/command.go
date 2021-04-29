package command

// Program is a sequence of commands.
type Program []Any

// Any command
type Any interface{}

// L type command
type L struct {
	Symbol string
}

// A type command
type A struct {
	Address int
	Symbol  string
	Static  bool
}

// C type command
type C struct {
	D Dest
	C string
	J string
}

// Dest part of C type command
type Dest struct {
	A bool
	D bool
	M bool
}
