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
			return m, cmdEndSession(m)

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
			m.state = StateRenameProject

		case "a", "A":
			m.state = StateAddProject
			m.temp = ""

		case "/", "f", "F":
			panic("Filter not yet implemented!")

		case "enter", " ":
			makeTemp(m.directory, m.projects[m.cursor].name)
			m.projects[m.cursor].time = time.Now()
			sortProjects(m.projects)
			m.cursor = 0

			return m, cmdEndSession(m)
		}

	case endMessage:
		return m, tea.Quit
	}

	return m, nil
}

func (m model) UpdateRenameProjectState(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m model) UpdateAddProjectState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEscape:
			m.state = StateList
		case tea.KeyEnter:
			err := os.Mkdir(filepath.Join(m.directory, m.temp), 0755)
			if err == nil {
				m.projects = append(m.projects, project{
					name:    m.temp,
					time:    time.Now(),
					visible: true,
				})

				sortProjects(m.projects)
			}

			m.state = StateList
		case tea.KeyBackspace:
			length := len(m.temp)
			if length > 0 {
				m.temp = m.temp[:length-1]
			}
		case tea.KeyRunes:
			m.temp += string(msg.Runes)
		}
	}

	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateList:
		return m.UpdateListState(msg)
	case StateAddProject:
		return m.UpdateAddProjectState(msg)
	case StateRenameProject:
		return m.UpdateRenameProjectState(msg)
	case StateFilterList:
		return m, tea.Quit
	}

	return m, tea.Quit
}

func (m model) View() string {
	// Title
	s := titleStyle.Render(" Test Projects ")
	s += "\n\n"

	if m.state == StateAddProject {
		s += selectedStyle.Render("Add Project: ") + renameStyle.Render(m.temp)
		s += "\n\n"
	}

	// TODO: add paging so only so list of projects does not extend past the terminal view
	for i, p := range m.projects {
		if !p.visible {
			continue
		}

		if m.cursor == i {
			switch m.state {
			case StateList:
				s += selectedStyle.Render("> "+p.name) + "\n"
			case StateRenameProject:
				s += selectedStyle.Render("> ") + renameStyle.Render(p.name) + "\n"
			case StateAddProject:
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
