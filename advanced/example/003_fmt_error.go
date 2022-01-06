package main

import (
	"fmt"
)

func main() {
	err := fmt.Errorf("这个是 fmt.Errorf() 创建的错误，错误编码为：%d", 500)
	fmt.Printf("err2 错误类型：%T，错误为：%v\n", err, err)
}
