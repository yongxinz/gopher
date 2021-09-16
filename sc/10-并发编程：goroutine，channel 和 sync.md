![](https://github.com/yongxinz/gopher/blob/main/sc/pic/10_%E5%B9%B6%E5%8F%91%E7%BC%96%E7%A8%8B.png)

**原文链接：** [Go 专栏｜并发编程：goroutine，channel 和 sync](https://mp.weixin.qq.com/s/VG4CSfT2OfxA6nfygWLSyw)

优雅的并发编程范式，完善的并发支持，出色的并发性能是 Go 语言区别于其他语言的一大特色。

在当今这个多核时代，并发编程的意义不言而喻。使用 Go 开发并发程序，操作起来非常简单，语言级别提供关键字 `go` 用于启动协程，并且在同一台机器上可以启动成千上万个协程。

下面就来详细介绍。

### goroutine

Go 语言的并发执行体称为 goroutine，使用关键词 `go` 来启动一个 goroutine。

`go` 关键词后面必须跟一个函数，可以是有名函数，也可以是无名函数，函数的返回值会被忽略。

`go` 的执行是非阻塞的。

先来看一个例子：

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	go spinner(100 * time.Millisecond)
	const n = 45
	fibN := fib(n)
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN) // Fibonacci(45) = 1134903170
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}
```

从执行结果来看，成功计算出了斐波那契数列的值，说明程序在 `spinner` 处并没有阻塞，而且 `spinner` 函数还一直在屏幕上打印提示字符，说明程序正在执行。

当计算完斐波那契数列的值，`main` 函数打印结果并退出，`spinner` 也跟着退出。

再来看一个例子，循环执行 10 次，打印两个数的和：

```go
package main

import "fmt"

func Add(x, y int) {
	z := x + y
	fmt.Println(z)
}

func main() {
	for i := 0; i < 10; i++ {
		go Add(i, i)
	}
}
```

有问题了，屏幕上什么都没有，为什么呢？

这就要看 Go 程序的执行机制了。当一个程序启动时，只有一个 goroutine 来调用 `main` 函数，称为主 goroutine。新的 goroutine 通过 `go` 关键词创建，然后并发执行。当 `main` 函数返回时，不会等待其他 goroutine 执行完，而是直接暴力结束所有 goroutine。

那有没有办法解决呢？当然是有的，请往下看。

### channel

一般写多进程程序时，都会遇到一个问题：进程间通信。常见的通信方式有信号，共享内存等。goroutine 之间的通信机制是通道 channel。

使用 `make` 创建通道：

```go
ch := make(chan int) // ch 的类型是 chan int
```

通道支持三个主要操作：`send`，`receive` 和 `close`。

```go
ch <- x // 发送
x = <-ch // 接收
<-ch // 接收，丢弃结果

close(ch) // 关闭
```

#### 无缓冲 channel

`make` 函数接受两个参数，第二个参数是可选参数，表示通道容量。不传或者传 0 表示创建了一个无缓冲通道。

无缓冲通道上的发送操作将会阻塞，直到另一个 goroutine 在对应的通道上执行接收操作。相反，如果接收先执行，那么接收 goroutine 将会阻塞，直到另一个 goroutine 在对应通道上执行发送。

所以，无缓冲通道是一种同步通道。

下面我们使用无缓冲通道把上面例子中出现的问题解决一下。

```go
package main

import "fmt"

func Add(x, y int, ch chan int) {
	z := x + y
	ch <- z
}

func main() {

	ch := make(chan int)
	for i := 0; i < 10; i++ {
		go Add(i, i, ch)
	}

	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}
}
```

可以正常输出结果。

主 goroutine 会阻塞，直到读取到通道中的值，程序继续执行，最后退出。

#### 缓冲 channel

创建一个容量是 5 的缓冲通道：

```go
ch := make(chan int, 5)
```

缓冲通道的发送操作在通道尾部插入一个元素，接收操作从通道的头部移除一个元素。如果通道满了，发送会阻塞，直到另一个 goroutine 执行接收。相反，如果通道是空的，接收会阻塞，直到另一个 goroutine 执行发送。

有没有感觉，其实缓冲通道和队列一样，把操作都解耦了。

#### 单向 channel

类型 `chan<- int` 是一个只能发送的通道，类型 `<-chan int` 是一个只能接收的通道。

任何双向通道都可以用作单向通道，但反过来不行。

还有一点需要注意，`close` 只能用在发送通道上，如果用在接收通道会报错。

看一个单向通道的例子：

```go
package main

import "fmt"

