package main

import (
	"fmt"
	"strconv"
)

func compileNumber(i int) {
	fmt.Printf(".text\n\t")
	fmt.Printf(".global intfn\n")
	fmt.Printf("intfn:\n\t")
	fmt.Printf("mov $%d, %%rax\n\t", i)
	fmt.Printf("ret\n")
}

func main() {
	var str string
	fmt.Scan(&str)

	i, ierr := strconv.Atoi(str)
	if ierr == nil {
		compileNumber(i)
	}
}
