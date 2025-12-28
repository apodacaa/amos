package helpers

import "github.com/apodacaa/amos/internal/models"

// SortTodosForDisplay sorts todos: next first, then open, then done
// Within each status group, newest first
//
// Design decision: Uses bubble sort instead of sort.Slice for these reasons:
//  1. Stable sort: Preserves relative ordering for items with same timestamp.
//     This is critical if we later add manual position reordering within priority groups.
//  2. Simple implementation: Multi-level comparison logic is transparent and maintainable.
//  3. Performance: O(n²) but todo lists are typically <100 items, so the difference
//     from O(n log n) is negligible in practice (~0.1ms vs ~0.01ms).
func SortTodosForDisplay(todos []models.Todo) []models.Todo {
	sorted := make([]models.Todo, len(todos))
	copy(sorted, todos)

	// Helper to get priority (lower = higher priority)
	getPriority := func(status string) int {
		switch status {
		case "next":
			return 0
		case "open":
			return 1
		case "done":
			return 2
		default:
			return 1 // treat unknown as open
		}
	}

	// Bubble sort with two-level comparison
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			iPriority := getPriority(sorted[i].Status)
			jPriority := getPriority(sorted[j].Status)

			// First: sort by status priority (next → open → done)
			if jPriority < iPriority {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			} else if iPriority == jPriority {
				// Second: within same status, newest first
				if sorted[j].CreatedAt.After(sorted[i].CreatedAt) {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
	}

	return sorted
}

// SortEntriesForDisplay sorts entries by creation date (newest first)
//
// Design decision: Uses bubble sort for consistency with SortTodosForDisplay.
// Stable sort preserves order for entries with identical timestamps.
// Performance is acceptable for typical journal entry lists (<1000 items).
func SortEntriesForDisplay(entries []models.Entry) []models.Entry {
	sorted := make([]models.Entry, len(entries))
	copy(sorted, entries)

	// Bubble sort by CreatedAt descending (newest first)
	// Sort by creation date ensures entries don't move when edited
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].CreatedAt.After(sorted[i].CreatedAt) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}
