package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	data, err := ioutil.ReadFile("./text.txt")
	if err != nil {
		fmt.Println("read error")
		os.Exit(1)
	}
	fmt.Println(string(data))
}
