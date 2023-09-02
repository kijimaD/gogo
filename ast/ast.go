package ast

import (
	"bytes"
	"strings"

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
	GetCtype() token.Ctype
}

// 構文解析器が生成する全てのASTのルートノードになる
// 全ての有効なプログラムは、ひと続きの文の集まり
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

type Var struct {
	Token token.Token
	Pos   int
	Ctype token.Ctype
}

func (v *Var) ExpressionNode()       {}
func (v *Var) TokenLiteral() string  { return v.Token.Literal }
func (v *Var) String() string        { return v.Token.Literal }
func (v *Var) GetCtype() token.Ctype { return v.Ctype }

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

func (sl *StringLiteral) ExpressionNode()       {}
func (sl *StringLiteral) TokenLiteral() string  { return sl.Token.Literal }
func (sl *StringLiteral) String() string        { return "\"" + sl.Token.Literal + "\"" }
func (sl *StringLiteral) GetCtype() token.Ctype { return token.CTYPE_STR }

type CharLiteral struct {
	Token token.Token
	Value rune
}

func (cl *CharLiteral) ExpressionNode()       {}
func (cl *CharLiteral) TokenLiteral() string  { return cl.Token.Literal }
func (cl *CharLiteral) String() string        { return `'` + cl.Token.Literal + `'` }
func (cl *CharLiteral) GetCtype() token.Ctype { return token.CTYPE_CHAR }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) ExpressionNode()       {}
func (il *IntegerLiteral) TokenLiteral() string  { return il.Token.Literal }
func (il *IntegerLiteral) String() string        { return il.Token.Literal }
func (il *IntegerLiteral) GetCtype() token.Ctype { return token.CTYPE_INT }

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
	Ctype    token.Ctype
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
func (ie *InfixExpression) GetCtype() token.Ctype { return ie.Ctype }

// int a = 1;
type DeclStatement struct {
	Token token.Token
	Name  *Var
	Value Expression
	Pos   int
	Ctype token.Ctype
}

func (de *DeclStatement) statementNode()       {}
func (de *DeclStatement) TokenLiteral() string { return de.Token.Literal }
func (de *DeclStatement) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(de.TokenLiteral() + " " + de.Name.Token.Literal)
	out.WriteString(" = ")
	out.WriteString(de.Value.String())
	out.WriteString(")")

	return out.String()
}

// f(20, 5)
type FuncallExpression struct {
	Token    token.Token // "("
	Function Expression
	Args     []Expression
}

func (fe *FuncallExpression) ExpressionNode()      {}
func (fe *FuncallExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *FuncallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range fe.Args {
		args = append(args, a.String())
	}
	out.WriteString(fe.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
func (fe *FuncallExpression) GetCtype() token.Ctype { return token.CTYPE_INT } // TODO: とりあえず返り値がintしかないのでハードコーディング
