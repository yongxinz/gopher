gRPC 这项技术真是太棒了，接口约束严格，性能还高，在 k8s 和很多微服务框架中都有应用。

作为一名程序员，学就对了。

之前用 Python 写过一些 gRPC 服务，现在准备用 Go 来感受一下原汁原味的 gRPC 程序开发。

本文的特点是直接用代码说话，通过开箱即用的完整代码，来介绍 gRPC 的各种使用方法。

代码已经上传到 [GitHub](https://github.com/yongxinz/go-example)，下面正式开始。

### 介绍

gRPC 是 Google 公司基于 Protobuf 开发的跨语言的开源 RPC 框架。gRPC 基于 HTTP/2 协议设计，可以基于一个 HTTP/2 链接提供多个服务，对于移动设备更加友好。

### 入门

首先来看一个最简单的 gRPC 服务，第一步是定义 proto 文件，因为 gRPC 也是 C/S 架构，这一步相当于明确接口规范。

**proto**

```proto
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

使用 protoc-gen-go 内置的 gRPC 插件生成 gRPC 代码：

```
protoc --go_out=plugins=grpc:. helloworld.proto
```

执行完这个命令之后，会在当前目录生成一个 helloworld.pb.go 文件，文件中分别定义了服务端和客户端的接口：

```go
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GreeterClient interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

// GreeterServer is the server API for Greeter service.
type GreeterServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}
```

接下来就是写服务端和客户端的代码，分别实现对应的接口。

**server**

```go
package main

import (
	"context"
	"fmt"
	"grpc-server/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type greeter struct {
}

func (*greeter) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	fmt.Println(req)
	reply := &proto.HelloReply{Message: "hello"}
	return reply, nil
}

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
```

**client**

```go
package main

import (
	"context"
	"fmt"
	"grpc-client/proto"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &proto.HelloRequest{Name: "zhangsan"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.Message)
}
```

这样就完成了最基础的 gRPC 服务的开发，接下来我们就在这个「基础模板」上不断丰富，学习更多特性。

### 流方式

接下来看看流的方式，顾名思义，数据可以源源不断的发送和接收。

流的话分单向流和双向流，这里我们直接通过双向流来举例。

**proto**

```proto
service Greeter {
    // Sends a greeting
    rpc SayHello (HelloRequest) returns (HelloReply) {}
    // Sends stream message
    rpc SayHelloStream (stream HelloRequest) returns (stream HelloReply) {}
}
```

增加一个流函数 `SayHelloStream`，通过 `stream` 关键词来指定流特性。

需要重新生成 helloworld.pb.go 文件，这里不再多说。

**server**

```go
func (*greeter) SayHelloStream(stream proto.Greeter_SayHelloStreamServer) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		fmt.Println("Recv: " + args.Name)
		reply := &proto.HelloReply{Message: "hi " + args.Name}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}
```

在「基础模板」上增加 `SayHelloStream` 函数，其他都不需要变。

**client**

```go
client := proto.NewGreeterClient(conn)

// 流处理
stream, err := client.SayHelloStream(context.Background())
if err != nil {
	log.Fatal(err)
}

// 发送消息
go func() {
	for {
		if err := stream.Send(&proto.HelloRequest{Name: "zhangsan"}); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}
}()

// 接收消息
for {
	reply, err := stream.Recv()
	if err != nil {
		if err == io.EOF {
			break
		}
		log.Fatal(err)
	}
	fmt.Println(reply.Message)
}
```

通过一个 goroutine 发送消息，主程序的 `for` 循环接收消息。

执行程序会发现，服务端和客户端都不断有打印输出。

### 验证器

接下来是验证器，这个需求是很自然会想到的，因为涉及到接口之间的请求，那么对参数进行适当的校验是很有必要的。

在这里我们使用 protoc-gen-govalidators 和 go-grpc-middleware 来实现。

先安装：

```
go get github.com/mwitkow/go-proto-validators/protoc-gen-govalidators

