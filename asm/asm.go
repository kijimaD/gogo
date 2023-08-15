package asm

import (
	"fmt"
	"log"

	"github.com/kijimaD/gogo/ast"
	"github.com/kijimaD/gogo/object"
	"github.com/kijimaD/gogo/parser"
	"github.com/kijimaD/gogo/token"
)

var Vpos = 1

// 変数の位置をアセンブリコードの中で正しいメモリアドレスに変換するための定数。1つの変数は4バイトに格納されている
const varWidth = 4

func emitBinop(env *object.Environment, i ast.InfixExpression) {
	var op string
	switch i.Operator {
	case token.PLUS:
		op = "add"
	case token.MINUS:
		op = "sub"
	case token.ASTERISK:
		op = "imul"
	case token.SLASH:
		op = ""
	default:
		log.Fatal("invalid operand:", op)
	}

	if i.Operator == token.SLASH {
		EmitExpr(env, i.Left)
		fmt.Printf("push %%rax\n\t")
		EmitExpr(env, i.Right)
		fmt.Printf("mov %%eax, %%ebx\n\t")
		fmt.Printf("pop %%rax\n\t")
		fmt.Printf("mov $0, %%edx\n\t")
		fmt.Printf("idiv %%ebx\n\t")
	} else {
		EmitExpr(env, i.Right)
		fmt.Printf("push %%rax\n\t")
		EmitExpr(env, i.Left)
		fmt.Printf("pop %%rbx\n\t")
		fmt.Printf("%s %%ebx, %%eax\n\t", op)
	}
}

func EvalDeclStmt(e *object.Environment, ds *ast.DeclStatement) {
	obj := &object.String{Value: ds.Name.Token.Literal, Pos: Vpos}
	e.Set(ds.Name.Token.Literal, obj)
	fmt.Printf("mov %%eax, -%d(%%rbp)\n\t", Vpos*varWidth)
	Vpos++
}

func evalIdentifier(e *object.Environment, ident *ast.Identifier) {
	result, ok := e.Get(ident.Token.Literal)
	if !ok {
		log.Fatal("not exist variable: ", ident.Token.Literal)
	}
	fmt.Printf("mov %%eax, -%d(%%rbp)\n\t", result.CurPos()*varWidth)
}

// 定義した文字列にデータラベルをつける
func EmitDataSection(p *parser.Parser) {
	if len(p.Strs) == 0 {
		return
	}
	for i, str := range p.Strs {
		fmt.Printf("\t.data\n")
		fmt.Printf(".s%d:\n\t", i)
		fmt.Printf(".string \"")
		fmt.Printf(`%s`, str)
		fmt.Printf("\"\n")
	}
	fmt.Printf("\t")
}

func EmitExpr(env *object.Environment, node ast.Node) {
	switch n := node.(type) {
	case *ast.IntegerLiteral:
		fmt.Printf("mov $%d, %%eax\n\t", int(n.Value))
	case *ast.StringLiteral:
		fmt.Printf("lea .s%d(%%rip), %%rax\n\t", n.ID)
	case *ast.Identifier:
		evalIdentifier(env, n)
	case *ast.InfixExpression:
		emitBinop(env, *n)
	}
}
