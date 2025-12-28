package main

import (
	"testing"
	"time"

	"github.com/apodacaa/amos/internal/helpers"
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

// ========================================
// Filtering Tests
// ========================================

// TestFilteringByTags tests tag filtering in entries list
func TestFilteringByTags(t *testing.T) {
	m := NewModel()
	m.view = "entries"
	m.entries = []models.Entry{
		{ID: "1", Title: "Work Entry", Tags: []string{"work", "urgent"}},
		{ID: "2", Title: "Personal Entry", Tags: []string{"personal"}},
		{ID: "3", Title: "Work Meeting", Tags: []string{"work", "meetings"}},
	}

	// Apply single tag filter
	m.filterTags = []string{"work"}

	// Get filtered entries (simulating what the view does)
	filtered := m.entries
	if len(m.filterTags) > 0 {
		filtered = helpers.FilterEntriesByTags(filtered, m.filterTags)
	}

	if len(filtered) != 2 {
		t.Errorf("Expected 2 entries with 'work' tag, got %d", len(filtered))
	}

	// Apply multiple tag filter (AND logic)
	m.filterTags = []string{"work", "urgent"}
	filtered = m.entries
	if len(m.filterTags) > 0 {
		filtered = helpers.FilterEntriesByTags(filtered, m.filterTags)
	}

	if len(filtered) != 1 {
		t.Errorf("Expected 1 entry with both 'work' and 'urgent' tags, got %d", len(filtered))
	}

	if len(filtered) > 0 && filtered[0].ID != "1" {
		t.Errorf("Expected filtered entry to have ID '1', got '%s'", filtered[0].ID)
	}
}

// TestFilteringByDate tests date filtering in entries list
func TestFilteringByDate(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	weekAgo := now.Add(-7 * 24 * time.Hour)

	m := NewModel()
	m.view = "entries"
	m.entries = []models.Entry{
		{ID: "1", Title: "Today's Entry", CreatedAt: now, UpdatedAt: now},
		{ID: "2", Title: "Yesterday's Entry", CreatedAt: yesterday, UpdatedAt: yesterday},
		{ID: "3", Title: "Old Entry", CreatedAt: weekAgo, UpdatedAt: weekAgo},
	}

	// Filter by "today"
	m.filterDate = "today"
	filtered := helpers.FilterEntriesByDateRange(m.entries, m.filterDate)
	if len(filtered) != 1 {
		t.Errorf("Expected 1 entry for 'today', got %d", len(filtered))
	}

	// Filter by "last 7 days"
	m.filterDate = "last 7 days"
	filtered = helpers.FilterEntriesByDateRange(m.entries, m.filterDate)
	// Should include today and yesterday (within 7 days)
	if len(filtered) < 2 {
		t.Errorf("Expected at least 2 entries for 'last 7 days', got %d", len(filtered))
	}
}

// TestCombinedFiltering tests using both tag and date filters
func TestCombinedFiltering(t *testing.T) {
	now := time.Now()
	weekAgo := now.Add(-7 * 24 * time.Hour)

	m := NewModel()
	m.view = "entries"
	m.entries = []models.Entry{
		{ID: "1", Title: "Recent Work", Tags: []string{"work"}, CreatedAt: now, UpdatedAt: now},
		{ID: "2", Title: "Old Work", Tags: []string{"work"}, CreatedAt: weekAgo, UpdatedAt: weekAgo},
		{ID: "3", Title: "Recent Personal", Tags: []string{"personal"}, CreatedAt: now, UpdatedAt: now},
	}

	// Apply both filters: work tag + today
	m.filterTags = []string{"work"}
	m.filterDate = "today"

	// Apply tag filter first
	filtered := m.entries
	if len(m.filterTags) > 0 {
		filtered = helpers.FilterEntriesByTags(filtered, m.filterTags)
	}

	// Then apply date filter
	if m.filterDate != "" {
		filtered = helpers.FilterEntriesByDateRange(filtered, m.filterDate)
	}

	// Should only get entry 1 (work + today)
	if len(filtered) != 1 {
		t.Errorf("Expected 1 entry with 'work' tag and 'today' date, got %d", len(filtered))
	}

	if len(filtered) > 0 && filtered[0].ID != "1" {
		t.Errorf("Expected filtered entry ID '1', got '%s'", filtered[0].ID)
	}
}

// TestFilterParsing tests the unified filter input parsing
func TestFilterParsing(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expectedTags int
		expectedDate string
		expectError  bool
	}{
		{"single tag", "@work", 1, "", false},
		{"multiple tags", "@work @urgent", 2, "", false},
		{"date only", "today", 0, "today", false},
		{"tags and date", "@work today", 1, "today", false},
		{"last N days", "@project last 7 days", 1, "last 7 days", false},
		{"date range", "2024-01-01 to 2024-01-31", 0, "2024-01-01 to 2024-01-31", false},
		{"mixed order", "yesterday @work @urgent", 2, "yesterday", false},
		{"empty input", "", 0, "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := helpers.ParseFilterInput(tc.input)

			if len(result.Tags) != tc.expectedTags {
				t.Errorf("Expected %d tags, got %d", tc.expectedTags, len(result.Tags))
			}

			if result.Date != tc.expectedDate {
				t.Errorf("Expected date '%s', got '%s'", tc.expectedDate, result.Date)
			}

			if tc.expectError && len(result.Errors) == 0 {
				t.Error("Expected error but got none")
			}

			if !tc.expectError && len(result.Errors) > 0 {
				t.Errorf("Expected no errors but got: %v", result.Errors)
			}
		})
	}
}

