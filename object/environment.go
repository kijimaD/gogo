package object

type Environment struct {
	store  map[string]Object
	VarPos int
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object), VarPos: 1}
}

func (e *Environment) Set(ident string, obj Object) {
	e.store[ident] = obj
	e.VarPos++
}

func (e *Environment) Get(ident string) (Object, bool) {
	obj, ok := e.store[ident]
	return obj, ok
}
