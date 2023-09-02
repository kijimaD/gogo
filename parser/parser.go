package parser

import (
	"fmt"
	"log"
	"strconv"

	"github.com/kijimaD/gogo/ast"
	"github.com/kijimaD/gogo/lexer"
	"github.com/kijimaD/gogo/object"
	"github.com/kijimaD/gogo/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token // 現在のトークン
	peekToken token.Token // 次のトークン

	// 構文解析関数は中置もしくは前置のマップどちらかにある
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	Strs   []string // 定義済みの文字列一覧。ラベルの定義に使う。スタックに入っているので、位置が必要
	errors []string
	Env    *object.Environment // パーサーから移動させたほうがいいかもしれない
}

func (p *Parser) Errors() []string {
	return p.errors
}

type (
	// どちらの関数もast.Expressionを返す。これが欲しいもの

	// 前置構文解析関数 ++1
	// 前置演算子には「左側」が存在しない
	prefixParseFn func() ast.Expression

	// 中置構文解析関数 n + 1
	// 引数は中置演算子の「左側」
	infixParseFn func(ast.Expression) ast.Expression
)

const (
	_int = iota
	LOWEST
	SUM
	PRODUCT
	CALL
)

var precedences = map[token.TokenType]int{
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   CALL,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		Strs:   []string{},
		errors: []string{},
		Env:    object.NewEnvironment(),
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.CHAR, p.parseCharLiteral)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.IDENT, p.parseIdent)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	// 2つトークンを読み込む。curTokenとpeekTokenの両方がセットされる
	p.nextToken()
	p.nextToken()

	return p
}

// 次のトークンに進む
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// peekTokenの型をチェックし、その型が正しい場合に限ってnextTokenを読んで、トークンを進める
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t,
		p.peekToken.Type,
	)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		switch p.curToken.Type {
		case token.ILLEGAL:
			p.errors = append(p.errors, "illegal token is detected!")
		case token.SEMICOLON:
		default:
			stmt := p.parseStatement()
			if stmt != nil {
				program.Statements = append(program.Statements, stmt)
			}
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// 文をパースする
// 文は代入とか、ifの実行文とか(条件部分は式)、返り値がないもの
func (p *Parser) parseStatement() ast.Statement {
	if p.curToken.Type == token.IDENT && p.isCtypeKeyword() {
		if decl := p.parseDeclStatement(); decl != nil {
			return decl
		}
	}
	return p.parseExpressionStatement()
}

// 式文
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	return stmt
}

// int a = 1
func (p *Parser) parseDeclStatement() *ast.DeclStatement {
	declstmt := &ast.DeclStatement{Token: p.curToken}

	ctype, err := p.getDeclCtype()
	if err != nil {
		p.errors = append(p.errors, "failed get ident type")
	}
	declstmt.Ctype = ctype

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	declstmt.Name = &ast.Var{Token: p.curToken}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	declstmt.Value = p.parseExpression(LOWEST)
	declstmt.Pos = p.Env.VarPos

	obj := &object.String{Value: declstmt.Name.Token.Literal, Pos: p.Env.VarPos} // TODO: とりあえずstring。ctype型によって変える
	p.Env.Set(declstmt.Name.Token.Literal, obj)

	return declstmt
}

// 式をパースする。現在位置に対応したパース関数を適用してASTを返す
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// 前置構文
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		msg := fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	leftExp := prefix()

	// 次のトークンの優先度が高く中置構文に対応してるなら、中置構文としてパースする
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()            // 中置関数の演算子のトークンに移動
		leftExp = infix(leftExp) // 中置関数の演算子をパースする
	}

	return leftExp
}

func (p *Parser) parseStringLiteral() ast.Expression {
	id := len(p.Strs)
	strlit := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal, ID: id}
	p.Strs = append(p.Strs, strlit.Value)
	return strlit
}

func (p *Parser) parseCharLiteral() ast.Expression {
	runes := []rune(p.curToken.Literal)
	a := ast.CharLiteral{Token: p.curToken, Value: runes[0]}
	return &a
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// token.identifierをパースする。変数の場合は値が存在するかチェックする
func (p *Parser) parseIdent() ast.Expression {
	if !p.peekTokenIs(token.LPAREN) {
		_, ok := p.Env.Get(p.curToken.Literal)
		if !ok {
			msg := fmt.Sprintf("not exist variable: %s", p.curToken.Literal)
			p.errors = append(p.errors, msg)
		}
	}
	// 前置関数と中置関数の仕組みで、処理しているトークンが関数呼び出しの場合はここの返り値は使われることがない
	a := &ast.Var{Token: p.curToken}
	return a
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken, // 現在のトークンは中置演算子の演算子
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()                                    // 中置演算子の右の引数に進む
	expression.Right = p.parseExpression(precedence) // 右側を評価する

	ctype, err := p.resultType(left, expression.Right)
	if err != nil {
		log.Fatal("type error")
	}
	expression.Ctype = ctype

	return expression
}

// 関数呼び出しは"("を真ん中とする中置構文。
// <f><(><args>
// 前  中  後
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.FuncallExpression{Token: p.curToken, Function: function}
	exp.Args = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	// リストの終端が来たら、次に進んで終了
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	// 次のトークンがコンマのときだけ繰り返すので、リストの最後の要素で止まる
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // [1<,> 2]
		p.nextToken() // [1, <2>]
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// 型宣言がどの型かを判定する
func (p *Parser) getDeclCtype() (token.Ctype, error) {
	if p.curToken.Type != token.IDENT {
		return token.CTYPE_VOID, fmt.Errorf("%s is not ident", p.curToken.Type)
	}

	switch p.curToken.Literal {
	case "int":
		return token.CTYPE_INT, nil
	case "char":
		return token.CTYPE_CHAR, nil
	case "string":
		return token.CTYPE_STR, nil
	default:
		return token.CTYPE_VOID, nil
	}
}

// 識別子が型か判定する
func (p *Parser) isCtypeKeyword() bool {
	ctype, err := p.getDeclCtype()
	if err != nil {
		log.Fatal(err)
	}
	return ctype != token.CTYPE_VOID
}

func (p *Parser) resultType(a ast.Expression, b ast.Expression) (token.Ctype, error) {
	small := a
	big := b
	incompatibleErr := fmt.Errorf("incompatible operands: %s and %s for %c", p.curToken, small, big)

	if a.GetCtype() > b.GetCtype() {
		small = b
		big = a
	}

	switch small.GetCtype() {
	case token.CTYPE_VOID:
		return token.CTYPE_VOID, incompatibleErr
	case token.CTYPE_INT:
		switch big.GetCtype() {
		case token.CTYPE_INT:
			return token.CTYPE_INT, nil
		case token.CTYPE_CHAR:
			return token.CTYPE_INT, nil
		case token.CTYPE_STR:
			return token.CTYPE_VOID, incompatibleErr
		default:
			log.Fatal("unknown type")
		}
	case token.CTYPE_CHAR:
		switch big.GetCtype() {
		case token.CTYPE_CHAR:
			return token.CTYPE_INT, nil
		case token.CTYPE_STR:
			return token.CTYPE_VOID, incompatibleErr
		default:
			log.Fatal("unknown type")
		}
	case token.CTYPE_STR:
		log.Fatal("unknown type")
	default:
		log.Fatal("unknown type")
	}

	return token.CTYPE_VOID, incompatibleErr
}
