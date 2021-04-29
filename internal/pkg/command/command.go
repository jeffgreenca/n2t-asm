package command

type Type int

const (
	UNUSED Type = iota
	L
	C
	A
)

// TODO more satisfying command object that can be used
// effectively by both parser and assembler
// with less casting and also not using the empty interface!

// specific types of commands have different underlying data (easy)
// but also have different supported operations (how to do that in go elegantly?)
// for example, the A command needs to support:
//   check if it is Final
//   if not final, get its Symbol (if that is in some other table, do some stuff)
//   set its Address to the real location (resolve Symbol) and set Final (not really needed?)
//   convert the final Address to a binary instruction
// whereas L command:
//   read its Symbol (that's all)
// whereas C command:
//   rather complex logic, but all encapsulated!

// so really the only problem is the A command that requires external data, the symbol table,
// and might need its internal representation updated based on some extern during assembly.
// everything else could just be the string representation of itself, although that isn't the most
// beautiful thing I suppose.

// other idea is to implement all the methods stubs, throw NotImplemented where it doesn't apply,
// so they all conform to the standard interface

/*
type Command interface {
	Symbol() string, error
	SetSymbol(string) error
}
func (cmd *cmdL) Symbol() string, error { return cmd.Symbol, nil }
func (cmd *cmdA) Symbol() string, error { return cmd.Symbol, nil }
func (cmd *cmdC) Symbol() string, error { return "", errors.New("not implemented") }
*/

/*
type Command interface {
	String() string, error
}
// string does different computations, but how do you also do SetSymbol type thing and IsFinal without cast
*/

// or, do some really specific things per command
/*
func (cmd *cmdL) UpdateSymbol(blah)
*/
