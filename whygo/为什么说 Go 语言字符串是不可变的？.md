**原文链接：** [为什么说 Go 语言字符串是不可变的？](https://mp.weixin.qq.com/s/AOb6AjKwyTwLeAUou0AU-Q)

最近有读者留言说，平时在写代码的过程中，是会对字符串进行修改的，但网上都说 Go 语言字符串是不可变的，这是为什么呢？

这个问题本身并不困难，但对于新手来说确实容易产生困惑，今天就来回答一下。

首先来看看它的底层结构：

```go
type stringStruct struct {
    str unsafe.Pointer
    len int
}
```

和切片的结构很像，只不过少了一个表示容量的 `cap` 字段。

*   `str`：指向一个 `[]byte` 类型的指针
*   `len`：字符串的长度

所以，当我们定义一个字符串：

```go
s := "Hello World"
```

那么它在内存中存储是这样的：

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/string.drawio.png)

当我们在程序中对字符串进行重新赋值时，比如这样：

```go
s := "Hello World"

s = "Hello AlwaysBeta"
```

底层的存储就变成了这样：

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/string.drawio%20\(1\).png)

Go 实际上是重新创建了一个 `[]byte{}` 切片，然后让指针指向了新的地址。

更直接一点，我们直接修改字符串中的单个字符，比如：

```go
s := "Hello World"
s[0] = 'h'
```

这样做的话，会直接报错：

```go
cannot assign to s[0] (strings are immutable)
```

如果一定要这么做的话，需要对字符串进行一个转换，转换成 `[]byte` 类型，修改之后再转换回 `string` 类型：

```go
s := "Hello World"
sBytes := []byte(s)
sBytes[0] = 'h'
s = string(sBytes)
```

这样就可以了。

以上就是本文的全部内容，如果觉得还不错的话欢迎**点赞**，**转发**和**关注**，感谢支持。

***

**推荐阅读：**

*   [Go 语言 map 如何顺序读取？](https://mp.weixin.qq.com/s/iScSgfpSE2y14GH7JNRJSA)
*   [Go 语言 map 是并发安全的吗？](https://mp.weixin.qq.com/s/4mDzMdMbunR_p94Du65QOA)
*   [Go 语言切片是如何扩容的？](https://mp.weixin.qq.com/s/VVM8nqs4mMGdFyCNJx16_g)
*   [Go 语言数组和切片的区别](https://mp.weixin.qq.com/s/esaAmAdmV4w3_qjtAzTr4A)
*   [Go 语言 new 和 make 关键字的区别](https://mp.weixin.qq.com/s/NBDkI3roHgNgW1iW4e_6cA)
*   [为什么 Go 不支持 \[\]T 转换为 \[\]interface](https://mp.weixin.qq.com/s/cwDEgnicK4jkuNpzulU2bw)
*   [为什么 Go 语言 struct 要使用 tags](https://mp.weixin.qq.com/s/L7-TJ-CzYfuVrIBWP7Ebaw)

