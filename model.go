package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	width         int
	height        int
	minIndex      int
	projectIndex  int
	projectOffset int
	state         AppState
	directory     string
	temp          string
	filterPattern string
	projects      []project
}

func NewModel(directory string) model {
	return model{
		projectOffset: 4,
		state:         StateList,
		directory:     directory,
		projects:      generateProjects(directory),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) ChangeState(state AppState) {
	m.state = state

	switch m.state {
	case StateAddProject, StateFilterList:
		m.projectOffset = 6
	default:
		m.projectOffset = 4
	}
}

func (m model) UpdateListState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, cmdEndSession(m)

		case "up", "k":
			if m.projectIndex > 0 {
				m.projectIndex--

				if m.minIndex > m.projectIndex {
					m.minIndex--
				}
			}

		case "down", "j":
			if m.projectIndex < len(m.projects)-1 {
				m.projectIndex++

				if m.minIndex+m.height-m.projectOffset-1 < m.projectIndex {
					m.minIndex++
				}
			}

		case "r", "R":
			m.temp = strings.Clone(m.projects[m.projectIndex].name)
			m.ChangeState(StateRenameProject)

		case "a", "A":
			m.temp = ""
			m.ChangeState(StateAddProject)

		case "/", "f", "F":
			m.ChangeState(StateFilterList)

		case "enter", " ":
			makeTemp(m.directory, m.projects[m.projectIndex].name)
			m.projects[m.projectIndex].time = time.Now()
			sortProjects(m.projects)
			m.projectIndex = 0

			return m, cmdEndSession(m)
		}
	}

	return m, nil
}

func (m model) UpdateRenameProjectState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			m.projects[m.projectIndex].name = m.temp
			m.ChangeState(StateList)
		case tea.KeyEnter:
			if len(m.projects[m.projectIndex].name) == 0 {
				m.projects[m.projectIndex].name = m.temp
			} else {
				err := os.Rename(
					filepath.Join(m.directory, m.temp),
					filepath.Join(m.directory, m.projects[m.projectIndex].name),
				)

				if err != nil {
					m.projects[m.projectIndex].name = m.temp
				} else {
					m.projects[m.projectIndex].time = time.Now()
					sortProjects(m.projects)
					m.projectIndex = 0
				}
			}

			m.ChangeState(StateList)
		case tea.KeyBackspace:
			length := len(m.projects[m.projectIndex].name)
			if length > 0 {
				m.projects[m.projectIndex].name = m.projects[m.projectIndex].name[:length-1]
			}
		case tea.KeyRunes:
			m.projects[m.projectIndex].name += string(msg.Runes)
		}
	}

	return m, nil
}

func (m model) UpdateAddProjectState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			m.ChangeState(StateList)
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

			m.ChangeState(StateList)
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

// simple filter version: https://github.com/forrestthewoods/lib_fts/blob/master/code/fts_fuzzy_match.h
//
// This works, but it is not as good as it could be and, more importantly, the project list view is messed up
func (m *model) filter() {
	for i := range len(m.projects) {
		name := m.projects[i].name
		indexP := 0
		indexS := 0

		for indexP < len(m.filterPattern) && indexS < len(name) {
			if unicode.ToLower(rune(m.filterPattern[indexP])) == unicode.ToLower(rune(name[indexS])) {
				indexP++
			}

			indexS++
		}

		m.projects[i].visible = indexP == len(m.filterPattern)
	}
}

func (m model) UpdateFilterState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			for i := range len(m.projects) {
				m.projects[i].visible = true
			}

			m.ChangeState(StateList)
		case tea.KeyEnter:
			m.ChangeState(StateList)
		case tea.KeyBackspace:
			length := len(m.filterPattern)
			if length > 0 {
				m.filterPattern = m.filterPattern[:length-1]
			}

			m.filter()
		case tea.KeyRunes:
			m.filterPattern += string(msg.Runes)
			m.filter()
		}
	}

	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, cmdEndSession(m)
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case endMessage:
		return m, tea.Quit
	}

	switch m.state {
	case StateList:
		return m.UpdateListState(msg)
	case StateAddProject:
		return m.UpdateAddProjectState(msg)
	case StateRenameProject:
		return m.UpdateRenameProjectState(msg)
	case StateFilterList:
		return m.UpdateFilterState(msg)
	}

	return m, cmdEndSession(m)
}

func (m model) View() string {
	if m.height < 8 {
		return errorStyle.Render("Please make your terminal taller...") + "\n"
	}

	if m.width < 20 {
		return errorStyle.Render("Please make your terminal wider...") + "\n"
	}

	// Title
	s := titleStyle.Render(" Test Projects ")
	s += "\n\n"

	// render add project if necessary
	switch m.state {
	case StateAddProject:
		s += selectedStyle.Render("Add Project: ") + renameStyle.Render(m.temp)
		s += "\n\n"
	case StateFilterList:
		s += selectedStyle.Render("Filter: ") + renameStyle.Render(m.filterPattern)
		s += "\n\n"
	}

	loopMax := min(m.minIndex+m.height-m.projectOffset, len(m.projects))
	for i := m.minIndex; i < loopMax; i++ {
		p := &m.projects[i]
		if !p.visible {
			continue
		}

		if m.projectIndex == i {
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
	s += "\nq quit, a add project, r rename project, / to filter"
	return s
}
