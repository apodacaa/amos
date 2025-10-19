package main

import "github.com/apodacaa/amos/internal/models"

// saveCompleteMsg is sent when save operation completes
type saveCompleteMsg struct {
	err error
}

// entriesLoadedMsg is sent when entries are loaded
type entriesLoadedMsg struct {
	entries []models.Entry
	err     error
}

// todosLoadedMsg is sent when todos are loaded
type todosLoadedMsg struct {
	todos []models.Todo
	err   error
}

// todoToggledMsg is sent when todo status is toggled
type todoToggledMsg struct {
	err error
}

// todoMovedMsg is sent when todo position is changed
type todoMovedMsg struct {
	err error
}
