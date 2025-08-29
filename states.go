package main

type AppState int

const (
	StateList = iota
	StateAddProject
	StateRenameProject
	StateFilterList
)