go get github.com/grpc-ecosystem/go-grpc-middleware
```

接下来修改 proto 文件：

**proto**

```proto
import "github.com/mwitkow/go-proto-validators@v0.3.2/validator.proto";

message HelloRequest {
    string name = 1 [
        (validator.field) = {regex: "^[z]{2,5}$"}
    ];
}
```

在这里对 `name` 参数进行校验，需要符合正则的要求才可以正常请求。

还有其他验证规则，比如对数字大小进行验证等，这里不做过多介绍。

接下来生成 *.pb.go 文件：

```
protoc  \
    --proto_path=${GOPATH}/pkg/mod \
    --proto_path=${GOPATH}/pkg/mod/github.com/gogo/protobuf@v1.3.2 \
    --proto_path=. \
    --govalidators_out=. --go_out=plugins=grpc:.\
    *.proto
```

执行成功之后，目录下会多一个 helloworld.validator.pb.go 文件。

这里需要特别注意一下，使用之前的简单命令是不行的，需要使用多个 `proto_path` 参数指定导入 proto 文件的目录。

官方给了两种依赖情况，一个是 google protobuf，一个是 gogo protobuf。我这里使用的是第二种。

即使使用上面的命令，也有可能会遇到这个报错：

```
Import "github.com/mwitkow/go-proto-validators/validator.proto" was not found or had errors
```

但不要慌，大概率是引用路径的问题，一定要看好自己的安装版本，以及在 `GOPATH` 中的具体路径。

最后是服务端代码改造：

引入包：

```
grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
```

然后在初始化的时候增加验证器功能：

```go
server := grpc.NewServer(
	grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_validator.UnaryServerInterceptor(),
		),
	),
	grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_validator.StreamServerInterceptor(),
		),
	),
)
```

启动程序之后，我们再用之前的客户端代码来请求，会收到报错：

```
2021/10/11 18:32:59 rpc error: code = InvalidArgument desc = invalid field Name: value 'zhangsan' must be a string conforming to regex "^[z]{2,5}$"
exit status 1
```

因为 `name: zhangsan` 是不符合服务端正则要求的，但是如果传参 `name: zzz`，就可以正常返回了。

### Token 认证

终于到认证环节了，先看 Token 认证方式，然后再介绍证书认证。

先改造服务端，有了上文验证器的经验，那么可以采用同样的方式，写一个拦截器，然后在初始化 server 时候注入。

**认证函数：**

```go
func Auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("missing credentials")
	}

	var user string
	var password string

	if val, ok := md["user"]; ok {
		user = val[0]
	}
	if val, ok := md["password"]; ok {
		password = val[0]
	}

	if user != "admin" || password != "admin" {
		return grpc.Errorf(codes.Unauthenticated, "invalid token")
	}

	return nil
}
```

`metadata.FromIncomingContext` 从上下文读取用户名和密码，然后和实际数据进行比较，判断是否通过认证。

**拦截器：**

```go
var authInterceptor grpc.UnaryServerInterceptor
authInterceptor = func(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	//拦截普通方法请求，验证 Token
	err = Auth(ctx)
	if err != nil {
		return
	}
	// 继续处理请求
	return handler(ctx, req)
}
```

**初始化：**

```go
server := grpc.NewServer(
	grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			authInterceptor,
			grpc_validator.UnaryServerInterceptor(),
		),
	),
	grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_validator.StreamServerInterceptor(),
		),
	),
)
```

除了上文的验证器，又多了 Token 认证拦截器 `authInterceptor`。

最后是客户端改造，客户端需要实现 `PerRPCCredentials` 接口。

```go
type PerRPCCredentials interface {
    // GetRequestMetadata gets the current request metadata, refreshing
    // tokens if required. This should be called by the transport layer on
    // each request, and the data should be populated in headers or other
    // context. If a status code is returned, it will be used as the status
    // for the RPC. uri is the URI of the entry point for the request.
    // When supported by the underlying implementation, ctx can be used for
    // timeout and cancellation.
    // TODO(zhaoq): Define the set of the qualified keys instead of leaving
    // it as an arbitrary string.
    GetRequestMetadata(ctx context.Context, uri ...string) (
        map[string]string,    error,
    )
    // RequireTransportSecurity indicates whether the credentials requires
    // transport security.
    RequireTransportSecurity() bool
}
```

`GetRequestMetadata` 方法返回认证需要的必要信息，`RequireTransportSecurity` 方法表示是否启用安全链接，在生产环境中，一般都是启用的，但为了测试方便，暂时这里不启用了。

**实现接口：**

```go
type Authentication struct {
	User     string
	Password string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (
	map[string]string, error,
) {
	return map[string]string{"user": a.User, "password": a.Password}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return false
}
```

**连接：**

```go
conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))
```

好了，现在我们的服务就有 Token 认证功能了。如果用户名或密码错误，客户端就会收到：

```
2021/10/11 20:39:35 rpc error: code = Unauthenticated desc = invalid token
exit status 1
```

如果用户名和密码正确，则可以正常返回。

### 单向证书认证

证书认证分两种方式：

1. 单向认证
2. 双向认证

先看一下单向认证方式：

#### 生成证书

首先通过 openssl 工具生成自签名的 SSL 证书。

**1、生成私钥：**

```
openssl genrsa -des3 -out server.pass.key 2048
```

**2、去除私钥中密码：**

```
openssl rsa -in server.pass.key -out server.key
```

**3、生成 csr 文件：**

```
openssl req -new -key server.key -out server.csr -subj "/C=CN/ST=beijing/L=beijing/O=grpcdev/OU=grpcdev/CN=example.grpcdev.cn"
```

**4、生成证书：**

```
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
```

再多说一句，分别介绍一下 X.509 证书包含的三个文件：key，csr 和 crt。

- **key：** 服务器上的私钥文件，用于对发送给客户端数据的加密，以及对从客户端接收到数据的解密。
- **csr：** 证书签名请求文件，用于提交给证书颁发机构（CA）对证书签名。
- **crt：** 由证书颁发机构（CA）签名后的证书，或者是开发者自签名的证书，包含证书持有人的信息，持有人的公钥，以及签署者的签名等信息。

#### gRPC 代码

证书有了之后，剩下的就是改造程序了，首先是服务端代码。

```go
// 证书认证-单向认证
creds, err := credentials.NewServerTLSFromFile("keys/server.crt", "keys/server.key")
if err != nil {
	log.Fatal(err)
	return
}

