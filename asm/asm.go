package asm

import (
	"fmt"
	"log"

	"github.com/kijimaD/gogo/ast"
	"github.com/kijimaD/gogo/object"
)

var Vpos = 1

func emitBinop(env *object.Environment, i ast.InfixExpression) {
	var op string
	switch i.Operator {
	case "+":
		op = "add"
	case "-":
		op = "sub"
	case "*":
		op = "imul"
	case "/":
		op = ""
	default:
		log.Fatal("invalid operand:", op)
	}

	EmitExpr(env, i.Left)
	fmt.Printf("push %%rax\n\t")
	EmitExpr(env, i.Right)

	if i.Operator == "/" {
		fmt.Printf("mov %%eax, %%ebx\n\t")
		fmt.Printf("pop %%rax\n\t")
		fmt.Printf("mov $0, %%edx\n\t")
		fmt.Printf("idiv %%ebx\n\t")
	} else {
		fmt.Printf("pop %%rbx\n\t")
		fmt.Printf("%s %%ebx, %%eax\n\t", op)
	}
}

func evalDeclStmt(e *object.Environment, ds *ast.DeclStatement) {
	obj := &object.String{Value: ds.Name.Value, Pos: Vpos}
	e.Set(ds.Name.Value, obj)
	fmt.Printf("mov %%eax, -%d(%%rbp)\n\t", Vpos*4)
	Vpos++
}

func evalIdentifier(e *object.Environment, ident *ast.Identifier) {
	result, ok := e.Get(ident.Value)
	if !ok {
		log.Fatal("not exist variable: ", ident.Value)
	}
	fmt.Printf("mov %%eax, -%d(%%rbp)\n\t", result.CurPos()*4)
}

func EmitExpr(env *object.Environment, node ast.Node) {
	switch n := node.(type) {
	case *ast.IntegerLiteral:
		fmt.Printf("mov $%d, %%eax\n\t", int(n.Value))
	case *ast.StringLiteral:
		fmt.Printf("lea .s%d(%%rip), %%rax\n\t", 1) // string IDを入れる
	case *ast.Identifier:
		resultobj, _ := env.Get(n.Value)
		fmt.Printf("mov -%d(%%rbp), %%eax\n\t", resultobj.CurPos()*4)
	case *ast.DeclStatement:
		evalDeclStmt(env, n)
	case *ast.InfixExpression:
		emitBinop(env, *n)
	}
}
