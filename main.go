package main

import (
	"C"
	"fmt"
)
import (
	"bufio"
	"log"
	"os"

	"github.com/kijimaD/gogo/ast"
	"github.com/kijimaD/gogo/lexer"
	"github.com/kijimaD/gogo/parser"
)

func emitIntexpr(e ast.Node) {
	switch ast := e.(type) {
	case *ast.IntegerLiteral:
		fmt.Printf("mov $%s, %%eax\n\t", ast.String())
	case *ast.InfixExpression:
		emitBinop(*ast)
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

func emitBinop(i ast.InfixExpression) {
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
		emitIntexpr(i.Left)
		fmt.Printf("push %%rax\n\t")
		emitIntexpr(i.Right)
		fmt.Printf("mov %%eax, %%ebx\n\t")
		fmt.Printf("pop %%rax\n\t")
		fmt.Printf("mov $0, %%edx\n\t")
		fmt.Printf("idiv %%ebx\n\t")
	} else {
		emitIntexpr(i.Right)
		fmt.Printf("push %%rax\n\t")
		emitIntexpr(i.Left)
		fmt.Printf("pop %%rbx\n\t")
		fmt.Printf("%s %%ebx, %%eax\n\t", op)
	}
}

func compile(e ast.Node) {
	switch ast := e.(type) {
	case *ast.StringLiteral:
		emitString(*ast)
	default:
		fmt.Printf(".text\n\t")
		fmt.Printf(".global intfn\n")
		fmt.Printf("intfn:\n\t")
		emitIntexpr(e)
		fmt.Printf("ret\n")
	}
}

func main() {
	var str string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str = scanner.Text()
	}

	l := lexer.New(str)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			log.Fatal(err)
		}
	}
	stmt, _ := prog.Statements[0].(*ast.ExpressionStatement)
	exp := stmt.Expression
	compile(exp)
}