server := grpc.NewServer(grpc.Creds(creds))
```

只有几行代码需要修改，很简单，接下来是客户端。

由于是单向认证，不需要为客户端单独生成证书，只需要把服务端的 crt 文件拷贝到客户端对应目录下即可。

```go
// 证书认证-单向认证
creds, err := credentials.NewClientTLSFromFile("keys/server.crt", "example.grpcdev.cn")
if err != nil {
	log.Fatal(err)
	return
}
conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
```

好了，现在我们的服务就支持单向证书认证了。

但是还没完，这里可能会遇到一个问题：

```
2021/10/11 21:32:37 rpc error: code = Unavailable desc = connection error: desc = "transport: authentication handshake failed: x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0"
exit status 1
```

原因是 Go 1.15 开始[废弃了 CommonName](https://golang.org/doc/go1.15#commonname)，推荐使用 SAN 证书。如果想要兼容之前的方式，可以通过设置环境变量的方式支持，如下：

```
export GODEBUG="x509ignoreCN=0"
```

但是需要注意，从 Go 1.17 开始，环境变量就不再生效了，必须通过 SAN 方式才行。所以，为了后续的 Go 版本升级，还是早日支持为好。

### 双向证书认证

最后来看看双向证书认证。

#### 生成带 SAN 的证书

还是先生成证书，但这次有一点不一样，我们需要生成带 SAN 扩展的证书。

什么是 SAN？

SAN（Subject Alternative Name）是 SSL 标准 x509 中定义的一个扩展。使用了 SAN 字段的 SSL 证书，可以扩展此证书支持的域名，使得一个证书可以支持多个不同域名的解析。

将默认的 OpenSSL 配置文件拷贝到当前目录。

Linux 系统在：

```
/etc/pki/tls/openssl.cnf
```

Mac 系统在：

```
/System/Library/OpenSSL/openssl.cnf
```

修改临时配置文件，找到 `[ req ]` 段落，然后将下面语句的注释去掉。

```
req_extensions = v3_req # The extensions to add to a certificate request
```

接着添加以下配置：

```
[ v3_req ]
# Extensions to add to a certificate request

basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = www.example.grpcdev.cn
```

`[ alt_names ]` 位置可以配置多个域名，比如：

```
[ alt_names ]
DNS.1 = www.example.grpcdev.cn
DNS.2 = www.test.grpcdev.cn
```

为了测试方便，这里只配置一个域名。

**1、生成 ca 证书：**

```go
openssl genrsa -out ca.key 2048

openssl req -x509 -new -nodes -key ca.key -subj "/CN=example.grpcdev.com" -days 5000 -out ca.pem
```

**2、生成服务端证书：**

```
# 生成证书
openssl req -new -nodes \
    -subj "/C=CN/ST=Beijing/L=Beijing/O=grpcdev/OU=grpcdev/CN=www.example.grpcdev.cn" \
    -config <(cat openssl.cnf \
        <(printf "[SAN]\nsubjectAltName=DNS:www.example.grpcdev.cn")) \
    -keyout server.key \
    -out server.csr
    
# 签名证书
openssl x509 -req -days 365000 \
    -in server.csr -CA ca.pem -CAkey ca.key -CAcreateserial \
    -extfile <(printf "subjectAltName=DNS:www.example.grpcdev.cn") \
    -out server.pem
```

**3、生成客户端证书：**

```
# 生成证书
openssl req -new -nodes \
    -subj "/C=CN/ST=Beijing/L=Beijing/O=grpcdev/OU=grpcdev/CN=www.example.grpcdev.cn" \
    -config <(cat openssl.cnf \
        <(printf "[SAN]\nsubjectAltName=DNS:www.example.grpcdev.cn")) \
    -keyout client.key \
    -out client.csr

# 签名证书
openssl x509 -req -days 365000 \
    -in client.csr -CA ca.pem -CAkey ca.key -CAcreateserial \
    -extfile <(printf "subjectAltName=DNS:www.example.grpcdev.cn") \
    -out client.pem
```

#### gRPC 代码

接下来开始修改代码，先看服务端：

```go
// 证书认证-双向认证
// 从证书相关文件中读取和解析信息，得到证书公钥、密钥对
cert, _ := tls.LoadX509KeyPair("cert/server.pem", "cert/server.key")
// 创建一个新的、空的 CertPool
certPool := x509.NewCertPool()
ca, _ := ioutil.ReadFile("cert/ca.pem")
// 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
certPool.AppendCertsFromPEM(ca)
// 构建基于 TLS 的 TransportCredentials 选项
creds := credentials.NewTLS(&tls.Config{
	// 设置证书链，允许包含一个或多个
	Certificates: []tls.Certificate{cert},
	// 要求必须校验客户端的证书。可以根据实际情况选用以下参数
	ClientAuth: tls.RequireAndVerifyClientCert,
	// 设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
	ClientCAs: certPool,
})
```

再看客户端：

```go
// 证书认证-双向认证
// 从证书相关文件中读取和解析信息，得到证书公钥、密钥对
cert, _ := tls.LoadX509KeyPair("cert/client.pem", "cert/client.key")
// 创建一个新的、空的 CertPool
certPool := x509.NewCertPool()
ca, _ := ioutil.ReadFile("cert/ca.pem")
// 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
certPool.AppendCertsFromPEM(ca)
// 构建基于 TLS 的 TransportCredentials 选项
creds := credentials.NewTLS(&tls.Config{
	// 设置证书链，允许包含一个或多个
	Certificates: []tls.Certificate{cert},
	// 要求必须校验客户端的证书。可以根据实际情况选用以下参数
	ServerName: "www.example.grpcdev.cn",
	RootCAs:    certPool,
})
```

大功告成。

### Python 客户端

前面已经说了，gRPC 是跨语言的，那么，本文最后我们用 Python 写一个客户端，来请求 Go 服务端。

使用最简单的方式来实现：

proto 文件就使用最开始的「基础模板」的 proto 文件：

```proto
syntax = "proto3";

