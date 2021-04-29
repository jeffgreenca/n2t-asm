package command

type L struct {
	Symbol string
}

type C struct {
	D Dest
	C string
	J string
}

type A struct {
	Address int
	Symbol  string
	Static  bool
}

type Dest struct {
	A bool
	D bool
	M bool
}

type Any interface{}
