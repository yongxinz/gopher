package main

import (
	"fmt"
)

func main() {
	// 从 0 值开始，逐项加 1
	const (
		a0 = iota // 0
		a1 = iota // 1
		a2 = iota // 2
	)
	fmt.Println(a0, a1, a2)

	// 简写，表达式相同，可以省略后面的
	const (
		b0 = iota // 0
		b1        // 1
		b2        // 2
	)
	fmt.Println(b0, b1, b2)

	const (
		b         = iota      // 0
		c float32 = iota * 10 // 10
		d         = iota      // 2
	)
	fmt.Println(b, c, d)

	// iota 在每个 const 开头被重置为 0
	const x = iota // 0
	fmt.Println(x)

	// 同上
	const y = iota // 0
	fmt.Println(y)

	// 枚举
	const (
		Sunday    = iota // 0
		Monday           // 1
		Tuesday          // 2
		Wednesday        // 3
		Thursday         // 4
		Friday           // 5
		Saturday         // 6
	)
	fmt.Println(Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday)
}
