**原文链接：** [Go 语言 map 是并发安全的吗？](https://mp.weixin.qq.com/s/4mDzMdMbunR_p94Du65QOA)

Go 语言中的 map 是一个非常常用的数据结构，它允许我们快速地存储和检索键值对。然而，在并发场景下使用 map 时，还是有一些问题需要注意的。

本文将探讨 Go 语言中的 map 是否是并发安全的，并提供三种方案来解决并发问题。

先来回答一下题目的问题，答案就是**并发不安全**。

看一段代码示例，当两个 goroutine 同时对同一个 map 进行写操作时，会发生什么？

```go
package main

import "sync"

func main() {
    m := make(map[string]int)
    m["foo"] = 1

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        for i := 0; i < 1000; i++ {
            m["foo"]++
        }
        wg.Done()
    }()

    go func() {
        for i := 0; i < 1000; i++ {
            m["foo"]++
        }
        wg.Done()
    }()

    wg.Wait()
}
```

在这个例子中，我们可以看到，两个 goroutine 将尝试同时对 map 进行写入。运行这个程序时，我们将看到一个错误：

```go
fatal error: concurrent map writes
```

也就是说，在并发场景下，这样操作 map 是不行的。

## 为什么是不安全的

因为它**没有内置的锁机制**来保护多个 goroutine 同时对其进行读写操作。

当多个 goroutine 同时对同一个 map 进行读写操作时，就会出现数据竞争和不一致的结果。

就像上例那样，当两个 goroutine 同时尝试更新同一个键值对时，最终的结果可能取决于哪个 goroutine 先完成了更新操作。这种不确定性可能会导致程序出现错误或崩溃。

Go 语言团队没有将 map 设计成并发安全的，是因为这样会增加程序的开销并降低性能。

如果 map 内置了锁机制，那么每次访问 map 时都需要进行加锁和解锁操作，这会增加程序的运行时间并降低性能。

此外，并不是所有的程序都需要在并发场景下使用 map，因此将锁机制内置到 map 中会对那些不需要并发安全的程序造成不必要的开销。

在实际使用过程中，开发人员可以根据程序的需求来选择是否需要保证 map 的并发安全性，从而在性能和安全性之间做出权衡。

## 如何并发安全

接下来介绍三种并发安全的方式：

1.  读写锁
2.  分片加锁
3.  sync.Map

### 加读写锁

第一种方法是使用**读写锁**，这是最容易想到的一种方式。在读操作时加读锁，在写操作时加写锁。

```go
package main

import (
    "fmt"
    "sync"
)

type SafeMap struct {
    sync.RWMutex
    Map map[string]string
}

func NewSafeMap() *SafeMap {
    sm := new(SafeMap)
    sm.Map = make(map[string]string)
    return sm
}

func (sm *SafeMap) ReadMap(key string) string {
    sm.RLock()
    value := sm.Map[key]
    sm.RUnlock()
    return value
}

func (sm *SafeMap) WriteMap(key string, value string) {
    sm.Lock()
    sm.Map[key] = value
    sm.Unlock()
}

func main() {
    safeMap := NewSafeMap()

    var wg sync.WaitGroup

    // 启动多个goroutine进行写操作
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            safeMap.WriteMap(fmt.Sprintf("name%d", i), fmt.Sprintf("John%d", i))
        }(i)
    }

    wg.Wait()

    // 启动多个goroutine进行读操作
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            fmt.Println(safeMap.ReadMap(fmt.Sprintf("name%d", i)))
        }(i)
    }

    wg.Wait()
}
```

在这个示例中，我们定义了一个 `SafeMap` 结构体，它包含一个 `sync.RWMutex` 和一个 `map[string]string`。

定义了两个方法：`ReadMap` 和 `WriteMap`。在 `ReadMap` 方法中，我们使用读锁来保护对 map 的读取操作。在 `WriteMap` 方法中，我们使用写锁来保护对 map 的写入操作。

在 `main` 函数中，我们启动了多个 goroutine 来进行读写操作，这些操作都是安全的。

### 分片加锁

上例中通过对整个 map 加锁来实现需求，但相对来说，锁会大大降低程序的性能，那如何优化呢？其中一个优化思路就是降低锁的粒度，不对整个 map 进行加锁。

这种方法是**分片加锁**，将这个 map 分成 n 块，每个块之间的读写操作都互不干扰，从而降低冲突的可能性。

```go
package main

import (
    "fmt"
    "sync"
)

const N = 16

type SafeMap struct {
    maps  [N]map[string]string
    locks [N]sync.RWMutex
}

func NewSafeMap() *SafeMap {
    sm := new(SafeMap)
    for i := 0; i < N; i++ {
        sm.maps[i] = make(map[string]string)
    }
    return sm
}

func (sm *SafeMap) ReadMap(key string) string {
    index := hash(key) % N
    sm.locks[index].RLock()
    value := sm.maps[index][key]
    sm.locks[index].RUnlock()
    return value
}

func (sm *SafeMap) WriteMap(key string, value string) {
    index := hash(key) % N
    sm.locks[index].Lock()
    sm.maps[index][key] = value
    sm.locks[index].Unlock()
}

func hash(s string) int {
    h := 0
    for i := 0; i < len(s); i++ {
        h = 31*h + int(s[i])
    }
    return h
}

func main() {
    safeMap := NewSafeMap()

    var wg sync.WaitGroup

    // 启动多个goroutine进行写操作
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            safeMap.WriteMap(fmt.Sprintf("name%d", i), fmt.Sprintf("John%d", i))
        }(i)
    }

    wg.Wait()

    // 启动多个goroutine进行读操作
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            fmt.Println(safeMap.ReadMap(fmt.Sprintf("name%d", i)))
        }(i)
    }

    wg.Wait()
}
```

在这个示例中，我们定义了一个 `SafeMap` 结构体，它包含一个长度为 `N` 的 map 数组和一个长度为 `N` 的锁数组。

定义了两个方法：`ReadMap` 和 `WriteMap`。在这两个方法中，我们都使用了一个 `hash` 函数来计算 `key` 应该存储在哪个 map 中。然后再对这个 map 进行读写操作。

在 `main` 函数中，我们启动了多个 goroutine 来进行读写操作，这些操作都是安全的。

有一个开源项目 [orcaman/concurrent-map](https://github.com/orcaman/concurrent-map) 就是通过这种思想来做的，感兴趣的同学可以看看。

### sync.Map

最后，在内置的 **sync 包**中（Go 1.9+）也有一个线程安全的 map，通过将读写分离的方式实现了某些特定场景下的性能提升。

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var m sync.Map
    var wg sync.WaitGroup

    // 启动多个goroutine进行写操作
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            m.Store(fmt.Sprintf("name%d", i), fmt.Sprintf("John%d", i))
        }(i)
    }

    wg.Wait()

    // 启动多个goroutine进行读操作
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            v, _ := m.Load(fmt.Sprintf("name%d", i))
            fmt.Println(v.(string))
        }(i)
    }

    wg.Wait()
}
```

有了官方的支持，代码瞬间少了很多，使用起来方便多了。

在这个示例中，我们使用了内置的 `sync.Map` 类型来存储键值对，使用 `Store` 方法来存储键值对，使用 `Load` 方法来获取键值对。

在 `main` 函数中，我们启动了多个 goroutine 来进行读写操作，这些操作都是安全的。

## 总结

Go 语言中的 map 本身并不是并发安全的。

在多个 goroutine 同时访问同一个 map 时，可能会出现并发不安全的现象。这是因为 Go 语言中的 map 并没有内置锁来保护对map的访问。

尽管如此，我们仍然可以使用一些方法来实现 map 的并发安全。

一种方法是使用读写锁，在读操作时加读锁，在写操作时加写锁。

另一种方法是分片加锁，将这个 map 分成 n 块，每个块之间的读写操作都互不干扰，从而降低冲突的可能性。

此外，在内置的 sync 包中（Go 1.9+）也有一个线程安全的 map，它通过将读写分离的方式实现了某些特定场景下的性能提升。

以上就是本文的全部内容，如果觉得还不错的话欢迎**点赞**，**转发**和**关注**，感谢支持。

***

**参考文章：**

*   <https://zhuanlan.zhihu.com/p/356739568>

**推荐阅读：**

*   [Go 语言切片是如何扩容的？](https://mp.weixin.qq.com/s/VVM8nqs4mMGdFyCNJx16_g)
*   [Go 语言数组和切片的区别](https://mp.weixin.qq.com/s/esaAmAdmV4w3_qjtAzTr4A)
*   [Go 语言 new 和 make 关键字的区别](https://mp.weixin.qq.com/s/NBDkI3roHgNgW1iW4e_6cA)
*   [为什么 Go 不支持 \[\]T 转换为 \[\]interface](https://mp.weixin.qq.com/s/cwDEgnicK4jkuNpzulU2bw)
*   [为什么 Go 语言 struct 要使用 tags](https://mp.weixin.qq.com/s/L7-TJ-CzYfuVrIBWP7Ebaw)