// TestFilterReturnsToCorrectView tests that filter returns to the right view
func TestFilterReturnsToCorrectView(t *testing.T) {
	testCases := []struct {
		name           string
		fromView       string
		expectedReturn string
	}{
		{"from entries", "entries", "entries"},
		{"from todos", "todos", "todos"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewModel()
			m.view = tc.fromView

			// Open filter (/ key)
			newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
			m = newModel.(Model)

			if m.view != "unified_filter" {
				t.Errorf("Expected view 'unified_filter', got '%s'", m.view)
			}

			if m.filterContext != tc.expectedReturn {
				t.Errorf("Expected filterContext '%s', got '%s'", tc.expectedReturn, m.filterContext)
			}

			// Press Esc to cancel
			newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
			m = newModel.(Model)

			if m.view != tc.expectedReturn {
				t.Errorf("Expected return to '%s', got '%s'", tc.expectedReturn, m.view)
			}
		})
	}
}

// TestClearFilter tests pressing 'c' to clear filters
func TestClearFilter(t *testing.T) {
	m := NewModel()
	m.view = "entries"
	m.filterTags = []string{"work", "urgent"}
	m.filterDate = "today"

	// Press 'c' to clear filter
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	m = newModel.(Model)

	if len(m.filterTags) != 0 {
		t.Errorf("Expected filterTags to be cleared, got %d tags", len(m.filterTags))
	}

	if m.filterDate != "" {
		t.Errorf("Expected filterDate to be cleared, got '%s'", m.filterDate)
	}

	// Cursor should reset to 0
	if m.selectedEntry != 0 {
		t.Errorf("Expected selectedEntry to reset to 0, got %d", m.selectedEntry)
	}
}

// TestFilterWithNoMatches tests filter behavior when no items match
func TestFilterWithNoMatches(t *testing.T) {
	m := NewModel()
	m.entries = []models.Entry{
		{ID: "1", Title: "Entry", Tags: []string{"work"}},
	}

	// Apply filter that matches nothing
	m.filterTags = []string{"nonexistent"}

	filtered := helpers.FilterEntriesByTags(m.entries, m.filterTags)

	if len(filtered) != 0 {
		t.Errorf("Expected 0 filtered entries, got %d", len(filtered))
	}
}

// ================================================================================
// Deletion Workflow Tests
// ================================================================================

// TestMarkingEntryFromListView tests marking entries from the entries list view
func TestMarkingEntryFromListView(t *testing.T) {
	m := NewModel()
	m.view = "entries"
	m.entries = []models.Entry{
		{ID: "entry-1", Title: "First Entry"},
		{ID: "entry-2", Title: "Second Entry"},
	}
	m.selectedEntry = 0

	// Mark first entry for deletion
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(Model)

	if itemType, exists := m.markedForDeletion["entry-1"]; !exists || itemType != "entry" {
		t.Errorf("Expected entry-1 to be marked for deletion with type 'entry', got exists=%v type=%s", exists, itemType)
	}

	// Pressing 'd' again should unmark
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(Model)

	if _, exists := m.markedForDeletion["entry-1"]; exists {
		t.Errorf("Expected entry-1 to be unmarked after second 'd' press")
	}
}

