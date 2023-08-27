package asm

import (
	"fmt"
	"log"

	"github.com/kijimaD/gogo/ast"
	"github.com/kijimaD/gogo/parser"
	"github.com/kijimaD/gogo/token"
)

var varPos = 1

// 変数の位置をアセンブリコードの中で正しいメモリアドレスに変換するための定数。1つの変数は4バイトに格納されている
const varWidth = 4

var regs = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func emitBinop(i ast.InfixExpression) {
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
		emitExpr(i.Left)
		fmt.Printf("push %%rax\n\t")
		emitExpr(i.Right)
		fmt.Printf("mov %%eax, %%ebx\n\t")
		fmt.Printf("pop %%rax\n\t")
		fmt.Printf("mov $0, %%edx\n\t")
		fmt.Printf("idiv %%ebx\n\t")
	} else {
		emitExpr(i.Right)
		fmt.Printf("push %%rax\n\t")
		emitExpr(i.Left)
		fmt.Printf("pop %%rbx\n\t")
		fmt.Printf("%s %%ebx, %%eax\n\t", op)
	}
}

func emitDeclStmt(ds *ast.DeclStatement) {
	fmt.Printf("mov %%eax, -%d(%%rbp)\n\t", ds.Pos*varWidth)
}

// TODO: varに変える
func evalIdentifier(ident *ast.Identifier) {
	fmt.Printf("mov %%eax, -%d(%%rbp)\n\t", ident.Pos*varWidth)
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

func EmitStmt(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		exp := s.Expression
		emitExpr(exp)
	case *ast.DeclStatement:
		exp := s.Value
		emitExpr(exp)
		emitDeclStmt(s)
	default:
		log.Fatal("not support statement:", s)
	}
}

func emitExpr(node ast.Node) {
	switch n := node.(type) {
	case *ast.IntegerLiteral:
		fmt.Printf("mov $%d, %%eax\n\t", int(n.Value))
	case *ast.StringLiteral:
		fmt.Printf("lea .s%d(%%rip), %%rax\n\t", n.ID)
	case *ast.CharLiteral:
		fmt.Printf("mov $%d, %%eax\n\t", n.Value)
	case *ast.Identifier:
		evalIdentifier(n)
	case *ast.InfixExpression:
		emitBinop(*n)
	case *ast.FuncallExpression:
		for i := 1; i < len(n.Args); i++ {
			fmt.Printf("push %%%s\n\t", regs[i])
		}
		for i := 0; i < len(n.Args); i++ {
			emitExpr(n.Args[i])
			fmt.Printf("push %%rax\n\t")
		}
		for i := len(n.Args) - 1; i >= 0; i-- {
			fmt.Printf("pop %%%s\n\t", regs[i])
		}
		fmt.Printf("mov $0, %%eax\n\t")
		fmt.Printf("call %s\n\t", n.Function.String())
		for i := len(n.Args) - 1; i > 0; i-- {
			fmt.Printf("pop %%%s\n\t", regs[i])
		}
	}
}
