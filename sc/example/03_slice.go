package main

import "fmt"

func main() {
	// 基于数组创建切片
	var array = [...]int{1, 2, 3, 4, 5, 6, 7, 8}

	s1 := array[3:6]
	s2 := array[:5]
	s3 := array[4:]
	s4 := array[:]

	fmt.Printf("s1: %v\n", s1) // s1: [4 5 6]
	fmt.Printf("s2: %v\n", s2) // s2: [1 2 3 4 5]
	fmt.Printf("s3: %v\n", s3) // s3: [5 6 7 8]
	fmt.Printf("s4: %v\n", s4) // s4: [1 2 3 4 5 6 7 8]

	// 使用 make 创建切片
	// len: 10, cap: 10
	a := make([]int, 10)
	// len: 10, cap: 15
	b := make([]int, 10, 15)

	fmt.Printf("a: %v, len: %d, cap: %d\n", a, len(a), cap(a))
	fmt.Printf("b: %v, len: %d, cap: %d\n", b, len(b), cap(b))

	// 切片遍历
	for i, n := range s1 {
		fmt.Println(i, n)
	}

	// 比较
	var s []int
	fmt.Println(len(s) == 0, s == nil) // true true
	s = nil
	fmt.Println(len(s) == 0, s == nil) // true true
	s = []int(nil)
	fmt.Println(len(s) == 0, s == nil) // true true
	s = []int{}
	fmt.Println(len(s) == 0, s == nil) // true false

	// 追加
	s5 := append(s4, 9)
	fmt.Printf("s5: %v\n", s5) // s5: [1 2 3 4 5 6 7 8 9]
	s6 := append(s4, 10, 11)
	fmt.Printf("s6: %v\n", s6) // s5: [1 2 3 4 5 6 7 8 10 11]

	// 追加另一个切片
	s7 := []int{12, 13}
	s7 = append(s7, s6...)
	fmt.Printf("s7: %v\n", s7) // s7: [12 13 1 2 3 4 5 6 7 8 10 11]

	// 复制
	s8 := []int{1, 2, 3, 4, 5}
	s9 := []int{5, 4, 3}
	s10 := []int{6}

	copy(s8, s9)
	fmt.Printf("s8: %v\n", s8) // s8: [5 4 3 4 5]
	copy(s10, s9)
	fmt.Printf("s10: %v\n", s10) // s10: [5]

	// 传参
	modify(s9)
	fmt.Println("main: ", s9) // main:  [30 4 3]
}

func modify(a []int) {
	a[0] = 30
	fmt.Println("modify: ", a) // modify:  [30 4 3]
}
