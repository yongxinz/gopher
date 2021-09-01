![](https://github.com/yongxinz/gopher/blob/main/sc/pic/04-%E5%AD%97%E5%85%B8%E5%92%8C%E7%BB%93%E6%9E%84%E4%BD%93.png)

**原文链接：** [Go 专栏｜复合数据类型：字典 map 和 结构体 struct](https://mp.weixin.qq.com/s/1unl6K9xHxy4V3KukORC3A)

楼下新开了一家重庆砂锅肥肠，扩音喇叭一直在放：正宗的老重庆砂锅肥肠，麻辣可口，老巴适了。

正不正宗不知道，反正听口音，我以为我回东北了。

本篇介绍复合数据类型的最后一篇：字典和结构体。内容很重要，编程时用的也多，需要熟练掌握才行。

> 本文所有代码基于 go1.16.6 编写。

### 字典

字典是一种非常常用的数据结构，Go 中用关键词 map 表示，类型是 `map[K]V`。`K` 和 `V` 分别是字典的键和值的数据类型，其中键必须支持相等运算符，比如数字，字符串等。

#### 创建字典

有两种方式可以创建字典，第一种是直接使用字面量创建；第二种使用内置函数 `make`。

字面量方式创建：

```Go
// 字面量方式创建
var m = map[string]int{"a": 1, "b": 2}
fmt.Println(m) // map[a:1 b:2]
```

使用 `make` 创建：

```Go
// 使用 make 创建
m1 := make(map[string]int)
fmt.Println(m1)
```

还可以初始化字典的长度。在已知字典长度的情况下，直接指定长度可以提升程序的执行效率。

```Go
// 指定长度
m2 := make(map[string]int, 10)
fmt.Println(m2)
```

字典的零值是 nil，对值是 nil 的字典赋值会报错。

```Go
// 零值是 nil
var m3 map[string]int
fmt.Println(m3 == nil, len(m3) == 0) // true true
// nil 赋值报错
// m3["a"] = 1
// fmt.Println(m3)	// panic: assignment to entry in nil map
```

#### 使用字典

**赋值：**

```Go
// 赋值
m["c"] = 3
m["d"] = 4
fmt.Println(m) // map[a:1 b:2 c:3 d:4]
```

**取值：**

```Go
// 取值
fmt.Println(m["a"], m["d"]) // 1 4
fmt.Println(m["k"])         // 0
```

即使在 Key 不存在的情况下，也是不报错的。而是返回对应类型的零值。

**删除元素：**

```Go
// 删除
delete(m, "c")
delete(m, "f") // key 不存在也不报错
fmt.Println(m) // map[a:1 b:2 d:4]
```

**获取长度：**

```Go
// 获取长度
fmt.Println(len(m)) // 3
```

**判断键是否存在：**

```Go
// 判断键是否存在
if value, ok := m["d"]; ok {
	fmt.Println(value) // 4
}
```

和 Python 对比起来看，这个用起来就很爽。

**遍历：**

```Go
// 遍历
for k, v := range m {
	fmt.Println(k, v)
}
```

#### 引用类型

map 是引用类型，所以在函数间传递时，也不会制造一个映射的副本，这点和切片类似，都很高效。

```Go
package main

import "fmt"

func main() {
	...

	// 传参
	modify(m)
	fmt.Println("main: ", m) // main:  map[a:1 b:2 d:4 e:10]
}

func modify(a map[string]int) {
	a["e"] = 10
	fmt.Println("modify: ", a) //	modify:  map[a:1 b:2 d:4 e:10]
}
```

### 结构体

结构体是一种聚合类型，包含零个或多个任意类型的命名变量，每个变量叫做结构体的成员。

#### 创建结构体

首先使用 `type` 来自定义一个结构体类型 `user`，里面有两个成员变量，分别是：`name` 和 `age`。

```Go
// 声明结构体
type user struct {
	name string
	age  int
}
```

结构体的初始化有两种方式：

第一种是按照声明字段的顺序逐个赋值，这里需要注意，字段的顺序要严格一致。

```Go
// 初始化
u1 := user{"zhangsan", 18}
fmt.Println(u1) // {zhangsan 18}
```

这样做的缺点很明显，如果字段顺便变了，那么凡是涉及到这个结构初始化的部分都要跟着变。

所以，更推荐使用第二种方式，按照字段名字来初始化。

```Go
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
```

未初始化的字段会赋值相应类型的零值。

#### 使用结构体

使用点号 `.` 来访问和赋值成员变量。

```Go
// 访问结构体成员
fmt.Println(u.name, u.age) // zhangsan 18
u.name = "lisi"
fmt.Println(u.name, u.age) // lisi 18
```

如果结构体的成员变量是可比较的，那么结构体也是可比较的。

```Go
// 结构体比较
u2 := user{
	age:  18,
	name: "zhangsan",
}
fmt.Println(u1 == u)  // false
fmt.Println(u1 == u2) // true
```

#### 结构体嵌套

现在我们已经定义一个 `user` 结构体了，假设我们再定义两个结构体 `admin` 和 `leader`，如下：

```Go
type admin struct {
	name    string
	age     int
	isAdmin bool
}

type leader struct {
	name     string
	age      int
	isLeader bool
}
```

那么问题就来了，有两个字段 `name` 和 `age` 被重复定义了多次。

懒是程序员的必修课。有没有什么办法可以复用这两个字段呢？答案就是结构体嵌套。

使用嵌套方式优化后变成了这样：

```Go
type admin struct {
	u       user
	isAdmin bool
}

type leader struct {
	u        user
	isLeader bool
}
```

代码看起来简洁了很多。

#### 匿名成员

但这样依然不是很完美，每次访问嵌套结构体的成员变量时还是有点麻烦。

```Go
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
```

这个时候就需要匿名成员登场了，不指定名称，只指定类型。

```Go
type admin1 struct {
	user
	isAdmin bool
}
```

通过这种方式可以省略掉中间变量，直接访问我们需要的成员变量。

```Go
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
```

### 总结

本文介绍了字典和结构体，两种很常用的数据类型。虽然篇幅不长，但基本操作也都包括，写代码肯定是没有问题的。更底层的原理和更灵活的用法就需要大家自己去探索和发现了。

当然，我也会在写完基础专栏之后，分享一些更深层的文章，欢迎大家关注，交流。

到目前为止，数据类型就都介绍完了。

先是学习了基础数据类型，包括整型，浮点型，复数类型，布尔型和字符串型。然后是复合数据类型，包括数组，切片，字典和结构体。

这些都是 Go 的基础，一定要多多练习，熟练掌握。文中的代码我都已经上传到 Github 了，有需要的同学可以点击文末地址，自行下载。


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
