package main

import (
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("content/hello.md")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(string(data))
}
