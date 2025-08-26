package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func makeTemp(cmd string) {
	dir, _ := os.UserHomeDir()
	err := os.WriteFile(filepath.Join(dir, ".gotest"), []byte(cmd), 0644)
	if err != nil {
		fmt.Printf("Failed to change directory: %s", err)
	}
}
