package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	ioReaderData := strings.NewReader("Hello AlwaysBeta")
	fmt.Printf("%T", ioReaderData)

	// creates a bytes.Buffer and read from io.Reader
	buf := &bytes.Buffer{}
	buf.ReadFrom(ioReaderData)

	// retrieve a byte slice from bytes.Buffer
	data := buf.Bytes()

	// only read the left bytes from 6
	fmt.Println(string(data[6:]))
}
