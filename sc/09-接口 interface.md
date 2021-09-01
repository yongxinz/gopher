![](https://github.com/yongxinz/gopher/blob/main/sc/pic/09_%E6%8E%A5%E5%8F%A3interface.png)

**原文链接：** [Go 专栏｜接口 interface](https://mp.weixin.qq.com/s/g7ngRIxxbd-M8K_sL_M4KQ)

Duck Typing，鸭子类型，在维基百科里是这样定义的：

> If it looks like a duck, swims like a duck, and quacks like a duck, then it probably is a duck.

翻译过来就是：如果某个东西长得像鸭子，游泳像鸭子，嘎嘎叫像鸭子，那它就可以被看成是一只鸭子。

它是动态编程语言的一种对象推断策略，它更关注对象能做什么，而不是对象的类型本身。

例如：在动态语言 Python 中，定义一个这样的函数：

```python
def hello_world(duck):
    duck.say_hello()
```

当调用此函数的时候，可以传入任意类型，只要它实现了 `say_hello()` 就可以。如果没实现，运行过程中会出现错误。

Go 语言作为一门静态语言，它通过接口的方式完美支持鸭子类型。

### 接口类型

之前介绍的类型都是具体类型，而接口是一种抽象类型，是多个方法声明的集合。在 Go 中，只要目标类型实现了接口要求的所有方法，我们就说它实现了这个接口。

先来看一个例子：

```go
package main

import "fmt"

// 定义接口，包含 Eat 方法
type Duck interface {
	Eat()
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
}
```

使用 `type` 关键词定义接口：

```go
type Duck interface {
	Eat()
}
```

接口包含了一个 `Eat()` 方法，然后定义两个结构体类型 `Cat` 和 `Dog`，分别实现了 `Eat` 方法。

```go
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
```

遍历接口切片，通过接口类型可以直接调用对应方法：

```go
s := []Duck{
	&Cat{},
	&Dog{},
}
for _, n := range s {
	n.Eat()
}

// 输出
// cat eat
// dog eat
```

### 接口赋值

接口赋值分两种情况：

1. 将对象实例赋值给接口
2. 将一个接口赋值给另一个接口

下面来分别说说：

#### 将对象实例赋值给接口

还是用上面的例子，因为 `Cat` 实现了 `Eat` 接口，所以可以直接将 `Cat` 实例赋值给接口。

```go
var c Duck = &Cat{}
c.Eat()
```

在这里一定要传结构体指针，如果直接传结构体会报错：

```go
var c Duck = Cat{}
c.Eat()
```

```
# command-line-arguments
./09_interface.go:25:6: cannot use Cat{} (type Cat) as type Duck in assignment:
	Cat does not implement Duck (Eat method has pointer receiver)
```

但是如果反过来呢？比如使用结构体来实现接口，使用结构体指针来赋值：

```go
// 定义 Cat 结构体，并实现 Eat 方法
type Cat struct{}

func (c Cat) Eat() {
	fmt.Println("cat eat")
}

var c Duck = &Cat{}
c.Eat() // cat eat
```

没有问题，可以正常执行。 

#### 将一个接口赋值给另一个接口

还是上面的例子，可以直接将 `c` 的值直接赋值给 `d`：

```go
var c Duck = &Cat{}
c.Eat()

var d Duck = c
d.Eat()
```

再来，我再定义一个接口 `Duck1`，这个接口包含两个方法 `Eat` 和 `Walk`，然后结构体 `Dog` 实现两个方法，但是 `Cat` 只实现 `Eat` 方法。

```go
type Duck1 interface {
	Eat()
	Walk()
}

// 定义 Dog 结构体，并实现 Eat 方法
type Dog struct{}

func (d *Dog) Eat() {
	fmt.Println("dog eat")
}

func (d *Dog) Walk() {
	fmt.Println("dog walk")
}
```

那么在赋值时，使用 `Duck1` 赋值给 `Duck` 是可以的，反过来就会报错。

```go
var c1 Duck1 = &Dog{}
var c2 Duck = c1
c2.Eat()
```

所以，已经初始化的接口变量 `c1` 直接赋值给另一个接口变量 `c2`，要求 `c2` 的方法集是 `c1` 的方法集的子集。

### 空接口

具有 0 个方法的接口称为空接口，它表示为 `interface {}`。由于空接口有 0 个方法，所以所有类型都实现了空接口。

```go
func main() {
	// interface 形参
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

func test(i interface{}) {
	fmt.Printf("Type = %T, value = %v\n", i, i)
}
```

### 类型断言

类型断言是作用在接口值上的操作，语法如下：

```go
x.(T)
```

其中 `x` 是接口类型的表达式，`T` 是断言类型。

作用是判断操作数的动态类型是否满足指定的断言类型。

有两种情况：

1. `T` 是具体类型
2. `T` 是接口类型

下面来分别举例说明：

#### 具体类型

类型断言会检查 `x` 的动态类型是否为 `T`，如果是，则输出 `x` 的值；如果不是，程序直接 `panic`。

```go
func main() {
	// 类型断言
	var n interface{} = 55
	assert(n) // 55
	var n1 interface{} = "hello"
	assert(n1) // panic: interface conversion: interface {} is string, not int
}

func assert(i interface{}) {
	s := i.(int)
	fmt.Println(s)
}
```

#### 接口类型

类型断言会检查 `x` 的动态类型是否满足接口类型 `T`，如果满足，则输出 `x` 的值，这个值可能是绑定实例的副本，也可能是指针的副本；如果不满足，程序直接 `panic`。

```go
func main() {
	// 类型断言
	assertInterface(c) // &{}
}

func assertInterface(i interface{}) {
	s := i.(Duck)
	fmt.Println(s)
}
```

如果有两个接收值，那么断言不会在失败时崩溃，而是会多返回一个布尔值，一般命名为 `ok`，来表示断言是否成功。

```go
func main() {
	// 类型断言
	var n1 interface{} = "hello"
	assertFlag(n1)
}

func assertFlag(i interface{}) {
	if s, ok := i.(int); ok {
		fmt.Println(s)
	}
}
```

### 类型查询

语法类似类型断言，只需将 `T` 直接用关键词 `type` 替代。

作用主要有两个：

1. 查询一个接口变量绑定的底层变量类型
2. 查询一个接口变量的底层变量是否还实现了其他接口

```go
func main() {
	// 类型查询
	SearchType(50)         // Int: 50
	SearchType("zhangsan") // String: zhangsan
	SearchType(c)          // dog eat
	SearchType(50.1)       // Unknown type
}

func SearchType(i interface{}) {
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
```

### 总结

本文从鸭子类型引出 Go 的接口，然后用一个例子简单展示了接口类型的用法，接着又介绍了接口赋值，空接口，类型断言和类型查询。

相信通过本篇文章大家能对接口有了整体的概念，并掌握了基本用法。

---

文章中的脑图和源码都上传到了 GitHub，有需要的同学可自行下载。

**地址：** https://github.com/yongxinz/gopher/tree/main/sc

关注公众号 **AlwaysBeta**，回复「**goebook**」领取 Go 编程经典书籍。

<center class="half">
    <img src="https://github.com/yongxinz/gopher/blob/main/alwaysbeta.JPG" width="300"/>
</center>

**Go 专栏文章列表：**

1. [开发环境搭建以及开发工具 VS Code 配置](<https://github.com/yongxinz/gopher/blob/main/sc/00-%E5%BC%80%E5%8F%91%E7%8E%AF%E5%A2%83%E6%90%AD%E5%BB%BA%E4%BB%A5%E5%8F%8A%E5%BC%80%E5%8F%91%E5%B7%A5%E5%85%B7%20VS%20Code%20%E9%85%8D%E7%BD%AE.md>)

2. [变量和常量的声明与赋值](<https://github.com/yongxinz/gopher/blob/main/sc/01-%E5%8F%98%E9%87%8F%E5%92%8C%E5%B8%B8%E9%87%8F%E7%9A%84%E5%A3%B0%E6%98%8E%E4%B8%8E%E8%B5%8B%E5%80%BC.md>)

3. [基础数据类型：整数、浮点数、复数、布尔值和字符串](<https://github.com/yongxinz/gopher/blob/main/sc/02-%E5%9F%BA%E7%A1%80%E6%95%B0%E6%8D%AE%E7%B1%BB%E5%9E%8B%EF%BC%9A%E6%95%B4%E6%95%B0%E3%80%81%E6%B5%AE%E7%82%B9%E6%95%B0%E3%80%81%E5%A4%8D%E6%95%B0%E3%80%81%E5%B8%83%E5%B0%94%E5%80%BC%E5%92%8C%E5%AD%97%E7%AC%A6%E4%B8%B2.md>)

4. [复合数据类型：数组和切片 slice](<https://github.com/yongxinz/gopher/blob/main/sc/03-%E5%A4%8D%E5%90%88%E6%95%B0%E6%8D%AE%E7%B1%BB%E5%9E%8B%EF%BC%9A%E6%95%B0%E7%BB%84%E5%92%8C%E5%88%87%E7%89%87%20slice.md>)

5. [复合数据类型：字典 map 和 结构体 struct](<https://github.com/yongxinz/gopher/blob/main/sc/04-%E5%A4%8D%E5%90%88%E6%95%B0%E6%8D%AE%E7%B1%BB%E5%9E%8B%EF%BC%9A%E5%AD%97%E5%85%B8%20map%20%E5%92%8C%20%E7%BB%93%E6%9E%84%E4%BD%93%20struct.md>)
6. [流程控制，一网打尽](<https://github.com/yongxinz/gopher/blob/main/sc/05-%E6%B5%81%E7%A8%8B%E6%8E%A7%E5%88%B6%EF%BC%8C%E4%B8%80%E7%BD%91%E6%89%93%E5%B0%BD.md>)
7. [函数那些事](<https://github.com/yongxinz/gopher/blob/main/sc/06-%E5%87%BD%E6%95%B0%E9%82%A3%E4%BA%9B%E4%BA%8B.md>)
8. [错误处理：defer，panic 和 recover](<https://github.com/yongxinz/gopher/blob/main/sc/07-%E9%94%99%E8%AF%AF%E5%A4%84%E7%90%86%EF%BC%9Adefer%EF%BC%8Cpanic%20%E5%92%8C%20recover.md>)
9. [说说方法](<https://github.com/yongxinz/gopher/blob/main/sc/08-%E8%AF%B4%E8%AF%B4%E6%96%B9%E6%B3%95.md>)