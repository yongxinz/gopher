**原文链接：** [Go 语言切片是如何扩容的？](https://mp.weixin.qq.com/s/VVM8nqs4mMGdFyCNJx16_g)

在 Go 语言中，有一个很常用的数据结构，那就是切片（Slice）。

切片是一个拥有相同类型元素的可变长度的序列，它是基于数组类型做的一层封装。它非常灵活，支持自动扩容。

切片是一种引用类型，它有三个属性：**指针**，**长度**和**容量**。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/slict1.png)

底层源码定义如下：

```go
type slice struct {
    array unsafe.Pointer
    len   int
    cap   int
}
```

1. **指针：** 指向 slice 可以访问到的第一个元素。
2. **长度：** slice 中元素个数。
3. **容量：** slice 起始元素到底层数组最后一个元素间的元素个数。

比如使用 `make([]byte, 5)` 创建一个切片，它看起来是这样的：

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/slice2.png)

## 声明和初始化

切片的使用还是比较简单的，这里举一个例子，直接看代码吧。

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

## 扩容时机

当切片的长度超过其容量时，切片会自动扩容。这通常发生在使用 `append` 函数向切片中添加元素时。

扩容时，Go 运行时会分配一个新的底层数组，并将原始切片中的元素复制到新数组中。然后，原始切片将指向新数组，并更新其长度和容量。

需要注意的是，由于**扩容会分配新数组并复制元素，因此可能会影响性能**。如果你知道要添加多少元素，可以使用 `make` 函数预先分配足够大的切片来避免频繁扩容。

接下来看看 `append` 函数，签名如下：

```go
func Append(slice []int, items ...int) []int
```

`append` 函数参数长度可变，可以追加多个值，还可以直接追加一个切片。使用起来比较简单，分别看两个例子：

**追加多个值：**

```go
package main

import "fmt"

func main() {
    s := []int{1, 2, 3}
    fmt.Println("初始切片:", s)

    s = append(s, 4, 5, 6)
    fmt.Println("追加多个值后的切片:", s)
}
```

输出结果为：

```go
初始切片: [1 2 3]
追加多个值后的切片: [1 2 3 4 5 6]
```

再来看一下直接**追加一个切片：**

```go
package main

import "fmt"

func main() {
    s1 := []int{1, 2, 3}
    fmt.Println("初始切片:", s1)

    s2 := []int{4, 5, 6}
    s1 = append(s1, s2...)
    fmt.Println("追加另一个切片后的切片:", s1)
}
```

输出结果为：

```go
初始切片: [1 2 3]
追加另一个切片后的切片: [1 2 3 4 5 6]
```

再来看一个**发生扩容**的例子：

```go
package main

import "fmt"

func main() {
    s := make([]int, 0, 3) // 创建一个长度为0，容量为3的切片
    fmt.Printf("初始状态: len=%d cap=%d %v\n", len(s), cap(s), s)

    for i := 1; i <= 5; i++ {
        s = append(s, i) // 向切片中添加元素
        fmt.Printf("添加元素%d: len=%d cap=%d %v\n", i, len(s), cap(s), s)
    }
}
```

输出结果为：

```go
初始状态: len=0 cap=3 []
添加元素1: len=1 cap=3 [1]
添加元素2: len=2 cap=3 [1 2]
添加元素3: len=3 cap=3 [1 2 3]
添加元素4: len=4 cap=6 [1 2 3 4]
添加元素5: len=5 cap=6 [1 2 3 4 5]
```

在这个例子中，我们创建了一个长度为 `0`，容量为 `3` 的切片。然后，我们使用 `append` 函数向切片中添加 `5` 个元素。

当我们添加第 `4` 个元素时，切片的长度超过了其容量。此时，切片会自动扩容。新的容量是原始容量的两倍，即 `6`。

表面现象已经看到了，接下来，我们就深入到源码层面，看看切片的扩容机制到底是什么样的。

## 源码分析

在 Go 语言的源码中，切片扩容通常是在进行切片的 `append` 操作时触发的。在进行 `append` 操作时，如果切片容量不足以容纳新的元素，就需要对切片进行扩容，此时就会调用 `growslice` 函数进行扩容。

`growslice` 函数定义在 Go 语言的 runtime 包中，它的调用是在编译后的代码中实现的。具体来说，当执行 `append` 操作时，编译器会将其转换为类似下面的代码：

