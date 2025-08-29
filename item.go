package main

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/list"
)

type project struct {
	name   string
	prefix string
}

func (i project) Title() string {
	return i.name
}

func (i project) Description() string { return "" }

func (i project) FilterValue() string {
	return i.name
}

func generateProjects(directory string) []list.Item {
	projects := []list.Item{}

	entries, err := os.ReadDir(directory)
	if err != nil {
		fmt.Printf("Error reading test directory: %s", err.Error())
		os.Exit(1)
	}

	modificationDate := make(map[string]time.Time)

	for _, e := range entries {
		if e.IsDir() {
			if stats, err := os.Stat(directory); err != nil {
				fmt.Fprintf(os.Stderr, "Error getting directory stats %s: %v", directory, err)
			} else {
				modificationDate[e.Name()] = stats.ModTime()
				projects = append(projects, project{e.Name(), stats.ModTime().GoString()})
			}
		}
	}

	slices.SortFunc(projects, func(a, b list.Item) int {
		return -modificationDate[a.FilterValue()].Compare(modificationDate[b.FilterValue()])
	})

	return projects
}
