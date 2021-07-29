package main

import "fmt"

func main() {
	// 字面量方式创建
	var m = map[string]int{"a": 1, "b": 2}
	fmt.Println(m) // map[a:1 b:2]

	// 使用 make 创建
	m1 := make(map[string]int)
	fmt.Println(m1)

	// 指定长度
	m2 := make(map[string]int, 10)
	fmt.Println(m2)

	// 零值是 nil
	var m3 map[string]int
	fmt.Println(m3 == nil, len(m3) == 0) // true true
	// nil 赋值报错
	// m3["a"] = 1
	// fmt.Println(m3)	// panic: assignment to entry in nil map

	// 赋值
	m["c"] = 3
	m["d"] = 4
	fmt.Println(m) // map[a:1 b:2 c:3 d:4]

	// 取值
	fmt.Println(m["a"], m["d"]) // 1 4
	fmt.Println(m["k"])         // 0

	// 删除
	delete(m, "c")
	delete(m, "f") // key 不存在也不报错
	fmt.Println(m) // map[a:1 b:2 d:4]

	// 获取长度
	fmt.Println(len(m)) // 3

	// 判断键是否存在
	if value, ok := m["d"]; ok {
		fmt.Println(value) // 4
	}

	// 遍历
	for k, v := range m {
		fmt.Println(k, v)
	}

	// 传参
	modify(m)
	fmt.Println("main: ", m) // main:  map[a:1 b:2 d:4 e:10]
}

func modify(a map[string]int) {
	a["e"] = 10
	fmt.Println("modify: ", a) //	modify:  map[a:1 b:2 d:4 e:10]
}