func counter(out chan<- int) {
	for x := 0; x < 10; x++ {
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for v := range in {
		out <- v * v
	}
	close(out)
}

func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

func main() {
	n := make(chan int)
	s := make(chan int)

	go counter(n)
	go squarer(s, n)
	printer(s)

}
```

### sync

sync 包提供了两种锁类型：`sync.Mutex` 和 `sync.RWMutex`，前者是互斥锁，后者是读写锁。

当一个 goroutine 获取了 `Mutex` 后，其他 goroutine 不管读写，只能等待，直到锁被释放。

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mutex sync.Mutex
	wg := sync.WaitGroup{}

	// 主 goroutine 先获取锁
	fmt.Println("Locking  (G0)")
	mutex.Lock()
	fmt.Println("locked (G0)")

	wg.Add(3)
	for i := 1; i < 4; i++ {
		go func(i int) {
			// 由于主 goroutine 先获取锁，程序开始 5 秒会阻塞在这里
			fmt.Printf("Locking (G%d)\n", i)
			mutex.Lock()
			fmt.Printf("locked (G%d)\n", i)

			time.Sleep(time.Second * 2)
			mutex.Unlock()
			fmt.Printf("unlocked (G%d)\n", i)

			wg.Done()
		}(i)
	}

	// 主 goroutine 5 秒后释放锁
	time.Sleep(time.Second * 5)
	fmt.Println("ready unlock (G0)")
	mutex.Unlock()
	fmt.Println("unlocked (G0)")

	wg.Wait()
}
```

`RWMutex` 属于经典的单写多读模型，当读锁被占用时，会阻止写，但不阻止读。而写锁会阻止写和读。

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var rwMutex sync.RWMutex
	wg := sync.WaitGroup{}

	Data := 0
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func(t int) {
			// 第一次运行后，写解锁。
			// 循环到第二次时，读锁定后，goroutine 没有阻塞，同时读成功。
			fmt.Println("Locking")
			rwMutex.RLock()
			defer rwMutex.RUnlock()
			fmt.Printf("Read data: %v\n", Data)
			wg.Done()
			time.Sleep(2 * time.Second)
		}(i)
		go func(t int) {
			// 写锁定下是需要解锁后才能写的
			rwMutex.Lock()
			defer rwMutex.Unlock()
			Data += t
			fmt.Printf("Write Data: %v %d \n", Data, t)
			wg.Done()
			time.Sleep(2 * time.Second)
		}(i)
	}

	wg.Wait()
}
```

### 总结

并发编程算是 Go 的特色，也是核心功能之一了，涉及的知识点其实是非常多的，本文也只是起到一个抛砖引玉的作用而已。

本文开始介绍了 goroutine 的简单用法，然后引出了通道的概念。

通道有三种：

1. 无缓冲通道
2. 缓冲通道
3. 单向通道

最后介绍了 Go 中的锁机制，分别是 sync 包提供的 `sync.Mutex`（互斥锁） 和 `sync.RWMutex`（读写锁）。

goroutine 博大精深，后面的坑还是要慢慢踩的。

---

文章中的脑图和源码都上传到了 GitHub，有需要的同学可自行下载。

**地址：** https://github.com/yongxinz/gopher/tree/main/sc

关注公众号 **AlwaysBeta**，回复「**goebook**」领取 Go 编程经典书籍。

<center class="half">
    <img src="https://github.com/yongxinz/gopher/blob/main/alwaysbeta.JPG" width="300"/>
</center>

**Go 专栏文章列表：**

1. [Go 专栏｜开发环境搭建以及开发工具 VS Code 配置](https://mp.weixin.qq.com/s/x1OW--3mwSTjgB2HaKGVVA)
2. [Go 专栏｜变量和常量的声明与赋值](https://mp.weixin.qq.com/s/cIceTj02bGa0BYqu-JN1Bg)
3. [Go 专栏｜基础数据类型：整数、浮点数、复数、布尔值和字符串](https://mp.weixin.qq.com/s/aotpxglSGRFfl6A1xPN-dw)
4. [Go 专栏｜复合数据类型：数组和切片 slice](https://mp.weixin.qq.com/s/MnjIeJPUAA6n48o4yns3hg)
5. [Go 专栏｜复合数据类型：字典 map 和 结构体 struct](https://mp.weixin.qq.com/s/1unl6K9xHxy4V3KukORC3A)
6. [Go 专栏｜流程控制，一网打尽](https://mp.weixin.qq.com/s/TbjT1dmTvwiKCzzbWc23kA)
7. [Go 专栏｜函数那些事](https://mp.weixin.qq.com/s/RKpyVrhtSk9pXMWNVpWYjQ)
8. [Go 专栏｜错误处理：defer，panic 和 recover](https://mp.weixin.qq.com/s/qYZXfAifBxwl1cDDaP0FNA)
9. [Go 专栏｜说说方法](https://mp.weixin.qq.com/s/qvFipY0pnmqxok6CVKquvg)
10. [Go 专栏｜接口 interface](https://mp.weixin.qq.com/s/g7ngRIxxbd-M8K_sL_M4KQ)