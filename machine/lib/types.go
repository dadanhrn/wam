package lib

const (
	HEAP_STR      = "STR"
	HEAP_STR_DATA = "STR_DATA"
	HEAP_REF      = "REF"
)

type HeapCell struct {
	Type string
	Data interface{}
}

type HeapFunctor struct {
	Identifier string
	Arity      int
}
