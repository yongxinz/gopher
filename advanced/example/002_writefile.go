package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	fileName := "./text.txt"
	s := "Hello AlwaysBeta"
	err := ioutil.WriteFile(fileName, []byte(s), 0777)
	fmt.Println(err)
}
