package main

import "fmt"

func main() {
	var a int = 10
	var b int32 = 20

	// fmt.Println(a + b)	// 报错 invalid operation: a + b (mismatched types int and int32)
	// 需要强制类型转换
	fmt.Println(a + int(b)) // 输出 30

	// 浮点型转整型
	var c float32 = 10.23
	fmt.Println(int(c)) // 输出 10

	// 取模
	fmt.Println(5 % 3)   // 输出 2
	fmt.Println(-5 % -3) // 输出 -2

	// 除法
	fmt.Println(5 / 3)     // 输出 1
	fmt.Println(5.0 / 3.0) // 输出 1.6666666666666667

	// 比较运算
	var i int32
	var j int64
	i, j = 1, 2

	// if i == j { // 报错 invalid operation: i == j (mismatched types int32 and int64)
	// 	fmt.Println("i and j are equal.")
	// }
	if i == 1 || j == 2 {
		fmt.Println("equal.")
	}

	// 复数
	var x complex64 = 3 + 5i
	var y complex128 = complex(3.5, 10)
	// 分别打印实部和虚部
	fmt.Println(real(x), imag(x)) // 输出 3 5
	fmt.Println(real(y), imag(y)) // 输出 3.5 10

	// 布尔
	ok := true
	fmt.Println(ok)

	// 类型转换
	// var e bool
	// e = bool(1)	// 报错  cannot convert 1 (type untyped int) to type bool

	m := 1
	// if m { // 报错 non-bool m (type int) used as if condition
	// 	fmt.Println("is true")
	// }
	if m == 1 {
		fmt.Println("m is 1")
	}
}
