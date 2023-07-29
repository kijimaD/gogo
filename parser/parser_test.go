package parser

import (
	"testing"

	"github.com/kijimaD/gogo/ast"
	"github.com/kijimaD/gogo/lexer"
	"github.com/kijimaD/gogo/token"
	"github.com/stretchr/testify/assert"
)

func TestParseProgram(t *testing.T) {
	l := lexer.NewLexer(`"hi" "all"`)
	p := New(l)

	result := p.ParseProgram()
	assert.Equal(t, 2, len(result.Statements))
}

func TestParseExpression(t *testing.T) {
	l := lexer.NewLexer(`"hi"`)
	p := New(l)

	actual := p.parseExpression()
	expect := &ast.StringLiteral{Token: token.Token{Type: "STRING", Literal: "hi"}, Value: "hi"}
	assert.Equal(t, expect, actual)
}

func TestNextToken(t *testing.T) {
	l := lexer.NewLexer(`"hello world"`)
	p := New(l)

	expectCur := token.Token{Type: "STRING", Literal: "hello world"}
	assert.Equal(t, expectCur, p.curToken)
	expectPeek := token.Token{Type: "EOF", Literal: ""}
	assert.Equal(t, expectPeek, p.peekToken)
}
