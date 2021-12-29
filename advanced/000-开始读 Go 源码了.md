**原文链接：** [开始读 Go 源码了](https://mp.weixin.qq.com/s/iPM-mPOepRuDqkBtcnG1ww)

学完 Go 的基础知识已经有一段时间了，那么接下来应该学什么呢？有几个方向可以考虑，比如说 Web 开发，网络编程等。

在下一阶段的学习之前，写了一个[开源项目｜Go 开发的一款分布式唯一 ID 生成系统](https://mp.weixin.qq.com/s/tCGYTlB4nJH1ClViFQJ6Cw)，如果你对这个项目感兴趣的话，可以在 [GitHub](https://github.com/yongxinz/id-maker) 上拿到源码。

在写项目的过程中，发现一个问题。实现功能是没问题的，但不知道自己写的代码是不是符合 Go 的风格，是不是够优雅。所以我觉得相比于继续学习应用开发，不如向底层前进，打好基础，打好写 Go 代码的基础。

所以，我决定开始读 Go 标准库源码，Go 一共有 150+ 标准库，想要全部读完的话不是不可能，但绝对是一项大工程，希望自己能坚持下去。

为什么从 Go 标准库的源码开始读呢？因为最近也看了一些 Go 底层原理的书，说实话，像 goroutine 调度，gc 垃圾回收这些内容，根本就看不懂。这要是一上来就读这部分代码，恐怕直接就放弃 Go 语言学习了。

而标准库就不一样了，有一部分代码根本不涉及底层原理，实现也相对简单，同时又能对 Go 的理念加深理解，作为入门再好不过了。然后再由简入深，循序渐进，就像打怪升级一样，一步一步征服 Go。

说了这么多，那到底应该怎么读呢？我想到了一些方法：

- 看官方标准库文档。
- 看网上其他人的技术文章。
- 写一些例子来练习如何使用。
- 如果可以的话，自己实现标准库的功能。
- 将自己的阅读心得总结输出。

可以通过上面的一种或几种方法相结合，然后再不断阅读不断总结，最终找到一个完全适合自己的方法。

下面是我总结的一些标准库及功能介绍：

- `archive/tar` 和 `/zip-compress`：压缩（解压缩）文件功能。
- `fmt`-`io`-`bufio`-`path/filepath`-`flag`：
  - `fmt`：提供格式化输入输出功能。 
  - `io`：提供基本输入输出功能，大多数是围绕系统功能的封装。
  - `bufio`：缓冲输入输出功能的封装。
  - `path/filepath`：用来操作在当前系统中的目标文件名路径。
  - `flag`：提供对命令行参数的操作。
- `strings`-`strconv`-`unicode`-`regexp`-`bytes`：
  - `strings`：提供对字符串的操作。
  - `strconv`：提供将字符串转换为基础类型的功能。
  - `unicode`：为 unicode 型的字符串提供特殊的功能。
  - `regexp`：正则表达式功能。
  - `bytes`：提供对字符型分片的操作。
  - `index/suffixarray`：子字符串快速查询。
- `math`-`math/cmath`-`math/big`-`math/rand-sort`：
  - `math`：基本的数学函数。
  - `math/cmath`：对复数的操作。
  - `math/rand`：伪随机数生成。
  - `sort`：为数组排序和自定义集合。
  - `math/big`：大数的实现和计算。
- `container`-`/list`-`/ring`-`/heap`：
  - `list`：双链表。
  - `ring`：环形链表。
  - `heap`：堆。
- `compress/bzip2`-`/flate`-`/gzip`-`/lzw`-`zlib`：
  - `compress/bzip2`：实现 bzip2 的解压。
  - `flate`：实现 deflate 的数据压缩格式，如 RFC 1951 所述。
  - `gzip`：实现 gzip 压缩文件的读写。
  - `lzw`：Lempel Ziv Welch 压缩数据格式实现。
  - `zlib`：实现 zlib 数据压缩格式的读写。
- `context`：用来简化对于处理单个请求的多个 goroutine 之间与请求域的数据、取消信号、截止时间等相关操作。
- `crypto`-`crypto/md5`-`crypto/sha1`：
  - `crypto`：常用密码常数的集合。
  - `crypto/md5`：MD5 加密。
  - `crypto/sha1`：SHA1 加密。
- `errors`：实现操作出错的方法。
- `expvar`：为公共变量提供标准化的接口。
- `hash`：所有散列函数实现的通用接口。
- `html`：HTML 文本转码转义功能。
- `sort`：提供用于对切片和用户定义的集合进行排序的原始函数。
- `unsafe`：包含了一些打破 Go 语言「类型安全」的命令，一般程序不会使用，可用在 C/C++ 程序的调用中。
- `syscall`-`os`-`os/exec`：
  - `syscall`：提供了操作系统底层调用的基本接口。
  - `os`：提供给我们一个平台无关性的操作系统功能接口，采用类 Unix 设计，隐藏了不同操作系统间差异，让不同的文件系统和操作系统对象表现一致。
  - `os/exec`：提供了运行外部操作系统命令和程序的方式。
- `time`-`log`：
  - `time`：日期和时间的基本操作。
  - `log`：记录程序运行时产生的日志。
- `encoding/json`-`encoding/xml`-`text/template`：
  - `encoding/json`：读取并解码和写入并编码 JSON 数据。
  - `encoding/xml`：简单的 XML1.0 解析器。
  - `text/template`：生成像 HTML 一样的数据与文本混合的数据驱动模板。
- `net`-`net/http`：
  - `net`：网络数据的基本操作。
  - `http`：提供了一个可扩展的 HTTP 服务器和基础客户端，解析 HTTP 请求和回复。
- `runtime`：Go 程序运行时的交互操作，例如垃圾回收和协程创建。
- `reflect`：实现通过程序运行时反射，让程序操作任意类型的变量。

这里仅仅列举了一部分标准库，更全面的标准库列表大家可以直接看官网。

那么问题来了，这么多库从何下手呢？

我这里做一个简单的分类，由于水平有限，只能做一些简单的梳理，然后大家可以结合自己的实际情况来做选择。

有些库涉及到非常专业的知识，投入产出比可能会比较低。比如 `archive`、`compress` 以及 `crypto`，涉及到压缩算法以及加密算法的知识。

有些库属于工具类，比如 `bufio`、`bytes`、`strings`、`path`、`strconv` 等，这些库不涉及领域知识，阅读起来比较容易。

有些库属于与操作系统打交道的，比如 `os`，`net`、`sync` 等，学习这些库需要对操作系统有明确的认识。

`net` 下的很多子包与网络协议相关，比如 `net/http`，涉及 `http` 报文的解析，需要对网络协议比较了解。

如果想要深入了解语言的底层原理，则需要阅读 `runtime` 库。

要想快速入门，并且了解语言的设计理念，建议阅读 `io` 以及 `fmt` 库，阅读后会对接口的设计理解更深。

我已经看了一些源码，虽然过程痛苦，但确实非常有用。前期可能理解起来比较困难，用的时间长一些，但形成固定套路之后，会越来越熟悉，用的时间也会更少，理解也会更深刻。

后续我还会继续总结输出，请大家持续关注，让我们学起来。

---

**开源项目：**

- [https://github.com/yongxinz/id-maker](https://github.com/yongxinz/id-maker)

**往期文章：**

- [开源项目｜Go 开发的一款分布式唯一 ID 生成系统](https://mp.weixin.qq.com/s/tCGYTlB4nJH1ClViFQJ6Cw)
- [测试小姐姐问我 gRPC 怎么用，我直接把这篇文章甩给了她](https://mp.weixin.qq.com/s/qdI2JqpMq6t2KN1byHaNCQ)
- [gRPC，爆赞](https://mp.weixin.qq.com/s/1Xbca4Dv0akonAZerrChgA)