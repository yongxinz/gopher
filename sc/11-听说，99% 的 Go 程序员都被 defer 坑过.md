**原文链接：** [听说，99% 的 Go 程序员都被 defer 坑过](https://mp.weixin.qq.com/s/1T6Z74Wri27Ap8skeJiyWQ)

先声明：我被坑过。

之前写 Go 专栏时，写过一篇文章：[Go 专栏｜错误处理：defer，panic 和 recover](https://mp.weixin.qq.com/s/qYZXfAifBxwl1cDDaP0FNA)。有小伙伴留言说：道理都懂，但还是不知道怎么用，而且还总出现莫名奇妙的问题。

出问题就对了，这个小东西坏的很，一不留神就出错。

所以，面对这种情况，我们今天就不讲道理了。直接把我珍藏多年的代码一把梭，凭借多年踩坑经历和写 BUG 经验，我要站着把这个坑迈过去。

<p style="text-align:center;color:#1e819e;font-size:1.2em;font-weight: bold;">一、</p>

先来一个简单的例子热热身：

```go
package main

import (
    "fmt"
)

func main() {
    defer func() {
        fmt.Println("first")
    }()

    defer func() {
        fmt.Println("second")
    }()

    fmt.Println("done")
}
```

输出：

```
done
second
first
```

这个比较简单，`defer` 语句的执行顺序是按调用 `defer` 语句的倒序执行。

<p style="text-align:center;color:#1e819e;font-size:1.2em;font-weight: bold;">二、</p>

看看这段代码有什么问题？

```go
for _, filename := range filenames {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close()
}
```

这段代码其实很危险，很可能会用尽所有文件描述符。因为 `defer` 语句不到函数的最后一刻是不会执行的，也就是说文件始终得不到关闭。所以切记，一定不要在 `for` 循环中使用 `defer` 语句。

那怎么优化呢？可以将循环体单独写一个函数，这样每次循环的时候都会调用关闭函数。

如下：

```go
for _, filename := range filenames {
    if err := doFile(filename); err != nil {
        return err
    }
}

func doFile(filename string) error {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close()
}
```

<p style="text-align:center;color:#1e819e;font-size:1.2em;font-weight: bold;">三、</p>

看看这三个函数的输出结果是什么？

```go
package main

import (
	"fmt"
)

func a() (r int) {
	defer func() {
		r++
	}()
	return 0
}

func b() (r int) {
	t := 5
	defer func() {
		t = t + 5
	}()
	return t
}

func c() (r int) {
	defer func(r int) {
		r = r + 5
	}(r)
	return 1
}

func main() {
	fmt.Println("a = ", a())
	fmt.Println("b = ", b())
	fmt.Println("c = ", c())
}
```

公布答案：

```
a =  1
b =  5
c =  1
```

你答对了吗？

说实话刚开始看到这个结果时，我是相当费解，完全不知道怎么回事。

但可以看到，这三个函数都有一个共同特点，它们都有一个命名返回值，并且都在函数中引用了这个返回值。

引用的方式分两种：分别是闭包和函数参数。

先看 `a()` 函数：

闭包通过 `r++` 修改了外部变量，返回值变成了 1。

相当于：

```go
func aa() (r int) {
	r = 0
	// 在 return 之前，执行 defer 函数
	func() {
		r++
	}()
	return
}
```

再看 `b()` 函数：

闭包内修改的只是局部变量 `t`，而外部变量 `t` 不受影响，所以还是返回 5。

相当于：

```go
func bb() (r int) {
	t := 5
	// 赋值
	r = t
	// 在 return 之前，执行 defer 函数
	// defer 函数没有对返回值 r 进行修改，只是修改了变量 t
	func() {
		t = t + 5
	}()
	return
}
```

最后是 `c` 函数：

参数传递是值拷贝，实参不受影响，所以还是返回 1。

相当于：

```go
func cc() (r int) {
	// 赋值
	r = 1
	// 这里修改的 r 是函数形参的值
	// 值拷贝，不影响实参值
	func(r int) {
		r = r + 5
	}(r)
	return
}
```

那么，为了避免写出这么令人意外的代码，最好在定义函数时就不要使用命名返回值。或者如果使用了，就不要在 `defer` 中引用。

再看下面两个例子：

```go
func d() int {
	r := 0
	defer func() {
		r++
	}()
	return r
}

func e() int {
	r := 0
	defer func(i int) {
		i++
	}(r)
	return 0
}
```

```
d =  0
e =  0
```

返回值符合预期，再也不用绞尽脑汁猜了。

<p style="text-align:center;color:#1e819e;font-size:1.2em;font-weight: bold;">四、</p>

`defer` 表达式的函数如果在 `panic` 后面，则这个函数无法被执行。

```go
func main() {
    panic("a")
    defer func() {
        fmt.Println("b")
    }()
}
```

输出如下，`b` 没有打印出来。

```
panic: a

goroutine 1 [running]:
main.main()
	xxx.go:87 +0x4ce
exit status 2
```

而如果 `defer` 在前，则可以执行。

```go
func main() {
	defer func() {
		fmt.Println("b")
	}()
	panic("a")
}
```

输出：

```go
b
panic: a

goroutine 1 [running]:
main.main()
    xxx.go:90 +0x4e7
exit status 2
```

<p style="text-align:center;color:#1e819e;font-size:1.2em;font-weight: bold;">五、</p>

看看下面这段代码的执行顺序：

```go
func G() {
	defer func() {
		fmt.Println("c")
	}()

	F()
	fmt.Println("继续执行")
}

func F() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获异常:", err)
		}
		fmt.Println("b")
	}()
	panic("a")
}

func main() {
	G()
}
```

顺序如下：

1. 调用 `G()` 函数；
2. 调用 `F()` 函数；
3. `F()` 中遇到 `panic`，立刻终止，不执行 `panic` 之后的代码；
4. 执行 `F()` 中 `defer` 函数，遇到 `recover` 捕获错误，继续执行 `defer` 中代码，然后返回；
5. 执行 `G()` 函数后续代码，最后执行 `G()` 中 `defer` 函数。

输出：

```
捕获异常: a
b
继续执行
c
```

<p style="text-align:center;color:#1e819e;font-size:1.2em;font-weight: bold;">五、</p>

看看下面这段代码的执行顺序：

```go
func G() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获异常:", err)
		}
		fmt.Println("c")
	}()

	F()
	fmt.Println("继续执行")
}

func F() {
	defer func() {
		fmt.Println("b")
	}()
	panic("a")
}

func main() {
	G()
}
```

顺序如下：

1. 调用 `G()` 函数；
2. 调用 `F()` 函数；
3. `F()` 中遇到 `panic`，立刻终止，不执行 `panic` 之后的代码；
4. 执行 `F()` 中 `defer` 函数，由于没有 `recover`，则将 `panic` 抛到 `G()` 中；
5. `G()` 收到 `panic` 则不会执行后续代码，直接执行 `defer` 函数；
6. `defer` 中捕获 `F()` 抛出的异常 `a`，然后继续执行，最后退出。

输出：

```
b
捕获异常: a
c
```

<p style="text-align:center;color:#1e819e;font-size:1.2em;font-weight: bold;">六、</p>

看看下面这段代码的执行顺序：

```go
func G() {
	defer func() {
		fmt.Println("c")
	}()

	F()
	fmt.Println("继续执行")
}

func F() {
	defer func() {
		fmt.Println("b")
	}()
	panic("a")
}

func main() {
	G()
}
```

顺序如下：

1. 调用 `G()` 函数；
2. 调用 `F()` 函数；
3. `F()` 中遇到 `panic`，立刻终止，不执行 `panic` 之后的代码；
4. 执行 `F()` 中 `defer` 函数，由于没有 `recover`，则将 `panic` 抛到 `G()` 中；
5. `G()` 收到 `panic` 则不会执行后续代码，直接执行 `defer` 函数；
6. 由于没有 `recover`，直接抛出 `F()` 抛过来的异常 `a`，然后退出。

输出：

```
b
c
panic: a

goroutine 1 [running]:
main.F()
	xxx.go:90 +0x5b
main.G()
	xxx.go:82 +0x48
main.main()
	xxx.go:107 +0x4a5
exit status 2
```

<p style="text-align:center;color:#1e819e;font-size:1.2em;font-weight: bold;">七、</p>

看看下面这段代码的执行顺序：

```go
func G() {
	defer func() {
		// goroutine 外进行 recover
		if err := recover(); err != nil {
			fmt.Println("捕获异常:", err)
		}
		fmt.Println("c")
	}()

	// 创建 goroutine 调用 F 函数
	go F()
	time.Sleep(time.Second)
}

func F() {
	defer func() {
		fmt.Println("b")
	}()
	// goroutine 内部抛出panic
	panic("a")
}

func main() {
	G()
}
```

顺序如下：

1. 调用 `G()` 函数；
2. 通过 goroutine 调用 `F()` 函数；
3. `F()` 中遇到 `panic`，立刻终止，不执行 `panic` 之后的代码；
4. 执行 `F()` 中 `defer` 函数，由于没有 `recover`，则将 `panic` 抛到 `G()` 中；
5. 由于 goroutine 内部没有进行 `recover`，则 goroutine 外部函数，也就是 `G()` 函数是没办法捕获的，程序直接崩溃退出。

输出：

```
b
panic: a

goroutine 6 [running]:
main.F()
	xxx.go:96 +0x5b
created by main.G
	xxx.go:87 +0x57
exit status 2
```

<p style="text-align:center;color:#1e819e;font-size:1.2em;font-weight: bold;">八、</p>

最后再说一个 `recover` 的返回值问题：

```go
defer func() {
	if err := recover(); err != nil {
		fmt.Println("捕获异常:", err.Error())
	}
}()
panic("a")
```

`recover` 返回的是 `interface {}` 类型，而不是 `error` 类型，所以这样使用的话会报错： 

```go
err.Error undefined (type interface {} is interface with no methods)
```

可以这样来转换一下：

```go
defer func() {
	if err := recover(); err != nil {
		fmt.Println("捕获异常:", fmt.Errorf("%v", err).Error())
	}
}()
panic("a")
```

或者直接打印结果：

```go
defer func() {
	if err := recover(); err != nil {
		fmt.Println("捕获异常:", err)
	}
}()
panic("a")
```

输出：

```
捕获异常: a
```

以上就是本文的全部内容，其实写过其他的语言的同学都知道，关闭文件句柄，释放锁等操作是很容易忘的。而 Go 语言通过 `defer` 很好地解决了这个问题，但在使用过程中还是要小心。

本文总结了一些容踩坑的点，希望能够帮助大家少写 BUG，如果大家觉得有用的话，欢迎点赞和转发。

---

文章中的脑图和源码都上传到了 [GitHub](https://github.com/yongxinz/gopher)，有需要的同学可自行下载。

**源码地址：** 

- [https://github.com/yongxinz/gopher/tree/main/sc](https://github.com/yongxinz/gopher/tree/main/sc)

**推荐阅读：**

- [gRPC，爆赞](https://mp.weixin.qq.com/s/1Xbca4Dv0akonAZerrChgA)
- [使用 grpcurl 通过命令行访问 gRPC 服务](https://mp.weixin.qq.com/s/GShwcGCopXVmxCKnYf5FhA)
- [推荐三个实用的 Go 开发工具](https://mp.weixin.qq.com/s/3GLMLhegB3wF5_62mpmePA)
- [被 Docker 日志坑惨了](https://mp.weixin.qq.com/s/3Tkc15dTCEDUAZaZ88pcSQ)

**参考：**

- 《Go 语言核心编程》
- [https://www.jianshu.com/p/63e3d57f285f](https://www.jianshu.com/p/63e3d57f285f)