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

// ApplyTodoFilters applies date and tag filters to todos in a single call.
// This centralizes the common filtering pattern used throughout the codebase.
//
// Design: Date filter first, then tags.
// Returns the original list if all filters are empty.
func ApplyTodoFilters(todos []models.Todo, dateFilter string, tagFilters []string) []models.Todo {
	// Apply date filter
	filtered := FilterTodosByDateRange(todos, dateFilter)

	// Apply tag filter to date-filtered results
	filtered = FilterTodosByTags(filtered, tagFilters)

	return filtered
}
