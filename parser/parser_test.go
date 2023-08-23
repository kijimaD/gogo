package parser

import (
	"strings"
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

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	t.Helper()
	opExp, ok := exp.(*ast.InfixExpression)
	assert.True(t, ok)
	assert.Equal(t, operator, opExp.Operator)
	assert.Equal(t, opExp.Left.String(), left)
	assert.Equal(t, opExp.Right.String(), right)

	return true
}

// Test body ================

func TestParsePrim(t *testing.T) {
	tests := []struct {
		input  string
		expect interface{}
		length int
	}{
		{
			`1`,
			`1`,
			1,
		},
		{
			`"1"`,
			`"1"`,
			1,
		},
		{
			`1;1`,
			`11`,
			2,
		},
		{
			`1;
2;
3;`,
			`123`,
			3,
		},
		{
			`"hello" "world"`,
			`"hello""world"`,
			2,
		},
		{
			`"hello"; "world"`,
			`"hello""world"`,
			2,
		},
		{
			`"hello"
				  "world"`,
			`"hello""world"`,
			2,
		},
		{
			`'h'`,
			`'h'`,
			1,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		pg := p.ParseProgram()

		assert.Equal(t, tt.length, len(pg.Statements))

		checkParserErrors(t, p)
		actual := pg.String()
		assert.Equal(t, tt.expect, actual)
	}
}

// ILLEGALトークンがあるとerrorsが入る
func TestParseProgramIllegal(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`"unbalance quote`},
		{`42a`},        // 数値から始まる識別子
		{`1+`},         // 中置演算子の右側がない
		{`'MULTIPLE'`}, // charリテラルに複数の文字
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		_ = p.ParseProgram()
		assertParserErrors(t, p)
	}

}

func TestParseInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  string
		operator   string
		rightValue string
	}{
		{"5 + 5", "5", "+", "5"},
		{"5 - 5", "5", "-", "5"},
		{"5 * 5", "5", "*", "5"},
		{"5 / 5", "5", "/", "5"},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program h.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		ok = testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)

		assert.True(t, ok)
	}
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect interface{}
	}{
		{
			input:  `"hi"`,
			expect: `"hi"`,
		},
		{
			input:  `    "hi"`,
			expect: `"hi"`,
		},
		{
			input:  `1`,
			expect: `1`,
		},
		{
			input:  `    1`,
			expect: `1`,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		actual := p.parseExpression(LOWEST)
		checkParserErrors(t, p)
		assert.Equal(t, tt.expect, actual.String())
	}
}

// 優先順位が正しいか
func TestParsePrecedence(t *testing.T) {
	tests := []struct {
		input     string
		expect    interface{}
		expectLen int
	}{
		{
			`1 + 2 - 3`,
			`((1 + 2) - 3)`,
			1,
		},
		{
			`1 - 2 + 3`,
			`((1 - 2) + 3)`,
			1,
		},
		{
			`1 + 2 * 3`,
			`(1 + (2 * 3))`,
			1,
		},
		{
			`1 + 2 / 3`,
			`(1 + (2 / 3))`,
			1,
		},
		{
			`1 * 2 / 3`,
			`((1 * 2) / 3)`,
			1,
		},
		{
			`1 + 2; 3 / 4`,
			`(1 + 2)(3 / 4)`,
			2,
		},
		{
			`1; 2 * 3`,
			`1(2 * 3)`,
			2,
		},
		{
			`int a = 1`,
			`(int a = 1)`,
			1,
		},
		{
			`int a = 1+2`,
			`(int a = (1 + 2))`,
			1,
		},
		{
			`int a = 1+2*3`,
			`(int a = (1 + (2 * 3)))`,
			1,
		},
		{
			`int a = 1; a+2*2`,
			`(int a = 1)(a + (2 * 2))`,
			2,
		},
		{
			`string a = "str"`,
			`(string a = "str")`,
			1,
		},
		{
			`f(1, 2)`,
			`f(1, 2)`,
			1,
		},
		{
			`f(a, b, c, d, e)`,
			`f(a, b, c, d, e)`,
			1,
		},
		{
			`f(1+1, 2)`,
			`f((1 + 1), 2)`,
			1,
		},
		{
			`f1(1, 2)`,
			`f1(1, 2)`,
			1,
		},
		{
			`sum5(1, 2)`,
			`sum5(1, 2)`,
			1,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		pg := p.ParseProgram()
		checkParserErrors(t, p)
		actual := pg.String()

		assert.Equal(t, tt.expect, actual)
		assert.Equal(t, tt.expectLen, len(pg.Statements))
	}
}

func TestNextToken(t *testing.T) {
	l := lexer.New(`"hello world"`)
	p := New(l)

	expectCur := token.Token{Type: "STRING", Literal: "hello world"}
	assert.Equal(t, expectCur, p.curToken)
	expectPeek := token.Token{Type: "EOF", Literal: ""}
	assert.Equal(t, expectPeek, p.peekToken)
}

func TestParseExpressionList(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{`(1,2,3)`, "1, 2, 3"},
		{`(1, 2, 3)`, "1, 2, 3"},
		{`(1+1,2,3)`, "(1 + 1), 2, 3"},
		{`(1,2,3`, ""},
		{`1,2,3)`, ""},
		{`()`, ""},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		list := p.parseExpressionList(token.RPAREN)

		results := []string{}
		for _, e := range list {
			results = append(results, e.String())
		}

		assert.Equal(t, tt.expect, strings.Join(results, ", "))
	}
}
