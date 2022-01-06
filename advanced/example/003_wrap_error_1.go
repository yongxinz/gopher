package main

import (
	"fmt"
)

type myError struct{}

func (e myError) Error() string {
	return "Error happended"
}

func main() {
	e1 := myError{}
	e2 := fmt.Errorf("E2: %w", e1)
	e3 := fmt.Errorf("E3: %w", e2)
	fmt.Println(e2)
	fmt.Println(e3)
}
