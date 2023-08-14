**原文链接：** [如何实现计数器限流？](https://mp.weixin.qq.com/s/CTemkZ2aKPCPTuQiDJri0Q)

上一篇文章 [go-zero 是如何做路由管理的？](https://mp.weixin.qq.com/s/uTJ1En-BXiLvH45xx0eFsA) 介绍了路由管理，这篇文章来说说限流，主要介绍计数器限流算法，具体的代码实现，我们还是来分析微服务框架 go-zero 的源码。

在微服务架构中，一个服务可能需要频繁地与其他服务交互，而过多的请求可能导致性能下降或系统崩溃。为了确保系统的稳定性和高可用性，限流算法应运而生。

限流算法允许在给定时间段内，对服务的请求流量进行控制和调整，以防止资源耗尽和服务过载。

计数器限流算法主要有两种实现方式，分别是：

1. 固定窗口计数器
2. 滑动窗口计数器

下面分别来介绍。

## 固定窗口计数器

算法概念如下：

- 将时间划分为多个窗口；
- 在每个窗口内每有一次请求就将计数器加一；
- 如果计数器超过了限制数量，则本窗口内所有的请求都被丢弃当时间到达下一个窗口时，计数器重置。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/periodlimit1.png)

固定窗口计数器是最为简单的算法，但这个算法有时会让通过请求量允许为限制的两倍。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/periodlimit2.png)

考虑如下情况：限制 1 秒内最多通过 5 个请求，在第一个窗口的最后半秒内通过了 5 个请求，第二个窗口的前半秒内又通过了 5 个请求。这样看来就是在 1 秒内通过了 10 个请求。

## 滑动窗口计数器

算法概念如下：

- 将时间划分为多个区间；
- 在每个区间内每有一次请求就将计数器加一维持一个时间窗口，占据多个区间；
- 每经过一个区间的时间，则抛弃最老的一个区间，并纳入最新的一个区间；
- 如果当前窗口内区间的请求计数总和超过了限制数量，则本窗口内所有的请求都被丢弃。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/periodlimit3.png)

滑动窗口计数器是通过将窗口再细分，并且按照时间滑动，这种算法避免了固定窗口计数器带来的双倍突发请求，但时间区间的精度越高，算法所需的空间容量就越大。

## go-zero 实现

go-zero 实现的是**固定窗口**的方式，计算一段时间内对同一个资源的访问次数，如果超过指定的 `limit`，则拒绝访问。当然如果在一段时间内访问不同的资源，每一个资源访问量都不超过 `limit`，此种情况是不会拒绝的。

而在一个分布式系统中，存在多个微服务提供服务。所以当瞬间的流量同时访问同一个资源，如何让计数器在分布式系统中正常计数？

这里要解决的一个主要问题就是计算的原子性，保证多个计算都能得到正确结果。

通过以下两个方面来解决：

- 使用 redis 的 `incrby` 做资源访问计数
- 采用 lua script 做整个窗口计算，保证计算的原子性

接下来先看一下 lua script 的源码：

```go
// core/limit/periodlimit.go

const periodScript = `local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCRBY", KEYS[1], 1)
if current == 1 then
    redis.call("expire", KEYS[1], window)
end
if current < limit then
    return 1
elseif current == limit then
    return 2
else
    return 0
end`
```

主要就是使用 `INCRBY` 命令来实现，第一次请求需要给 key 加上一个过期时间，到达过期时间之后，key 过期被清楚，重新计数。

限流器初始化：

```go
type (
    // PeriodOption defines the method to customize a PeriodLimit.
    PeriodOption func(l *PeriodLimit)

    // A PeriodLimit is used to limit requests during a period of time.
    PeriodLimit struct {
        period     int  // 窗口大小，单位 s
        quota      int  // 请求上限
        limitStore *redis.Redis
        keyPrefix  string   // key 前缀
        align      bool
    }
)

// NewPeriodLimit returns a PeriodLimit with given parameters.
func NewPeriodLimit(period, quota int, limitStore *redis.Redis, keyPrefix string,
    opts ...PeriodOption) *PeriodLimit {
    limiter := &PeriodLimit{
        period:     period,
        quota:      quota,
        limitStore: limitStore,
        keyPrefix:  keyPrefix,
    }

    for _, opt := range opts {
        opt(limiter)
    }

    return limiter
}
```

调用限流：

```go
// key 就是需要被限制的资源标识
func (h *PeriodLimit) Take(key string) (int, error) {
    return h.TakeCtx(context.Background(), key)
}

// TakeCtx requests a permit with context, it returns the permit state.
func (h *PeriodLimit) TakeCtx(ctx context.Context, key string) (int, error) {
    resp, err := h.limitStore.EvalCtx(ctx, periodScript, []string{h.keyPrefix + key}, []string{
        strconv.Itoa(h.quota),
        strconv.Itoa(h.calcExpireSeconds()),
    })
    if err != nil {
        return Unknown, err
    }

    code, ok := resp.(int64)
    if !ok {
        return Unknown, ErrUnknownCode
    }

    switch code {
    case internalOverQuota: // 超过上限
        return OverQuota, nil
    case internalAllowed:   // 未超过，允许访问
        return Allowed, nil
    case internalHitQuota:  // 正好达到限流上限
        return HitQuota, nil
    default:
        return Unknown, ErrUnknownCode
    }
}
```

上文已经介绍了，固定时间窗口会有临界突发问题，并不是那么严谨，下篇文章我们来介绍令牌桶限流。

以上就是本文的全部内容，如果觉得还不错的话欢迎**点赞**，**转发**和**关注**，感谢支持。

***

**参考文章：**

- https://juejin.cn/post/6895928148521648141
- https://juejin.cn/post/7051406419823689765
- https://www.infoq.cn/article/Qg2tX8fyw5Vt-f3HH673

**推荐阅读：**

- [go-zero 是如何做路由管理的？](https://mp.weixin.qq.com/s/uTJ1En-BXiLvH45xx0eFsA)
