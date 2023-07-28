**原文链接：** [Go 语言 context 都能做什么？](https://mp.weixin.qq.com/s/7IliODEUt3JpEuzL8K_sOg)

很多 Go 项目的源码，在读的过程中会发现一个很常见的参数 `ctx`，而且基本都是作为函数的第一个参数。

为什么要这么写呢？这个参数到底有什么用呢？带着这样的疑问，我研究了这个参数背后的故事。

开局一张图：

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/IMG_8342.PNG)

核心是 `Context` 接口：

```go
// A Context carries a deadline, cancelation signal, and request-scoped values
// across API boundaries. Its methods are safe for simultaneous use by multiple
// goroutines.
type Context interface {
    // Done returns a channel that is closed when this Context is canceled
    // or times out.
    Done() <-chan struct{}

    // Err indicates why this context was canceled, after the Done channel
    // is closed.
    Err() error

    // Deadline returns the time when this Context will be canceled, if any.
    Deadline() (deadline time.Time, ok bool)

    // Value returns the value associated with key or nil if none.
    Value(key interface{}) interface{}
}
```

包含四个方法：

*   `Done()`：返回一个 channel，当 times out 或者调用 cancel 方法时。
*   `Err()`：返回一个错误，表示取消 ctx 的原因。
*   `Deadline()`：返回截止时间和一个 bool 值。
*   `Value()`：返回 key 对应的值。

有四个结构体实现了这个接口，分别是：`emptyCtx`, `cancelCtx`, `timerCtx` 和 `valueCtx`。

其中 `emptyCtx` 是空类型，暴露了两个方法：

```go
func Background() Context
func TODO() Context
```

一般情况下，会使用 `Background()` 作为根 ctx，然后在其基础上再派生出子 ctx。要是不确定使用哪个 ctx，就使用 `TODO()`。

另外三个也分别暴露了对应的方法：

```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
func WithValue(parent Context, key, val interface{}) Context
```

## 遵循规则

在使用 Context 时，要遵循以下四点规则：

1.  不要将 Context 放入结构体，而是应该作为第一个参数传入，命名为 `ctx`。
2.  即使函数允许，也不要传入 `nil` 的 Context。如果不知道用哪种 Context，可以使用 `context.TODO()`。
3.  使用 Context 的 Value 相关方法只应该用于在程序和接口中传递和请求相关的元数据，不要用它来传递一些可选的参数。
4.  相同的 Context 可以传递给不同的 goroutine；Context 是并发安全的。

## WithCancel

```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
```

`WithCancel` 返回带有新 `Done` 通道的父级副本。当调用返回的 `cancel` 函数或关闭父上下文的 `Done` 通道时，返回的 `ctx` 的 `Done` 通道将关闭。

取消此上下文会释放与其关联的资源，因此在此上下文中运行的操作完成后，代码应立即调用 `cancel`。

举个例子：

这段代码演示了如何使用可取消上下文来防止 goroutine 泄漏。在函数结束时，由 `gen` 启动的 goroutine 将返回而不会泄漏。

```go
package main

import (
    "context"
    "fmt"
)

func main() {
    // gen generates integers in a separate goroutine and
    // sends them to the returned channel.
    // The callers of gen need to cancel the context once
    // they are done consuming generated integers not to leak
    // the internal goroutine started by gen.
    gen := func(ctx context.Context) <-chan int {
        dst := make(chan int)
        n := 1
        go func() {
            for {
                select {
                case <-ctx.Done():
                    return // returning not to leak the goroutine
                case dst <- n:
                    n++
                }
            }
        }()
        return dst
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // cancel when we are finished consuming integers

    for n := range gen(ctx) {
        fmt.Println(n)
        if n == 5 {
            break
        }
    }
}
```

输出：

```go
1
2
3
4
5
```

## WithDeadline

```go
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)
```

`WithDeadline` 返回父上下文的副本，并将截止日期调整为不晚于 `d`。如果父级的截止日期已经早于 `d`，则 `WithDeadline(parent, d)` 在语义上等同于 `parent`。

当截止时间到期、调用返回的取消函数时或当父上下文的 `Done` 通道关闭时，返回的上下文的 `Done` 通道将关闭。

取消此上下文会释放与其关联的资源，因此在此上下文中运行的操作完成后，代码应立即调用取消。

举个例子：

这段代码传递具有截止时间的上下文，来告诉阻塞函数，它应该在到达截止时间时立刻退出。

