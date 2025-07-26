package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Testing pre-commit hooks")
}

func testFunction() {
	x := 1
	y := 2
	fmt.Printf("x: %d, y: %d\n", x, y)
}
