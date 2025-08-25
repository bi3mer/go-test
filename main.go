package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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

type listKeyMap struct {
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	makeProject      key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		makeProject: key.NewBinding(
			key.WithKeys("a", "+"),
			key.WithHelp("a/+", "make a project"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	testDirectory string
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
		}
	}

	return model{
		testDirectory: directory,
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
			// m.delegateKeys.remove.SetEnabled(true)
			// newItem := m.itemGenerator.next()
			// insCmd := m.list.InsertItem(0, newItem)
			// statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added " + newItem.Title()))
			// return m, tea.Batch(insCmd, statusCmd)
			return m, m.list.NewStatusMessage("Make project not implemented yet...")
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

func main() {
	// ===========================================================================
	// Set up and get user test directory
	// ===========================================================================
	const TEST_DIRECTORY = "gotestdir"
	directory := os.Getenv(TEST_DIRECTORY)
	if directory == "" {
		fmt.Println("Test directory must be set (e.g. `export gotestdir=~/tests`).")
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
	makeTemp(".")

	if _, err := tea.NewProgram(NewModel(directory), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
