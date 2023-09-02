// 変数宣言では、スタックから取り出すときに何番目にあるかの情報が必要である

package object

import (
	"fmt"

	"github.com/kijimaD/gogo/token"
)

const (
	INTEGER_OBJ = "INTEGER"
	STRING_OBJ  = "STRING"
	CHAR_OBJ    = "CHAR"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
	CurPos() int // スタックの何番目にあるか?
	GetCtype() token.Ctype
}

type Integer struct {
	Value int64
	Pos   int
}

func (i *Integer) Type() ObjectType      { return INTEGER_OBJ }
func (i *Integer) Inspect() string       { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) CurPos() int           { return i.Pos }
func (i *Integer) GetCtype() token.Ctype { return token.CTYPE_INT }

type String struct {
	Value string
	Pos   int
}

func (s *String) Type() ObjectType      { return STRING_OBJ }
func (s *String) Inspect() string       { return s.Value }
func (s *String) CurPos() int           { return s.Pos }
func (s *String) GetCtype() token.Ctype { return token.CTYPE_STR }

type Char struct {
	Value int64
	Pos   int
}

func (c *Char) Type() ObjectType      { return CHAR_OBJ }
func (c *Char) Inspect() string       { return fmt.Sprintf("%d", c.Value) }
func (c *Char) CurPos() int           { return c.Pos }
func (c *Char) GetCtype() token.Ctype { return token.CTYPE_CHAR }
