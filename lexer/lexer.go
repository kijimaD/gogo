package lexer

import (
	"fmt"

	"github.com/kijimaD/gogo/token"
)

type Lexer struct {
	input        string
	position     int // 現在検査中のバイトchの位置
	readPosition int // 入力における次の位置
	ch           byte
}

// ソースコード文字列を引数に取り、初期化する
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// 現在位置の文字を読み込む
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipSpace()

	switch l.ch {
	case '"':
		tok.Type = token.STRING
		err, literal := l.readString()
		if err != nil {
			tok = newToken(token.ILLEGAL, l.ch)
		}
		tok.Literal = literal
	case '\'':
		tok.Type = token.CHAR
		if l.ch == '\'' {
			l.readChar() // 左のシングルクォートを飛ばす <'>a'
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
		tok.Literal = string(l.ch)
		l.readChar() // 本体を飛ばす '<a>'
		if l.ch != '\'' {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case 0:
		// 終端文字
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT

			// 数字の次の文字がスペース区切りなしで非数字だったら文法エラー -- 42a など
			// これはパーサーでやることではないような気もする
			// 評価時にきめることではないのか?
			if isLetter(l.ch) {
				tok = newToken(token.ILLEGAL, l.ch)
				l.readChar()
				return tok
			}

			return tok // readNumberは "1+2"で1にあったとき現在値を+に進めているので、この関数の最終行で1文字余計に進めないようにreturnが必要
		} else if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.IDENT
			// TODO: ここでIDENTを組み込みのものか判断すればよさそう
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tok
}

// 文字列をすべて読んで、次の非文字列の領域に現在地を進める
func (l *Lexer) readString() (error, string) {
	startPos := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}

		// ダブルクォートがペアにならずに終端するとエラー
		if l.ch == 0 {
			return fmt.Errorf("unexpected EOF"), ""
		}
	}
	return nil, l.input[startPos:l.position]
}

// 数字をすべて読んで、次の非数字の領域に現在地を進める
// "1+2" 1で実行したとき、現在地を+にすすめる
func (l *Lexer) readNumber() string {
	startPos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

// identの最初の文字はアルファベットでないといけない。2文字目からは数字が使える
func (l *Lexer) readIdentifier() string {
	startPos := l.position
	if isLetter(l.ch) {
		for isLetter(l.ch) || isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[startPos:l.position]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
