package main

import (
	"database/sql"
	"fmt"
)

func foo() error {
	return fmt.Errorf("foo err, %v", sql.ErrNoRows)
}

func bar() error {
	return foo()
}

func main() {
	// err := errors.New("这是 errors.New() 创建的错误")
	// fmt.Printf("err 错误类型：%T，错误为：%v\n", err, err)

	err := bar()
	if err == sql.ErrNoRows {
		fmt.Printf("data not found, %+v\n", err)
		return
	}
	if err != nil {
		fmt.Println("Unknown error")
	}
}
