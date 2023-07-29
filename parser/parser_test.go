package parser

import (
	"testing"

	"github.com/kijimaD/gogo/lexer"
)

func TestParseProgram(t *testing.T) {
	l := lexer.NewLexer("hi")
	p := New(l)
	p.ParseProgram()
}
