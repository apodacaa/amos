package helpers

import "github.com/apodacaa/amos/internal/models"

// ApplyEntryFilters applies both date and tag filters to entries in a single call.
// This centralizes the common filtering pattern used throughout the codebase.
//
// Design: Date filter is applied first, then tag filter (order doesn't matter for AND logic).
// Returns the original list if both filters are empty.
func ApplyEntryFilters(entries []models.Entry, dateFilter string, tagFilters []string) []models.Entry {
	// Apply date filter first
	filtered := FilterEntriesByDateRange(entries, dateFilter)

	// Apply tag filter to date-filtered results
	filtered = FilterEntriesByTags(filtered, tagFilters)

	return filtered
}

// ApplyTodoFilters applies visibility, date, and tag filters to todos in a single call.
// This centralizes the common filtering pattern used throughout the codebase.
//
// Design: Visibility filter first (removes done todos if hidden), then date, then tags.
// Returns the original list if all filters are empty/disabled.
func ApplyTodoFilters(todos []models.Todo, dateFilter string, tagFilters []string, hideCompleted bool) []models.Todo {
	// Apply visibility filter first (hide completed if enabled)
	filtered := FilterTodosByVisibility(todos, hideCompleted)

	// Apply date filter
	filtered = FilterTodosByDateRange(filtered, dateFilter)

	// Apply tag filter to date-filtered results
	filtered = FilterTodosByTags(filtered, tagFilters)

	return filtered
}

// FilterTodosByVisibility filters out done todos if hideCompleted is true.
// Returns the original slice if hideCompleted is false.
func FilterTodosByVisibility(todos []models.Todo, hideCompleted bool) []models.Todo {
	if !hideCompleted {
		return todos
	}
	filtered := make([]models.Todo, 0, len(todos))
	for _, todo := range todos {
		if todo.Status != "done" {
			filtered = append(filtered, todo)
		}
	}
	return filtered
}
