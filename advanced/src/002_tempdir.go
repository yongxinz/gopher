package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	dir, err := ioutil.TempDir("./", "Test")
	if err != nil {
		fmt.Println(err)
	}
	defer os.Remove(dir) // 用完删除
	fmt.Printf("%s\n", dir)
}
