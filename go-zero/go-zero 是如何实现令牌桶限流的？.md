**原文链接：** [go-zero 是如何实现令牌桶限流的？](https://mp.weixin.qq.com/s/--AdUcwOQyP6r5W8ziVwUg)

上一篇文章介绍了 [如何实现计数器限流？](https://mp.weixin.qq.com/s/CTemkZ2aKPCPTuQiDJri0Q)主要有两种实现方式，分别是固定窗口和滑动窗口，并且分析了 go-zero 采用固定窗口方式实现的源码。

但是采用固定窗口实现的限流器会有两个问题：

1. 会出现请求量超出限制值两倍的情况
2. 无法很好处理流量突增问题

这篇文章来介绍一下令牌桶算法，可以很好解决以上两个问题。

## 工作原理

算法概念如下：

- 令牌以固定速率生成；
- 生成的令牌放入令牌桶中存放，如果令牌桶满了则多余的令牌会直接丢弃，当请求到达时，会尝试从令牌桶中取令牌，取到了令牌的请求可以执行；
- 如果桶空了，那么尝试取令牌的请求会被直接丢弃。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/tokenlimit.png)

令牌桶算法既能够将所有的请求平均分布到时间区间内，又能接受服务器能够承受范围内的突发请求，因此是目前使用较为广泛的一种限流算法。

## 源码实现

源码分析我们还是以 go-zero 项目为例，首先来看生成令牌的部分，依然是使用 Redis 来实现。

```go
// core/limit/tokenlimit.go

// 生成 token 速率
script = `local rate = tonumber(ARGV[1])
// 通容量
local capacity = tonumber(ARGV[2])
// 当前时间戳
local now = tonumber(ARGV[3])
// 请求数量
local requested = tonumber(ARGV[4])
// 需要多少秒才能把桶填满
local fill_time = capacity/rate
// 向下取整，ttl 为填满时间 2 倍
local ttl = math.floor(fill_time*2)
// 当前桶剩余容量，如果为 nil，说明第一次使用，赋值为桶最大容量
local last_tokens = tonumber(redis.call("get", KEYS[1]))
if last_tokens == nil then
    last_tokens = capacity
end

// 上次请求时间戳，如果为 nil 则赋值 0
local last_refreshed = tonumber(redis.call("get", KEYS[2]))
if last_refreshed == nil then
    last_refreshed = 0
end

// 距离上一次请求的时间跨度
local delta = math.max(0, now-last_refreshed)
// 距离上一次请求的时间跨度能生成的 token 数量和桶内剩余 token 数量的和
// 与桶容量比较，取二者的小值
local filled_tokens = math.min(capacity, last_tokens+(delta*rate))
// 判断请求数量和桶内 token 数量的大小
local allowed = filled_tokens >= requested
// 被请求消耗掉之后，更新剩余 token 数量
local new_tokens = filled_tokens
if allowed then
    new_tokens = filled_tokens - requested
end

// 更新 redis token
redis.call("setex", KEYS[1], ttl, new_tokens)
// 更新 redis 刷新时间
redis.call("setex", KEYS[2], ttl, now)

return allowed`
```

Redis 中主要保存两个 key，分别是 token 数量和刷新时间。

核心思想就是比较两次请求时间间隔内生成的 token 数量 + 桶内剩余 token 数量，和请求量之间的大小，如果满足则允许，否则则不允许。

限流器初始化：

```go
// A TokenLimiter controls how frequently events are allowed to happen with in one second.
type TokenLimiter struct {
    // 生成 token 速率
    rate           int
    // 桶容量
    burst          int
    store          *redis.Redis
    // 桶 key
    tokenKey       string
    // 桶刷新时间 key
    timestampKey   string
    rescueLock     sync.Mutex
    // redis 健康标识
    redisAlive     uint32
    // redis 健康监控启动状态
    monitorStarted bool
    // 内置单机限流器
    rescueLimiter  *xrate.Limiter
}

// NewTokenLimiter returns a new TokenLimiter that allows events up to rate and permits
// bursts of at most burst tokens.
func NewTokenLimiter(rate, burst int, store *redis.Redis, key string) *TokenLimiter {
    tokenKey := fmt.Sprintf(tokenFormat, key)
    timestampKey := fmt.Sprintf(timestampFormat, key)

    return &TokenLimiter{
        rate:          rate,
        burst:         burst,
        store:         store,
        tokenKey:      tokenKey,
        timestampKey:  timestampKey,
        redisAlive:    1,
        rescueLimiter: xrate.NewLimiter(xrate.Every(time.Second/time.Duration(rate)), burst),
    }
}
```

其中有一个变量 `rescueLimiter`，这是一个进程内的限流器。如果 Redis 发生故障了，那么就使用这个，算是一个保障，尽量避免系统被突发流量拖垮。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/go-zero-limit.png)

