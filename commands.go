package main

import (
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/apodacaa/amos/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
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

// saveEntry saves the current entry and extracts todos
func (m Model) saveEntry() tea.Cmd {
	return func() tea.Msg {
		content := m.textarea.Value()

		// Parse content into title and body
		title, body := helpers.ParseEntryContent(content)

		// Extract tags from title and body
		tags := helpers.ExtractTags(title + " " + body)

		// Extract todos from content
		todoTitles := helpers.ExtractTodos(content)

		// Load existing todos for this entry to avoid duplicates
		existingTodosByTitle := make(map[string]string) // title -> ID
		allTodos, err := storage.LoadTodos()
		if err == nil {
			// Build map of existing todos for this entry
			for _, existingID := range m.currentEntry.TodoIDs {
				for _, todo := range allTodos {
					if todo.ID == existingID {
						existingTodosByTitle[todo.Title] = todo.ID
						break
					}
				}
			}
		}

		// Create todo IDs list (preserve existing + add new)
		todoIDs := make([]string, 0, len(todoTitles))

		// Create and save only NEW todos
		for _, todoTitle := range todoTitles {
			// Check if this todo already exists
			if existingID, exists := existingTodosByTitle[todoTitle]; exists {
				// Reuse existing todo ID
				todoIDs = append(todoIDs, existingID)
			} else {
				// Create new todo
				todo := models.Todo{
					ID:        uuid.New().String(),
					Title:     todoTitle,
					Status:    "open",
					Tags:      helpers.ExtractTags(todoTitle), // Extract tags from todo title
					CreatedAt: time.Now(),
					EntryID:   &m.currentEntry.ID, // Link to this entry
				}

				// Save new todo
				if err := storage.SaveTodo(todo); err != nil {
					return saveCompleteMsg{err: err}
				}

				todoIDs = append(todoIDs, todo.ID)
			}
		}

		// Update current entry
		m.currentEntry.Title = title
		m.currentEntry.Body = body
		m.currentEntry.Tags = tags
		m.currentEntry.TodoIDs = todoIDs
		m.currentEntry.Timestamp = time.Now()

		// Save entry to storage
		err = storage.SaveEntry(m.currentEntry)

		return saveCompleteMsg{entry: m.currentEntry, err: err}
	}
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
				// Count linked todos for this entry
				for _, e := range m.entries {
					if e.ID == id {
						linkedTodoCount += len(e.TodoIDs)
						break
					}
				}
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

// saveConfigCmd saves user configuration
func saveConfigCmd(config models.Config) tea.Cmd {
	return func() tea.Msg {
		err := storage.SaveConfig(config)
		return configSavedMsg{err: err}
	}
}
