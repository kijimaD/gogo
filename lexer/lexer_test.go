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
			name:   "ダブルクォートがなく数字でもない",
			input:  `naked string`,
			expect: token.ILLEGAL,
		},
		{
			name:   "ダブルクォートのペアが合わない",
			input:  `"not pair`,
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
