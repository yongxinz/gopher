package main

import "fmt"

// 声明结构体
type user struct {
	name string
	age  int
}

type admin struct {
	u       user
	isAdmin bool
}

type leader struct {
	u        user
	isLeader bool
}

type admin1 struct {
	user
	isAdmin bool
}

func main() {
	// 初始化
	u1 := user{"zhangsan", 18}
	fmt.Println(u1) // {zhangsan 18}

	// 更好的方式
	// u := user{
	// 	age: 20,
	// }
	// fmt.Println(u)	// { 20}
	u := user{
		name: "zhangsan",
		age:  18,
	}
	fmt.Println(u) // {zhangsan 18}

	// 访问结构体成员
	fmt.Println(u.name, u.age) // zhangsan 18
	u.name = "lisi"
	fmt.Println(u.name, u.age) // lisi 18

	// 结构体比较
	u2 := user{
		age:  18,
		name: "zhangsan",
	}
	fmt.Println(u1 == u)  // false
	fmt.Println(u1 == u2) // true

	// 结构体嵌套
	a := admin{
		u:       u,
		isAdmin: true,
	}
	fmt.Println(a) // {{lisi 18} true}
	a.u.name = "wangwu"
	fmt.Println(a.u.name)  // wangwu
	fmt.Println(a.u.age)   // 18
	fmt.Println(a.isAdmin) // true

	l := leader{
		u:        u,
		isLeader: false,
	}
	fmt.Println(l) // {{lisi 18} false}

	// 匿名成员
	a1 := admin1{
		user:    u,
		isAdmin: true,
	}
	a1.age = 20
	a1.isAdmin = false

	fmt.Println(a1)         // {{lisi 20} false}
	fmt.Println(a1.name)    // lisi
	fmt.Println(a1.age)     // 20
	fmt.Println(a1.isAdmin) // false
}
