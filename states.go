package main

type AppState int

const (
	StateList = iota
	StateAdd
	StateRename
	StateFilter
)
