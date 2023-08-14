**原文链接：** [go-zero 是如何做路由管理的？](https://mp.weixin.qq.com/s/uTJ1En-BXiLvH45xx0eFsA)

go-zero 是一个微服务框架，包含了 web 和 rpc 两大部分。

而对于 web 框架来说，路由管理是必不可少的一部分，那么本文就来探讨一下 go-zero 的路由管理是怎么做的，具体采用了哪种技术方案。

## 路由管理方案

路由管理方案有很多种，具体应该如何选择，应该根据使用场景，以及实现的难易程度做综合分析，下面介绍常见的三种方案。

注意这里只是做一个简单的概括性对比，更加详细的内容可以看这篇文章：[HTTP Router 算法演进](https://mp.weixin.qq.com/s/Ec2KyQ1ObyJuAOSvFa5xXg)。

### 标准库方案

最简单的方案就是直接使用 `map[string]func()` 作为路由的数据结构，键为具体的路由，值为具体的处理方法。

```go
// 路由管理数据结构

type ServeMux struct {
    mu    sync.RWMutex          // 对象操作读写锁
    m     map[string]muxEntry   // 存储路由映射关系
}
```

这种方案优点就是实现简单，性能较高；缺点也很明显，占用内存更高，更重要的是不够灵活。

### Trie Tree

Trie Tree 也称为字典树或前缀树，是一种用于高效存储和检索、用于从某个集合中查到某个特定 key 的数据结构。 

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/trie1.jpeg)

Trie Tree 时间复杂度低，和一般的树形数据结构相比，Trie Tree 拥有更快的前缀搜索和查询性能。

和查询时间复杂度为 `O(1)` 常数的哈希算法相比，Trie Tree 支持前缀搜索，并且可以节省哈希函数的计算开销和避免哈希值碰撞的情况。

最后，Trie Tree 还支持对关键字进行字典排序。

### Radix Tree

Radix Tree（基数树）是一种特殊的数据结构，用于高效地存储和搜索字符串键值对，它是一种基于前缀的树状结构，通过将相同前缀的键值对合并在一起来减少存储空间的使用。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/RadixTree.png)

Radix Tree 通过合并公共前缀来降低存储空间的开销，避免了 Trie Tree 字符串过长和字符集过大时导致的存储空间过多问题，同时公共前缀优化了路径层数，提升了插入、查询、删除等操作效率。

比如 Gin 框架使用的开源组件 HttpRouter 就是采用这个方案。

## go-zero 路由规则

在使用 go-zero 开发项目时，定义路由需要遵守如下规则：

1. 路由必须以 `/` 开头
2. 路由节点必须以 `/` 分隔
3. 路由节点中可以包含 `:`，但是 `:` 必须是路由节点的第一个字符，`:` 后面的节点值必须要在结请求体中有 `path tag` 声明，用于接收路由参数
4. 路由节点可以包含字母、数字、下划线、中划线

接下来就让我们深入到源码层面，相信看过源码之后，你就会更懂这些规则的意义了。

## go-zero 源码实现

