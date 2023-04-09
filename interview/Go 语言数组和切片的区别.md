**原文链接：** [Go 语言数组和切片的区别](https://mp.weixin.qq.com/s/esaAmAdmV4w3_qjtAzTr4A)

在 Go 语言中，数组和切片看起来很像，但其实它们又有很多的不同之处，这篇文章就来说说它们到底有哪些不同。

另外，这个问题在面试中也经常会被问到，属于入门级题目，看过文章之后，相信你会有一个很好的答案。

## 数组

数组是同一种数据类型元素的集合，数组在定义时需要指定长度和元素类型。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/array.png)

例如：`[4]int` 表示一个包含四个整数的数组，数组的大小是固定的。并且长度是其类型的一部分（`[4]int` 和 `[5]int` 是不同的、不兼容的类型）。

数组元素可以通过索引来访问，比如表达式 `s[n]` 表示访问第 `n` 个元素，索引从零开始。

### 声明以及初始化

```go
func main() {
    var nums [3]int   // 声明并初始化为默认零值
    var nums1 = [4]int{1, 2, 3, 4}  // 声明同时初始化
    var nums2 = [...]int{1, 2, 3, 4, 5} // ...可以表示后面初始化值的长度
    fmt.Println(nums)    // [0 0 0]
    fmt.Println(nums1)   // [1 2 3 4]
    fmt.Println(nums2)   // [1 2 3 4 5]
}
```

### 函数参数

如果数组作为函数的参数，那么实际传递的是一份数组的拷贝，而不是数组的指针。这也就意味着，在函数中修改数组的元素是不会影响到原始数组的。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/arrayparams.png)

```go
package main

import (
    "fmt"
)

func Add(numbers [5]int) {
    for i := 0; i < len(numbers); i++ {
        numbers[i] = numbers[i] + 1
    }
    fmt.Println("numbers in Add:", numbers) // [2 3 4 5 6]
}

func main() {
    // declare and initialize the array
    var numbers [5]int
    for i := 0; i < len(numbers); i++ {
        numbers[i] = i + 1
    }

    Add(numbers)
    fmt.Println("numbers in main:", numbers) // [1 2 3 4 5]
}
```

## 切片

数组的使用场景相对有限，切片才更加常用。

切片（Slice）是一个拥有相同类型元素的可变长度的序列。它是基于数组类型做的一层封装。它非常灵活，支持自动扩容。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/slice.png)

切片是一种引用类型，它有三个属性：**指针**，**长度**和**容量**。

1. 指针：指向 slice 可以访问到的第一个元素。
2. 长度：slice 中元素个数。
3. 容量：slice 起始元素到底层数组最后一个元素间的元素个数。

底层源码定义如下：

```go
type slice struct {
    array unsafe.Pointer
    len   int
    cap   int
}
```

### 声明以及初始化

```go
func main() {
    var nums []int  // 声明切片
    fmt.Println(len(nums), cap(nums)) // 0 0
    nums = append(nums, 1)   // 初始化
    fmt.Println(len(nums), cap(nums)) // 1 1

    nums1 := []int{1,2,3,4}    // 声明并初始化
    fmt.Println(len(nums1), cap(nums1))    // 4 4

    nums2 := make([]int,3,5)   // 使用make()函数构造切片
    fmt.Println(len(nums2), cap(nums2))    // 3 5
}
```

### 函数参数

当切片作为函数参数时，和数组是不同的，如果一个函数接受一个切片参数，它对切片元素所做的更改将对调用者可见，类似于将指针传递给了底层数组。

```go
package main

import (
    "fmt"
)

func Add(numbers []int) {
    for i := 0; i < len(numbers); i++ {
        numbers[i] = numbers[i] + 1
    }
    fmt.Println("numbers in Add:", numbers) // [2 3 4 5 6]
}

func main() {
    var numbers []int
    for i := 0; i < 5; i++ {
        numbers = append(numbers, i+1)
    }

    Add(numbers)

    fmt.Println("numbers in main:", numbers) // [2 3 4 5 6]
}
```

再看一下上面的例子，把参数由数组变成切片，`Add` 函数中的修改会影响到 `main` 函数。

## 总结

最后来总结一下，面试时也可以这么来回答：

1. 数组是一个长度固定的数据类型，其长度在定义时就已经确定，不能动态改变；切片是一个长度可变的数据类型，其长度在定义时可以为空，也可以指定一个初始长度。
2. 数组的内存空间是在定义时分配的，其大小是固定的；切片的内存空间是在运行时动态分配的，其大小是可变的。
3. 当数组作为函数参数时，函数操作的是数组的一个副本，不会影响原始数组；当切片作为函数参数时，函数操作的是切片的引用，会影响原始切片。
4. 切片还有容量的概念，它指的是分配的内存空间。

以上就是本文的全部内容，如果觉得还不错的话欢迎**点赞**，**转发**和**关注**，感谢支持。

***

**参考文章：**

- https://go.dev/doc/effective_go#arrays
- https://go.dev/blog/slices-intro
- https://levelup.gitconnected.com/go-programming-array-vs-slice-5902b7fdd436

**推荐阅读：**

- [Go 语言 new 和 make 关键字的区别](https://mp.weixin.qq.com/s/NBDkI3roHgNgW1iW4e_6cA)
- [为什么 Go 不支持 []T 转换为 []interface](https://mp.weixin.qq.com/s/cwDEgnicK4jkuNpzulU2bw)
- [为什么 Go 语言 struct 要使用 tags](https://mp.weixin.qq.com/s/L7-TJ-CzYfuVrIBWP7Ebaw)