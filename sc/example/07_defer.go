package main

import (
	"fmt"
)

func main() {
	defer func() {
		fmt.Println("first")
	}()

	defer func() {
		fmt.Println("second")
	}()

	fmt.Println("done")

	fmt.Println(triple(4)) // 12
}

func double(x int) (result int) {
	defer func() {
		fmt.Printf("double(%d) = %d\n", x, result)
	}()

	return x + x
}

func triple(x int) (result int) {
	defer func() {
		result += x
	}()

	return double(x)
}
