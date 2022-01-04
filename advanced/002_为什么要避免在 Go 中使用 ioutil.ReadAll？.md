**原文链接：** [为什么要避免在 Go 中使用 ioutil.ReadAll？](https://mp.weixin.qq.com/s/e2A3ME4vhOK2S3hLEJtPsw)

`ioutil.ReadAll` 主要的作用是从一个 `io.Reader` 中读取所有数据，直到结尾。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/002-ioutil-readall-github.png)

在 GitHub 上搜索 `ioutil.ReadAll`，类型选择 Code，语言选择 Go，一共得到了 637307 条结果。

这说明 `ioutil.ReadAll` 还是挺受欢迎的，主要也是用起来确实方便。

但是当遇到大文件时，这个函数就会暴露出两个明显的缺点：

1. 性能问题，文件越大，性能越差。
2. 文件过大的话，可能直接撑爆内存，导致程序崩溃。

为什么会这样呢？这篇文章就通过源码来分析背后的原因，并试图给出更好的解决方案。

下面我们正式开始。

### ioutil.ReadAll

首先，我们通过一个例子看一下 `ioutil.ReadAll` 的使用场景。比如说，使用 `http.Client` 发送 `GET` 请求，然后再读取返回内容：

```go
func main() {
	res, err := http.Get("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	
	robots, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", robots)
}
```

`http.Get()` 返回的数据，存储在 `res.Body` 中，通过 `ioutil.ReadAll` 将其读取出来。

表面上看这段代码没有什么问题，但仔细分析却并非如此。想要探究其背后的原因，就只能靠源码说话。

`ioutil.ReadAll` 的源码如下：

```go
// src/io/ioutil/ioutil.go

func ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
```

Go 1.16 版本开始，直接调用 `io.ReadAll()` 函数，下面再看看 `io.ReadAll()` 的实现：

```go
// src/io/io.go

func ReadAll(r Reader) ([]byte, error) {
    // 创建一个 512 字节的 buf
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// 如果 buf 满了，则追加一个元素，使其重新分配内存
			b = append(b, 0)[:len(b)]
		}
		// 读取内容到 buf
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		// 遇到结尾或者报错则返回
		if err != nil {
			if err == EOF {
				err = nil
			}
			return b, err
		}
	}
}
```

我给代码加上了必要的注释，这段代码的执行主要分三个步骤：

1. 创建一个 512 字节的 `buf`；
2. 不断读取内容到 `buf`，当 `buf` 满的时候，会追加一个元素，促使其重新分配内存；
3. 直到结尾或报错，则返回；

知道了执行步骤，但想要分析其性能问题，还需要了解 Go 切片的扩容策略，如下：

1. 如果期望容量大于当前容量的两倍就会使用期望容量；
2. 如果当前切片的长度小于 1024 就会将容量翻倍；
3. 如果当前切片的长度大于 1024 就会每次增加 25% 的容量，直到新容量大于期望容量；

也就是说，如果待拷贝数据的容量小于 512 字节的话，性能不受影响。但如果超过 512 字节，就会开始切片扩容。数据量越大，扩容越频繁，性能受影响越大。

如果数据量足够大的话，内存可能就直接撑爆了，这样的话影响就大了。

那有更好的替换方案吗？当然是有的，我们接着往下看。

### io.Copy

可以使用 `io.Copy` 函数来代替，源码定义如下：

```go
src/io/io.go

func Copy(dst Writer, src Reader) (written int64, err error) {
	return copyBuffer(dst, src, nil)
}
```

其功能是直接从 `src` 读取数据，并写入到 `dst`。

和 `ioutil.ReadAll` 最大的不同就是没有把所有数据一次性都取出来，而是不断读取，不断写入。

具体实现 `Copy` 的逻辑在 `copyBuffer` 函数中实现：

```go
// src/io/io.go

func copyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
	// 如果源实现了 WriteTo 方法，则直接调用 WriteTo
	if wt, ok := src.(WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// 同样的，如果目标实现了 ReaderFrom 方法，则直接调用 ReaderFrom
	if rt, ok := dst.(ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	// 如果 buf 为空，则创建 32KB 的 buf
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	// 循环读取数据并写入
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
```

此函数执行步骤如下：

1. 如果源实现了 `WriteTo` 方法，则直接调用 `WriteTo` 方法；
2. 同样的，如果目标实现了 ReaderFrom 方法，则直接调用 ReaderFrom 方法；
3. 如果 `buf` 为空，则创建 32KB 的 `buf`；
4. 最后就是循环 `Read` 和 `Write`；

对比之后就会发现，`io.Copy` 函数不会一次性读取全部数据，也不会频繁进行切片扩容，显然在数据量大时是更好的选择。

### ioutil 其他函数

再看看 `ioutil` 包的其他函数：

- `func ReadDir(dirname string) ([]os.FileInfo, error)`
- `func ReadFile(filename string) ([]byte, error)`
- `func WriteFile(filename string, data []byte, perm os.FileMode) error`
- `func TempFile(dir, prefix string) (f *os.File, err error)`
- `func TempDir(dir, prefix string) (name string, err error)`
- `func NopCloser(r io.Reader) io.ReadCloser`

