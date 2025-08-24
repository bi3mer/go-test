package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
)

type Model struct {
	project_date_prefix string
	directory           string
	text_input          textinput.Model
}

func NewModel() Model {
	text_input := textinput.New()
	text_input.Placeholder = "New project name..."
	text_input.Focus()
	text_input.CharLimit = 156
	text_input.Width = 20
	text_input.Prompt = ""

	t := time.Now()

	return Model{
		fmt.Sprintf("%d-%d-%d-", t.Year(), t.Month(), t.Day()),
		"",
		text_input,
	}
}

func (model Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		return m, nil
	}

	m.text_input, cmd = m.text_input.Update(msg)
	return m, cmd
}

func (model Model) View() string {
	s := model.directory + "\n\n"
	s += model.project_date_prefix
	s += model.text_input.View() + "\n"

	return s
}

func main() {
	model := NewModel()

	// ===========================================================================
	// Set up and get user test directory
	// ===========================================================================
	model.directory = os.Getenv(TEST_DIRECTORY)
	if model.directory == "" {
		fmt.Println("Test directory must be set (e.g. `export gotestdir=~/tests`).")
		return
	}

	if strings.HasPrefix(model.directory, "~/") {
		home, error := os.UserHomeDir()
		if error != nil {
			fmt.Println("Unexpected error getting user home dir: ", error)
			return
		}

		model.directory = filepath.Join(home, model.directory[2:])
	}

	if os.MkdirAll(model.directory, os.ModePerm) != nil {
		fmt.Print(`Error automatically creating the directory. Update the environemnt"
variable 'gotestdir' and/or make the directory yourself.`)

		return
	}

	// ===========================================================================
	// Start bubbletea command line application
	// ===========================================================================
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
