package main

func main() {
	type Test struct{}
	v := Test{}
	Print(v)
}

func Print(v interface{}) {
	println(v)
}
