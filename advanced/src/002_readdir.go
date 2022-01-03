package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	dirName := "../"
	fileInfos, _ := ioutil.ReadDir(dirName)
	fmt.Println(len(fileInfos))
	for i := 0; i < len(fileInfos); i++ {
		fmt.Printf("%T\n", fileInfos[i])
		fmt.Println(i, fileInfos[i].Name(), fileInfos[i].IsDir())

	}
}
