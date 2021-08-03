package main

import (
	"fmt"
)

func main() {
	i := 1
	// 只有条件
	for i <= 3 {
		fmt.Println(i)
		i = i + 1
	}

	// 有变量初始化和条件
	for j := 7; j <= 9; j++ {
		fmt.Println(j)
	}

	// 死循环
	for {
		fmt.Println("loop")
		break
	}

	// 遍历数组
	a := [...]int{10, 20, 30, 40}
	for i := range a {
		fmt.Println(i)
	}
	for i, v := range a {
		fmt.Println(i, v)
	}

	// 遍历切片
	s := []string{"a", "b", "c"}
	for i := range s {
		fmt.Println(i)
	}
	for i, v := range s {
		fmt.Println(i, v)
	}

	// 遍历字典
	m := map[string]int{"a": 10, "b": 20, "c": 30}
	for k := range m {
		fmt.Println(k)
	}
	for k, v := range m {
		fmt.Println(k, v)
	}
}