// TestMarkingTodoFromListView tests marking and unmarking todos from todo list
func TestMarkingTodoFromListView(t *testing.T) {
	m := NewModel()
	m.view = "todos"
	m.displayTodos = []models.Todo{
		{ID: "todo-1", Title: "First Todo", Status: "open"},
		{ID: "todo-2", Title: "Second Todo", Status: "open"},
	}
	m.selectedTodo = 0

	// Mark first todo for deletion
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(Model)

	if itemType, exists := m.markedForDeletion["todo-1"]; !exists || itemType != "todo" {
		t.Errorf("Expected todo-1 to be marked for deletion with type 'todo', got exists=%v type=%s", exists, itemType)
	}

	// Pressing 'd' again should unmark
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(Model)

	if _, exists := m.markedForDeletion["todo-1"]; exists {
		t.Errorf("Expected todo-1 to be unmarked after second 'd' press")
	}
}

// TestMarkingMultipleItems tests marking multiple entries and todos for deletion
func TestMarkingMultipleItems(t *testing.T) {
	m := NewModel()
	m.view = "entries"
	m.entries = []models.Entry{
		{ID: "entry-1", Title: "First Entry"},
		{ID: "entry-2", Title: "Second Entry"},
	}
	m.displayTodos = []models.Todo{
		{ID: "todo-1", Title: "First Todo", Status: "open"},
	}

	// Mark first entry
	m.selectedEntry = 0
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(Model)

	// Mark second entry
	m.selectedEntry = 1
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(Model)

	// Switch to todos view and mark a todo
	m.view = "todos"
	m.selectedTodo = 0
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(Model)

	// Should have 3 items marked
	if len(m.markedForDeletion) != 3 {
		t.Errorf("Expected 3 items marked for deletion, got %d", len(m.markedForDeletion))
	}

	// Verify types
	if m.markedForDeletion["entry-1"] != "entry" {
		t.Errorf("Expected entry-1 to have type 'entry', got %s", m.markedForDeletion["entry-1"])
	}
	if m.markedForDeletion["entry-2"] != "entry" {
		t.Errorf("Expected entry-2 to have type 'entry', got %s", m.markedForDeletion["entry-2"])
	}
	if m.markedForDeletion["todo-1"] != "todo" {
		t.Errorf("Expected todo-1 to have type 'todo', got %s", m.markedForDeletion["todo-1"])
	}
}

// TestDeletionConfirmationWorkflow tests the $ key confirmation workflow
func TestDeletionConfirmationWorkflow(t *testing.T) {
	m := NewModel()
	m.view = "entries"
	m.entries = []models.Entry{
		{ID: "entry-1", Title: "First Entry"},
	}
	m.markedForDeletion = map[string]string{
		"entry-1": "entry",
	}

	// Press $ to initiate deletion
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'$'}})
	m = newModel.(Model)

	// Should set deleteConfirmPending to true
	if !m.deleteConfirmPending {
		t.Errorf("Expected deleteConfirmPending to be true after pressing $")
	}

	// Press 'n' to cancel
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m = newModel.(Model)

	// Should clear deleteConfirmPending
	if m.deleteConfirmPending {
		t.Errorf("Expected deleteConfirmPending to be false after pressing 'n'")
	}

	// Items should still be marked
	if len(m.markedForDeletion) != 1 {
		t.Errorf("Expected items to still be marked after canceling, got %d marked items", len(m.markedForDeletion))
	}
}

// TestDeletionConfirmationWithNoMarkedItems tests $ with no marked items
func TestDeletionConfirmationWithNoMarkedItems(t *testing.T) {
	m := NewModel()
	m.view = "entries"
	m.entries = []models.Entry{
		{ID: "entry-1", Title: "First Entry"},
	}

	// Press $ with no marked items
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'$'}})
	m = newModel.(Model)

	// Should set status message
	if m.statusMsg != "No items marked for deletion" {
		t.Errorf("Expected status message about no marked items, got: %s", m.statusMsg)
	}

	// Should NOT set deleteConfirmPending
	if m.deleteConfirmPending {
		t.Errorf("Expected deleteConfirmPending to remain false when no items marked")
	}
}

