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

	// 字符串
	s1 := "hello"
	s2 := "world"

	// 原始字符串
	s := `row1\r\n
	row2`
	fmt.Println(s)

	// 字符串拼接
	s3 := s1 + s2
	fmt.Println(s3)
	// 取字符串长度
	fmt.Println(len(s3))
	// 取单个字符
	fmt.Println(s3[4])
	// 字符串切片
	fmt.Println(s3[2:4])
	fmt.Println(s3[:4])
	fmt.Println(s3[2:])
	fmt.Println(s3[:])

	// 修改报错
	// s3[0] = "H"	// cannot assign to s3[0] (strings are immutable)

	s4 := "hello 世界"

	// 遍历字节数组
	for i := 0; i < len(s4); i++ {
		fmt.Println(i, s4[i])
	}

	// 遍历 rune 数组
	for i, v := range s4 {
		fmt.Println(i, v)
	}
}
