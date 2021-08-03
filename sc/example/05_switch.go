package main

import (
	"fmt"
	"time"
)

func main() {
	i := 2
	fmt.Print("write ", i, " as ")
	switch i {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two") // write 2 as two
		fallthrough
	case 3:
		fmt.Println("three") // three
	case 4, 5, 6:
		fmt.Println("four, five, six")
	}

	switch num := 9; num {
	case 1:
		fmt.Println("one")
	default:
		fmt.Println("nine") // nine
	}

	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("it's the weekend")
	default:
		fmt.Println("it's a weekday") // it's a weekday
	}

	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("it's before noon")
	default:
		fmt.Println("it's after noon") // it's after noon
	}
}
