package main

import "github.com/apodacaa/amos/internal/models"

// saveCompleteMsg is sent when save operation completes
type saveCompleteMsg struct {
	entry models.Entry
	err   error
}

// finalizeCompleteMsg is sent when entry finalization completes (todo extraction on exit)
type finalizeCompleteMsg struct {
	entry   models.Entry
	err     error
	skipped bool // true if entry was empty and not saved
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

// statusTimeoutMsg is sent when status message should be cleared
type statusTimeoutMsg struct{}

// deleteCompleteMsg is sent when delete operation completes successfully
type deleteCompleteMsg struct {
	entryCount      int
	todoCount       int
	linkedTodoCount int
}

// deleteErrorMsg is sent when delete operation fails
type deleteErrorMsg struct {
	err     error
	message string
}

// configSavedMsg is sent when config save operation completes
type configSavedMsg struct {
	err error
}

// updateCheckCompleteMsg is sent when update check completes
type updateCheckCompleteMsg struct {
	latestVersion   string // Latest version from GitHub (e.g., "v1.5.0")
	updateAvailable bool   // Whether an update is available
}
