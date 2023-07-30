package parser

import (
	"testing"

	"github.com/kijimaD/gogo/ast"
	"github.com/kijimaD/gogo/lexer"
	"github.com/kijimaD/gogo/token"
	"github.com/stretchr/testify/assert"
)

// Helper ================

// エラーがあった場合にテストを失敗させる
func checkParserErrors(t *testing.T, p *Parser) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

// エラーがない場合にテストを失敗させる
func assertParserErrors(t *testing.T, p *Parser) {
	t.Helper()
	if len(p.Errors()) != 0 {
		return
	}
	t.Errorf("expected parser error, but not error!")
	t.FailNow()
}

// Test body ================

func TestParseProgram(t *testing.T) {
	l := lexer.NewLexer(`"hi" "all"`)
	p := New(l)

	result := p.ParseProgram()
	checkParserErrors(t, p)
	assert.Equal(t, 2, len(result.Statements))
}

// ILLEGALトークンがあるとerrorsが入る
func TestParseProgramIllegal(t *testing.T) {
	l := lexer.NewLexer(`illegal`)
	p := New(l)

	p.ParseProgram()
	assertParserErrors(t, p)
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect interface{}
	}{
		{
			name:   "文字列をパースする",
			input:  `"hi"`,
			expect: &ast.StringLiteral{Token: token.Token{Type: "STRING", Literal: "hi"}, Value: "hi"},
		},
		{
			name:   "空白をスルーする(文字列)",
			input:  `    "hi"`,
			expect: &ast.StringLiteral{Token: token.Token{Type: "STRING", Literal: "hi"}, Value: "hi"},
		},
		{
			name:   "整数をパースする",
			input:  `1`,
			expect: &ast.IntegerLiteral{Token: token.Token{Type: "INT", Literal: "1"}, Value: 1},
		},
		{
			name:   "空白をスルーする(整数)",
			input:  `    1`,
			expect: &ast.IntegerLiteral{Token: token.Token{Type: "INT", Literal: "1"}, Value: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := New(l)
			actual := p.parseExpression()
			checkParserErrors(t, p)
			assert.Equal(t, tt.expect, actual)

		})
	}
}

func TestNextToken(t *testing.T) {
	l := lexer.NewLexer(`"hello world"`)
	p := New(l)

	expectCur := token.Token{Type: "STRING", Literal: "hello world"}
	assert.Equal(t, expectCur, p.curToken)
	expectPeek := token.Token{Type: "EOF", Literal: ""}
	assert.Equal(t, expectPeek, p.peekToken)
}
