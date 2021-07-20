package main

import (
	"fmt"
	"os"
)

// 全局变量
var gg = "global"

func main() {

	// 声明单个变量
	var a = "initial"
	fmt.Println(a)

	// 声明多个变量
	var b, c int = 1, 2
	fmt.Println(b, c)

	// 以组的方式声明多个变量
	var (
		b1, c1 int
		b2, c2 = 3, 4
	)
	fmt.Println(b1, c1)
	fmt.Println(b2, c2)

	// 声明布尔值变量
	var d = true
	fmt.Println(d)

	// 没有初始值，会赋默认零值
	var v1 int
	var v2 string
	var v3 bool
	var v4 [10]int // 数组
	var v5 []int   // 数组切片
	var v6 struct {
		e int
	}
	var v7 *int           // 指针
	var v8 map[string]int // map，key 为 string 类型，value 为 int 类型
	var v9 func(e int) int
	fmt.Println(v1, v2, v3, v4, v5, v6, v7, v8, v9)

	// 短变量声明方式
	f := "short"
	fmt.Println(f)

	// 声明赋值多个变量
	g, h := 5, "alwaysbeta"
	fmt.Println(g, h)

	var i int
	// i := 100 // 报错 no new variables on left side of :=
	i, j := 100, 101 // 有新值 j，不报错
	fmt.Println(i, j)

	// 指针
	k := 6
	l := &k         // l 为整型指针，指向 k
	fmt.Println(*l) // 输出 6
	*l = 7
	fmt.Println(k) // 输出 7

	// 使用内置函数 new 声明
	var p = new(int)
	fmt.Println(*p) // 输出整型默认值 0
	*p = 8
	fmt.Println(*p) // 输出 8

	// 变量赋值
	var m, n int
	m = 9
	n = 10
	fmt.Println(m, n)

	// 多重赋值
	m, n = n, m
	fmt.Println(m, n)

	// 空标识符
	r := [5]int{1, 2, 3, 4, 5}
	for _, v := range r {
		// fmt.Println(i, v)
		// fmt.Println(v)	// 定义 i 但不用会报错 i declared but not used
		fmt.Println(v) // 忽略索引
	}

	// 作用域
	fmt.Println(gg) // 输出 global
	gg = "local"
	fmt.Println(gg) // 输出 local

	// 条件分支下的作用域
	// if f, err := os.Open("./00_hello.go"); err != nil {
	// 	fmt.Println(err)
	// }
	// f.Close()	// 报错 f.Close undefined (type string has no field or method Close)

	// 正确写法
	file, err := os.Open("00_hello.go")
	if err != nil {
		fmt.Println(err)
	}
	file.Close()
}
