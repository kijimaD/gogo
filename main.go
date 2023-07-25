package main

import "fmt"

func main() {
	var i int
	fmt.Scan(&i)

	fmt.Printf(".global mymain\n")
	fmt.Printf("mymain:\n")
	fmt.Printf("\tmov $%d, %%eax\n", i)
	fmt.Printf("\tret\n")
}
