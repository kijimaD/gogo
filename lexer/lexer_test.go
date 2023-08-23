package lexer

import (
	"testing"

	"github.com/kijimaD/gogo/token"
	"github.com/stretchr/testify/assert"
)

func TestReadChar(t *testing.T) {
	l := New("hi")
	assert.Equal(t, uint8('h'), l.ch)
	l.readChar()
	assert.Equal(t, uint8('i'), l.ch)
	l.readChar()
	assert.Equal(t, uint8(0), l.ch)
}

func TestReadString(t *testing.T) {
	l := New(`"hello" "world"`)
	_, actual := l.readString()
	expect := "hello"
	assert.Equal(t, expect, actual)

	l.readChar()
	l.readChar()

	_, actual = l.readString()
	expect = "world"
	assert.Equal(t, expect, actual)
}

// ダブルクォートのペアがあっていない場合はエラー
func TestReadStringFail(t *testing.T) {
	l := New(`"hello`)
	err, actual := l.readString()

	expect := ``
	assert.Equal(t, expect, actual)
	assert.Error(t, err)
}

func TestReadNumber1(t *testing.T) {
	l := New("12 34")
	actual := l.readNumber()
	expect := "12"
	assert.Equal(t, expect, actual)

	l.readChar()

	actual = l.readNumber()
	expect = "34"
	assert.Equal(t, expect, actual)
}

func TestReadNumber2(t *testing.T) {
	l := New("12a")
	actual := l.readNumber()
	expect := "12"
	assert.Equal(t, expect, actual)
}

func TestSkipSpace(t *testing.T) {
	l := New(`   123`)
	assert.Equal(t, uint8(' '), l.ch)
	l.skipSpace()
	assert.Equal(t, uint8('1'), l.ch)
}

func TestIllegal(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect token.TokenType
	}{
		{
			name:   "ダブルクォートのペアが合わない",
			input:  `"not pair`,
			expect: token.ILLEGAL,
		},
		{
			name:   "シングルクォートのペアが合わない",
			input:  `'not pair`,
			expect: token.ILLEGAL,
		},
		{
			name:   "シングルクォートに複数の文字",
			input:  `'MULTIPLE CHAR'`,
			expect: token.ILLEGAL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			actual := l.NextToken()
			if actual.Type != tt.expect {
				t.Errorf("got %s want %s", actual, tt.expect)
			}
		})
	}
}

func TestNextToken(t *testing.T) {
	input := `1 + 2;
3 * 4;
5 / 6;
7 - 8;
a = 1;
"hello";
'c';
	1  ;
a;
abc;
a + b;
42a;
a42;
f(1);
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "1"},
		{token.PLUS, "+"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},

		{token.INT, "3"},
		{token.ASTERISK, "*"},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.SLASH, "/"},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},

		{token.INT, "7"},
		{token.MINUS, "-"},
		{token.INT, "8"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},

		{token.STRING, "hello"},
		{token.SEMICOLON, ";"},

		{token.CHAR, "c"},
		{token.SEMICOLON, ";"},

		{token.INT, "1"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "a"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "abc"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},

		{token.ILLEGAL, "a"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "a42"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "f"},
		{token.LPAREN, "("},
		{token.INT, "1"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)
	for _, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}

func TestIsLetter(t *testing.T) {
	assert.True(t, isLetter('a'))
	assert.True(t, isLetter('B'))
	assert.False(t, isLetter('1'))
}