```go
package main

import (
    "context"
    "fmt"
    "time"
)

const shortDuration = 1 * time.Millisecond

func main() {
    d := time.Now().Add(shortDuration)
    ctx, cancel := context.WithDeadline(context.Background(), d)

    // Even though ctx will be expired, it is good practice to call its
    // cancellation function in any case. Failure to do so may keep the
    // context and its parent alive longer than necessary.
    defer cancel()

    select {
    case <-time.After(1 * time.Second):
        fmt.Println("overslept")
    case <-ctx.Done():
        fmt.Println(ctx.Err())
    }
}
```

输出：

```go
context deadline exceeded
```

## WithTimeout

```go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```

`WithTimeout` 返回 `WithDeadline(parent, time.Now().Add(timeout))`。

取消此上下文会释放与其关联的资源，因此在此上下文中运行的操作完成后，代码应立即调用取消。

举个例子：

这段代码传递带有超时的上下文，以告诉阻塞函数应在超时后退出。

```go
package main

import (
    "context"
    "fmt"
    "time"
)

const shortDuration = 1 * time.Millisecond

func main() {
    // Pass a context with a timeout to tell a blocking function that it
    // should abandon its work after the timeout elapses.
    ctx, cancel := context.WithTimeout(context.Background(), shortDuration)
    defer cancel()

    select {
    case <-time.After(1 * time.Second):
        fmt.Println("overslept")
    case <-ctx.Done():
        fmt.Println(ctx.Err()) // prints "context deadline exceeded"
    }

}
```

输出：

```go
context deadline exceeded
```

## WithValue

```go
func WithValue(parent Context, key, val any) Context
```

`WithValue` 返回父级的副本，其中与 `key` 关联的值为 `val`。

其中键必须是可比较的，并且不应是字符串类型或任何其他内置类型，以避免使用上下文的包之间发生冲突。 `WithValue` 的用户应该定义自己的键类型。

为了避免分配给 `interface{}`，上下文键通常具有具体的 `struct{}` 类型。或者，导出的上下文键变量的静态类型应该是指针或接口。

举个例子：

这段代码演示了如何将值传递到上下文以及如何检索它（如果存在）。

```go
package main

import (
    "context"
    "fmt"
)

func main() {
    type favContextKey string

    f := func(ctx context.Context, k favContextKey) {
        if v := ctx.Value(k); v != nil {
            fmt.Println("found value:", v)
            return
        }
        fmt.Println("key not found:", k)
    }

    k := favContextKey("language")
    ctx := context.WithValue(context.Background(), k, "Go")

    f(ctx, k)
    f(ctx, favContextKey("color"))
}
```

输出：

```go
found value: Go
key not found: color
```

本文的大部分内容，包括代码示例都是翻译自官方文档，代码都是经过验证可以执行的。如果有不是特别清晰的地方，可以直接去读官方文档。

以上就是本文的全部内容，如果觉得还不错的话欢迎**点赞**，**转发**和**关注**，感谢支持。

***

**官方文档：**

*   <https://pkg.go.dev/context@go1.20.5>

**源码分析：**

*   <https://mritd.com/2021/06/27/golang-context-source-code/>
*   <https://www.qtmuniao.com/2020/07/12/go-context/>
*   <https://seekload.net/2021/11/28/go-context.html>

**推荐阅读：**

*   [Go 语言 map 如何顺序读取？](https://mp.weixin.qq.com/s/iScSgfpSE2y14GH7JNRJSA)
*   [Go 语言 map 是并发安全的吗？](https://mp.weixin.qq.com/s/4mDzMdMbunR_p94Du65QOA)
*   [Go 语言切片是如何扩容的？](https://mp.weixin.qq.com/s/VVM8nqs4mMGdFyCNJx16_g)
*   [Go 语言数组和切片的区别](https://mp.weixin.qq.com/s/esaAmAdmV4w3_qjtAzTr4A)
*   [Go 语言 new 和 make 关键字的区别](https://mp.weixin.qq.com/s/NBDkI3roHgNgW1iW4e_6cA)
*   [为什么 Go 不支持 \[\]T 转换为 \[\]interface](https://mp.weixin.qq.com/s/cwDEgnicK4jkuNpzulU2bw)
*   [为什么 Go 语言 struct 要使用 tags](https://mp.weixin.qq.com/s/L7-TJ-CzYfuVrIBWP7Ebaw)