```go
slice = append(slice, elem)
```

在上述代码中，如果切片容量不足以容纳新的元素，则会调用 `growslice` 函数进行扩容。所以 `growslice` 函数的调用是**由编译器在生成的机器码中实现的，而不是在源代码中显式调用的**。

切片扩容策略有两个阶段，go1.18 之前和之后是不同的，这一点在 go1.18 的 release notes 中有说明。

下面我用 go1.17 和 go1.18 两个版本来分开说明。先通过一段测试代码，直观感受一下两个版本在扩容上的区别。

```go
package main

import "fmt"

func main() {
    s := make([]int, 0)

    oldCap := cap(s)

    for i := 0; i < 2048; i++ {
        s = append(s, i)

        newCap := cap(s)

        if newCap != oldCap {
            fmt.Printf("[%d -> %4d] cap = %-4d  |  after append %-4d  cap = %-4d\n", 0, i-1, oldCap, i, newCap)
            oldCap = newCap
        }
    }
}
```

上述代码先创建了一个空的 slice，然后在一个循环里不断往里面 `append` 新元素。

然后记录容量的变化，每当容量发生变化的时候，记录下老的容量，添加的元素，以及添加完元素之后的容量。

这样就可以观察，新老 slice 的容量变化情况，从而找出规律。

运行结果（**1.17 版本**）：

```go
[0 ->   -1] cap = 0     |  after append 0     cap = 1   
[0 ->    0] cap = 1     |  after append 1     cap = 2   
[0 ->    1] cap = 2     |  after append 2     cap = 4   
[0 ->    3] cap = 4     |  after append 4     cap = 8   
[0 ->    7] cap = 8     |  after append 8     cap = 16  
[0 ->   15] cap = 16    |  after append 16    cap = 32  
[0 ->   31] cap = 32    |  after append 32    cap = 64  
[0 ->   63] cap = 64    |  after append 64    cap = 128 
[0 ->  127] cap = 128   |  after append 128   cap = 256 
[0 ->  255] cap = 256   |  after append 256   cap = 512 
[0 ->  511] cap = 512   |  after append 512   cap = 1024
[0 -> 1023] cap = 1024  |  after append 1024  cap = 1280
[0 -> 1279] cap = 1280  |  after append 1280  cap = 1696
[0 -> 1695] cap = 1696  |  after append 1696  cap = 2304
```

运行结果（**1.18 版本**）：

```go
[0 ->   -1] cap = 0     |  after append 0     cap = 1
[0 ->    0] cap = 1     |  after append 1     cap = 2   
[0 ->    1] cap = 2     |  after append 2     cap = 4   
[0 ->    3] cap = 4     |  after append 4     cap = 8   
[0 ->    7] cap = 8     |  after append 8     cap = 16  
[0 ->   15] cap = 16    |  after append 16    cap = 32  
[0 ->   31] cap = 32    |  after append 32    cap = 64  
[0 ->   63] cap = 64    |  after append 64    cap = 128 
[0 ->  127] cap = 128   |  after append 128   cap = 256 
[0 ->  255] cap = 256   |  after append 256   cap = 512 
[0 ->  511] cap = 512   |  after append 512   cap = 848 
[0 ->  847] cap = 848   |  after append 848   cap = 1280
[0 -> 1279] cap = 1280  |  after append 1280  cap = 1792
[0 -> 1791] cap = 1792  |  after append 1792  cap = 2560
```

根据上面的结果还是能看到区别的，具体扩容策略下面边看源码边说明。

### go1.17

扩容调用的是 `growslice` 函数，我复制了其中计算新容量部分的代码。

```go
// src/runtime/slice.go

func growslice(et *_type, old slice, cap int) slice {
    // ...

    newcap := old.cap
    doublecap := newcap + newcap
    if cap > doublecap {
        newcap = cap
    } else {
        if old.cap < 1024 {
            newcap = doublecap
        } else {
            // Check 0 < newcap to detect overflow
            // and prevent an infinite loop.
            for 0 < newcap && newcap < cap {
                newcap += newcap / 4
            }
            // Set newcap to the requested cap when
            // the newcap calculation overflowed.
            if newcap <= 0 {
                newcap = cap
            }
        }
    }

    // ...

    return slice{p, old.len, newcap}
}
```

