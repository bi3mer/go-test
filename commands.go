package main

import tea "github.com/charmbracelet/bubbletea"

type endMessage struct{}

func cmdEndSession(m model) tea.Cmd {
	return func() tea.Msg {
		saveProjects(m.projects, m.directory)
		return endMessage{}
	}
}
