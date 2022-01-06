package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

func main() {
	//返回*strings.Reader
	reader := strings.NewReader("Hello AlwaysBeta")
	r := ioutil.NopCloser(reader)
	defer r.Close()

	fmt.Println(reflect.TypeOf(reader))
	data, _ := ioutil.ReadAll(reader)
	fmt.Println(string(data))
}
