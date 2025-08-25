package main

import (
	"fmt"
	"os"
)

func makeTemp(directory string) {
	err := os.WriteFile("temp", []byte(directory), 0644)
	if err != nil {
		fmt.Printf("Failed to change directory: %s", err)
	}
}
