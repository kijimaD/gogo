package main

type Lexer struct {
	input        string
	position     int // 現在検査中のバイトchの位置
	readPosition int // 入力における次の位置
	ch           byte
}

// ソースコード文字列を引数に取り、初期化する
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// 現在の文字を読み込む
func (l *Lexer) Next() {
	l.skipSpace()

	switch l.ch {
	case '"':
		s := l.readString()
		compileString(s)
	default:
		if isDigit(l.ch) {
			s := l.readNumber()
			compileNumber(s)
			return // readNumberは "1+2"で1にあったとき現在値を+に進めるので、この関数の最終行でまた進めないようにreturnが必要
		}
	}

	l.readChar()
}

// 次の1文字を読んでinput文字列の現在位置を進める
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCIIコードの"NUL"文字に対応している
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// 文字列をすべて読んで、次の非文字列の領域に現在地を進める
func (l *Lexer) readString() string {
	startPos := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[startPos:l.position]
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

func (l *Lexer) skipSpace() {
	for l.ch == ' ' {
		l.readChar()
	}
}

// 数字か判定する
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isSpace(ch byte) bool {
	return ch == ' '
}
