package parser

import (
	"fmt"
	"strconv"

	"github.com/kijimaD/gogo/ast"
	"github.com/kijimaD/gogo/lexer"
	"github.com/kijimaD/gogo/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token // 現在のトークン
	peekToken token.Token // 次のトークン

	// 構文解析関数は中置もしくは前置のマップどちらかにある
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	errors []string
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
)

var precedences = map[token.TokenType]int{
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.IDENT, p.parseIdent)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)

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
	switch p.curToken.Type {
	case token.IDENT:
		// TODO: 型をテーブルから参照する
		if p.curToken.Literal == token.CTYPE_VOID ||
			p.curToken.Literal == token.CTYPE_INT ||
			p.curToken.Literal == token.CTYPE_CHAR ||
			p.curToken.Literal == token.CTYPE_STR {
			decl := p.parseDeclStatement()
			if decl != nil {
				return decl
			}
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
	stmt := &ast.DeclStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
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

var Strings = []string{}

func (p *Parser) parseStringLiteral() ast.Expression {
	id := len(Strings)
	strlit := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal, ID: id}
	Strings = append(Strings, strlit.Value)
	return strlit
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

func (p *Parser) parseIdent() ast.Expression {
	ident := &ast.Identifier{Token: p.curToken}
	return ident
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken, // 現在のトークンは中置演算子の演算子
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()                                    // 中置演算子の右の引数に進む
	expression.Right = p.parseExpression(precedence) // 右側を評価する

	return expression
}
