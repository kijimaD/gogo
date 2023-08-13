package asm

import (
	"fmt"
	"log"

	"github.com/kijimaD/gogo/ast"
	"github.com/kijimaD/gogo/object"
)

var Vpos = 1

func emitIntexpr(env *object.Environment, n ast.Node) {
	switch nt := n.(type) {
	case *ast.IntegerLiteral:
		fmt.Printf("mov $%s, %%eax\n\t", nt.String())
	case *ast.Identifier:
		evalIdentifier(env, nt)
	case *ast.InfixExpression:
		emitBinop(env, *nt)
	default:
		log.Fatal("not cover type:")
	}
}

func emitString(e ast.StringLiteral) {
	fmt.Printf("\t.data\n")
	fmt.Printf(".mydata:\n\t")
	fmt.Printf(".string \"%s\"\n\t", e.String())
	fmt.Printf(".text\n\t")
	fmt.Printf(".global stringfn\n")
	fmt.Printf("stringfn:\n\t")
	fmt.Printf("lea .mydata(%%rip), %%rax\n\t")
	fmt.Printf("ret\n")
}

func emitBinop(e *object.Environment, i ast.InfixExpression) {
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

	if i.Operator == "/" {
		emitIntexpr(e, i.Left)
		fmt.Printf("push %%rax\n\t")
		emitIntexpr(e, i.Right)
		fmt.Printf("mov %%eax, %%ebx\n\t")
		fmt.Printf("pop %%rax\n\t")
		fmt.Printf("mov $0, %%edx\n\t")
		fmt.Printf("idiv %%ebx\n\t")
	} else {
		emitIntexpr(e, i.Right)
		fmt.Printf("push %%rax\n\t")
		emitIntexpr(e, i.Left)
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

func Compile(env *object.Environment, n ast.Node) {
	switch node := n.(type) {
	case *ast.StringLiteral:
		emitString(*node)
	case *ast.InfixExpression, *ast.IntegerLiteral:
		fmt.Printf(".text\n\t")
		fmt.Printf(".global intfn\n")
		fmt.Printf("intfn:\n\t")
		emitIntexpr(env, node)
		fmt.Printf("ret\n")
	case *ast.Identifier:
		resultobj, _ := env.Get(node.Value)
		fmt.Printf("mov -%d(%%rbp), %%eax\n\t", resultobj.CurPos()*4)
	case *ast.DeclStatement:
		evalDeclStmt(env, node)
	}
}
