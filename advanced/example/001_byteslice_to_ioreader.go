package main

import (
	"bytes"
	"fmt"
	"log"
)

func main() {
	data := []byte("Hello AlwaysBeta")

	// byte slice to bytes.Reader, which implements the io.Reader interface
	reader := bytes.NewReader(data)

	// read the data from reader
	buf := make([]byte, len(data))
	if _, err := reader.Read(buf); err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(buf))
}
