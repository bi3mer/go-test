package main

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	makeProject      key.Binding
	renameProject    key.Binding
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
		renameProject: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename a project"),
		),
	}
}
