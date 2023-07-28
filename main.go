package main

import (
	"C"
	"fmt"
)
import (
	"bufio"
	"os"
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

func main() {
	var str string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str = scanner.Text()
	}

	l := NewLexer(str)
	l.Next()
}
