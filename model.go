package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	state         AppState
	testDirectory string
	textInput     textinput.Model
	list          list.Model
	keys          *listKeyMap
	delegateKeys  *delegateKeyMap
}

func NewModel(directory string) model {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	projects := generateProjects(directory)

	delegate := NewItemDelegate(delegateKeys, directory)
	projectsList := list.New(projects, delegate, 0, 0)
	projectsList.Title = "Test Projects"
	projectsList.Styles.Title = titleStyle
	projectsList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.makeProject,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
			listKeys.renameProject,
		}
	}

	textInput := textinput.New()
	textInput.Placeholder = "New project name..."
	textInput.CharLimit = 156
	textInput.Width = 60
	textInput.Prompt = ""

	// projectDatePrefix: fmt.Sprintf("%d-%d-%d-", t.Year(), t.Month(), t.Day()),

	return model{
		testDirectory: directory,
		textInput:     textInput,
		list:          projectsList,
		keys:          listKeys,
		delegateKeys:  delegateKeys,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch m.state {
	case StateList:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			h, v := appStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)

		case tea.KeyMsg:
			// Don't match any of the keys below if we're actively filtering.
			if m.list.FilterState() == list.Filtering {
				break
			}

			switch {
			case key.Matches(msg, m.keys.togglePagination):
				m.list.SetShowPagination(!m.list.ShowPagination())
				return m, nil

			case key.Matches(msg, m.keys.toggleHelpMenu):
				m.list.SetShowHelp(!m.list.ShowHelp())
				return m, nil

			case key.Matches(msg, m.keys.makeProject):
				m.state = StateAdd
				m.textInput.Placeholder = "New project name..."
				m.textInput.Focus()
				return m, nil

			case key.Matches(msg, m.keys.renameProject):
				m.textInput.Placeholder = "Rename project..."
				m.state = StateRename
				m.textInput.SetValue(m.list.SelectedItem().FilterValue())
				m.textInput.Focus()

				return m, nil
			}
		}

		newListModel, cmd := m.list.Update(msg)
		m.list = newListModel
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	case StateAdd:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				newProjectDirectory := filepath.Join(m.testDirectory, m.textInput.Value())
				os.Mkdir(newProjectDirectory, 0755)
				makeTemp(m.testDirectory, newProjectDirectory)
				return m, tea.Quit
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyEscape:
				m.state = StateList
				return m, nil
			}
		}

		newAddProject, cmd := m.textInput.Update(msg)
		m.textInput = newAddProject
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	case StateRename:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				userInput := strings.TrimSpace(m.textInput.Value())
				if len(userInput) == 0 {
					return m, tea.Batch(cmds...) // TODO: use status message?
				}

				oldName := filepath.Join(m.testDirectory, m.list.SelectedItem().FilterValue())
				newName := filepath.Join(m.testDirectory, userInput)
				os.Rename(oldName, newName)

				stats, err := os.Stat(newName)
				if err == nil {
					newItem := project{userInput, ""}
					m.list.Items()[m.list.Index()] = newItem
				} else {
					newItem := project{m.textInput.Value(), stats.ModTime().GoString()}
					m.list.Items()[m.list.Index()] = newItem
				}

				m.state = StateList

				return m, tea.Batch(cmds...)
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyEscape:
				m.state = StateList
				return m, tea.Batch(cmds...)
			}
		}

		newAddProject, cmd := m.textInput.Update(msg)
		m.textInput = newAddProject
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

	fmt.Println("Error: entered unknown app state: %i", m.state)
	return m, tea.Quit
}

func (m model) View() string {
	switch m.state {
	case StateList:
		return appStyle.Render(m.list.View())
	case StateAdd:
		return "Create Project\n\n  " + m.textInput.View()
	case StateRename:
		return "Rename Project\n\n  " + m.textInput.View()
	}

	return fmt.Sprintf("Error: entered unknown app state: %d\n", m.state)
}
