package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

type project struct {
	name    string
	time    time.Time
	visible bool
}

func generateProjects(directory string) []project {
	projects := []project{}

	// read .gotestdb to find projects already made
	dbPath := filepath.Join(directory, ".gotestdb")
	if stat, statErr := os.Stat(dbPath); statErr == nil && !stat.IsDir() {
		dbFile, err := os.Open(dbPath)

		if err == nil {
			defer dbFile.Close()
			scanner := bufio.NewScanner(dbFile)
			for scanner.Scan() {
				lineData := strings.Split(scanner.Text(), ",")
				projectTime, timeErr := time.Parse(time.StampNano, lineData[1])
				if timeErr == nil {
					projects = append(projects, project{
						name:    lineData[0],
						time:    projectTime,
						visible: false,
					})
				}
			}
		}
	}

	// read the directory to see if the user has made any projects without the cli and
	// add them to the list for future use
	entries, err := os.ReadDir(directory)
	if err != nil {
		fmt.Printf("Error reading test directory: %s", err.Error())
		os.Exit(1)
	}

	for _, e := range entries {
		projectName := e.Name()
		if e.IsDir() {
			found := false
			for i := 0; i < len(projects); i++ {
				p := &projects[i]
				if p.name == projectName {
					p.visible = true
					found = true
					break
				}
			}

			if !found {
				projects = append(projects, project{
					name:    projectName,
					time:    time.Now(),
					visible: true,
				})
			}
		}
	}

	// remove any project where visible is false, because this means that there is no
	// corresponding directory in the users test directory
	filtered := projects[:0]
	for _, p := range projects {
		if p.visible {
			filtered = append(filtered, p)
		}
	}

	sortProjects(filtered)
	return filtered
}

func saveProjects(projects []project, directory string) {
	file, err := os.Create(filepath.Join(directory, ".gotestdb"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating .gotestdb: %v", err)
	}
	defer file.Close()

	for _, p := range projects {
		file.WriteString(fmt.Sprintf("%s,%s\n", p.name, p.time.Format(time.StampNano))) // skip error
	}
}

func sortProjects(projects []project) {
	slices.SortFunc(projects, func(a, b project) int {
		return -a.time.Compare(b.time)
	})
}
