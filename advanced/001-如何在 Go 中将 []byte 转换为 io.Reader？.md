**原文链接：** [如何在 Go 中将 []byte 转换为 io.Reader？](https://mp.weixin.qq.com/s/nFkob92GOs6Gp75pxA5wCQ)

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/001-byte-slice-to-io-reader.png)

在 stackoverflow 上看到一个问题，题主进行了一个网络请求，接口返回的是 `[]byte`。如果想要将其转换成 `io.Reader`，需要怎么做呢？

这个问题解决起来并不复杂，简单几行代码就可以轻松将其转换成功。不仅如此，还可以再通过几行代码反向转换回来。

下面听我慢慢给你吹，首先直接看两段代码。

### []byte 转 io.Reader

```go
package main

import (
	"bytes"
	"fmt"
	"log"
)

func main() {
	data := []byte("Hello AlwaysBeta")

	// byte slice to bytes.Reader, which implements the io.Reader interface
	reader := bytes.NewReader(data)

	// read the data from reader
	buf := make([]byte, len(data))
	if _, err := reader.Read(buf); err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(buf))
}
```

输出：

```
Hello AlwaysBeta
```

这段代码先将 `[]byte` 数据转换到 `reader` 中，然后再从 `reader` 中读取数据，并打印输出。

### io.Reader 转 []byte

```go
package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	ioReaderData := strings.NewReader("Hello AlwaysBeta")

	// creates a bytes.Buffer and read from io.Reader
	buf := &bytes.Buffer{}
	buf.ReadFrom(ioReaderData)

	// retrieve a byte slice from bytes.Buffer
	data := buf.Bytes()

	// only read the left bytes from 6
	fmt.Println(string(data[6:]))
}
```

输出：

```
AlwaysBeta
```

这段代码先创建了一个 `reader`，然后读取数据到 `buf`，最后打印输出。

以上两段代码就是 `[]byte` 和 `io.Reader` 互相转换的过程。对比这两段代码不难发现，都有 `NewReader` 的身影。而且在转换过程中，都起到了关键作用。

那么问题来了，这个 `NewReader` 到底是什么呢？接下来我们通过源码来一探究竟。

### 源码解析

Go 的 `io` 包提供了最基本的 IO 接口，其中 `io.Reader` 和 `io.Writer` 两个接口最为关键，很多原生结构都是围绕这两个接口展开的。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/io-reader-writer.png)

下面就来分别说说这两个接口：

#### Reader 接口

`io.Reader` 表示一个读取器，它将数据从某个资源读取到传输缓冲区。在缓冲区中，数据可以被流式传输和使用。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/io-reader.png)

接口定义如下：

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

`Read()` 方法将 `len(p)` 个字节读取到 `p` 中。它返回读取的字节数 `n`，以及发生错误时的错误信息。

举一个例子：

```go
package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	reader := strings.NewReader("Clear is better than clever")
	p := make([]byte, 4)

	for {
		n, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF:", n)
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(n, string(p[:n]))
	}
}
```

输出：

```
4 Clea
4 r is
4  bet
4 ter
4 than
4  cle
3 ver
EOF: 0
```

这段代码从 `reader` 不断读取数据，每次读 4 个字节，然后打印输出，直到结尾。

最后一次返回的 n 值有可能小于缓冲区大小。

#### Writer 接口

`io.Writer` 表示一个编写器，它从缓冲区读取数据，并将数据写入目标资源。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/io-writer.drawio.png)

```go
type Writer interface {
   Write(p []byte) (n int, err error)
}
```

`Write` 方法将 `len(p)` 个字节从 `p` 中写入到对象数据流中。它返回从 `p` 中被写入的字节数 `n`，以及发生错误时返回的错误信息。

举一个例子：

