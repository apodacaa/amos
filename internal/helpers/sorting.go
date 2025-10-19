package helpers

import "github.com/apodacaa/amos/internal/models"

// SortTodosForDisplay sorts todos the same way as the UI display:
// open first, then by position (lower first), then newest first
func SortTodosForDisplay(todos []models.Todo) []models.Todo {
	sorted := make([]models.Todo, len(todos))
	copy(sorted, todos)

	// Bubble sort with three-level comparison
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			// First: open todos before done todos
			if sorted[i].Status == "done" && sorted[j].Status == "open" {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			} else if sorted[i].Status == sorted[j].Status {
				// Second: sort by position (lower position first)
				if sorted[i].Position != sorted[j].Position {
					if sorted[j].Position < sorted[i].Position {
						sorted[i], sorted[j] = sorted[j], sorted[i]
					}
				} else {
					// Third: sort by created date (newest first)
					if sorted[j].CreatedAt.After(sorted[i].CreatedAt) {
						sorted[i], sorted[j] = sorted[j], sorted[i]
					}
				}
			}
		}
	}

	return sorted
}

// NormalizeTodoPositions renumbers todo positions to 0, 1, 2, 3...
// This is used after sorting to ensure positions are sequential
func NormalizeTodoPositions(todos []models.Todo) []models.Todo {
	normalized := make([]models.Todo, len(todos))
	copy(normalized, todos)

	for i := range normalized {
		normalized[i].Position = i
	}

	return normalized
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
