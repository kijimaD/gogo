package main

import (
	"C"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func compileNumber(i int) {
	fmt.Printf(".text\n\t")
	fmt.Printf(".global intfn\n")
	fmt.Printf("intfn:\n\t")
	fmt.Printf("mov $%d, %%rax\n\t", i)
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

	i, ierr := strconv.Atoi(str)
	if ierr == nil {
		compileNumber(i)
	}
	if str != "" {
		compileString(str)
	}
}
