**原文链接：** [Go 语言 new 和 make 关键字的区别](https://mp.weixin.qq.com/s/NBDkI3roHgNgW1iW4e_6cA)

本篇文章来介绍一道非常常见的面试题，到底有多常见呢？可能很多面试的开场白就是由此开始的。那就是 new 和 make 这两个内置函数的区别。

其实这个问题本身并不复杂，简单来说就是，new 只分配内存，而 make 只能用于 slice、map 和 chan 的初始化，下面我们就来详细介绍一下。

## new

new 是一个内置函数，它会分配一段内存，并返回指向该内存的指针。

其函数签名如下：

### 源码

```go
// The new built-in function allocates memory. The first argument is a type,
// not a value, and the value returned is a pointer to a newly
// allocated zero value of that type.
func new(Type) *Type
```

从上面的代码可以看出，new 函数只接受一个参数，这个参数是一个类型，并且返回一个指向该类型内存地址的指针。

同时 new 函数会把分配的内存置为零，也就是类型的零值。

### 使用

使用 new 函数为变量分配内存空间：

```go
p1 := new(int)
fmt.Printf("p1 --> %#v \n ", p1) //(*int)(0xc42000e250) 
fmt.Printf("p1 point to --> %#v \n ", *p1) //0

var p2 *int
i := 0
p2 = &i
fmt.Printf("p2 --> %#v \n ", p2) //(*int)(0xc42000e278) 
fmt.Printf("p2 point to --> %#v \n ", *p2) //0
```

上面的代码是等价的，`new(int)` 将分配的空间初始化为 int 的零值，也就是 0，并返回 int 的指针，这和直接声明指针并初始化的效果是相同的。

当然，new 函数不仅能够为系统默认的数据类型分配空间，自定义类型也可以使用 new 函数来分配空间，如下所示：

```go
type Student struct {
   name string
   age int
}
var s *Student
s = new(Student) //分配空间
s.name = "zhangsan"
fmt.Println(s)
```

这就是 new 函数，它返回的永远是类型的指针，指针指向分配类型的内存地址。需要注意的是，new 函数只会分配内存空间，但并不会初始化该内存空间。

## make

make 也是用于内存分配的，但是和 new 不同，它只用于 slice、map 和 chan 的内存创建，而且它返回的类型就是这三个类型本身，而不是他们的指针类型。因为这三种类型本身就是引用类型，所以就没有必要返回他们的指针了。

其函数签名如下：

### 源码

```go
// The make built-in function allocates and initializes an object of type
// slice, map, or chan (only). Like new, the first argument is a type, not a
// value. Unlike new, make's return type is the same as the type of its
// argument, not a pointer to it. The specification of the result depends on
// the type:
// Slice: The size specifies the length. The capacity of the slice is
// equal to its length. A second integer argument may be provided to
// specify a different capacity; it must be no smaller than the
// length, so make([]int, 0, 10) allocates a slice of length 0 and
// capacity 10.
// Map: An empty map is allocated with enough space to hold the
// specified number of elements. The size may be omitted, in which case
// a small starting size is allocated.
// Channel: The channel's buffer is initialized with the specified
// buffer capacity. If zero, or the size is omitted, the channel is
// unbuffered.
func make(t Type, size ...IntegerType) Type
```

通过上面的代码可以看出 make 函数的 `t` 参数必须是 slice、map 和 chan 中的一个，并且返回值也是类型本身。

### 使用

下面用 slice 来举一个例子：

```go
var s1 []int
if s1 == nil {
    fmt.Printf("s1 is nil --> %#v \n ", s1) // []int(nil)
}

s2 := make([]int, 3)
if s2 == nil {
    fmt.Printf("s2 is nil --> %#v \n ", s2)
} else {
    fmt.Printf("s2 is not nill --> %#v \n ", s2)// []int{0, 0, 0}
}
```

slice 的零值是 `nil`，但使用 make 初始化之后，slice 内容被类型 int 的零值填充，如：`[]int{0, 0, 0}`。

map 和 chan 也是类似的，就不多说了。

## 总结

通过以上分析，总结一下 new 和 make 主要区别如下：

1. make 只能用来分配及初始化类型为 slice、map 和 chan 的数据。new 可以分配任意类型的数据；
2. new 分配返回的是指针，即类型 `*Type`。make 返回类型本身，即 `Type`；
3. new 分配的空间被清零。make 分配空间后，会进行初始化；


以上就是本文的全部内容，如果觉得还不错的话欢迎**点赞**，**转发**和**关注**，感谢支持。

***

**参考文章：**

- https://go.dev/doc/effective_go#allocation_new
- http://c.biancheng.net/view/5722.html
- https://sanyuesha.com/2017/07/26/go-make-and-new/

**推荐阅读：**

- [为什么 Go 不支持 []T 转换为 []interface](https://mp.weixin.qq.com/s/cwDEgnicK4jkuNpzulU2bw)
- [为什么 Go 语言 struct 要使用 tags](https://mp.weixin.qq.com/s/L7-TJ-CzYfuVrIBWP7Ebaw)