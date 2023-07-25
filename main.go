package main

import "fmt"

func main() {
	fmt.Printf("\t.globl main\n")
	fmt.Printf("main:\n")
	fmt.Printf("\tmovl\t$%d, %%eax\n", 0)
	fmt.Printf("\tret\n")
}
