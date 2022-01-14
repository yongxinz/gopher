**原文链接：** [Go Error 嵌套到底是怎么实现的？](https://mp.weixin.qq.com/s/nWb-0RTDG1Pg5ZmJZfbEPA)

Go Error 的设计哲学是 「Errors Are Values」。

这句话应该怎么理解呢？翻译起来挺难的。不过从源码的角度来看，好像更容易理解其背后的含义。

Go Error 源码很简单，寥寥几行：

```go
// src/builtin/builtin.go

type error interface {
	Error() string
}
```

`error` 是一个接口类型，只需要实现 `Error()` 方法即可。在 `Error()` 方法中，就可以返回自定义结构体的任意内容。

下面首先说说如何创建 `error`。

## 创建 Error

创建 `error` 有两种方式，分别是：

1. `errors.New()`；
2. `fmt.Errorf()`。

### errors.New()

`errors.New()` 的使用延续了 Go 的一贯风格，`New` 一下就可以了。

举一个例子：

```go
package main

import (
	"errors"
	"fmt"
)

func main() {
	err := errors.New("这是 errors.New() 创建的错误")
	fmt.Printf("err 错误类型：%T，错误为：%v\n", err, err)
}

/* 输出
err 错误类型：*errors.errorString，错误为：这是 errors.New() 创建的错误
*/
```

这段代码唯一让人困惑的地方可能就是错误类型了，但没关系。只要看一下源码，就瞬间迎刃而解。

源码如下：

```go
// src/errors/errors.go

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
```

可以看到，`errorString` 是一个结构体，实现了 `Error()` 方法，`New` 函数直接返回 `errorString` 指针。

这种用法很简单，但不实用。假如我还想返回程序的上下文信息，它就没辙了。

下面看第二种方式。

### fmt.Errorf()

还是先看一个例子：

```go
package main

import (
	"database/sql"
	"fmt"
)

func foo() error {
	return sql.ErrNoRows
}

func bar() error {
	return foo()
}

func main() {
	err := bar()
	if err == sql.ErrNoRows {
		fmt.Printf("data not found, %+v\n", err)
		return
	}
	if err != nil {
		fmt.Println("Unknown error")
	}
}

/* 输出
data not found, sql: no rows in result set
*/
```

这个例子输出了我们想要的结果，但是还不够。

一般情况下，我们会通过使用 `fmt.Errorf()` 函数，附加上我们想添加的文本信息，使返回内容更明确，处理起来更灵活。

所以，`foo()` 函数会改成下面这样：

```go
func foo() error {
   return fmt.Errorf("foo err, %v", sql.ErrNoRows)
}
```

这时问题就出现了，经过 `fmt.Errorf()` 的封装，原始 `error` 类型发生了改变，这就导致 `err == sql.ErrNoRows` 不再成立，返回信息变成了 `Unknown error`。

如果想根据返回的 `error` 类型做不同处理，就无法实现了。

因此，Go 1.13 为我们提供了 `wrapError` 来处理这个问题。

## Wrap Error

看一个例子：

```go
package main

import (
	"fmt"
)

type myError struct{}

func (e myError) Error() string {
	return "Error happended"
}

func main() {
	e1 := myError{}
	e2 := fmt.Errorf("E2: %w", e1)
	e3 := fmt.Errorf("E3: %w", e2)
	fmt.Println(e2)
	fmt.Println(e3)
}

/* output
E2: Error happended
E3: E2: Error happended
*/
```

乍一看好像好没什么区别，但背后的实现原理却并不相同。

Go 扩展了 `fmt.Errorf()` 函数，增加了一个 `%w` 标识符来创建 `wrapError`。

```go
// src/fmt/errors.go

func Errorf(format string, a ...interface{}) error {
	p := newPrinter()
	p.wrapErrs = true
	p.doPrintf(format, a)
	s := string(p.buf)
	var err error
	if p.wrappedErr == nil {
		err = errors.New(s)
	} else {
		err = &wrapError{s, p.wrappedErr}
	}
	p.free()
	return err
}
```

当使用 `w%` 时，函数会返回 `&wrapError{s, p.wrappedErr}`，`wrapError` 结构体定义如下：

```go
// src/fmt/errors.go

type wrapError struct {
	msg string
	err error
}

func (e *wrapError) Error() string {
	return e.msg
}

func (e *wrapError) Unwrap() error {
	return e.err
}
```

实现了 `Error()` 方法，说明它是一个 `error`，而 `Unwrap()` 方法是为了获取被封装的 `error`。

```go
// src/errors/wrap.go

func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}
```

它们之间的关系是这样的：

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/errors.png)

因此，我们可以使用 `w%` 将上文中的程序进行改造，使其内容输出更丰富。

如下：

```go
package main

import (
	"database/sql"
	"errors"
	"fmt"
)

func bar() error {
	if err := foo(); err != nil {
		return fmt.Errorf("bar failed: %w", foo())
	}
	return nil
}

func foo() error {
	return fmt.Errorf("foo failed: %w", sql.ErrNoRows)
}

func main() {
	err := bar()
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("data not found,  %+v\n", err)
		return
	}
	if err != nil {
		fmt.Println("Unknown error")
	}
}

/* output
data not found,  bar failed: foo failed: sql: no rows in result set
*/
```

终于有了让人满意的输出结果，每个函数都增加了必要的上下文信息，而且也符合对错误类型的判断。

`errors.Is()` 函数用来判断 `err` 以及其封装的 `error` 链中是否包含目标类型。这也就解决了上文提出的无法判断错误类型的问题。

## 后记

其实，Go 目前对 Error 的处理方式也是充满争议的。不过，官方团队正在积极和社区交流，提出改进方法。相信在不久的将来，一定会找到更好的解决方案。

现阶段来说，大部分团队可能会选择 `github.com/pkg/errors` 包来进行错误处理。如果感兴趣的话，可以学学看。

好了，本文就到这里吧。**关注我，带你通过问题读 Go 源码。**

---

**源码地址：**

- [https://github.com/yongxinz/gopher](https://github.com/yongxinz/gopher)

**推荐阅读：**

- [为什么要避免在 Go 中使用 ioutil.ReadAll？](https://mp.weixin.qq.com/s/e2A3ME4vhOK2S3hLEJtPsw)
- [如何在 Go 中将 []byte 转换为 io.Reader？](https://mp.weixin.qq.com/s/nFkob92GOs6Gp75pxA5wCQ)
- [开始读 Go 源码了](https://mp.weixin.qq.com/s/iPM-mPOepRuDqkBtcnG1ww)

**参考文章：**

- https://chasecs.github.io/posts/the-philosophy-of-go-error-handling/
- https://medium.com/@dche423/golang-error-handling-best-practice-cn-42982bd72672
- https://www.flysnow.org/2019/09/06/go1.13-error-wrapping.html