![](https://github.com/yongxinz/gopher/blob/main/sc/pic/00_GO%E5%AE%89%E8%A3%85%E4%B8%8E%E9%85%8D%E7%BD%AE.png)

**原文链接：** [Go 专栏｜开发环境搭建以及开发工具 VS Code 配置](https://mp.weixin.qq.com/s/x1OW--3mwSTjgB2HaKGVVA)

Go 专栏的第一篇，想学 Go 的同学们，走起～

### Go 安装

我的个人电脑是 Mac，然后工作主要使用 Linux，所以在这里主要介绍在这两个系统下的安装。

**下载地址：**

- Go 官网下载地址：https://golang.org/dl/
- Go 官方镜像站（推荐）：https://golang.google.cn/dl/

![](https://github.com/yongxinz/gopher/blob/main/sc/pic/00_go_dl.png)

直接安装最新版本 go1.16.6，后续文章都会在此版本下开发，测试。

#### Mac 下安装

可以通过 `brew` 方式安装，也可以直接在官网下载可执行文件，然后双击安装包，不停下一步就可以了。

![](https://github.com/yongxinz/gopher/blob/main/sc/pic/00_mac_install_go.png)

#### Linux 下安装

下载安装包：
```shell
$ wget https://golang.google.cn/dl/go1.16.6.linux-amd64.tar.gz
```

解压到 `/usr/local` 目录：
```shell
$ sudo tar -zxvf go1.16.6.linux-amd64.tar.gz -C /usr/local
```

然后配置环境变量，打开 `$HOME/.bash_profile` 文件，增加下面两行代码：

```shell
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin
```

最后使环境变量生效：

```shell
$ source $HOME/.bash_profile
```

安装完成后，在终端执行查看版本命令，如果能正确输出版本信息，那就说明安装成功了。

```shell
$ go version
go version go1.16.6 linux/amd64
```

### 配置环境变量

`GOROOT` 和 `GOPATH` 都是环境变量，其中 `GOROOT` 是我们安装 Go 开发包的路径，`GOPATH` 会有一个默认目录。

由于 go1.11 之后使用 go mod 来管理依赖包，不再强制我们必须把代码写在 `GOPATH/src` 目录下，所以使用默认即可，无需修改。

默认 `GOPROXY` 配置是 `GOPROXY=https://proxy.golang.org,direct`，由于国内访问不到，所以我们需要换一个 PROXY，这里推荐使用：

1. https://goproxy.io
2. https://goproxy.cn

配置 `GOPROXY`：

```shell
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct
```

go mod 先这样配置就可以了，后续再来写文章详细介绍。

### 开发工具 VS Code

开发工具可以根据自己的喜好来，可以用 Goland，VS Code，当然 Vim 也可以。

我比较喜欢 VS Code，插件丰富，而且免费。

官方下载地址：https://code.visualstudio.com/Download

安装 Go 插件，并重启：

![](https://github.com/yongxinz/gopher/blob/main/sc/pic/00_vs_code_install_go.png)

### 第一个 Go 程序

好了，一切准备就绪，让我们开始 Hello World 吧。

```go
// 00_hello.go

package main  // 声明 main 包

import "fmt"  // 导入内置 fmt 包

func main(){  // main函数，程序执行入口
	fmt.Println("Hello World!")  // 在终端打印 Hello World!
}
```

使用 `go build` 命令编译：
```shell
$ go build 00_hello.go
$ ls
00_hello    00_hello.go go.mod
```

可以看到在目录下生成了可执行文件 `00_hello`，然后运行一下试试：
```shell
$ ./00_hello
Hello World!
```

成功输出！

还可以直接使用 `go run` 命令来执行代码，在调试的时候更加方便。
```shell
$ go run 00_hello.go
Hello World!
```

我可真厉害，又学会了一门编程语言。


---
文章中的脑图和源码都上传到了 GitHub，有需要的同学可自行下载。

**地址：** https://github.com/yongxinz/gopher/tree/main/sc

关注公众号 **AlwaysBeta**，回复「**goebook**」领取 Go 编程经典书籍。

<center class="half">
    <img src="https://github.com/yongxinz/gopher/blob/main/alwaysbeta.JPG" width="300"/>
</center>
