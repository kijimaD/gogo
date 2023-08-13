package main

import (
	"C"
)
import (
	"bufio"
	"log"
	"os"

	"github.com/kijimaD/gogo/asm"
	"github.com/kijimaD/gogo/ast"
	"github.com/kijimaD/gogo/lexer"
	"github.com/kijimaD/gogo/parser"
)

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
	asm.Compile(exp)
}
