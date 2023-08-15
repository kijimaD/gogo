package ast

import (
	"bytes"

	"github.com/kijimaD/gogo/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	ExpressionNode()
}

// 構文解析器が生成する全てのASTのルートノードになる
// 全ての有効なmonekyプログラムは、ひと続きの文の集まり
type Program struct {
	Statements []Statement
}

// 文字列表示してデバッグしやすいようにする
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type Identifier struct {
	Token token.Token
}

func (i *Identifier) ExpressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Token.Literal }

type ExpressionStatement struct {
	Token      token.Token // 式の最初のトークン
	Expression Expression  // 式を保持
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type StringLiteral struct {
	Token token.Token
	Value string
	ID    int
}

func (sl *StringLiteral) ExpressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Token.Literal + "\"" }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) ExpressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) ExpressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// int a = 1;
type DeclStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (de *DeclStatement) statementNode()       {}
func (de *DeclStatement) TokenLiteral() string { return de.Token.Literal }
func (de *DeclStatement) String() string {
	var out bytes.Buffer
	out.WriteString(de.TokenLiteral() + " " + de.Name.Token.Literal)
	out.WriteString(" = ")
	out.WriteString(de.Value.String())

	return out.String()
}
