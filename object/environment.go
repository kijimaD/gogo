package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

func (e *Environment) Set(ident string, obj Object) {
	e.store[ident] = obj
}

func (e *Environment) Get(ident string) (Object, bool) {
	obj, ok := e.store[ident]
	return obj, ok
}