```go
package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	// 创建 Buffer 暂存空间，并将一个字符串写入 Buffer
	// 使用 io.Writer 的 Write 方法写入
	var buf bytes.Buffer
	buf.Write([]byte("hello world , "))

	// 用 Fprintf 将一个字符串拼接到 Buffer 里
	fmt.Fprintf(&buf, " welcome to golang !")

	// 将 Buffer 的内容输出到标准输出设备
	buf.WriteTo(os.Stdout)
}
```

输出：

```
hello world ,  welcome to golang !
```

`bytes.Buffer` 是一个结构体类型，用来暂存写入的数据，其实现了 `io.Writer` 接口的 `Write` 方法。

`WriteTo` 方法定义：

```go
func (b *Buffer) WriteTo(w io.Writer) (n int64, err error)
```

`WriteTo` 方法第一个参数是 `io.Writer` 接口类型。

### 转换原理

再说回文章开头的转换问题。

只要某个实例实现了接口 `io.Reader` 里的方法 `Read()` ，就满足了接口 `io.Reader`。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/io-bytes-strings.png)

`bytes` 和 `strings` 包都实现了 `Read()` 方法。

```go
// src/bytes/reader.go

// NewReader returns a new Reader reading from b.
func NewReader(b []byte) *Reader { return &Reader{b, 0, -1} }
```

```go
// src/strings/reader.go

// NewReader returns a new Reader reading from s.
// It is similar to bytes.NewBufferString but more efficient and read-only.
func NewReader(s string) *Reader { return &Reader{s, 0, -1} }
```

在调用 `NewReader` 的时候，会返回了对应的 `T.Reader` 类型，而它们都是通过 `io.Reader` 扩展而来的，所以也就实现了转换。

### 总结

在开发过程中，避免不了要进行一些 IO 操作，包括打印输出，文件读写，网络连接等。

在 Go 语言中，也提供了一系列标准库来应对这些操作，主要封装在以下几个包中：

- `io`：基本的 IO 操作接口。
- `io/ioutil`：封装了一些实用的 IO 函数。
- `fmt`：实现了 IO 格式化操作。
- `bufio`：实现了带缓冲的 IO。 
- `net.Conn`：网络读写。
- `os.Stdin`，`os.Stdout`：系统标准输入输出。
- `os.File`:系统文件操作。
- `bytes`：字节相关 IO 操作。

除了 `io.Reader` 和 `io.Writer` 之外，`io` 包还封装了很多其他基本接口，比如 `ReaderAt`，`WriterAt`，`ReaderFrom` 和 `WriterTo` 等，这里就不一一介绍了。这部分代码并不复杂，读起来很轻松，而且还能加深对接口的理解，推荐大家看看。

好了，本文就到这里吧。关注我，带你通过问题读 Go 源码。

---

**推荐阅读：**

- [开始读 Go 源码了](https://mp.weixin.qq.com/s/iPM-mPOepRuDqkBtcnG1ww)

**热情推荐：**

- [计算机经典书籍（含下载方式）](https://mp.weixin.qq.com/s?__biz=MzI3MjY1ODI2Ng==&mid=2247484320&idx=1&sn=4f9ef828917db8b9c23688902ca46477&chksm=eb2e7995dc59f0834030ad6bad95190a9e1f5b9d44da9e53922ef8c81919b8bc68fa0b9841fd&token=1764237540&lang=zh_CN#rd)
- **[技术博客](https://github.com/yongxinz/tech-blog)：** 硬核后端技术干货，内容包括 Python、Django、Docker、Go、Redis、ElasticSearch、Kafka、Linux 等。
- **[Go 程序员](https://github.com/yongxinz/gopher)：** Go 学习路线图，包括基础专栏，进阶专栏，源码阅读，实战开发，面试刷题，必读书单等一系列资源。
- **[面试题汇总](https://github.com/yongxinz/backend-interview)：** 包括 Python、Go、Redis、MySQL、Kafka、数据结构、算法、编程、网络等各种常考题。

**参考文章：**

- https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter01/01.1.html
- https://www.cnblogs.com/jiujuan/p/14005731.html
- https://segmentfault.com/a/1190000015591319