package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	f, err := ioutil.TempFile("./", "Test")
	if err != nil {
		fmt.Println(err)
	}
	defer os.Remove(f.Name()) // 用完删除
	fmt.Printf("%s\n", f.Name())
}
