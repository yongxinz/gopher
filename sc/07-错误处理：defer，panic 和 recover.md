![](https://github.com/yongxinz/gopher/blob/main/sc/pic/07_%E9%94%99%E8%AF%AF%E5%A4%84%E7%90%86.png)

**原文链接：** [Go 专栏｜错误处理：defer，panic 和 recover](https://mp.weixin.qq.com/s/qYZXfAifBxwl1cDDaP0FNA)

最近校招又开始了，我也接到了一些面试工作，当我问「你觉得自己有什么优势」时，十个人里有八个的回答里会有一条「精力充沛，能加班」。

怪不得国家都给认证了：新生代农民工。合着我们这根本就不是什么脑力劳动者，而是靠出卖体力的苦劳力。

好了，废话不多说，肝文还确实需要体力。

这篇来说说 Go 的错误处理。

### 错误处理

错误处理相当重要，合理地抛出并记录错误能在排查问题时起到事半功倍的作用。

Go 中有关于错误处理的标准模式，即 `error` 接口，定义如下：

```go
type error interface {
	Error() string
}
```

大部分函数，如果需要返回错误的话，基本都会将 `error` 作为多个返回值的最后一个，举个例子：

```go
package main

import "fmt"

func main() {
	n, err := echo(10)
	if err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(n)
	}
}

func echo(param int) (int, error) {
	return param, nil
}
```

我们也可以使用自定义的 `error` 类型，比如调用标准库的 `os.Stat` 方法，返回的错误就是自定义类型：

```go
type PathError struct {
	Op   string
	Path string
	Err  error
}

func (e *PathError) Error() string {
	return e.Op + " " + e.Path + ": " + e.Err.Error()
}
```

暂时看不懂也没有关系，等学会了接口之后，再回过头来看这段代码，应该就豁然开朗了。

### defer

延迟函数调用，`defer` 后边会接一个函数，但该函数不会立刻被执行，而是等到包含它的程序返回时（包含它的函数执行了 `return` 语句、运行到函数结尾自动返回、对应的 goroutine `panic`），`defer` 函数才会被执行。

通常用于资源释放、打印日志、异常捕获等。

```go
func main() {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    /**
     * 这里defer要写在err判断的后边而不是os.Open后边
     * 如果资源没有获取成功，就没有必要对资源执行释放操作
     * 如果err不为nil而执行资源执行释放操作，有可能导致panic
     */
    defer f.Close()
}
```

`defer` 语句经常成对出现，比如打开和关闭，连接和断开，加锁和解锁。

`defer` 语句在 `return` 语句之后执行。

```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println(triple(4)) // 12
}

func double(x int) (result int) {
	defer func() {
		fmt.Printf("double(%d) = %d\n", x, result)
	}()

	return x + x
}

func triple(x int) (result int) {
	defer func() {
		result += x
	}()

	return double(x)
}
```

切勿在 `for` 循环中使用 `defer` 语句，因为 `defer` 语句不到函数的最后一刻是不会执行的，所以下面这段代码很可能会用尽所有文件描述符。

```go
for _, filename := range filenames {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close()
}
```

一种解决办法是将循环体单独写一个函数，这样每次循环的时候都会调用关闭函数。

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

`defer` 语句的执行是按调用 `defer` 语句的倒序执行。

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

### panic 和 recover

一般情况下，在程序里记录错误日志，就可以帮助我们在碰到异常时快速定位问题。

但还有一些错误比较严重的，比如数组越界访问，程序会主动调用 `panic` 来抛出异常，然后程序退出。

如果不想程序退出的话，可以使用 `recover` 函数来捕获并恢复。

感觉挺不好理解的，但仔细想想其实和 `try-catch` 也没什么区别。

先来看看两个函数的定义：

```go
func panic(interface{})
func recover() interface{}
```

`panic` 参数类型是 `interface{}`，所以可以接收任意参数类型，比如：

```go
panic(404)
panic("network broken")
panic(Error("file not exists"))
```

`recover` 需要在 `defer` 函数中执行，举个例子：

```go
package main

import (
	"fmt"
)

func main() {
	G()
}

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
```

输出：

```
捕获异常: a
b
继续执行
c
```

`F()` 中抛出异常被捕获，`G()` 还可以正常继续执行。如果 `F()` 没有捕获的话，那么 `panic` 会向上传递，直接导致 `G()` 异常，然后程序直接退出。

还有一个场景就是我们自己在调试程序时，可以使用 `panic` 来中断程序，抛出异常，用于排查问题。

这个就不举例了，反正是我们自己调试，怎么爽怎么来就行了。

### 总结

错误处理在开发过程中至关重要，好的错误处理可以使程序更加健壮。而且将错误信息清晰地记录日志，在排查问题时非常有用。

Go 中使用 `error` 类型进行错误处理，还可以在此基础上自定义错误类型。

使用 `defer` 语句进行延迟调用，用来关闭或释放资源。

使用 `panic` 和 `recover` 来抛出错误和恢复。

使用 `panic` 一般有两种情况：

1. 程序遇到无法执行的错误时，主动调用 `panic` 结束运行；
2. 在调试程序时，主动调用 `panic` 结束运行，根据抛出的错误信息来定位问题。

为了程序的健壮性，可以使用 `recover` 捕获错误，恢复程序运行。

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