package main

import (
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/apodacaa/amos/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

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

// toggleTodo toggles the status of the currently selected todo
func (m Model) toggleTodo() tea.Cmd {
	return func() tea.Msg {
		// Get the sorted todo to toggle (same sorting as display)
		sorted := make([]models.Todo, len(m.todos))
		copy(sorted, m.todos)
		// Sort: open before done, then by position, then newest first
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i].Status == "done" && sorted[j].Status == "open" {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				} else if sorted[i].Status == sorted[j].Status {
					if sorted[i].Position != sorted[j].Position {
						if sorted[j].Position < sorted[i].Position {
							sorted[i], sorted[j] = sorted[j], sorted[i]
						}
					} else {
						if sorted[j].CreatedAt.After(sorted[i].CreatedAt) {
							sorted[i], sorted[j] = sorted[j], sorted[i]
						}
					}
				}
			}
		}

		if m.selectedTodo < 0 || m.selectedTodo >= len(sorted) {
			return todoToggledMsg{err: nil}
		}

		todoToToggle := sorted[m.selectedTodo]

		// Toggle status
		if todoToToggle.Status == "open" {
			todoToToggle.Status = "done"
		} else {
			todoToToggle.Status = "open"
		}

		// Save the updated todo
		err := storage.SaveTodo(todoToToggle)
		return todoToggledMsg{err: err}
	}
}

// moveTodo moves a todo up or down in priority (changes position)
func (m Model) moveTodo(direction string) tea.Cmd {
	return func() tea.Msg {
		// Get the sorted todos (same sorting as display)
		sorted := make([]models.Todo, len(m.todos))
		copy(sorted, m.todos)
		// Sort: open before done, then by position, then newest first
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i].Status == "done" && sorted[j].Status == "open" {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				} else if sorted[i].Status == sorted[j].Status {
					if sorted[i].Position != sorted[j].Position {
						if sorted[j].Position < sorted[i].Position {
							sorted[i], sorted[j] = sorted[j], sorted[i]
						}
					} else {
						if sorted[j].CreatedAt.After(sorted[i].CreatedAt) {
							sorted[i], sorted[j] = sorted[j], sorted[i]
						}
					}
				}
			}
		}

		if m.selectedTodo < 0 || m.selectedTodo >= len(sorted) {
			return todoMovedMsg{err: nil}
		}

		currentTodo := sorted[m.selectedTodo]

		// Determine target position based on direction
		var targetIdx int
		if direction == "up" {
			targetIdx = m.selectedTodo - 1
			if targetIdx < 0 {
				return todoMovedMsg{err: nil} // Already at top
			}
		} else { // down
			targetIdx = m.selectedTodo + 1
			if targetIdx >= len(sorted) {
				return todoMovedMsg{err: nil} // Already at bottom
			}
		}

		targetTodo := sorted[targetIdx]

		// Swap positions
		currentTodo.Position, targetTodo.Position = targetTodo.Position, currentTodo.Position

		// Save both todos
		err := storage.SaveTodo(currentTodo)
		if err != nil {
			return todoMovedMsg{err: err}
		}
		err = storage.SaveTodo(targetTodo)
		return todoMovedMsg{err: err}
	}
}

// deleteEntry deletes the currently selected entry
func (m Model) deleteEntry() tea.Cmd {
	return func() tea.Msg {
		// Get the sorted entry to delete
		sorted := make([]models.Entry, len(m.entries))
		copy(sorted, m.entries)
		// Sort by timestamp descending (newest first)
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[j].Timestamp.After(sorted[i].Timestamp) {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}

		if m.selectedEntry < 0 || m.selectedEntry >= len(sorted) {
			return entryDeletedMsg{err: nil}
		}

		entryToDelete := sorted[m.selectedEntry]

		// Load all entries
		entries, err := storage.LoadEntries()
		if err != nil {
			return entryDeletedMsg{err: err}
		}

		// Find and remove the entry
		newEntries := make([]models.Entry, 0, len(entries)-1)
		for _, e := range entries {
			if e.ID != entryToDelete.ID {
				newEntries = append(newEntries, e)
			}
		}

		// Save updated entries
		err = storage.SaveEntries(newEntries)
		return entryDeletedMsg{err: err}
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

		// Create todo IDs list
		todoIDs := make([]string, 0, len(todoTitles))

		// Load existing todos to determine max position
		existingTodos, err := storage.LoadTodos()
		if err != nil {
			return saveCompleteMsg{err: err}
		}

		// Find max position
		maxPosition := 0
		for _, t := range existingTodos {
			if t.Position > maxPosition {
				maxPosition = t.Position
			}
		}

		// Create and save todos
		for idx, todoTitle := range todoTitles {
			todo := models.Todo{
				ID:        uuid.New().String(),
				Title:     todoTitle,
				Status:    "open",
				Tags:      helpers.ExtractTags(todoTitle), // Extract tags from todo title
				CreatedAt: time.Now(),
				EntryID:   &m.currentEntry.ID,    // Link to this entry
				Position:  maxPosition + idx + 1, // Append to end
			}

			// Save each todo
			if err := storage.SaveTodo(todo); err != nil {
				return saveCompleteMsg{err: err}
			}

			todoIDs = append(todoIDs, todo.ID)
		}

		// Update current entry
		m.currentEntry.Title = title
		m.currentEntry.Body = body
		m.currentEntry.Tags = tags
		m.currentEntry.TodoIDs = todoIDs
		m.currentEntry.Timestamp = time.Now()

		// Save entry to storage
		err = storage.SaveEntry(m.currentEntry)

		return saveCompleteMsg{err: err}
	}
}
