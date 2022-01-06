package main

import (
	"database/sql"
	"errors"
	"fmt"
)

func bar() error {
	if err := foo(); err != nil {
		return fmt.Errorf("bar failed: %w", foo())
	}
	return nil
}

func foo() error {
	return fmt.Errorf("foo failed: %w", sql.ErrNoRows)
}

func main() {
	err := bar()
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("data not found,  %+v\n", err)
		return
	}
	if err != nil {
		fmt.Println("Unknown error")
	}
}