下面举例详细说明：

#### ReadDir

```go
// ReadDir 读取指定目录中的所有目录和文件（不包括子目录）。
// 返回读取到的文件信息列表和遇到的错误，列表是经过排序的。
func ReadDir(dirname string) ([]os.FileInfo, error)
```

**举例：**

```go
package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	dirName := "../"
	fileInfos, _ := ioutil.ReadDir(dirName)
	fmt.Println(len(fileInfos))
	for i := 0; i < len(fileInfos); i++ {
		fmt.Printf("%T\n", fileInfos[i])
		fmt.Println(i, fileInfos[i].Name(), fileInfos[i].IsDir())

	}
}
```

#### ReadFile

```go
// ReadFile 读取文件中的所有数据，返回读取的数据和遇到的错误
// 如果读取成功，则 err 返回 nil，而不是 EOF
func ReadFile(filename string) ([]byte, error)
```

**举例：**

```go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	data, err := ioutil.ReadFile("./test.txt")
	if err != nil {
		fmt.Println("read error")
		os.Exit(1)
	}
	fmt.Println(string(data))
}
```

#### WriteFile

```go
// WriteFile 向文件中写入数据，写入前会清空文件。
// 如果文件不存在，则会以指定的权限创建该文件。
// 返回遇到的错误。
func WriteFile(filename string, data []byte, perm os.FileMode) error
```

**举例：**

```go
package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	fileName := "./text.txt"
	s := "Hello AlwaysBeta"
	err := ioutil.WriteFile(fileName, []byte(s), 0777)
	fmt.Println(err)
}
```

#### TempFile

```go
// TempFile 在 dir 目录中创建一个以 prefix 为前缀的临时文件，并将其以读
// 写模式打开。返回创建的文件对象和遇到的错误。
// 如果 dir 为空，则在默认的临时目录中创建文件（参见 os.TempDir），多次
// 调用会创建不同的临时文件，调用者可以通过 f.Name() 获取文件的完整路径。
// 调用本函数所创建的临时文件，应该由调用者自己删除。
func TempFile(dir, prefix string) (f *os.File, err error)
```

**举例：**

```go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	f, err := ioutil.TempFile("./", "Test")
	if err != nil {
		fmt.Println(err)
	}
	defer os.Remove(f.Name()) // 用完删除
	fmt.Printf("%s\n", f.Name())
}
```

#### TempDir

```go
// TempDir 功能同 TempFile，只不过创建的是目录，返回目录的完整路径。
func TempDir(dir, prefix string) (name string, err error)
```

**举例：**

```go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	dir, err := ioutil.TempDir("./", "Test")
	if err != nil {
		fmt.Println(err)
	}
	defer os.Remove(dir) // 用完删除
	fmt.Printf("%s\n", dir)
}
```

#### NopCloser

```go
// NopCloser 将 r 包装为一个 ReadCloser 类型，但 Close 方法不做任何事情。
func NopCloser(r io.Reader) io.ReadCloser
```

这个函数的使用场景是这样的：

有时候我们需要传递一个 `io.ReadCloser` 的实例，而我们现在有一个 `io.Reader` 的实例，比如：`strings.Reader`。

这个时候 `NopCloser` 就派上用场了。它包装一个 `io.Reader`，返回一个 `io.ReadCloser`，相应的 `Close` 方法啥也不做，只是返回 `nil`。

**举例：**

```go
package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

func main() {
	//返回 *strings.Reader
	reader := strings.NewReader("Hello AlwaysBeta")
	r := ioutil.NopCloser(reader)
	defer r.Close()

	fmt.Println(reflect.TypeOf(reader))
	data, _ := ioutil.ReadAll(reader)
	fmt.Println(string(data))
}
```

### 总结

`ioutil` 提供了几个很实用的工具函数，背后实现逻辑也并不复杂。

本篇文章从一个问题入手，重点研究了 `ioutil.ReadAll` 函数。主要原因是在小数据量的情况下，这个函数并没有什么问题，但当数据量大时，它就变成了一颗定时炸弹。有可能会影响程序的性能，甚至会导致程序崩溃。

接下来给出对应的解决方案，在数据量大的情况下，最好使用 `io.Copy` 函数。

文章最后继续介绍了 `ioutil` 的其他几个函数，并给出了程序示例。相关代码都会上传到 [GitHub](https://github.com/yongxinz/gopher/tree/main/advanced/src)，需要的同学可以自行下载。

好了，本文就到这里吧。关注我，带你通过问题读 Go 源码。

---

**源码地址：**

- [https://github.com/yongxinz/gopher](https://github.com/yongxinz/gopher)

**推荐阅读：**

- [如何在 Go 中将 []byte 转换为 io.Reader？](https://mp.weixin.qq.com/s/nFkob92GOs6Gp75pxA5wCQ)
- [开始读 Go 源码了](https://mp.weixin.qq.com/s/iPM-mPOepRuDqkBtcnG1ww)

**参考文章：**

- https://haisum.github.io/2017/09/11/golang-ioutil-readall/
- https://juejin.cn/post/6977640348679929886
- https://zhuanlan.zhihu.com/p/76231663