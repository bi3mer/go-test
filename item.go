package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

type project struct {
	name string
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

	for _, e := range entries {
		if e.IsDir() {
			projects = append(projects, project{e.Name()})
		}
	}

	slices.SortFunc(projects, func(a, b list.Item) int {
		return -strings.Compare(a.FilterValue(), b.FilterValue())
	})

	return projects
}
