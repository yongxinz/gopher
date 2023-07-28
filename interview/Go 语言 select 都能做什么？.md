**原文链接：** [Go 语言 select 都能做什么？](https://mp.weixin.qq.com/s/YyyMzYxMi8I4HEaxzy4c7g)

在 Go 语言中，`select` 是一个关键字，用于监听和 `channel` 有关的 IO 操作。

通过 `select` 语句，我们可以同时监听多个 `channel`，并在其中任意一个 `channel` 就绪时进行相应的处理。

本文将总结一下 `select` 语句的常见用法，以及在使用过程中的注意事项。

## 基本语法

`select` 语句的基本语法如下：

```go
select {
case <-channel1:
    // 通道 channel1 就绪时的处理逻辑
case data := <-channel2:
    // 通道 channel2 就绪时的处理逻辑
default:
    // 当没有任何通道就绪时的默认处理逻辑
}
```

看到这个语法，很容易想到 `switch` 语句。

虽然 `select` 语句和 `switch` 语句在表面上有些相似，但它们的用途和功能是不同的。

`switch` 用于条件判断，而 `select` 用于通道操作。不能在 `select` 语句中使用任意类型的条件表达式，只能对通道进行操作。

## 使用规则

虽然语法简单，但是在使用过程中，还是有一些地方需要注意，我总结了如下四点：

1.  `select` 语句只能用于通道操作，用于在多个通道之间进行选择，以监听通道的就绪状态，而不是用于其他类型的条件判断。
2.  `select` 语句可以包含多个 `case` 子句，每个 `case` 子句对应一个通道操作。当其中任意一个通道就绪时，相应的 `case` 子句会被执行。
3.  如果多个通道都已经就绪，`select` 语句会随机选择一个通道来执行。这样确保了多个通道之间的公平竞争。
4.  `select` 语句的执行可能是阻塞的，也可能是非阻塞的。如果没有任何一个通道就绪且没有默认的 `default` 子句，`select` 语句会阻塞，直到有一个通道就绪。如果有 `default` 子句，且没有任何通道就绪，那么 `select` 语句会执行 `default` 子句，从而避免阻塞。

## 多路复用

`select` 最常见的用途之一，同时监听多个通道，并根据它们的就绪状态执行不同的操作。

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    c1 := make(chan string)
    c2 := make(chan string)

    go func() {
        time.Sleep(3 * time.Second)
        c1 <- "one"
    }()

    go func() {
        time.Sleep(3 * time.Second)
        c2 <- "two"
    }()

    select {
    case msg := <-c1:
        fmt.Println(msg)
    case msg := <-c2:
        fmt.Println(msg)
    }
}
```

执行上面的代码，程序会随机打印 `one` 或者 `two`，如果通道为空的话，程序就会一直阻塞在那里。

## 非阻塞通信

当通道中没有数据可读或者没有缓冲空间可写时，普通的读写操作将会阻塞。

但通过 `select` 语句，我们可以在没有数据就绪时执行默认的逻辑，避免程序陷入无限等待状态。

```go
package main

import (
    "fmt"
)

func main() {
    channel := make(chan int)

    select {
    case data := <-channel:
        fmt.Println("Received:", data)
    default:
        fmt.Println("No data available.")
    }
}
```

执行上面代码，程序会执行 `default` 分支。

输出：

```go
No data available.
```

## 超时处理

通过结合 `select` 和 `time.After` 函数，我们可以在指定时间内等待通道就绪，超过时间后执行相应的逻辑。

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    channel := make(chan int)

    select {
    case data := <-channel:
        fmt.Println("Received:", data)
    case <-time.After(3 * time.Second):
        fmt.Println("Timeout occurred.")
    }
}
```

执行上面代码，如果 `channel` 在 `3` 秒内没有数据可读，`select` 会执行 `time.After` 分支。

输出：

```go
Timeout occurred.
```

以上就是本文的全部内容，如果觉得还不错的话欢迎**点赞**，**转发**和**关注**，感谢支持。

***

**推荐阅读：**

*   [Go 语言 context 都能做什么？](https://mp.weixin.qq.com/s/7IliODEUt3JpEuzL8K_sOg)

