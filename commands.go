package main

import (
	"os/exec"
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/apodacaa/amos/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

// clearStatusAfterDelay returns a command that sends a statusTimeoutMsg after 3 seconds
func clearStatusAfterDelay() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return statusTimeoutMsg{}
	})
}

// loadEntries loads all entries from storage
func (m Model) loadEntries() tea.Cmd {
	return func() tea.Msg {
		entries, err := storage.LoadEntries()
		return entriesLoadedMsg{entries: entries, err: err}
	}
}

// loadTodos loads all todos from storage (async)
func (m Model) loadTodos() tea.Cmd {
	return func() tea.Msg {
		todos, err := storage.LoadTodos()
		return todosLoadedMsg{todos: todos, err: err}
	}
}

// loadEntriesAndTodos loads both entries and todos (for entry list view with todo stats)
func (m Model) loadEntriesAndTodos() tea.Cmd {
	return tea.Batch(m.loadEntries(), m.loadTodos())
}

// toggleTodoImmediate saves a todo immediately without reloading
func (m Model) toggleTodoImmediate(todo models.Todo) tea.Cmd {
	return func() tea.Msg {
		err := storage.SaveTodo(todo)
		// Return empty message - we don't reload, just save
		return todoToggledMsg{err: err}
	}
}

// launchEditor opens the external editor for entry creation/editing
// existingContent is the existing entry content (empty for new entries)
func launchEditor(existingContent string) tea.Cmd {
	// Create temp file with content or template
	tempFile, err := helpers.CreateTempFile(existingContent)
	if err != nil {
		return func() tea.Msg {
			return editorFinishedMsg{tempFile: "", err: err}
		}
	}

	// Get editor command
	editor := helpers.GetEditor()

	// Launch editor with ExecProcess (suspends TUI, returns control on exit)
	c := tea.ExecProcess(newEditorCmd(editor, tempFile), func(err error) tea.Msg {
		return editorFinishedMsg{tempFile: tempFile, err: err}
	})

	return c
}

// newEditorCmd creates an exec.Cmd for the given editor and file
func newEditorCmd(editor, filePath string) *exec.Cmd {
	return exec.Command(editor, filePath)
}

// saveTodo saves a standalone todo and returns to dashboard
func (m Model) saveTodo() tea.Cmd {
	return func() tea.Msg {
		// Save todo
		err := storage.SaveTodo(m.currentTodo)

		return saveCompleteMsg{err: err}
	}
}

// deleteMarkedCmd deletes all marked entries and todos
func (m Model) deleteMarkedCmd() tea.Cmd {
	return func() tea.Msg {
		// Separate entries and todos, count linked todos
		var entryIDs []string
		var todoIDs []string
		linkedTodoCount := 0

		for id, itemType := range m.markedForDeletion {
			if itemType == "entry" {
				entryIDs = append(entryIDs, id)
				// Count linked todos for this entry (single source of truth: use FilterTodosByEntry)
				linkedTodos := helpers.FilterTodosByEntry(m.todos, id)
				linkedTodoCount += len(linkedTodos)
			} else if itemType == "todo" {
				todoIDs = append(todoIDs, id)
			}
		}

		// Delete entries first (cascade deletes linked todos)
		for _, id := range entryIDs {
			err := storage.DeleteEntryCascade(id)
			if err != nil {
				return deleteErrorMsg{err: err, message: "Failed to delete entries"}
			}
		}

		// Delete standalone todos
		// Note: entry-linked todos will already be deleted by cascade
		for _, id := range todoIDs {
			// Try to delete - if it was already deleted by cascade, storage layer handles it
			storage.DeleteTodo(id)
		}

		return deleteCompleteMsg{
			entryCount:      len(entryIDs),
			todoCount:       len(todoIDs),
			linkedTodoCount: linkedTodoCount,
		}
	}
}

// switchWorkspaceCmd switches data directory and reloads entries/todos
func switchWorkspaceCmd(workspace storage.Workspace) tea.Cmd {
	return func() tea.Msg {
		if err := storage.SwitchWorkspace(workspace.DataDir); err != nil {
			return workspaceSwitchedMsg{err: err, name: workspace.Name}
		}
		entries, err := storage.LoadEntries()
		if err != nil {
			return workspaceSwitchedMsg{err: err, name: workspace.Name}
		}
		todos, err := storage.LoadTodos()
		if err != nil {
			return workspaceSwitchedMsg{err: err, name: workspace.Name}
		}
		return workspaceSwitchedMsg{
			entries: entries,
			todos:   todos,
			name:    workspace.Name,
		}
	}
}

// saveConfigCmd saves user configuration to system config
func saveConfigCmd(config storage.SystemConfig) tea.Cmd {
	return func() tea.Msg {
		err := storage.SaveSystemConfig(config)
		return configSavedMsg{err: err}
	}
}

// checkForUpdatesCmd performs asynchronous update check against GitHub
func checkForUpdatesCmd(currentVersion string) tea.Cmd {
	return func() tea.Msg {
		// Fetch latest version from GitHub API
		latestVersion := helpers.FetchLatestVersion()

		// If fetch failed or returned empty, return no update
		if latestVersion == "" {
			return updateCheckCompleteMsg{
				latestVersion:   "",
				updateAvailable: false,
			}
		}

		// Compare versions
		updateAvailable := helpers.IsUpdateAvailable(currentVersion, latestVersion)

		return updateCheckCompleteMsg{
			latestVersion:   latestVersion,
			updateAvailable: updateAvailable,
		}
	}
}
