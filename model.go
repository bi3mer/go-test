package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	directory string
	projects  []project
	cursor    int
}

func NewModel(directory string) model {
	return model{
		directory: directory,
		projects:  generateProjects(directory),
		cursor:    0,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.projects)-1 {
				m.cursor++
			}

		case "enter", " ":
			makeTemp(m.directory, m.projects[m.cursor].name)
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	// The header
	s := titleStyle.Render(" Test Projects ")
	s += "\n\n"

	// Iterate over our choices
	for i, p := range m.projects {
		if m.cursor == i {
			s += selectedStyle.Render("> "+p.name) + "\n"
		} else {
			s += defaultStyle.Render("  "+p.name) + "\n"
		}
	}

	s += "\nPress q to quit.\n"
	return s
}