在分配内存空间之前需要先确定新的切片容量，运行时根据切片的当前容量选择不同的策略进行扩容：

1. 如果期望容量大于当前容量的两倍就会使用期望容量；
2. 如果当前切片的长度小于 1024 就会将容量翻倍；
3. 如果当前切片的长度大于等于 1024 就会每次增加 25% 的容量，直到新容量大于期望容量；

### go1.18

```go
// src/runtime/slice.go

func growslice(et *_type, old slice, cap int) slice {
    // ...

    newcap := old.cap
    doublecap := newcap + newcap
    if cap > doublecap {
        newcap = cap
    } else {
        const threshold = 256
        if old.cap < threshold {
            newcap = doublecap
        } else {
            // Check 0 < newcap to detect overflow
            // and prevent an infinite loop.
            for 0 < newcap && newcap < cap {
                // Transition from growing 2x for small slices
                // to growing 1.25x for large slices. This formula
                // gives a smooth-ish transition between the two.
                newcap += (newcap + 3*threshold) / 4
            }
            // Set newcap to the requested cap when
            // the newcap calculation overflowed.
            if newcap <= 0 {
                newcap = cap
            }
        }
    }

    // ...

    return slice{p, old.len, newcap}
}
```

和之前版本的区别，主要在扩容阈值，以及这行代码：`newcap += (newcap + 3*threshold) / 4`。

在分配内存空间之前需要先确定新的切片容量，运行时根据切片的当前容量选择不同的策略进行扩容：

1. 如果期望容量大于当前容量的两倍就会使用期望容量；
2. 如果当前切片的长度小于阈值（默认 256）就会将容量翻倍；
3. 如果当前切片的长度大于等于阈值（默认 256），就会每次增加 25% 的容量，基准是 `newcap + 3*threshold`，直到新容量大于期望容量；

### 内存对齐

分析完两个版本的扩容策略之后，再看前面的那段测试代码，就会发现扩容之后的容量并不是严格按照这个策略的。

那是为什么呢？

实际上，`growslice` 的后半部分还有更进一步的优化（内存对齐等），靠的是 `roundupsize` 函数，在计算完 `newcap` 值之后，还会有一个步骤计算最终的容量：

```go
capmem = roundupsize(uintptr(newcap) * ptrSize)
newcap = int(capmem / ptrSize)
```

这个函数的实现就不在这里深入了，先挖一个坑，以后再来补上。

## 总结

切片扩容通常是在进行切片的 `append` 操作时触发的。在进行 `append` 操作时，如果切片容量不足以容纳新的元素，就需要对切片进行扩容，此时就会调用 `growslice` 函数进行扩容。

切片扩容分两个阶段，分为 go1.18 之前和之后：

**一、go1.18 之前：**

1. 如果期望容量大于当前容量的两倍就会使用期望容量；
2. 如果当前切片的长度小于 1024 就会将容量翻倍；
3. 如果当前切片的长度大于 1024 就会每次增加 25% 的容量，直到新容量大于期望容量；

**二、go1.18 之后：**

1. 如果期望容量大于当前容量的两倍就会使用期望容量；
2. 如果当前切片的长度小于阈值（默认 256）就会将容量翻倍；
3. 如果当前切片的长度大于等于阈值（默认 256），就会每次增加 25% 的容量，基准是 `newcap + 3*threshold`，直到新容量大于期望容量；

以上就是本文的全部内容，如果觉得还不错的话欢迎**点赞**，**转发**和**关注**，感谢支持。

***

**参考文章：**

- https://go.dev/doc/go1.18
- https://go.dev/blog/slices
- https://go.dev/blog/slices-intro
- https://golang.design/go-questions/slice/grow/
- https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-array-and-slice/

**推荐阅读：**

- [Go 语言数组和切片的区别](https://mp.weixin.qq.com/s/esaAmAdmV4w3_qjtAzTr4A)
- [Go 语言 new 和 make 关键字的区别](https://mp.weixin.qq.com/s/NBDkI3roHgNgW1iW4e_6cA)
- [为什么 Go 不支持 []T 转换为 []interface](https://mp.weixin.qq.com/s/cwDEgnicK4jkuNpzulU2bw)
- [为什么 Go 语言 struct 要使用 tags](https://mp.weixin.qq.com/s/L7-TJ-CzYfuVrIBWP7Ebaw)
