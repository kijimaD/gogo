package main

import (
	"C"
)
import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/kijimaD/gogo/asm"
	"github.com/kijimaD/gogo/lexer"
	"github.com/kijimaD/gogo/object"
	"github.com/kijimaD/gogo/parser"
)

func main() {
	var str string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str = scanner.Text()
	}

	env := object.NewEnvironment()

	l := lexer.New(str)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			log.Fatal(err)
		}
	}
	asm.EmitDataSection(p)
	fmt.Printf(".text\n\t")
	fmt.Printf(".global mymain\n")
	fmt.Printf("mymain:\n\t")

	for _, stmt := range prog.Statements {
		asm.EmitStmt(env, stmt)
	}
	fmt.Printf("ret\n")
}
