package main

import "fmt"

func Add(x, y int, ch chan int) {
	z := x + y
	ch <- z
}

func counter(out chan<- int) {
	for x := 0; x < 10; x++ {
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for v := range in {
		out <- v * v
	}
	close(out)
}

func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

func main() {
	n := make(chan int)
	s := make(chan int)

	go counter(n)
	go squarer(s, n)
	printer(s)

	ch := make(chan int)
	for i := 0; i < 10; i++ {
		go Add(i, i, ch)
	}

	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}
}