首先需要说明的是，底层数据结构使用的是二叉搜索树，还不是很了解的同学可以看这篇文章：[使用 Go 语言实现二叉搜索树](https://mp.weixin.qq.com/s/2wYRmG_AiiHYjLDEXg94Ag)

### 节点定义

先看一下节点定义：

```go
// core/search/tree.go

const (
    colon = ':'
    slash = '/'
)

type (
    // 节点
    node struct {
        item     interface{}
        children [2]map[string]*node
    }

    // A Tree is a search tree.
    Tree struct {
        root *node
    }
)
```

重点说一下 `children`，它是一个包含两个元素的数组，元素 `0` 存正常路由键，元素 `1` 存以 `:` 开头的路由键，这些是 url 中的变量，到时候需要替换成实际值。

举一个例子，有这样一个路由 `/api/:user`，那么 `api` 会存在 `children[0]`，`user` 会存在 `children[1]`。

具体可以看看这段代码：

```go
func (nd *node) getChildren(route string) map[string]*node {
    // 判断路由是不是以 : 开头
    if len(route) > 0 && route[0] == colon {
        return nd.children[1]
    }

    return nd.children[0]
}
```

### 路由添加

```go
// Add adds item to associate with route.
func (t *Tree) Add(route string, item interface{}) error {
    // 需要路由以 / 开头
    if len(route) == 0 || route[0] != slash {
        return errNotFromRoot
    }

    if item == nil {
        return errEmptyItem
    }

    // 把去掉 / 的路由作为参数传入
    err := add(t.root, route[1:], item)
    switch err {
    case errDupItem:
        return duplicatedItem(route)
    case errDupSlash:
        return duplicatedSlash(route)
    default:
        return err
    }
}


func add(nd *node, route string, item interface{}) error {
    if len(route) == 0 {
        if nd.item != nil {
            return errDupItem
        }

        nd.item = item
        return nil
    }

    // 继续判断，看看是不是有多个 /
    if route[0] == slash {
        return errDupSlash
    }

    for i := range route {
        // 判断是不是 /，目的就是去处两个 / 之间的内容
        if route[i] != slash {
            continue
        }

        token := route[:i]
        
        // 看看有没有子节点，如果有子节点，就在子节点下面继续添加
        children := nd.getChildren(token)
        if child, ok := children[token]; ok {
            if child != nil {
                return add(child, route[i+1:], item)
            }

            return errInvalidState
        }

        // 没有子节点，那么新建一个
        child := newNode(nil)
        children[token] = child
        return add(child, route[i+1:], item)
    }

    children := nd.getChildren(route)
    if child, ok := children[route]; ok {
        if child.item != nil {
            return errDupItem
        }

        child.item = item
    } else {
        children[route] = newNode(item)
    }

    return nil
}
```

主要部分代码都已经加了注释，其实这个过程就是树的构建，如果读过之前那篇文章，那这里还是比较好理解的。

### 路由查找

先来看一段 `match` 代码：

```go
func match(pat, token string) innerResult {
    if pat[0] == colon {
        return innerResult{
            key:   pat[1:],
            value: token,
            named: true,
            found: true,
        }
    }

    return innerResult{
        found: pat == token,
    }
}
```

这里有两个参数：

- `pat`：路由树中存储的路由
- `token`：实际请求的路由，可能包含参数值

还是刚才的例子 `/api/:user`，如果是 `api`，没有以 `:` 开头，那就不会走 `if` 逻辑。

接下来匹配 `:user` 部分，如果实际请求的 url 是 `/api/zhangsan`，那么会将 `user` 作为 `key`，`zhangsan` 作为 `value` 保存到结果中。

下面是搜索查找代码：

```go
// Search searches item that associates with given route.
func (t *Tree) Search(route string) (Result, bool) {
    // 第一步先判断是不是 / 开头
    if len(route) == 0 || route[0] != slash {
        return NotFound, false
    }

    var result Result
    ok := t.next(t.root, route[1:], &result)
    return result, ok
}

func (t *Tree) next(n *node, route string, result *Result) bool {
    if len(route) == 0 && n.item != nil {
        result.Item = n.item
        return true
    }

    for i := range route {
        // 和 add 里同样的提取逻辑
        if route[i] != slash {
            continue
        }

        token := route[:i]
        return n.forEach(func(k string, v *node) bool {
            r := match(k, token)
            if !r.found || !t.next(v, route[i+1:], result) {
                return false
            }
            // 如果 url 中有参数，会把键值对保存到结果中
            if r.named {
                addParam(result, r.key, r.value)
            }

            return true
        })
    }

    return n.forEach(func(k string, v *node) bool {
        if r := match(k, route); r.found && v.item != nil {
            result.Item = v.item
            if r.named {
                addParam(result, r.key, r.value)
            }

            return true
        }

        return false
    })
}
```

以上就是路由管理的大部分代码，整个文件也就 200 多行，逻辑也并不复杂，通读之后还是很有收获的。

大家如果感兴趣的话，可以找到项目更详细地阅读。也可以关注我，接下来还会分析其他模块的源码。

以上就是本文的全部内容，如果觉得还不错的话欢迎**点赞**，**转发**和**关注**，感谢支持。

***

**推荐阅读：**

*   [使用 Go 语言实现二叉搜索树](https://mp.weixin.qq.com/s/2wYRmG_AiiHYjLDEXg94Ag)
*   [HTTP Router 算法演进](https://mp.weixin.qq.com/s/Ec2KyQ1ObyJuAOSvFa5xXg)