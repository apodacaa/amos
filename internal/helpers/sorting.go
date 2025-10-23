package helpers

import "github.com/apodacaa/amos/internal/models"

// SortTodosForDisplay sorts todos: next first, then open, then done
// Within each status group, newest first
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

// SortEntriesForDisplay sorts entries by timestamp (newest first)
func SortEntriesForDisplay(entries []models.Entry) []models.Entry {
	sorted := make([]models.Entry, len(entries))
	copy(sorted, entries)

	// Bubble sort by timestamp descending (newest first)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Timestamp.After(sorted[i].Timestamp) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}