package proto;

// The greeting service definition.
service Greeter {
    // Sends a greeting
    rpc SayHello (HelloRequest) returns (HelloReply) {}
    // Sends stream message
    rpc SayHelloStream (stream HelloRequest) returns (stream HelloReply) {}
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

同样的，也需要通过命令行的方式生成 pb.py 文件：

```
python3 -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. ./*.proto
```

执行成功之后会在目录下生成 helloworld_pb2.py 和 helloworld_pb2_grpc.py 两个文件。

这个过程也可能会报错：

```
ModuleNotFoundError: No module named 'grpc_tools'
```

别慌，是缺少包，安装就好：

```
pip3 install grpcio
pip3 install grpcio-tools
```

最后看一下 Python 客户端代码：

```python
import grpc

import helloworld_pb2
import helloworld_pb2_grpc


def main():
    channel = grpc.insecure_channel("127.0.0.1:50051")
    stub = helloworld_pb2_grpc.GreeterStub(channel)
    response = stub.SayHello(helloworld_pb2.HelloRequest(name="zhangsan"))
    print(response.message)


if __name__ == '__main__':
    main()
```

这样，就可以通过 Python 客户端请求 Go 启的服务端服务了。

### 总结

本文通过实战角度出发，直接用代码说话，来说明 gRPC 的一些应用。

内容包括简单的 gRPC 服务，流处理模式，验证器，Token 认证和证书认证。

除此之外，还有其他值得研究的内容，比如超时控制，REST 接口和负载均衡等。以后还会抽时间继续完善剩下这部分内容。

本文中的代码都经过测试验证，可以直接执行，并且已经上传到 [GitHub](https://github.com/yongxinz/go-example/tree/main/grpc-example)，小伙伴们可以一遍看源码，一遍对照文章内容来学习。

---

**源码地址：** 

- [https://github.com/yongxinz/go-example/tree/main/grpc-example](https://github.com/yongxinz/go-example/tree/main/grpc-example)
- [https://github.com/yongxinz/gopher/tree/main/blog](https://github.com/yongxinz/gopher/tree/main/blog)

关注公众号 **AlwaysBeta**，回复「**goebook**」领取 Go 编程经典书籍。

<center class="half">
    <img src="https://github.com/yongxinz/gopher/blob/main/alwaysbeta.JPG" width="300"/>
</center>

**往期文章：**

- [推荐三个实用的 Go 开发工具](https://mp.weixin.qq.com/s/3GLMLhegB3wF5_62mpmePA)
- [被 Docker 日志坑惨了](https://mp.weixin.qq.com/s/3Tkc15dTCEDUAZaZ88pcSQ)
- [使用 grpcurl 通过命令行访问 gRPC 服务](https://mp.weixin.qq.com/s/GShwcGCopXVmxCKnYf5FhA)
- [这个 TCP 问题你得懂：Cannot assign requested address](https://mp.weixin.qq.com/s/-cThzr5N2w3IEYYf-duCDA)

**参考文章：**

- [https://github.com/mwitkow/go-proto-validators](https://github.com/mwitkow/go-proto-validators)
- [https://github.com/Bingjian-Zhu/go-grpc-example](https://github.com/Bingjian-Zhu/go-grpc-example)
- [http://gaodongfei.com/archives/start-grpc](http://gaodongfei.com/archives/start-grpc)
- [https://liaoph.com/openssl-san/](https://liaoph.com/openssl-san/)
- [https://www.cnblogs.com/jackluo/p/13841286.html](https://www.cnblogs.com/jackluo/p/13841286.html)
