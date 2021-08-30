package main

import "fmt"

// 定义接口，包含 Eat 方法
type Duck interface {
	Eat()
}

type Duck1 interface {
	Eat()
	Walk()
}

// 定义 Cat 结构体，并实现 Eat 方法
type Cat struct{}

func (c *Cat) Eat() {
	fmt.Println("cat eat")
}

// 定义 Dog 结构体，并实现 Eat 方法
type Dog struct{}

func (d *Dog) Eat() {
	fmt.Println("dog eat")
}

func (d *Dog) Walk() {
	fmt.Println("dog walk")
}

func main() {
	var c Duck = &Cat{}
	c.Eat()

	var d Duck = &Dog{}
	d.Eat()

	s := []Duck{
		&Cat{},
		&Dog{},
	}
	for _, n := range s {
		n.Eat()
	}

	var c1 Duck1 = &Dog{}
	var c2 Duck = c1
	c2.Eat()

	// 类型断言
	var n interface{} = 55
	assert(n) // 55
	var n1 interface{} = "hello"
	// assert(n1) // panic: interface conversion: interface {} is string, not int
	assertFlag(n1)

	assertInterface(c) // &{}

	// 类型查询
	searchType(50)         // Int: 50
	searchType("zhangsan") // String: zhangsan
	searchType(c1)         // dog eat
	searchType(50.1)       // Unknown type

	// 空接口
	s1 := "Hello World"
	i := 50
	strt := struct {
		name string
	}{
		name: "AlwaysBeta",
	}
	test(s1)
	test(i)
	test(strt)
}

func assert(i interface{}) {
	s := i.(int)
	fmt.Println(s)
}

func assertInterface(i interface{}) {
	s := i.(Duck)
	fmt.Println(s)
}

func assertFlag(i interface{}) {
	if s, ok := i.(int); ok {
		fmt.Println(s)
	}
}

func searchType(i interface{}) {
	switch v := i.(type) {
	case string:
		fmt.Printf("String: %s\n", i.(string))
	case int:
		fmt.Printf("Int: %d\n", i.(int))
	case Duck:
		v.Eat()
	default:
		fmt.Printf("Unknown type\n")
	}
}

func test(i interface{}) {
	fmt.Printf("Type = %T, value = %v\n", i, i)
}
