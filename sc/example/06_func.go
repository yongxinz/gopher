package main

import (
	"fmt"
)

func main() {
	funcA() // i am funcA

	// 函数签名
	fmt.Printf("%T\n", add) // func(int, int) int
	fmt.Printf("%T\n", sub) // func(int, int) int

	// 不定参数
	fmt.Println(funcSum(1, 2))    // 3
	fmt.Println(funcSum(1, 2, 3)) // 6

	// slice 参数
	s := []int{1, 2, 3, 4}
	fmt.Println(funcSum(s...)) // 10
	fmt.Println(funcSum1(s))   // 10

	// 多返回值
	fmt.Println(swap(1, 2)) // 2 1

	x, _ := swap(1, 2)
	fmt.Println(x) // 2

	// 匿名函数
	sum := func(a, b int) int { return a + b }
	fmt.Println(sum(1, 2)) // 3

	// 作为参数
	fmt.Println(funcSum2(sum, 3, 5)) // 8

	// 作为返回值
	f := wrap("add")
	fmt.Println(f(2, 4)) // 6

	// 直接调用
	fmt.Println(func(a, b int) int { return a + b }(4, 5)) // 9
}

func funcA() {
	fmt.Println("i am funcA")
}

func add(x int, y int) int {
	return x + y
}

func sub(x int, y int) (z int) {
	z = x - y
	return
}

// 简写形式
func add1(x, y int) int {
	return x + y
}

func sub1(x, y int) (z int) {
	z = x - y
	return
}

// 不定参数
func funcSum(args ...int) (ret int) {
	for _, arg := range args {
		ret += arg
	}
	return
}

// slice 参数
func funcSum1(args []int) (ret int) {
	for _, arg := range args {
		ret += arg
	}
	return
}

// 多返回值
func swap(x, y int) (int, int) {
	return y, x
}

// 匿名函数作为参数
func funcSum2(f func(int, int) int, x, y int) int {
	return f(x, y)
}

// 匿名函数作为返回值
func wrap(op string) func(int, int) int {
	switch op {
	case "add":
		return func(a, b int) int {
			return a + b
		}
	case "sub":
		return func(a, b int) int {
			return a + b
		}

	default:
		return nil
	}
}
