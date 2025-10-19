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

// moveTodo moves a todo up or down in priority (changes position)
func (m Model) moveTodo(direction string) tea.Cmd {
	return func() tea.Msg {
		// Load all todos fresh from storage
		allTodos, err := storage.LoadTodos()
		if err != nil {
			return todoMovedMsg{err: err}
		}

		if len(allTodos) < 2 {
			return todoMovedMsg{err: nil} // Nothing to move
		}

		// Sort same way as display and normalize positions
		sorted := helpers.SortTodosForDisplay(allTodos)
		sorted = helpers.NormalizeTodoPositions(sorted)

		if m.selectedTodo < 0 || m.selectedTodo >= len(sorted) {
			return todoMovedMsg{err: nil}
		}

		// Find target index
		targetIdx := m.selectedTodo
		if direction == "up" {
			targetIdx--
			if targetIdx < 0 {
				return todoMovedMsg{err: nil} // Already at top
			}
		} else { // down
			targetIdx++
			if targetIdx >= len(sorted) {
				return todoMovedMsg{err: nil} // Already at bottom
			}
		}

		// Swap positions
		sorted[m.selectedTodo].Position, sorted[targetIdx].Position = sorted[targetIdx].Position, sorted[m.selectedTodo].Position

		// Save both todos
		err = storage.SaveTodo(sorted[m.selectedTodo])
		if err != nil {
			return todoMovedMsg{err: err}
		}
		err = storage.SaveTodo(sorted[targetIdx])
		return todoMovedMsg{err: err}
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

// saveTodo saves a standalone todo and returns to dashboard
func (m Model) saveTodo() tea.Cmd {
	return func() tea.Msg {
		// Load existing todos to determine position
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

		// Set position for new todo (at the end)
		m.currentTodo.Position = maxPosition + 1

		// Save todo
		err = storage.SaveTodo(m.currentTodo)

		return saveCompleteMsg{err: err}
	}
}
