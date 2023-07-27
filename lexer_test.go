package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadChar(t *testing.T) {
	l := NewLexer("hi")
	assert.Equal(t, uint8('h'), l.ch)
	l.readChar()
	assert.Equal(t, uint8('i'), l.ch)
	l.readChar()
	assert.Equal(t, uint8(0), l.ch)
}

func TestReadString(t *testing.T) {
	l := NewLexer(`"hello" "world"`)
	actual := l.readString()
	expect := "hello"
	assert.Equal(t, expect, actual)

	l.readChar()
	l.readChar()

	actual = l.readString()
	expect = "world"
	assert.Equal(t, expect, actual)
}

func TestReadStringFail(t *testing.T) {
	l := NewLexer(`"hello`)
	actual := l.readString()
	expect := `hello`
	assert.Equal(t, expect, actual)
}

func TestReadNumber1(t *testing.T) {
	l := NewLexer("12 34")
	actual := l.readNumber()
	expect := "12"
	assert.Equal(t, expect, actual)

	l.readChar()

	actual = l.readNumber()
	expect = "34"
	assert.Equal(t, expect, actual)
}

// 暫定の挙動
func TestReadNumber2(t *testing.T) {
	l := NewLexer("12a")
	actual := l.readNumber()
	expect := "12"
	assert.Equal(t, expect, actual)
}
