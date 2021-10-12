**原文链接：** [使用 grpcurl 通过命令行访问 gRPC 服务](https://mp.weixin.qq.com/s/GShwcGCopXVmxCKnYf5FhA)

一般情况下测试 gRPC 服务，都是通过客户端来直接请求服务端。如果客户端还没准备好的话，也可以使用 [BloomRPC](https://appimage.github.io/BloomRPC/) 这样的 GUI 客户端。

如果环境不支持安装这种 GUI 客户端的话，那么有没有一种工具，类似于 `curl` 这样的，直接通过终端，在命令行发起请求呢？

答案肯定是有的，就是本文要介绍的 `grpcurl`。

### gRPC Server

首先来写一个简单的 gRPC Server：

helloworld.proto：

```
syntax = "proto3";

package proto;

// The greeting service definition.
service Greeter {
    // Sends a greeting
    rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
    string name = 1;
}

// The response message containing the greetings
message HelloReply {
    string message = 1;
}
```

main.go

```go
package main

import (
	"context"
	"fmt"
	"grpc-hello/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	// 注册 grpcurl 所需的 reflection 服务
	reflection.Register(server)
	// 注册业务服务
	proto.RegisterGreeterServer(server, &greeter{})

	fmt.Println("grpc server start ...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type greeter struct {
}

func (*greeter) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	fmt.Println(req)
	reply := &proto.HelloReply{Message: "hello"}
	return reply, nil
}

```

运行服务：

```
go run main.go

server start ...
```

### grpcurl 安装

这里我介绍三种方式：

#### Mac

```
brew install grpcurl
```

#### Docker

```
# Download image
docker pull fullstorydev/grpcurl:latest
# Run the tool
docker run fullstorydev/grpcurl api.grpc.me:443 list
```

#### go tool

如果有 Go 环境的话，可以通过 go tool 来安装：

```
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### grpcurl 使用

在使用 grpcurl 时，需要通过 `-cert` 和 `-key` 参数设置公钥和私钥文件，表示链接启用了 TLS 协议的服务。

对于没有启用 TLS 协议的 gRPC 服务，通过 `-plaintext` 参数忽略 TLS 证书的验证过程。

如果是 Unix Socket 协议，则需要指定 `-unix` 参数。

**查看服务列表：**

```
grpcurl -plaintext 127.0.0.1:50051 list
```

输出：

```
grpc.reflection.v1alpha.ServerReflection
proto.Greeter
```

**查看某个服务的方法列表：**

```
grpcurl -plaintext 127.0.0.1:50051 list proto.Greeter
```

输出：

```
proto.Greeter.SayHello
```

**查看方法定义：**

```
grpcurl -plaintext 127.0.0.1:50051 describe proto.Greeter.SayHello
```

输出：

```
proto.Greeter.SayHello is a method:
rpc SayHello ( .proto.HelloRequest ) returns ( .proto.HelloReply );
```

**查看请求参数：**

```
grpcurl -plaintext 127.0.0.1:50051 describe proto.HelloRequest
```

输出：

```
proto.HelloRequest is a message:
message HelloRequest {
  string name = 1;
}
```

**请求服务：**

```
grpcurl -d '{"name": "zhangsan"}' -plaintext 127.0.0.1:50051 proto.Greeter.SayHello
```

输出：

```
{
  "message": "hello"
}
```

`-d` 参数后面也可以跟 `@`，表示从标准输入读取 json 参数，一般用于输入比较复杂的 json 数据，也可以用于测试流方法。

```
grpcurl -d @ -plaintext 127.0.0.1:50051 proto.Greeter.SayHello
```

### 可能遇到的错误

可能会遇到三个报错：

**1、gRPC Server 未启用 TLS：**

报错信息：

```
Failed to dial target host "127.0.0.1:50051": tls: first record does not look like a TLS handshake
```

**解决：**

请求时增加参数：`-plaintext`，参考上面的命令。

**2、服务没有启动 reflection 反射服务**

报错信息：

```
Failed to list services: server does not support the reflection API
```

**解决：**

这行代码是关键，一定要包含：

```go
// 注册 grpcurl 所需的 reflection 服务
reflection.Register(server)
```

**3、参数格式错误：**

报错信息：

```
Error invoking method "greet.Greeter/SayHello": error getting request data: invalid character 'n' looking for beginning of object key string
```

**解决：**

`-d` 后面参数为 json 格式，并且需要使用 `''` 包裹起来。

### 总结

用这个工具做一些简单的测试还是相当方便的，上手也简单。只要掌握文中提到的几条命令，基本可以涵盖大部分的测试需求了。

---

**扩展阅读：**

1. https://appimage.github.io/BloomRPC/
2. https://github.com/fullstorydev/grpcurl

文章中的脑图和源码都上传到了 GitHub，有需要的同学可自行下载。

**地址：** https://github.com/yongxinz/gopher/tree/main/blog

关注公众号 **AlwaysBeta**，回复「**goebook**」领取 Go 编程经典书籍。

<center class="half">
    <img src="https://github.com/yongxinz/gopher/blob/main/alwaysbeta.JPG" width="300"/>
</center>

**往期文章列表：**

1. [被 Docker 日志坑惨了](https://mp.weixin.qq.com/s/3Tkc15dTCEDUAZaZ88pcSQ)
2. [推荐三个实用的 Go 开发工具](https://mp.weixin.qq.com/s/3GLMLhegB3wF5_62mpmePA)
3. [这个 TCP 问题你得懂：Cannot assign requested address](https://mp.weixin.qq.com/s/-cThzr5N2w3IEYYf-duCDA)


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
11. [Go 专栏｜并发编程：goroutine，channel 和 sync](https://mp.weixin.qq.com/s/VG4CSfT2OfxA6nfygWLSyw)