package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(0, 0)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 0)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

func main() {
	// ===========================================================================
	// Set up and get user test directory
	// ===========================================================================
	const TEST_DIRECTORY = "gotestdir"
	directory := os.Getenv(TEST_DIRECTORY)
	if directory == "" {
		fmt.Println("Test directory must be set (e.g. `export gotestdir=\"~/tests\"`).")
		return
	}

	if strings.HasPrefix(directory, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Unexpected error getting user home dir: ", err)
			return
		}

		directory = filepath.Join(home, directory[2:])
	}

	if os.MkdirAll(directory, os.ModePerm) != nil {
		fmt.Print(`Error automatically creating the directory. Update the environemnt"
variable 'gotestdir' and/or make the directory yourself.`)

		return
	}

	// ===========================================================================
	// Start the app
	// ===========================================================================
	makeTemp(directory, "")
	if _, err := tea.NewProgram(NewModel(directory), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
