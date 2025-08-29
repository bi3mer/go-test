package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	state     AppState
	directory string
	temp      string
	projects  []project
	cursor    int
}

func NewModel(directory string) model {
	return model{
		state:     StateList,
		directory: directory,
		projects:  generateProjects(directory),
		cursor:    0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) UpdateListState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, endSession(m)

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.projects)-1 {
				m.cursor++
			}

		case "r", "R":
			m.temp = strings.Clone(m.projects[m.cursor].name)
			m.state = StateRename

		case "a", "A":
			panic("Add project not implemented!")

		case "/", "f", "F":
			panic("Filter not yet implemented!")

		case "enter", " ":
			makeTemp(m.directory, m.projects[m.cursor].name)
			m.projects[m.cursor].time = time.Now()
			sortProjects(m.projects)
			m.cursor = 0

			return m, endSession(m)
		}

	case endMessage:
		return m, tea.Quit
	}

	return m, nil
}

func (m model) UpdateRenameState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEscape:
			m.projects[m.cursor].name = m.temp
			m.state = StateList
		case tea.KeyEnter:
			if len(m.projects[m.cursor].name) == 0 {
				m.projects[m.cursor].name = m.temp
			} else {
				err := os.Rename(
					filepath.Join(m.directory, m.temp),
					filepath.Join(m.directory, m.projects[m.cursor].name),
				)

				if err != nil {
					m.projects[m.cursor].name = m.temp
				} else {
					m.projects[m.cursor].time = time.Now()
					sortProjects(m.projects)
					m.cursor = 0
				}
			}

			m.state = StateList
		case tea.KeyBackspace:
			length := len(m.projects[m.cursor].name)
			if length > 0 {
				m.projects[m.cursor].name = m.projects[m.cursor].name[:length-1]
			}
		case tea.KeyRunes:
			m.projects[m.cursor].name += string(msg.Runes)
		}
	}

	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateList:
		return m.UpdateListState(msg)
	case StateAdd:
		return m, nil
	case StateRename:
		return m.UpdateRenameState(msg)
	}

	return m, tea.Quit
}

func (m model) View() string {
	// Title
	s := titleStyle.Render(" Test Projects ")
	s += "\n\n"

	for i, p := range m.projects {
		if !p.visible {
			continue
		}

		if m.cursor == i {
			switch m.state {
			case StateList:
				s += selectedStyle.Render("> "+p.name) + "\n"
			case StateRename:
				s += selectedStyle.Render("> ") + renameStyle.Render(p.name) + "\n"
			case StateAdd:
				s += defaultStyle.Render("> "+p.name) + "\n"

			}
		} else {
			s += defaultStyle.Render("  "+p.name) + "\n"
		}
	}

	// header
	s += "\nq quit, a add project, r rename project, / to filter\n"
	return s
}
