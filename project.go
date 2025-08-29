package main

import (
	"fmt"
	"os"
	"slices"
	"time"
)

type project struct {
	name string
	time time.Time
}

func generateProjects(directory string) []project {
	projects := []project{}

	entries, err := os.ReadDir(directory)
	if err != nil {
		fmt.Printf("Error reading test directory: %s", err.Error())
		os.Exit(1)
	}

	for _, e := range entries {
		if e.IsDir() {
			projects = append(projects, project{e.Name(), time.Now()})
		}
	}

	slices.SortFunc(projects, func(a, b project) int {
		return a.time.Compare(b.time)
	})

	return projects
}

func saveProjects(projects []project, directory string) {
	// @TODO: save to $directory/.gotestdb
}
