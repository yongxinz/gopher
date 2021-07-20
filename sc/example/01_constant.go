package main

import (
	"fmt"
)

// 全局变量
const s string = "constant"

func main() {
	fmt.Println(s)

	// 无类型整型常量
	const n = 500000000

	// 用编译阶段即可计算出值的表达式来赋值
	const d = 3e20 / n
	fmt.Println(d)

	// 类型转换
	fmt.Println(int64(d))

	const Pi float64 = 3.14159265358979323846
	// 无类型浮点常量
	const zero = 0.0

	// 无类型整型和字符串常量
	const a, b, c = 3, 4, "foo"
	fmt.Println(a, b, c)

	// 多个常量
	const (
		size int64 = 1024
		eof        = -1 // 无类型整型常量
	)
	fmt.Println(size, eof)
}
