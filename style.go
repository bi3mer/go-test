package main

import "github.com/charmbracelet/lipgloss"

var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF10F0")).
	Background(lipgloss.Color("#333333"))

var selectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF10F0")).
	Padding(0)

var defaultStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#666666")).
	Padding(0)
