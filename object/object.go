package object

import "fmt"

const (
	INTEGER_OBJ = "INTEGER"
	STRING_OBJ  = "STRING"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
	CurPos() int // スタックの何番目にあるか?
}

type Integer struct {
	Value int64
	Pos   int
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) CurPos() int      { return i.Pos }

type String struct {
	Value string
	Pos   int
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) CurPos() int      { return s.Pos }
