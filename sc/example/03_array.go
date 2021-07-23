package main

import "fmt"

func main() {
	var a [3]int
	// 输出数组第一个元素
	fmt.Println(a[0]) // 0
	// 输出数组长度
	fmt.Println(len(a)) // 3

	// 数组字面量初始化
	var b [3]int = [3]int{1, 2, 3}
	var c [3]int = [3]int{1, 2}
	fmt.Println(b)    // [1 2 3]
	fmt.Println(c[2]) // 0

	// 使用 ...
	d := [...]int{1, 2, 3, 4, 5}
	fmt.Printf("%T\n", d) // [5]int

	// 指定索引位置初始化
	e := [4]int{5, 2: 10}
	f := [...]int{2, 4: 6}
	fmt.Println(e) // [5 0 10 0]
	fmt.Println(f) // [2 0 0 0 6]

	// 二维数组
	var g [4][2]int
	h := [4][2]int{{10, 11}, {20, 21}, {30, 31}, {40, 41}}
	// 声明并初始化外层数组中索引为 1 和 3 的元素
	i := [4][2]int{1: {20, 21}, 3: {40, 41}}
	// 声明并初始化外层数组和内层数组的单个元素
	j := [...][2]int{1: {0: 20}, 3: {1: 41}}
	fmt.Println(g, h, i, j)

	// 数组比较
	a1 := [2]int{1, 2}
	a2 := [...]int{1, 2}
	a3 := [2]int{1, 3}
	// a4 := [3]int{1, 2}
	fmt.Println(a1 == a2, a1 == a3, a2 == a3) // true false false
	// fmt.Println(a1 == a4)                     // invalid operation: a1 == a4 (mismatched types [2]int and [3]int)

	// 数组遍历
	for i, n := range e {
		fmt.Println(i, n)
	}

	// 数组复制
	x := [2]int{10, 20}
	y := x
	fmt.Printf("x: %p, %v\n", &x, x) // x: 0xc00012e020, [10 20]
	fmt.Printf("y: %p, %v\n", &y, y) // y: 0xc00012e030, [10 20]
	test(x)

	// 传参
	modify(x)
	fmt.Println("main: ", x) // main:  [10 20]
}

func test(a [2]int) {
	fmt.Printf("a: %p, %v\n", &a, a) // a: 0xc00012e060, [10 20]
}

func modify(a [2]int) {
	a[0] = 30
	fmt.Println("modify: ", a) // modify:  [30 20]
}
