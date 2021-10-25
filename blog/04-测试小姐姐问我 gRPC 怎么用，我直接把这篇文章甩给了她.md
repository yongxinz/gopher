**原文链接：** [测试小姐姐问我 gRPC 怎么用，我直接把这篇文章甩给了她](https://mp.weixin.qq.com/s/qdI2JqpMq6t2KN1byHaNCQ)

上篇文章 [gRPC，爆赞](https://mp.weixin.qq.com/s/1Xbca4Dv0akonAZerrChgA) 直接爆了，内容主要包括：简单的 gRPC 服务，流处理模式，验证器，Token 认证和证书认证。

在多个平台的阅读量都创了新高，在 oschina 更是获得了首页推荐，阅读量到了 1w+，这已经是我单篇阅读的高峰了。

看来只要用心写还是有收获的。

这篇咱们还是从实战出发，主要介绍 gRPC 的发布订阅模式，REST 接口和超时控制。

相关代码我会都上传到 [GitHub](https://github.com/yongxinz/go-example)，感兴趣的小伙伴可以去查看或下载。

### 发布和订阅模式

发布订阅是一个常见的设计模式，开源社区中已经存在很多该模式的实现。其中 docker 项目中提供了一个 pubsub 的极简实现，下面是基于 pubsub 包实现的本地发布订阅代码：

```go
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/moby/moby/pkg/pubsub"
)

func main() {
	p := pubsub.NewPublisher(100*time.Millisecond, 10)

	golang := p.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, "golang:") {
				return true
			}
		}
		return false
	})
	docker := p.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, "docker:") {
				return true
			}
		}
		return false
	})

	go p.Publish("hi")
	go p.Publish("golang: https://golang.org")
	go p.Publish("docker: https://www.docker.com/")
	time.Sleep(1)

	go func() {
		fmt.Println("golang topic:", <-golang)
	}()
	go func() {
		fmt.Println("docker topic:", <-docker)
	}()

	<-make(chan bool)
}

```

这段代码首先通过 `pubsub.NewPublisher` 创建了一个对象，然后通过 `p.SubscribeTopic` 实现订阅，`p.Publish` 来发布消息。

执行效果如下：

```
docker topic: docker: https://www.docker.com/
golang topic: golang: https://golang.org
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan receive]:
main.main()
	/Users/zhangyongxin/src/go-example/grpc-example/pubsub/server/pubsub.go:43 +0x1e7
exit status 2
```

订阅消息可以正常打印。

但有一个死锁报错，是因为这条语句 `<-make(chan bool)` 引起的。但是如果没有这条语句就不能正常打印订阅消息。

这里就不是很懂了，有没有大佬知道，欢迎留言，求指导。

接下来就用 gRPC 和 pubsub 包实现发布订阅模式。

需要实现四个部分：
1. **proto** 文件；
2. **服务端：** 用于接收订阅请求，同时也接收发布请求，并将发布请求转发给订阅者；
3. **订阅客户端：** 用于从服务端订阅消息，处理消息；
4. **发布客户端：** 用于向服务端发送消息。

#### proto 文件

首先定义 proto 文件：

```
syntax = "proto3";

package proto;
 
message String {
    string value = 1;
}
 
service PubsubService {
    rpc Publish (String) returns (String);
    rpc SubscribeTopic (String) returns (stream String);
    rpc Subscribe (String) returns (stream String);
}
```

定义三个方法，分别是一个发布 `Publish` 和两个订阅 `Subscribe` 和 `SubscribeTopic`。

`Subscribe` 方法接收全部消息，而 `SubscribeTopic` 根据特定的 `Topic` 接收消息。

#### 服务端

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"server/proto"
	"strings"
	"time"

	"github.com/moby/moby/pkg/pubsub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type PubsubService struct {
	pub *pubsub.Publisher
}

func (p *PubsubService) Publish(ctx context.Context, arg *proto.String) (*proto.String, error) {
	p.pub.Publish(arg.GetValue())
	return &proto.String{}, nil
}

func (p *PubsubService) SubscribeTopic(arg *proto.String, stream proto.PubsubService_SubscribeTopicServer) error {
	ch := p.pub.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, arg.GetValue()) {
				return true
			}
		}
		return false
	})

	for v := range ch {
		if err := stream.Send(&proto.String{Value: v.(string)}); nil != err {
			return err
		}
	}
	return nil
}

