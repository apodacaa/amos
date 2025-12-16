package main

import (
	"testing"

	"github.com/apodacaa/amos/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

// TestNavigationFromEntriesToEntryView tests pressing Enter in entries list
func TestNavigationFromEntriesToEntryView(t *testing.T) {
	m := NewModel()
	m.view = "entries"
	m.entries = []models.Entry{
		{ID: "1", Title: "Test Entry", Body: "Test body"},
		{ID: "2", Title: "Second Entry", Body: "Second body"},
	}
	m.selectedEntry = 0

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newModel.(Model)

	if m.view != "view_entry" {
		t.Errorf("Expected view to be 'view_entry', got '%s'", m.view)
	}

	if m.viewingEntry.ID != "1" {
		t.Errorf("Expected viewingEntry ID to be '1', got '%s'", m.viewingEntry.ID)
	}
}

// TestNavigationFromEntryViewToEntriesList tests pressing Esc in entry view
func TestNavigationFromEntryViewToEntriesList(t *testing.T) {
	m := NewModel()
	m.view = "view_entry"
	m.viewingEntry = models.Entry{ID: "1", Title: "Test Entry"}

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = newModel.(Model)

	if m.view != "entries" {
		t.Errorf("Expected view to be 'entries', got '%s'", m.view)
	}
}

// TestNavigationToNewEntry tests pressing 'n' to create new entry
func TestNavigationToNewEntry(t *testing.T) {
	testCases := []struct {
		name     string
		fromView string
	}{
		{"from entries list", "entries"},
		{"from entry view", "view_entry"},
		{"from todo list", "todos"},
		{"from todo view", "view_todo"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewModel()
			m.view = tc.fromView

			newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			m = newModel.(Model)

			if m.view != "entry" {
				t.Errorf("Expected view to be 'entry', got '%s'", m.view)
			}

			if m.editingMode {
				t.Error("Expected editingMode to be false for new entry")
			}
		})
	}
}

// TestNavigationToTodoList tests pressing 't' to jump to todos
func TestNavigationToTodoList(t *testing.T) {
	testCases := []struct {
		name     string
		fromView string
	}{
		{"from entries list", "entries"},
		{"from entry view", "view_entry"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewModel()
			m.view = tc.fromView

			newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
			m = newModel.(Model)

			if m.view != "todos" {
				t.Errorf("Expected view to be 'todos', got '%s'", m.view)
			}
		})
	}
}

// TestNavigationToEntriesList tests pressing 'e' to jump to entries
func TestNavigationToEntriesList(t *testing.T) {
	testCases := []struct {
		name     string
		fromView string
	}{
		{"from todo list", "todos"},
		{"from todo view", "view_todo"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewModel()
			m.view = tc.fromView

			newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			m = newModel.(Model)

			if m.view != "entries" {
				t.Errorf("Expected view to be 'entries', got '%s'", m.view)
			}
		})
	}
}

// TestNavigationToFilter tests pressing '/' to open filter
func TestNavigationToFilter(t *testing.T) {
	testCases := []struct {
		name            string
		fromView        string
		expectedContext string
	}{
		{"from entries list", "entries", "entries"},
		{"from todo list", "todos", "todos"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewModel()
			m.view = tc.fromView

			newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
			m = newModel.(Model)

			if m.view != "unified_filter" {
				t.Errorf("Expected view to be 'unified_filter', got '%s'", m.view)
			}

			if m.filterContext != tc.expectedContext {
				t.Errorf("Expected filterContext to be '%s', got '%s'", tc.expectedContext, m.filterContext)
			}
		})
	}
}

// TestNavigationToHelp tests pressing '?' to open help
func TestNavigationToHelp(t *testing.T) {
	testCases := []struct {
		name     string
		fromView string
	}{
		{"from entries list", "entries"},
		{"from entry view", "view_entry"},
		{"from todo list", "todos"},
		{"from todo view", "view_todo"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewModel()
			m.view = tc.fromView

			newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			m = newModel.(Model)

			if m.view != "help" {
				t.Errorf("Expected view to be 'help', got '%s'", m.view)
			}

			if m.previousView != tc.fromView {
				t.Errorf("Expected previousView to be '%s', got '%s'", tc.fromView, m.previousView)
			}
		})
	}
}

// TestNavigationWithEmptyEntries tests navigation with no entries
func TestNavigationWithEmptyEntries(t *testing.T) {
	m := NewModel()
	m.view = "entries"
	m.entries = []models.Entry{} // Empty list
	m.selectedEntry = 0

	// Pressing Enter should not change view when no entries
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newModel.(Model)

	if m.view != "entries" {
		t.Errorf("Expected view to remain 'entries' with empty list, got '%s'", m.view)
	}
}

// TestNavigationCursorBounds tests that cursor stays within bounds
func TestNavigationCursorBounds(t *testing.T) {
	m := NewModel()
	m.view = "entries"
	m.entries = []models.Entry{
		{ID: "1", Title: "Entry 1"},
		{ID: "2", Title: "Entry 2"},
		{ID: "3", Title: "Entry 3"},
	}
	m.selectedEntry = 0

	// Press 'j' multiple times (down key)
	for i := 0; i < 5; i++ {
		newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		m = newModel.(Model)
	}

	// Cursor should not exceed length - 1
	if m.selectedEntry >= len(m.entries) {
		t.Errorf("Cursor exceeded bounds: %d >= %d", m.selectedEntry, len(m.entries))
	}

	// Press 'k' multiple times (up key)
	for i := 0; i < 5; i++ {
		newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
		m = newModel.(Model)
	}

	// Cursor should not go below 0
	if m.selectedEntry < 0 {
		t.Errorf("Cursor went below 0: %d", m.selectedEntry)
	}
}

// TestEditingExistingEntry tests pressing 'i' in entry view to edit
func TestEditingExistingEntry(t *testing.T) {
	m := NewModel()
	m.view = "view_entry"
	m.viewingEntry = models.Entry{
		ID:    "test-id",
		Title: "Test Entry",
		Body:  "Test body with @tag",
		Tags:  []string{"tag"},
	}

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	m = newModel.(Model)

	if m.view != "entry" {
		t.Errorf("Expected view to be 'entry', got '%s'", m.view)
	}

	if !m.editingMode {
		t.Error("Expected editingMode to be true when editing existing entry")
	}

	if m.originalEntryID != "test-id" {
		t.Errorf("Expected originalEntryID to be 'test-id', got '%s'", m.originalEntryID)
	}
}

// TestMarkingEntryForDeletion tests pressing 'd' to mark entry
func TestMarkingEntryForDeletion(t *testing.T) {
	m := NewModel()
	m.view = "view_entry"
	m.viewingEntry = models.Entry{ID: "test-id", Title: "Test"}
	m.markedForDeletion = make(map[string]string)

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(Model)

	// Check if entry was marked (key is just the ID, value is "entry")
	if entryType, exists := m.markedForDeletion["test-id"]; !exists || entryType != "entry" {
		t.Errorf("Expected entry test-id to be marked for deletion with type 'entry', got exists=%v type=%s", exists, entryType)
	}

	// Press 'd' again to unmark
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(Model)

	if _, exists := m.markedForDeletion["test-id"]; exists {
		t.Error("Expected entry to be unmarked after pressing 'd' again")
	}
}

// TestQuitFromEntryList tests pressing 'q' to quit
func TestQuitFromEntryList(t *testing.T) {
	m := NewModel()
	m.view = "entries"

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	// Should return tea.Quit command
	if cmd == nil {
		t.Error("Expected quit command to be returned")
	}
}
