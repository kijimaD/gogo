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

func compileNumber(s string) {
	fmt.Printf(".text\n\t")
	fmt.Printf(".global intfn\n")
	fmt.Printf("intfn:\n\t")
	fmt.Printf("mov $%s, %%rax\n\t", s)
	fmt.Printf("ret\n")
}

func compileString(s string) {
	fmt.Printf("\t.data\n")
	fmt.Printf(".mydata:\n\t")
	fmt.Printf(".string \"%s\"\n\t", s)
	fmt.Printf(".text\n\t")
	fmt.Printf(".global stringfn\n")
	fmt.Printf("stringfn:\n\t")
	fmt.Printf("lea .mydata(%%rip), %%rax\n\t")
	fmt.Printf("ret\n")
}

func printQuote(s string) {
	for _, c := range s {
		if c == '"' || c == '\\' {
			fmt.Print("\\")
		}
		fmt.Printf("%c", c)
	}
}

func main() {
	var str string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str = scanner.Text()
	}

	l := lexer.NewLexer(str)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			log.Fatal(err)
		}
	}
	stmt, _ := prog.Statements[0].(*ast.ExpressionStatement)
	switch node := stmt.Expression.(type) {
	case *ast.IntegerLiteral:
		compileNumber(node.String())
	case *ast.StringLiteral:
		compileString(node.String())
	}

}