func (p *PubsubService) Subscribe(arg *proto.String, stream proto.PubsubService_SubscribeServer) error {
	ch := p.pub.Subscribe()

	for v := range ch {
		if err := stream.Send(&proto.String{Value: v.(string)}); nil != err {
			return err
		}
	}
	return nil
}

func NewPubsubService() *PubsubService {
	return &PubsubService{pub: pubsub.NewPublisher(100*time.Millisecond, 10)}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 简单调用
	server := grpc.NewServer()
	// 注册 grpcurl 所需的 reflection 服务
	reflection.Register(server)
	// 注册业务服务
	proto.RegisterPubsubServiceServer(server, NewPubsubService())

	fmt.Println("grpc server start ...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

对比之前的发布订阅程序，其实这里是将 `*pubsub.Publisher` 作为了 gRPC 的结构体 `PubsubService` 的一个成员。

然后还是按照 gRPC 的开发流程，实现结构体对应的三个方法。

最后，在注册服务时，将 `NewPubsubService()` 服务注入，实现本地发布订阅功能。

#### 订阅客户端

```go
package main

import (
	"client/proto"
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewPubsubServiceClient(conn)
	stream, err := client.Subscribe(
		context.Background(), &proto.String{Value: "golang:"},
	)
	if nil != err {
		log.Fatal(err)
	}

	go func() {
		for {
			reply, err := stream.Recv()
			if nil != err {
				if io.EOF == err {
					break
				}
				log.Fatal(err)
			}
			fmt.Println("sub1: ", reply.GetValue())
		}
	}()

	streamTopic, err := client.SubscribeTopic(
		context.Background(), &proto.String{Value: "golang:"},
	)
	if nil != err {
		log.Fatal(err)
	}

	go func() {
		for {
			reply, err := streamTopic.Recv()
			if nil != err {
				if io.EOF == err {
					break
				}
				log.Fatal(err)
			}
			fmt.Println("subTopic: ", reply.GetValue())
		}
	}()

	<-make(chan bool)
}
```

新建一个 `NewPubsubServiceClient` 对象，然后分别实现 `client.Subscribe` 和 `client.SubscribeTopic` 方法，再通过 goroutine 不停接收消息。

#### 发布客户端

```go
package main

import (
	"client/proto"
	"context"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := proto.NewPubsubServiceClient(conn)

	_, err = client.Publish(
		context.Background(), &proto.String{Value: "golang: hello Go"},
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Publish(
		context.Background(), &proto.String{Value: "docker: hello Docker"},
	)
	if nil != err {
		log.Fatal(err)
	}

}
```

新建一个 `NewPubsubServiceClient` 对象，然后通过 `client.Publish` 方法发布消息。

当代码全部写好之后，我们开三个终端来测试一下：

**终端1** 上启动服务端：

```
go run main.go
```

**终端2** 上启动订阅客户端：

```
go run sub_client.go
```

**终端3** 上执行发布客户端：

```
go run pub_client.go
```

这样，在 **终端2** 上就有对应的输出了：

```
subTopic:  golang: hello Go
sub1:  golang: hello Go
sub1:  docker: hello Docker
```

也可以再多开几个订阅终端，那么每一个订阅终端上都会有相同的内容输出。

**源码地址：** [GitHub](https://github.com/yongxinz/go-example/tree/main/grpc-example/pubsub)

### REST 接口

gRPC 一般用于集群内部通信，如果需要对外提供服务，大部分都是通过 REST 接口的方式。开源项目 grpc-gateway 提供了将 gRPC 服务转换成 REST 服务的能力，通过这种方式，就可以直接访问 gRPC API 了。

但我觉得，实际上这么用的应该还是比较少的。如果提供 REST 接口的话，直接写一个 HTTP 服务会方便很多。

#### proto 文件

第一步还是创建一个 proto 文件：

```
syntax = "proto3";

package proto;

import "google/api/annotations.proto";

message StringMessage {
  string value = 1;
}

service RestService {
    rpc Get(StringMessage) returns (StringMessage) {
        option (google.api.http) = {
            get: "/get/{value}"
        };
    }
    rpc Post(StringMessage) returns (StringMessage) {
        option (google.api.http) = {
            post: "/post"
            body: "*"
        };
    }
}
```

定义一个 REST 服务 `RestService`，分别实现 `GET` 和 `POST` 方法。

安装插件：

```
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
```

生成对应代码：

```
protoc -I/usr/local/include -I. \
    -I$GOPATH/pkg/mod \
    -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
    --grpc-gateway_out=. --go_out=plugins=grpc:.\
    --swagger_out=. \
    helloworld.proto
```

`--grpc-gateway_out` 参数可生成对应的 gw 文件，`--swagger_out` 参数可生成对应的 API 文档。

在我这里生成的两个文件如下：

```
helloworld.pb.gw.go
helloworld.swagger.json
```

#### REST 服务

```go
package main

import (
	"context"
	"log"
	"net/http"

	"rest/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	err := proto.RegisterRestServiceHandlerFromEndpoint(
		ctx, mux, "localhost:50051",
		[]grpc.DialOption{grpc.WithInsecure()},
	)
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":8080", mux)
}
```

这里主要是通过实现 gw 文件中的 `RegisterRestServiceHandlerFromEndpoint` 方法来连接 gRPC 服务。

#### gRPC 服务

```go
package main

import (
	"context"
	"net"

	"rest/proto"

	"google.golang.org/grpc"
)

type RestServiceImpl struct{}

func (r *RestServiceImpl) Get(ctx context.Context, message *proto.StringMessage) (*proto.StringMessage, error) {
	return &proto.StringMessage{Value: "Get hi:" + message.Value + "#"}, nil
}

func (r *RestServiceImpl) Post(ctx context.Context, message *proto.StringMessage) (*proto.StringMessage, error) {
	return &proto.StringMessage{Value: "Post hi:" + message.Value + "@"}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	proto.RegisterRestServiceServer(grpcServer, new(RestServiceImpl))
	lis, _ := net.Listen("tcp", ":50051")
	grpcServer.Serve(lis)
}
```

gRPC 服务的实现方式还是和以前一样。

以上就是全部代码，现在来测试一下：

启动三个终端：

**终端1** 启动 gRPC 服务：

```
go run grpc_service.go
```

**终端2** 启动 REST 服务：

```
go run rest_service.go
```

**终端3** 来请求 REST 服务：

```
$ curl localhost:8080/get/gopher
{"value":"Get hi:gopher"}

$ curl localhost:8080/post -X POST --data '{"value":"grpc"}'
{"value":"Post hi:grpc"}
```

**源码地址：** [GitHub](https://github.com/yongxinz/go-example/tree/main/grpc-example/rest)

### 超时控制

最后一部分介绍一下超时控制，这部分内容是非常重要的。

一般的 WEB 服务 API，或者是 Nginx 都会设置一个超时时间，超过这个时间，如果还没有数据返回，服务端可能直接返回一个超时错误，或者客户端也可能结束这个连接。

如果没有这个超时时间，那是相当危险的。所有请求都阻塞在服务端，会消耗大量资源，比如内存。如果资源耗尽的话，甚至可能会导致整个服务崩溃。

那么，在 gRPC 中怎么设置超时时间呢？主要是通过上下文 `context.Context` 参数，具体来说就是 `context.WithDeadline` 函数。

#### proto 文件

创建最简单的 proto 文件，这个不多说。

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

#### 客户端

```go
package main

import (
	"client/proto"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	// 简单调用
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer conn.Close()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(3*time.Second)))
	defer cancel()

	client := proto.NewGreeterClient(conn)
	// 简单调用
	reply, err := client.SayHello(ctx, &proto.HelloRequest{Name: "zzz"})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Fatalln("client.SayHello err: deadline")
			}
		}

		log.Fatalf("client.SayHello err: %v", err)
	}
	fmt.Println(reply.Message)
}
```

通过下面的函数设置一个 3s 的超时时间：

```go
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(3*time.Second)))
defer cancel()
```

然后在响应错误中对超时错误进行检测。

#### 服务端

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime"
	"server/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type greeter struct {
}

func (*greeter) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	data := make(chan *proto.HelloReply, 1)
	go handle(ctx, req, data)
	select {
	case res := <-data:
		return res, nil
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "Client cancelled, abandoning.")
	}
}

func handle(ctx context.Context, req *proto.HelloRequest, data chan<- *proto.HelloReply) {
	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
		runtime.Goexit() //超时后退出该Go协程
	case <-time.After(4 * time.Second): // 模拟耗时操作
		res := proto.HelloReply{
			Message: "hello " + req.Name,
		}
		// //修改数据库前进行超时判断
		// if ctx.Err() == context.Canceled{
		// 	...
		// 	//如果已经超时，则退出
		// }
		data <- &res
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 简单调用
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
```

服务端增加一个 `handle` 函数，其中 `case <-time.After(4 * time.Second)` 表示 4s 之后才会执行其对应代码，用来模拟超时请求。

如果客户端超时时间超过 4s 的话，就会产生超时报错。

下面来模拟一下：

**服务端：**

```
$ go run main.go
grpc server start ...
2021/10/24 22:57:40 context deadline exceeded
```

**客户端：**

```
$ go run main.go
2021/10/24 22:57:40 client.SayHello err: deadline
exit status 1
```

**源码地址：** [GitHub](https://github.com/yongxinz/go-example/tree/main/grpc-example/deadline)

### 总结

本文主要介绍了 gRPC 的三部分实战内容，分别是：

1. 发布订阅模式
2. REST 接口
3. 超时控制

个人感觉，超时控制还是最重要的，在平时的开发过程中需要多多注意。

结合上篇文章，gRPC 的实战内容就写完了，代码全部可以执行，也都上传到了 [GitHub](https://github.com/yongxinz/go-example)。

大家如果有任何疑问，欢迎给我留言，如果感觉不错的话，也欢迎关注和转发。


---


**源码地址：** 

- [https://github.com/yongxinz/go-example](https://github.com/yongxinz/go-example)
- [https://github.com/yongxinz/gopher](https://github.com/yongxinz/gopher)

关注公众号 **AlwaysBeta**，回复「**goebook**」领取 Go 编程经典书籍。

<center class="half">
    <img src="https://github.com/yongxinz/gopher/blob/main/alwaysbeta.JPG" width="300"/>
</center>

**推荐阅读：**

- [gRPC，爆赞](https://mp.weixin.qq.com/s/1Xbca4Dv0akonAZerrChgA)
- [使用 grpcurl 通过命令行访问 gRPC 服务](https://mp.weixin.qq.com/s/GShwcGCopXVmxCKnYf5FhA)
- [听说，99% 的 Go 程序员都被 defer 坑过](https://mp.weixin.qq.com/s/1T6Z74Wri27Ap8skeJiyWQ)

**参考：**

- [https://chai2010.cn/advanced-go-programming-book/ch4-rpc/readme.html](https://chai2010.cn/advanced-go-programming-book/ch4-rpc/readme.html)
- [https://codeleading.com/article/94674952433/](https://codeleading.com/article/94674952433/)
- [https://juejin.cn/post/6844904017017962504](https://juejin.cn/post/6844904017017962504)
- [https://www.cnblogs.com/FireworksEasyCool/p/12702959.html](https://www.cnblogs.com/FireworksEasyCool/p/12702959.html)