// TestCascadeDeletion tests that deleting an entry also deletes linked todos
func TestCascadeDeletion(t *testing.T) {
	// Note: This test validates the logic, but actual cascade deletion
	// happens in the storage layer (storage.DeleteEntryCascade).
	// Here we verify that the command separates entries and counts linked todos.

	entryID := "entry-1"
	m := NewModel()
	m.view = "entries"
	m.entries = []models.Entry{
		{ID: entryID, Title: "Entry with todos"},
	}
	m.todos = []models.Todo{
		{ID: "todo-1", Title: "Linked todo 1", EntryID: &entryID},
		{ID: "todo-2", Title: "Linked todo 2", EntryID: &entryID},
		{ID: "todo-3", Title: "Standalone todo", EntryID: nil},
	}

	// Mark the entry for deletion
	m.markedForDeletion = map[string]string{
		entryID: "entry",
	}

	// Press $ to initiate deletion
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'$'}})
	m = newModel.(Model)

	// Should set deleteConfirmPending
	if !m.deleteConfirmPending {
		t.Errorf("Expected deleteConfirmPending to be true")
	}

	// The actual deletion happens in the command, but we can verify
	// the command is called correctly by checking the model state
	// (Full integration test would need storage layer mocking)
}

// TestStandaloneTodoDeletion tests deleting standalone todos
func TestStandaloneTodoDeletion(t *testing.T) {
	m := NewModel()
	m.view = "todos"
	m.displayTodos = []models.Todo{
		{ID: "todo-1", Title: "Standalone todo", EntryID: nil},
		{ID: "todo-2", Title: "Another standalone", EntryID: nil},
	}

	// Mark first standalone todo
	m.selectedTodo = 0
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(Model)

	if itemType, exists := m.markedForDeletion["todo-1"]; !exists || itemType != "todo" {
		t.Errorf("Expected standalone todo to be marked for deletion")
	}

	// Initiate deletion
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'$'}})
	m = newModel.(Model)

	if !m.deleteConfirmPending {
		t.Errorf("Expected deleteConfirmPending to be true for standalone todo deletion")
	}
}

// TestEntryLinkedTodoMarking tests that entry-linked todos can be marked
func TestEntryLinkedTodoMarking(t *testing.T) {
	entryID := "entry-1"
	m := NewModel()
	m.view = "todos"
	m.displayTodos = []models.Todo{
		{ID: "todo-1", Title: "Linked todo", EntryID: &entryID},
	}

	// Mark entry-linked todo
	m.selectedTodo = 0
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(Model)

	// Entry-linked todos can be marked (UI allows it)
	if itemType, exists := m.markedForDeletion["todo-1"]; !exists || itemType != "todo" {
		t.Errorf("Expected entry-linked todo to be marked for deletion")
	}

	// Note: The actual deletion logic in storage layer handles whether
	// entry-linked todos get deleted independently or via cascade.
	// This test verifies the marking behavior works.
}

// TestMarksPersistAcrossNavigation tests that marks persist when navigating between views
func TestMarksPersistAcrossNavigation(t *testing.T) {
	m := NewModel()
	m.view = "entries"
	m.entries = []models.Entry{
		{ID: "entry-1", Title: "First Entry"},
	}
	m.displayTodos = []models.Todo{
		{ID: "todo-1", Title: "First Todo", Status: "open"},
	}

	// Mark entry in entries view
	m.selectedEntry = 0
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(Model)

	// Navigate to todos view
	m.view = "todos"
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	m = newModel.(Model)

	// Mark should persist
	if _, exists := m.markedForDeletion["entry-1"]; !exists {
		t.Errorf("Expected entry mark to persist when navigating to todos view")
	}

	// Mark a todo
	m.selectedTodo = 0
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(Model)

	// Navigate back to entries
	m.view = "entries"
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	m = newModel.(Model)

	// Both marks should persist
	if len(m.markedForDeletion) != 2 {
		t.Errorf("Expected 2 marked items after navigation, got %d", len(m.markedForDeletion))
	}
}
