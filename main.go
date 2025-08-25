package main

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

const TEST_DIRECTORY = "gotestdir"

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func makeTemp(directory string) {
	err := os.WriteFile("temp", []byte(directory), 0644)
	if err != nil {
		fmt.Printf("Failed to change directory: %s", err)
	}
}

type Model struct {
	err                 *error
	project_date_prefix string
	directory           string
	text_input          textinput.Model
	projects            list.Model
}

func NewModel(directory string) Model {
	text_input := textinput.New()
	text_input.Placeholder = "New project name..."
	text_input.CharLimit = 156
	text_input.Width = 20
	text_input.Prompt = ""

	project_names := []list.Item{}
	entries, err := os.ReadDir(directory)
	if err != nil {
		fmt.Printf("Error reading test directory: %s", err.Error())
		os.Exit(1)
	}

	for _, e := range entries {
		if e.IsDir() {
			project_names = append(project_names, item(e.Name()))
		}
	}

	slices.SortFunc(project_names, func(a, b list.Item) int {
		return cmp.Compare(a.FilterValue(), b.FilterValue())
	})

	project_list := list.New(project_names, itemDelegate{}, 30, len(project_names)*3)
	project_list.FilterInput.Focus()

	t := time.Now()

	return Model{
		nil,
		fmt.Sprintf("%d-%d-%d-", t.Year(), t.Month(), t.Day()),
		directory,
		text_input,
		project_list,
	}
}

func (model Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if model.text_input.Focused() {
				// do something to make the directory and then change the directory
				// with some option for opening zed, nivm, etc. via an environment variable
				// model.directory = filepath.Join(model.directory, directory)
				return model, tea.Quit
			} else if model.projects.FilterInput.Focused() {
				i, ok := model.projects.SelectedItem().(item)
				if ok {
					makeTemp(filepath.Join(model.directory, string(i)))
				}

				return model, tea.Quit
			}

		case tea.KeyEscape, tea.KeyCtrlC:
			return model, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		model.err = &msg
		return model, nil
	}

	if model.text_input.Focused() {
		model.text_input, cmd = model.text_input.Update(msg)
		return model, cmd
	} else if model.projects.FilterInput.Focused() {
		model.projects, cmd = model.projects.Update(msg)
		return model, cmd
	}

	return model, nil
}

func (model Model) View() string {
	s := model.directory + "\n\n"
	if model.err == nil {
		// s += model.project_date_prefix
		// s += model.text_input.View() + "\n"
		s += model.projects.View()

	} else {
		s += "Error: " + (*model.err).Error()
	}

	return s
}

func main() {
	// ===========================================================================
	// Set up and get user test directory
	// ===========================================================================
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
	// Start bubbletea command line application
	// ===========================================================================
	makeTemp(".")
	model := NewModel(directory)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	fmt.Print("\033[H\033[2J") // clear the screen
}
