package main

import (
	"fmt"
	"time"
)

func a() (r int) {
	defer func() {
		r++
	}()
	return 0
}

func b() (r int) {
	t := 5
	defer func() {
		t = t + 5
	}()
	return t
}

func c() (r int) {
	defer func(r int) {
		r = r + 5
	}(r)
	return 1
}

func aa() (r int) {
	r = 0
	// 在 return 之前，执行 defer 函数
	func() {
		r++
	}()
	return
}

func bb() (r int) {
	t := 5
	// 赋值
	r = t
	// 在 return 之前，执行 defer 函数
	// defer 函数没有对返回值 r 进行修改，只是修改了变量 t
	func() {
		t = t + 5
	}()
	return
}

func cc() (r int) {
	// 赋值
	r = 1
	// 这里修改的 r 是函数形参的值
	// 值拷贝，不影响实参值
	func(r int) {
		r = r + 5
	}(r)
	return
}

func d() int {
	r := 0
	defer func() {
		r++
	}()
	return r
}

func e() int {
	r := 0
	defer func(i int) {
		i++
	}(r)
	return 0
}

func G() {
	defer func() {
		// goroutine 外进行 recover
		if err := recover(); err != nil {
			fmt.Println("捕获异常:", err)
		}
		fmt.Println("c")
	}()
	// 创建 goroutine 调用 F 函数
	go F()
	time.Sleep(time.Second)
}

func F() {
	defer func() {
		fmt.Println("b")
	}()
	// goroutine 内部抛出panic
	panic("a")
}

func main() {
	fmt.Println("a = ", a())
	fmt.Println("b = ", b())
	fmt.Println("c = ", c())
	fmt.Println("aa = ", aa())
	fmt.Println("bb = ", bb())
	fmt.Println("cc = ", cc())
	fmt.Println("d = ", d())
	fmt.Println("e = ", e())

	// defer func() {
	// 	fmt.Println("b")
	// }()
	// panic("a")

	// G()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获异常:", err)
		}
	}()
	panic("a")
}
