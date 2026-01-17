package main

import (
	"strings"
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

// saveEntry saves the current entry (keeps !todo markup intact)
func (m Model) saveEntry() tea.Cmd {
	return func() tea.Msg {
		content := m.textarea.Value()

		// Parse content into title and body
		title, body := helpers.ParseEntryContent(content)

		// Extract tags from title and body (include !todo lines)
		tags := helpers.ExtractTags(title + " " + body)

		// Update current entry (keep body as-is with !todo markup)
		m.currentEntry.Title = title
		m.currentEntry.Body = body // Keep !todo markup - extraction happens on exit
		m.currentEntry.Tags = tags
		m.currentEntry.UpdatedAt = time.Now() // Only update UpdatedAt, preserve CreatedAt

		// Save entry to storage
		err := storage.SaveEntry(m.currentEntry)

		return saveCompleteMsg{entry: m.currentEntry, err: err}
	}
}

// finalizeAndExit extracts todos, cleans entry text, and exits entry form
// This is called when the user presses ESC to exit the entry form
func (m Model) finalizeAndExit() tea.Cmd {
	return func() tea.Msg {
		// Always use m.currentEntry as source (last saved version)
		entry := m.currentEntry

		// Skip saving if entry has no meaningful content
		// This handles the case where user creates a new entry and immediately exits
		hasContent := strings.TrimSpace(entry.Title) != "" || strings.TrimSpace(entry.Body) != ""
		if !hasContent {
			return finalizeCompleteMsg{entry: entry, err: nil, skipped: true}
		}

		content := entry.Title + "\n" + entry.Body

		// Extract todos from SAVED content BEFORE cleaning
		todoTitles := helpers.ExtractTodos(content)

		// Remove !todo markup from body
		cleanedBody := helpers.RemoveTodoMarkup(entry.Body)

		// Extract tags from title and cleaned body
		tags := helpers.ExtractTags(entry.Title + " " + cleanedBody)

		// Create todos from extracted !todo lines
		for _, todoTitle := range todoTitles {
			now := time.Now()
			todo := models.Todo{
				ID:        uuid.New().String(),
				Title:     todoTitle,
				Status:    "open",
				Tags:      helpers.ExtractTags(todoTitle),
				CreatedAt: now,
				UpdatedAt: now,
				EntryID:   &entry.ID,
			}

			if err := storage.SaveTodo(todo); err != nil {
				return finalizeCompleteMsg{err: err}
			}
		}

		// Update entry with cleaned body
		entry.Body = cleanedBody
		entry.Tags = tags
		entry.UpdatedAt = time.Now()

		// Save cleaned entry
		err := storage.SaveEntry(entry)

		return finalizeCompleteMsg{entry: entry, err: err}
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