提供了四个可调用方法：

```go
// Allow is shorthand for AllowN(time.Now(), 1).
func (lim *TokenLimiter) Allow() bool {
    return lim.AllowN(time.Now(), 1)
}

// AllowCtx is shorthand for AllowNCtx(ctx,time.Now(), 1) with incoming context.
func (lim *TokenLimiter) AllowCtx(ctx context.Context) bool {
    return lim.AllowNCtx(ctx, time.Now(), 1)
}

// AllowN reports whether n events may happen at time now.
// Use this method if you intend to drop / skip events that exceed the rate.
// Otherwise, use Reserve or Wait.
func (lim *TokenLimiter) AllowN(now time.Time, n int) bool {
    return lim.reserveN(context.Background(), now, n)
}

// AllowNCtx reports whether n events may happen at time now with incoming context.
// Use this method if you intend to drop / skip events that exceed the rate.
// Otherwise, use Reserve or Wait.
func (lim *TokenLimiter) AllowNCtx(ctx context.Context, now time.Time, n int) bool {
    return lim.reserveN(ctx, now, n)
}
```

最终调用的都是 `reverveN` 方法：

```go
func (lim *TokenLimiter) reserveN(ctx context.Context, now time.Time, n int) bool {
    // 判断 Redis 健康状态，如果 Redis 故障，则使用进程内限流器
    if atomic.LoadUint32(&lim.redisAlive) == 0 {
        return lim.rescueLimiter.AllowN(now, n)
    }

    // 执行限流脚本
    resp, err := lim.store.EvalCtx(ctx,
        script,
        []string{
            lim.tokenKey,
            lim.timestampKey,
        },
        []string{
            strconv.Itoa(lim.rate),
            strconv.Itoa(lim.burst),
            strconv.FormatInt(now.Unix(), 10),
            strconv.Itoa(n),
        })
    // redis allowed == false
    // Lua boolean false -> r Nil bulk reply
    if err == redis.Nil {
        return false
    }
    if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
        logx.Errorf("fail to use rate limiter: %s", err)
        return false
    }
    if err != nil {
        logx.Errorf("fail to use rate limiter: %s, use in-process limiter for rescue", err)
        // 如果有异常的话，会启动进程内限流
        lim.startMonitor()
        return lim.rescueLimiter.AllowN(now, n)
    }

    code, ok := resp.(int64)
    if !ok {
        logx.Errorf("fail to eval redis script: %v, use in-process limiter for rescue", resp)
        lim.startMonitor()
        return lim.rescueLimiter.AllowN(now, n)
    }

    // redis allowed == true
    // Lua boolean true -> r integer reply with value of 1
    return code == 1
}
```

最后看一下进程内限流的启动与恢复：

```go
func (lim *TokenLimiter) startMonitor() {
    lim.rescueLock.Lock()
    defer lim.rescueLock.Unlock()

    // 需要加锁保护，如果程序已经启动了，直接返回，不要重复启动
    if lim.monitorStarted {
        return
    }

    lim.monitorStarted = true
    atomic.StoreUint32(&lim.redisAlive, 0)

    go lim.waitForRedis()
}

func (lim *TokenLimiter) waitForRedis() {
    ticker := time.NewTicker(pingInterval)
    // 更新监控进程的状态
    defer func() {
        ticker.Stop()
        lim.rescueLock.Lock()
        lim.monitorStarted = false
        lim.rescueLock.Unlock()
    }()

    for range ticker.C {
        // 对 redis 进行健康监测，如果 redis 服务恢复了
        // 则更新 redisAlive 标识，并退出 goroutine
        if lim.store.Ping() {
            atomic.StoreUint32(&lim.redisAlive, 1)
            return
        }
    }
}
```

以上就是本文的全部内容，如果觉得还不错的话欢迎**点赞**，**转发**和**关注**，感谢支持。

***

**参考文章：**

- https://juejin.cn/post/7052171117116522504
- https://www.infoq.cn/article/Qg2tX8fyw5Vt-f3HH673

**推荐阅读：**

- [如何实现计数器限流？](https://mp.weixin.qq.com/s/CTemkZ2aKPCPTuQiDJri0Q)
- [go-zero 是如何做路由管理的？](https://mp.weixin.qq.com/s/uTJ1En-BXiLvH45xx0eFsA)