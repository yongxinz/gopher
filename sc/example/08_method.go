package main

import (
	"fmt"
)

type Person struct {
	name string
}

type Point struct {
	x, y int
}

func main() {
	p := Person{name: "zhangsan"}

	// 调用方法
	fmt.Println(p.String()) // person name is zhangsan

	// 值接收者
	p.Modify()
	fmt.Println(p.String()) // person name is zhangsan

	// 指针接收者
	p.ModifyP()
	fmt.Println(p.String()) // person name is lisi
	// 等价于
	(&p).ModifyP()
	fmt.Println(p.String())

	(&p).Modify()
	fmt.Println(p.String())

	// 方法变量
	p1 := Point{1, 2}
	q1 := Point{3, 4}
	f := p1.Add
	fmt.Println(f(q1)) // {4 6}

	// 方法表达式
	f1 := Point.Add
	fmt.Println(f1(p1, q1)) // {4 6}
}

func (p Person) String() string {
	return "person name is " + p.name
}

// 值接收者
func (p Person) Modify() {
	p.name = "lisi"
}

// 指针接收者
func (p *Person) ModifyP() {
	p.name = "lisi"
}

func (p Point) Add(q Point) Point {
	return Point{p.x + q.x, p.y + q.y}
}
