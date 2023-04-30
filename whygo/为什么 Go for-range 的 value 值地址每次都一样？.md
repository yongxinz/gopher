**原文链接：** [为什么 Go for-range 的 value 值地址每次都一样？](https://mp.weixin.qq.com/s/OoJ42UVYe72492mRUGtdvA)

循环语句是一种常用的控制结构，在 Go 语言中，除了 `for` 关键字以外，还有一个 `range` 关键字，可以使用 `for-range`  循环迭代数组、切片、字符串、map 和 channel 这些数据类型。

但是在使用 `for-range` 循环迭代数组和切片的时候，是很容易出错的，甚至很多老司机一不小心都会在这里翻车。

具体是怎么翻的呢？我们接着看。

## 现象

先来看两段很有意思的代码：

### 无限循环

如果我们在遍历数组的同时向数组中添加元素，能否得到一个永远都不会停止的循环呢？

比如下面这段代码：

```go
func main() {
    arr := []int{1, 2, 3}
    for _, v := range arr {
        arr = append(arr, v)
    }
    fmt.Println(arr)
}
```

程序输出：

```go
$ go run main.go
1 2 3 1 2 3
```

上述代码的输出意味着循环只遍历了原始切片中的三个元素，我们在遍历切片时追加的元素并没有增加循环的执行次数，所以循环最终还是停了下来。

### 相同地址

第二个例子是使用 Go 语言经常会犯的一个错误。

当我们在遍历一个数组时，如果获取 `range` 返回变量的地址并保存到另一个数组或者哈希时，会遇到令人困惑的现象：

```go
func main() {
    arr := []int{1, 2, 3}
    newArr := []*int{}
    for _, v := range arr {
        newArr = append(newArr, &v)
    }
    for _, v := range newArr {
        fmt.Println(*v)
    }
}
```

程序输出：

```go
$ go run main.go
3 3 3
```

上述代码并没有输出 `1 2 3`，而是输出 `3 3 3`。

正确的做法应该是使用 `&arr[i]` 替代 `&v`，像这种编程中的细节是很容易出错的。

## 原因

具体原因也并不复杂，一句话就能解释。

对于数组、切片或字符串，每次迭代，`for-range` 语句都会将原始值的副本传递给迭代变量，而非原始值本身。

口说无凭，具体是不是这样，还得靠源码说话。

Go 编译器会将 `for-range` 语句转换成类似 C 语言的[**三段式循环**](https://github.com/golang/gofrontend/blob/e387439bfd24d5e142874b8e68e7039f74c744d7/go/statements.cc#L5384)结构，就像这样：

```go
// Arrange to do a loop appropriate for the type.  We will produce
//   for INIT ; COND ; POST {
//           ITER_INIT
//           INDEX = INDEX_TEMP
//           VALUE = VALUE_TEMP // If there is a value
//           original statements
//   }
```

迭代[**数组**](https://github.com/golang/gofrontend/blob/e387439bfd24d5e142874b8e68e7039f74c744d7/go/statements.cc#L5501)时，是这样：

```go
// The loop we generate:
//   len_temp := len(range)
//   range_temp := range
//   for index_temp = 0; index_temp < len_temp; index_temp++ {
//           value_temp = range_temp[index_temp]
//           index = index_temp
//           value = value_temp
//           original body
//   }
```

[**切片**](https://github.com/golang/gofrontend/blob/e387439bfd24d5e142874b8e68e7039f74c744d7/go/statements.cc#L5593)：

```go
//   for_temp := range
//   len_temp := len(for_temp)
//   for index_temp = 0; index_temp < len_temp; index_temp++ {
//           value_temp = for_temp[index_temp]
//           index = index_temp
//           value = value_temp
//           original body
//   }
```

从上面的代码片段，可以总结两点：

1. 在循环开始前，会将数组或切片赋值给一个新变量，在赋值过程中就发生了拷贝，迭代的实际上是副本，这也就解释了现象 1。
2. 在循环过程中，会将迭代元素赋值给一个临时变量，这又发生了拷贝。如果取地址的话，每次都是一样的，都是临时变量的地址。

以上就是本文的全部内容，如果觉得还不错的话欢迎**点赞**，**转发**和**关注**，感谢支持。

***

**参考文章：**

- <https://garbagecollected.org/2017/02/22/go-range-loop-internals/>
- <https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-for-range/>

**推荐阅读：**

- [为什么 Go 不支持 []T 转换为 []interface](https://mp.weixin.qq.com/s/cwDEgnicK4jkuNpzulU2bw)
- [为什么 Go 语言 struct 要使用 tags](https://mp.weixin.qq.com/s/L7-TJ-CzYfuVrIBWP7Ebaw)