package main

import (
	"C"
	"fmt"
)
import (
	"bufio"
	"os"

	"github.com/kijimaD/gogo/lexer"
	"github.com/kijimaD/gogo/token"
)

func compileNumber(s string) {
	fmt.Printf(".text\n\t")
	fmt.Printf(".global intfn\n")
	fmt.Printf("intfn:\n\t")
	fmt.Printf("mov $%s, %%rax\n\t", s)
	fmt.Printf("ret\n")
}

func compileString(s string) {

}

func printQuote(s string) {
	for _, c := range s {
		if c == '"' || c == '\\' {
			fmt.Print("\\")
		}
		fmt.Printf("%c", c)
	}
}

func emitString(tok token.Token) {
	fmt.Printf("\t.data\n")
	fmt.Printf(".mydata:\n\t")
	fmt.Printf(".string ")
	printQuote(tok.Literal)
	fmt.Printf("\n\t")
	fmt.Printf(".text\n\t")
	fmt.Printf(".global stringfn\n")
	fmt.Printf("stringfn:\n\t")
	fmt.Printf("lea .mydata(%%rip), %%rax\n\t")
	fmt.Printf("ret\n")
	return
}

func compile(token.Token) {

}

func main() {
	var str string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str = scanner.Text()
	}

	l := lexer.NewLexer(str)
	l.NextToken()
}
